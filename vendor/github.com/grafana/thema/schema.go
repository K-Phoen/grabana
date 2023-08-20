package thema

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
)

var (
	_ Schema                = &schemaDef{}
	_ TypedSchema[Assignee] = &unaryTypedSchema[Assignee]{}
)

var (
	pathSchDef   = cue.MakePath(cue.Hid("_#schema", "github.com/grafana/thema"))
	pathExamples = cue.MakePath(cue.Str("examples"))
	pathSch      = cue.MakePath(cue.Str("schema"))
	pathJoin     = cue.MakePath(cue.Hid("_join", "github.com/grafana/thema"))
)

// schemaDef represents a single #SchemaDef, with a backlink to its containing
// #Lineage.
type schemaDef struct {
	// ref holds a reference to the entire #SchemaDef object.
	ref cue.Value

	// def holds a reference to #SchemaDef._#schema. See those docs.
	def cue.Value

	// v is the version of this schema.
	v SyntacticVersion

	lin *baseLineage
}

// Examples returns the set of examples of this schema defined in the original
// lineage. The string key is the name given to the example.
func (sch *schemaDef) Examples() map[string]*Instance {
	examplesNode := sch.Underlying().LookupPath(pathExamples)
	it, err := examplesNode.Fields()
	if err != nil {
		panic(err)
	}

	examples := make(map[string]*Instance)
	for it.Next() {
		label := it.Selector().String()
		examples[label] = &Instance{
			valid: true,
			raw:   it.Value(),
			name:  label,
			sch:   sch,
		}
	}

	return examples
}

func (sch *schemaDef) rt() *Runtime {
	return sch.Lineage().Runtime()
}

// Validate checks that the provided data is valid with respect to the
// schema. If valid, the data is wrapped in an Instance and returned.
// Otherwise, a nil Instance is returned along with an error detailing the
// validation failure.
//
// While Validate takes a cue.Value, this is only to avoid having to trigger
// the translation internally; input values must be concrete. To use
// incomplete CUE values with Thema schemas, prefer working directly in CUE,
// or if you must, rely on Underlying().
func (sch *schemaDef) Validate(data cue.Value) (*Instance, error) {
	sch.rt().rl()
	defer sch.rt().ru()
	// TODO which approach is actually the right one, unify or subsume? ugh
	// err := sch.raw.Subsume(data, cue.All(), cue.Raw())
	// if err != nil {
	// 	return nil, err
	// 	// return nil, mungeValidateErr(err, sch)
	// }

	x := sch.def.Unify(data)

	// The cue.Concrete(true) option ensure that Concrete all values resulting
	// from the unification of the schema and data are concrete.
	// ie: every field defined by the schema has a concrete value associated to it,
	// and no required field was omitted.
	if err := x.Validate(cue.Concrete(true)); err != nil {
		return nil, mungeValidateErr(err, sch)
	}

	return &Instance{
		valid: true,
		raw:   data,
		sch:   sch,
		name:  "", // FIXME how are we getting this out?
	}, nil
}

// Successor returns the next schema in the lineage, or nil if it is the last schema.
func (sch *schemaDef) Successor() Schema {
	if s := sch.successor(); s != nil {
		return s
	}
	return nil
}

func (sch *schemaDef) successor() *schemaDef {
	if sch.lin.allv[len(sch.lin.allv)-1] == sch.v {
		return nil
	}

	succv := sch.lin.allv[searchSynv(sch.lin.allv, sch.v)+1]
	return sch.lin.schema(succv)
}

// Predecessor returns the previous schema in the lineage, or nil if it is the first schema.
func (sch *schemaDef) Predecessor() Schema {
	if s := sch.predecessor(); s != nil {
		return s
	}
	return nil
}

func (sch *schemaDef) predecessor() *schemaDef {
	if sch.v == synv() {
		return nil
	}

	predv := sch.lin.allv[searchSynv(sch.lin.allv, sch.v)-1]
	return sch.lin.schema(predv)
}

// LatestInMajor returns the Schema with the newest (largest) minor version
// within this Schema's major version. If the receiver Schema is the latest, it
// will return itself.
func (sch *schemaDef) LatestInMajor() Schema {
	return sch.lin.allsch[searchSynv(sch.lin.allv, SyntacticVersion{sch.v[0] + 1, 0})]
}

