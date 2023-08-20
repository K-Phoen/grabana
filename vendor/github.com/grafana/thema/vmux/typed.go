package vmux

import (
	"fmt"

	"github.com/grafana/thema"
)

// TypedMux is a version multiplexer that maps a []byte containing data at any
// schematized version to a [thema.TypedInstance] at a particular schematized version.
type TypedMux[T thema.Assignee] func(b []byte) (*thema.TypedInstance[T], thema.TranslationLacunas, error)

// NewTypedMux creates a [TypedMux] func from the provided [thema.TypedSchema].
//
// When the returned mux func is called, it will:
//
//   - Decode the input []byte using the provided [Decoder], then
//   - Pass the result to [thema.TypedSchema.ValidateTyped], then
//   - Call [thema.Instance.Translate] on the result, to the version of the provided [thema.TypedSchema], then
//   - Return the resulting [thema.TypedInstance], [thema.TranslationLacunas], and error
//
// The returned error may be from any of the above steps.
func NewTypedMux[T thema.Assignee](sch thema.TypedSchema[T], dec Decoder) TypedMux[T] {
	ctx := sch.Lineage().Underlying().Context()
	// Prepare no-match error string once for reuse
	vstring := allvstr(sch)

	return func(b []byte) (*thema.TypedInstance[T], thema.TranslationLacunas, error) {
		v, err := dec.Decode(ctx, b)
		if err != nil {
			// TODO wrap error for use with errors.Is
			return nil, nil, err
		}

		// Try the given schema first, on the premise that in general it's the
		// most likely one for an application to encounter
		tinst, err := sch.ValidateTyped(v)
		if err == nil {
			return tinst, nil, nil
		}

		// Walk in reverse order on the premise that, in general, newer versions are more
		// likely to be provided than older versions
		isch := latest(sch.Lineage())
		for ; isch != nil; isch = isch.Predecessor() {
			if isch.Version() == sch.Version() {
				continue
			}

			if inst, ierr := isch.Validate(v); ierr == nil {
				trinst, lac, err := inst.Translate(sch.Version())
				if err != nil {
					return nil, nil, err
				}

				// TODO perf: introduce a typed translator to avoid wastefully re-binding the go type every time
				tinst, err := thema.BindInstanceType(trinst, sch)
				if err != nil {
					panic(fmt.Errorf("unreachable, instance type should always be bindable: %w", err))
				}
				return tinst, lac, nil
			}
		}

		return nil, nil, fmt.Errorf("data invalid against all versions (%s), error against %s: %w", vstring, sch.Version(), err)
	}
}
