package thema

import (
	"fmt"

	"cuelang.org/go/cue"
	cerrors "cuelang.org/go/cue/errors"
	"cuelang.org/go/pkg/encoding/json"
	"github.com/cockroachdb/errors"

	terrors "github.com/grafana/thema/errors"
)

// BindInstanceType produces a TypedInstance, given an Instance and a
// TypedSchema derived from its Instance.Schema().
//
// The only possible error occurs if the TypedSchema is not derived from the
// Instance.Schema().
func BindInstanceType[T Assignee](inst *Instance, tsch TypedSchema[T]) (*TypedInstance[T], error) {
	// if !schemaIs(inst.Schema(), tsch) {
	// FIXME stop assuming underlying type UGH
	if !tsch.(*unaryTypedSchema[T]).is(inst.Schema()) {
		return nil, fmt.Errorf("typed schema is not derived from instance's schema")
	}

	return &TypedInstance[T]{
		Instance: inst,
		tsch:     tsch,
	}, nil
}

// Instance represents data that is a valid instance of a Thema [Schema].
//
// It is not possible to create a valid Instance directly. They can only be
// obtained by successful call to [Schema.Validate].
type Instance struct {
	// The CUE representation of the input data
	raw cue.Value
	// A name for the input data, primarily for use in error messages
	name string
	// The schema the data validated against/of which the input data is a valid instance
	sch Schema

	// simple flag the prevents external creation
	valid bool
}

func (i *Instance) check() {
	if !i.valid {
		panic("Instance is not valid; Instances must be created by a call to thema.Schema.Validate")
	}
}

// Hydrate returns a copy of the Instance with all default values specified by
// the schema included.
//
// NOTE hydration implementation is a WIP. If errors are encountered, the
// original input is returned unchanged.
func (i *Instance) Hydrate() *Instance {
	i.check()

	i.sch.Lineage().Runtime()
	ni, err := doHydrate(i.sch.Underlying(), i.raw)
	// FIXME For now, just no-op it if we error
	if err != nil {
		return i
	}

	return &Instance{
		valid: true,
		raw:   ni,
		name:  i.name,
		sch:   i.sch,
	}
}

// Dehydrate returns a copy of the Instance with all default values specified by
// the schema removed.
//
// NOTE dehydration implementation is a WIP. If errors are encountered, the
// original input is returned unchanged.
func (i *Instance) Dehydrate() *Instance {
	i.check()

	ni, _, err := doDehydrate(i.sch.Underlying(), i.raw)
	// FIXME For now, just no-op it if we error
	if err != nil {
		return i
	}

	return &Instance{
		valid: true,
		raw:   ni,
		name:  i.name,
		sch:   i.sch,
	}
}

// AsSuccessor translates the instance into the form specified by the successor
// schema.
func (i *Instance) AsSuccessor() (*Instance, TranslationLacunas, error) {
	i.check()
	// If it's a minor version upgrade, we can safely shortcut and just create
	// a new instance
	nsch := i.Schema().Successor()
	if nsch.Version()[0] == i.Schema().Version()[0] {
		ni := new(Instance)
		*ni = *i
		ni.sch = nsch
		return ni, nil, nil
	}
	return i.Translate(i.sch.Successor().Version())
}

// AsPredecessor translates the instance into the form specified by the predecessor
// schema.
func (i *Instance) AsPredecessor() (*Instance, TranslationLacunas, error) {
	i.check()
	return i.Translate(i.sch.Predecessor().Version())
}

// Underlying returns the cue.Value representing the data contained in the Instance.
func (i *Instance) Underlying() cue.Value {
	i.check()
	return i.raw
}

// Schema returns the [Schema] corresponding to this instance.
func (i *Instance) Schema() Schema {
	i.check()
	return i.sch
}

func (i *Instance) rt() *Runtime {
	return getLinLib(i.Schema().Lineage())
}

// TypedInstance represents data that is a valid instance of a Thema
// [TypedSchema].
//
// A TypedInstance is to a [TypedSchema] as an [Instance] is to a [Schema].
//
// It is not possible to create a valid TypedInstance directly. They can only be
// obtained by successful call to [TypedSchema.Validate].
type TypedInstance[T Assignee] struct {
	*Instance
	tsch TypedSchema[T]
}

// TypedSchema returns the [TypedSchema] corresponding to this instance.
//
// This method is identical to [Instance.Schema], except that it returns the already-typed variant.
func (inst *TypedInstance[T]) TypedSchema() TypedSchema[T] {
	inst.check()
	return inst.tsch
}

// Value returns a Go struct of this TypedInstance's generic [Assignee] type,
// populated with the data contained in this instance, including default values, etc.
//
// This method is similar to [json.Unmarshal] - it decodes serialized data into a standard Go type
// for working with in all the usual ways.
func (inst *TypedInstance[T]) Value() (T, error) {
	inst.check()

	t := inst.tsch.NewT()
	// TODO figure out correct pointer handling here
	err := inst.Instance.raw.Decode(&t)
	return t, err
}

// ValueP is the same as Value, but panics if an error is encountered.
func (inst *TypedInstance[T]) ValueP() T {
	inst.check()

	t, err := inst.Value()
	if err != nil {
		panic(fmt.Errorf("error decoding value: %w", err))
	}
	return t
}

