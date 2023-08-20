package thema

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
)

// AssignableTo indicates whether all valid instances of the provided Thema
// schema can be assigned to the provided Go type.
//
// If the provided T is a pointer, it will be dereferenced before verification.
// Double pointers (or any n-pointer > 1) are not allowed.
//
// The provided T must be struct-kinded, as it is a requirement that all Thema
// schemas are of base type struct.
//
//	type MyType struct {
//		MyField string `json:"myfield"`
//	}
//
//	AssignableTo(sch, &MyType{})
//
// Assignability rules are specified here: https://github.com/grafana/thema/blob/main/docs/invariants.md#go-assignability
func AssignableTo(sch Schema, T any) error {
	rt := sch.Lineage().Runtime()
	rt.rl()
	defer rt.ru()
	return assignable(sch.Underlying().LookupPath(pathSchDef), T)
}

// ErrPointerDepth indicates that a Go type having pointer indirection depth greater than 1, such as
//
//	**struct{ V: string })
//
// was provided to a Thema func that checks assignability, such as [BindType].
var ErrPointerDepth = errors.New("assignability does not support more than one level of pointer indirection")

const scalarKinds = cue.NullKind | cue.BoolKind |
	cue.IntKind | cue.FloatKind | cue.StringKind | cue.BytesKind