// Underlying returns the cue.Value that represents the underlying CUE #SchemaDef.
//
// The #SchemaDef is not directly helpful for most use cases. But useful values are easily
// accessed by calling [cue.Value.LookupPath] on the returned value:
//   - "schema": the literal schema definition provided by the user.
//   - "_#schema": the user-provided schema, unified with the lineage joinSchema and
//     recursively closed.
func (sch *schemaDef) Underlying() cue.Value {
	return sch.ref
}

// Version returns the schema's version number.
func (sch *schemaDef) Version() SyntacticVersion {
	return sch.v
}

// Lineage returns the lineage that contains this schema.
func (sch *schemaDef) Lineage() Lineage {
	return sch.lin
}

func (sch *schemaDef) _schema() {}

// BindType produces a [TypedSchema], given a [Schema] that is [AssignableTo]
// the [Assignee] type parameter T. T must be struct-kinded, and at most one
// level of pointer indirection is allowed.
//
// An error is returned if the provided Schema is not assignable to the given
// struct type.
func BindType[T Assignee](sch Schema, t T) (TypedSchema[T], error) {
	if err := AssignableTo(sch, t); err != nil {
		return nil, err
	}

	tsch := &unaryTypedSchema[T]{
		Schema: sch,
	}

	// Verify that there are no problematic errors emitted from decoding.
	if err := sch.Underlying().LookupPath(pathSchDef).Decode(t); err != nil {
		// Because assignability has already been established, the only errors here
		// _should_ be those arising from schema fields without concrete defaults. But
		// to avoid swallowing other error types, try to filter out those from the list
		// that aren't relevant for decoding purposes. This means we're choosing false
		// negatives over false positives.
		var actual errors.Error
		for _, e := range errors.Errors(err) {
			// TODO would love a better check, but CUE needs a better error architecture first
			if !strings.Contains(e.Error(), "cannot convert non-concrete") {
				actual = errors.Append(actual, e)
			}
		}

		if len(errors.Errors(actual)) > 0 {
			return nil, actual
		}
	}

	// It's now established that decoding the value is error-free. Ideally, there
	// would be some way of precomputing a trivially copyable value so that we
	// could avoid needing to call any reflection at runtime. However, even if the
	// T parameter is not a pointer type, it could contain a type with a pointer.
	// And there's no way to do that without reflection. So for now, the simplest
	// thing to do is just make a decode call newfn func itself.
	//
	// Given that the constraints on Thema assignable types are narrower than on
	// general Go types CUE can decode onto, there may be some opportunity for a
	// specialized implementation to improve performance - but we'll attempt that
	// iff performance is actually shown to be a problem.
	rt := getLinLib(sch.Lineage())
	tsch.newfn = func() T {
		nt := new(T)
		rt.rl()
		sch.Underlying().LookupPath(pathSchDef).Decode(nt) //nolint:gosec,errcheck
		rt.ru()
		return *nt
	}

	tsch.tlin = &unaryConvLineage[T]{
		Lineage: sch.Lineage(),
		tsch:    tsch,
	}

	return tsch, nil
}

func schemaIs(s1, s2 Schema) bool {
	// TODO will need something smarter here if/when we have more types representing schema
	vs1, is1 := s1.(*schemaDef)
	vs2, is2 := s2.(*schemaDef)
	if !is1 || !is2 {
		panic(fmt.Sprintf("TODO implement schema comparison handler for types %T and %T", s1, s2))
		return false
	}
	return vs1 == vs2
}

type unaryTypedSchema[T Assignee] struct {
	Schema
	newfn func() T
	tlin  ConvergentLineage[T]
}

func (sch *unaryTypedSchema[T]) NewT() T {
	return sch.newfn()
}

func (sch *unaryTypedSchema[T]) is(osch Schema) bool {
	return schemaIs(sch.Schema, osch)
}

func (sch *unaryTypedSchema[T]) ValidateTyped(data cue.Value) (*TypedInstance[T], error) {
	inst, err := sch.Schema.Validate(data)
	if err != nil {
		return nil, err
	}

	return &TypedInstance[T]{
		Instance: inst,
		tsch:     sch,
	}, nil
}

func (sch *unaryTypedSchema[T]) ConvergentLineage() ConvergentLineage[T] {
	return sch.tlin
}