// Translate transforms the provided [Instance] to an Instance of a different
// [Schema] from the same [Lineage]. A new *Instance is returned representing the
// transformed value, along with any lacunas accumulated along the way.
//
// Forward translation within a major version (e.g. 0.0 to 0.7) is trivial, as
// all those schema changes are established as backwards compatible by Thema's
// lineage invariants. In such cases, the lens is referred to as implicit, as
// the lineage author does not write it, with translation relying on simple
// unification. Lacunas cannot be emitted from such translations.
//
// Forward translation across major versions (e.g. 0.0 to 1.0), and all reverse
// translation regardless of sequence boundaries (e.g. 1.1 to either 1.0
// or 0.0), is nontrivial and relies on explicitly defined lenses, which
// introduce room for lacunas and author judgment.
//
// Thema translation is non-invertible by design. That is, Thema does not seek
// to generally guarantee that translating an instance from 0.0->1.0->0.0 will
// result in the exact original data. Input state preservation can be fully
// achieved in the program depending on Thema, so we avoid introducing
// complexity into Thema that is not essential for all use cases.
//
// Errors only occur in cases where lenses were written in an unexpected way -
// for example, not all fields were mapped over, and the resulting object is not
// concrete. All errors returned from this func will children of [terrors.ErrInvalidLens].
func (i *Instance) Translate(to SyntacticVersion) (*Instance, TranslationLacunas, error) {
	i.check()

	if len(i.Schema().Lineage().(*baseLineage).lensmap) > 0 {
		return i.translateGo(to)
	}

	// TODO define this in terms of AsSuccessor and AsPredecessor, rather than those in terms of this.
	newsch, err := i.Schema().Lineage().Schema(to)
	if err != nil {
		panic(fmt.Sprintf("no schema in lineage with version %v, cannot translate", to))
	}

	out, err := cueArgs{
		"inst": i.raw,
		"to":   to,
		"from": i.Schema().Version(),
		"lin":  i.Schema().Lineage().Underlying(),
	}.call("#Translate", i.rt())
	if err != nil {
		// This can't happen without a name change or an invariant violation
		panic(err)
	}

	if out.Err() != nil {
		return nil, nil, errors.Mark(out.Err(), terrors.ErrInvalidLens)
	}

	lac := make(multiTranslationLacunas, 0)
	out.LookupPath(cue.MakePath(cue.Str("lacunas"))).Decode(&lac)

	// Attempt to evaluate #Translate result to remove intermediate structures created by #Translate.
	// Otherwise, all the #Translate results are non-concrete, which leads to undesired effects.
	raw, _ := out.LookupPath(cue.MakePath(cue.Str("result"), cue.Str("result"))).Default()

	// Check that the result is concrete by trying to marshal/export it as JSON
	_, err = json.Marshal(raw)
	if err != nil {
		return nil, nil, errors.Mark(fmt.Errorf("lens produced a non-concrete result: %s", cerrors.Details(err, nil)), terrors.ErrLensIncomplete)
	}

	// Ensure the result is a valid instance of the target schema
	inst, err := newsch.Validate(raw)
	if err != nil {
		return nil, nil, errors.Mark(err, terrors.ErrLensResultIsInvalidData)
	}
	return inst, lac, err
}

func (i *Instance) translateGo(to SyntacticVersion) (*Instance, TranslationLacunas, error) {
	from := i.Schema().Version()
	if to == from {
		// TODO make sure this mirrors the pure CUE behavior
		return i, nil, nil
	}
	lensmap := i.Schema().Lineage().(*baseLineage).lensmap

	sch := i.Schema()
	ti := new(Instance)
	*ti = *i
	for sch.Version() != to {
		var nsch Schema
		if to.Less(from) {
			nsch = sch.Predecessor()
		} else {
			nsch = sch.Successor()
		}

		var rti *Instance
		var err error
		if to.Less(from) || sch.Version()[0] != nsch.Version()[0] {
			// Going backward, or crossing major version - need explicit lens
			mlid := lid(sch.Version(), nsch.Version())
			rti, err = lensmap[mlid].Mapper(ti, nsch)
			if err != nil {
				return nil, nil, fmt.Errorf("error executing %s migration: %w", mlid, err)
			}
			// Ensure that
			//  - the returned instance exists
			//  - the caller returned an instance of the expected schema version
			if rti == nil {
				return nil, nil, fmt.Errorf("lens returned a nil instance")
			}
			if rti.Schema().Version() != nsch.Version() {
				return nil, nil, fmt.Errorf("lens returned an instance of the wrong schema version: expected %v, got %v", nsch.Version(), rti.Schema().Version())
			}
		} else {
			// going up a minor version - neither errors nor lacunas are possible
			rti, _, err = ti.AsSuccessor()
			if err != nil {
				panic(fmt.Sprintf("unreachable - error on minor version upgrade: %s", err))
			}
		}
		*ti = *rti
		sch = nsch
	}

	return ti, nil, nil
}

type multiTranslationLacunas []struct {
	V   SyntacticVersion `json:"v"`
	Lac []Lacuna         `json:"lacunas"`
}

func (lac multiTranslationLacunas) AsList() []Lacuna {
	// FIXME This loses info, naturally - need to rework the lacuna types
	var l []Lacuna
	for _, v := range lac {
		l = append(l, v.Lac...)
	}
	return l
}