func assignable(sch cue.Value, T interface{}) error {
	v := reflect.ValueOf(T)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Ptr {
		return ErrPointerDepth
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("must provide struct-kinded type, got *%s", v.Kind())
	}

	ctx := sch.Context()
	gt := ctx.EncodeType(v.Interface())

	// None of the builtin CUE functions do _quite_ what we want here. In the
	// simple case, we might check subsumption of the Go type by the CUE
	// schema, but that falls down because bounds constraints in CUE may be
	// narrower than pure native Go types (e.g. string enums), and Go's type
	// system gives us no choice but to accept that case.
	//
	// We also can't check subsumption of the CUE schema by the Go type, because
	// the addition of `*null |` to every field with a Go pointer type means
	// that many CUE schema values won't be instances of the Go type. (This is
	// something that can be changed in the logic of EncodeType, hopefully.)
	//
	// More importantly, checking subsumption of the schema type by the Go type
	// will not verify that all the schema fields are satisfied/exist - only
	// that the ones that do exist in the Go type are also present and subsumed
	// in the schema type. And that's not even considering optional fields. This
	// makes builtin cue.Value.Subsume() insufficient for our purposes.
	//
	// Forcing the Go type to be closed would plausibly help with all of this,
	// except erroneous nulls. But the above considerations force us to roll our
	// own definition of assignability, at least for now.

	// Errors, keyed by string
	errs := make(assignErrs)

	type checkfn func(gval, sval cue.Value, p cue.Path)
	var check, checkstruct, checklist, checkscalar checkfn

	check = func(gval, sval cue.Value, p cue.Path) {
		sk, gk := sval.IncompleteKind(), gval.IncompleteKind()
		schemaHadNull := sk&cue.NullKind != 0
		sk &^= cue.NullKind
		gk &^= cue.NullKind

		ogval := gval
		// At least for now, we have to deal with these unhelpful *null
		// appearing in the encoding of pointer types.
		//
		// We can't do the same for the schema side, because this null stripper
		// relies on the fact that all actual Go type declarations will come across
		// as a single value, without any disjunctions.
		gval = stripLeadNull(gval)

		// strict equality _might_ be too restrictive? But it's better to start there
		if sk != gk && gk != (cue.TopKind^cue.NullKind) {
			errs[p.String()] = fmt.Errorf("%s: is kind %s in schema, but kind %s in Go type", p, sk, gk)
			return
		} else if gk == (cue.TopKind ^ cue.NullKind) {
			// Escape hatch for a Go interface{}/any
			return
		}
		op, _ := sval.Expr()

		switch sk {
		case cue.ListKind:
			checklist(gval, sval, p)
		case cue.FloatKind, cue.IntKind, cue.StringKind, cue.BytesKind, cue.BoolKind, cue.NumberKind:
			if schemaHadNull {
				checkscalar(ogval, sval, p)
			} else {
				checkscalar(gval, sval, p)
			}
		case cue.StructKind:
			if op == cue.OrOp {
				errs[p.String()] = fmt.Errorf("%s: contains disjunction over struct types, but Go type is not any", p)
				return
			}
			checkstruct(gval, sval, p)
		case cue.NullKind:
			errs[p.String()] = fmt.Errorf("%s: null is not permitted in schema; express optionality with ?", p)
		default:
			if sk&scalarKinds == sk {
				errs[p.String()] = fmt.Errorf(
					"%s: allows multiple basic CUE types %s, may only correspond to interface{}/any", p, sk,
				)
				return
			}
			panic(fmt.Sprintf("unhandled kind %s", sk))
		}
	}

	checkstruct = func(ogval, osval cue.Value, p cue.Path) {
		ss, gmap := structToSlice(osval), structToMap(ogval)

		// The returned cue.Value appears to differ depending on whether it's
		// accessed through an iterator vs. LookupPath. This matters(?) in the
		// context of doing things like comparing kinds.
		for _, vp := range ss {
			sval, p := vp.Value, cue.MakePath(append(p.Selectors(), vp.Path.Selectors()...)...)
			// Optional() gives us paths that are either optional or not, which
			// seems reasonable at least until we really formally define this
			// relation
			// gval, exists := gmap[vp.Path.Optional().String()].Value
			gvp, exists := gmap[vp.Path.String()]

			// TODO replace these one-offs with formalized error types
			if !exists {
				errs[p.String()] = fmt.Errorf("%s: field present in schema, absent from Go type", p)
				continue
			}

			gval := gvp.Value
			// Remove the paths from the map for leftover checking later
			delete(gmap, vp.Path.String())
			// delete(gmap, vp.Path.Optional().String())

			check(gval, sval, p)
		}

		for _, vp := range gmap {
			errs[vp.Path.String()] = fmt.Errorf("%s: field present in Go type, absent from schema", vp.Path.String())
		}
	}

	checklist = func(gval, sval cue.Value, p cue.Path) {
		// If the schema is an open list with a default value, sval.Len() will produce a bottom value because it can't
		// know which side of the disjunction to pick. There could be many branches on this disjunction, though. Only
		// allow two, and only if one is a default.
		if _, has := sval.Default(); has {
			// Bewildering: calling Expr() on a list disjunction with a default value appears to
			// just drop the marked/default value, and not indicate an OrOp at all. This is handy,
			// but unexpected behavior, and it feels dangerous to rely on.
			_, evals := sval.Expr()
			if len(evals) != 1 {
				errs[p.String()] = fmt.Errorf("%s: schema is a complex disjunction of list types, may only correspond to interface{}/any", p)
				return
			}
			sval = evals[0]
		}

		glen, slen := gval.Len(), sval.Len()
		// Ensure alignment of list openness/closedness
		if glen.IsConcrete() != slen.IsConcrete() {
			if slen.IsConcrete() {
				errs[p.String()] = fmt.Errorf("%s: list is closed in schema, Go type must be an array, not slice", p)
			} else {
				errs[p.String()] = fmt.Errorf("%s: list is open in schema, Go type must be a slice, not array", p)
			}
			return
		}

		if err := glen.Subsume(slen); err != nil {
			// should be unreachable?
			errs[p.String()] = fmt.Errorf(
				"%s: incompatible list lengths in schema (%s) and Go type (%s)", p, slen, glen,
			)
			return
		}

		// Vars capturing the values we'll eventually need to walk down into and check
		var los, log cue.Value

		// Whether the list is open or closed on the schema side, it may contain
		// a fixed set of values. If they exist, we have to ensure those
		// elements are the same, as Go's type system can't express type
		// variance across an array or list. Of course, checking "sameness" of
		// incomplete values isn't (?) trivial. Mutual subsume...
		iter, _ := sval.List()
		// Variables for retaining list iteration state
		var lastval cue.Value
		var nonempty bool
		if nonempty = iter.Next(); nonempty {
			lastval = iter.Value()
			for iter.Next() {
				los = iter.Value() // it's fine to just keep updating the reference
				// Failures indicate the CUE schema is unrepresentable in Go.
				// That's the kind of thing we'd likely prefer to know/have in
				// some more universal place.
				lerr, rerr := lastval.Subsume(los, cue.Schema()), los.Subsume(lastval, cue.Schema())
				if lerr != nil || rerr != nil {
					errs[p.String()] = fmt.Errorf("%s: schema is list with multiple types, not representable in Go", p)
					return
				}

				lastval = iter.Value()
			}
		}

		if glen.IsConcrete() {
			// It has already been established that all the list elements are
			// the same, so taking the first one is sufficient. Use an iter,
			// since we don't trust LookupPath.
			iter, _ = gval.List()
			if !iter.Next() {
				// Only happens if list is empty, which means list is empty on both sides.
				// Weird, but not illegal
				return
			}
			log = iter.Value()
		} else {
			los = sval.LookupPath(cue.MakePath(cue.AnyIndex))
			// If there were actual list elements, make sure the AnyIndex type is the same as them, too.
			if nonempty {
				lerr, rerr := lastval.Subsume(los, cue.Schema()), los.Subsume(lastval, cue.Schema())
				if lerr != nil || rerr != nil {
					errs[p.String()] = fmt.Errorf("%s: schema is list of multiple types; not representable in Go", p)
					return
				}
			}
			log = gval.LookupPath(cue.MakePath(cue.AnyIndex))
		}
		p = cue.MakePath(append(p.Selectors(), cue.AnyIndex)...)
		check(log, los, p)
	}

	checkscalar = func(gval, sval cue.Value, p cue.Path) {
		// Because the CUE types can have narrower bounds, and we're
		// really interested in whether all valid schema instances will
		// be assignable to the Go type, we have to see if the Go type
		// subsumes the schema, rather than the more intuitive check
		// that the schema subsumes the Go type.
		if err := gval.Subsume(sval, cue.Schema()); err != nil {
			errs[p.String()] = fmt.Errorf("%s: schema type %v not subsumed by Go type %v", p, sval, gval)
		}
	}

	// Walk down the whole struct tree
	check(gt, sch, cue.MakePath())

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func stripLeadNull(v cue.Value) cue.Value {
	if op, vals := v.Expr(); op == cue.OrOp {
		// Walk over the vals, because there may be more than one null (e.g. omitempty +
		// slice/map type).
		for i := 0; i < len(vals); i++ {
			if vals[i].Null() != nil {
				return vals[i]
			}
		}
	}
	return v
}

type valpath struct {
	Path  cue.Path
	Value cue.Value
}

type structSlice []valpath

func structToMap(v cue.Value) map[string]valpath {
	m := make(map[string]valpath)
	iter, err := v.Fields(cue.Optional(true))
	if err != nil {
		panic(err)
	}

	for iter.Next() {
		vp := valpath{
			Path:  cue.MakePath(iter.Selector()),
			Value: iter.Value(),
		}
		// fmt.Printf("sm %v %#v\n", iter.Selector(), iter.Value())
		m[iter.Selector().String()] = vp
		// m[iter.Selector().Optional().String()] = vp
	}

	return m
}

func structToSlice(v cue.Value) structSlice {
	var ss structSlice
	iter, err := v.Fields(cue.Optional(true))
	if err != nil {
		panic(err)
	}

	for iter.Next() {
		// fmt.Printf("ss %v %#v\n", iter.Selector(), iter.Value())
		ss = append(ss, valpath{
			Path:  cue.MakePath(iter.Selector()),
			Value: iter.Value(),
		})
	}

	return ss
}

type assignErrs map[string]error

func (m assignErrs) Error() string {
	var buf bytes.Buffer

	var i int
	for _, err := range m {
		if i == len(m)-1 {
			fmt.Fprint(&buf, err)
		} else {
			fmt.Fprint(&buf, err, "\n")
		}
		i++
	}

	return (&buf).String()
}
