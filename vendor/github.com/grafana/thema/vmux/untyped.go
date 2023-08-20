package vmux

import (
	"fmt"

	"github.com/grafana/thema"
)

// UntypedMux is a version multiplexer that maps a []byte containing data at any
// schematized version to a [thema.Instance] at a particular schematized version.
type UntypedMux func(b []byte) (*thema.Instance, thema.TranslationLacunas, error)

// NewUntypedMux creates an [UntypedMux] from the provided [thema.Schema].
//
// When the returned mux func is called, it will:
//
//   - Decode the input []byte using the provided [Decoder], then
//   - Pass the result to [thema.Schema.Validate], then
//   - Call [thema.Instance.Translate] on the result, to the version of the provided [thema.Schema], then
//   - Return the resulting [thema.Instance], [thema.TranslationLacunas], and error
//
// The returned error may be from any of the above steps.
func NewUntypedMux(sch thema.Schema, dec Decoder) UntypedMux {
	ctx := sch.Lineage().Underlying().Context()
	// Prepare no-match error string once for reuse
	vstring := allvstr(sch)

	return func(b []byte) (*thema.Instance, thema.TranslationLacunas, error) {
		v, err := dec.Decode(ctx, b)
		if err != nil {
			// TODO wrap error for use with errors.Is
			return nil, nil, err
		}

		// Try the given schema first, on the premise that in general it's the
		// most likely one for an application to encounter
		tinst, err := sch.Validate(v)
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
				return inst.Translate(sch.Version())
			}
		}

		return nil, nil, fmt.Errorf("data invalid against all versions (%s), error against %s: %w", vstring, sch.Version(), err)
	}
}
