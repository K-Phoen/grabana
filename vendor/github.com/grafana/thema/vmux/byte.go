package vmux

import "github.com/grafana/thema"

// ByteMux is a version multiplexer that maps a []byte containing data at any
// schematized version to a []byte containing data at a particular schematized version.
type ByteMux func(b []byte) ([]byte, thema.TranslationLacunas, error)

// NewByteMux creates a [ByteMux] func from the provided [thema.Schema].
//
// When the returned mux func is called, it will:
//
//   - Decode the input []byte using the provided [Codec], then
//   - Pass the result to [thema.Schema.Validate], then
//   - Call [thema.Instance.Translate] on the result, to the version of the provided [thema.Schema], then
//   - Encode the resulting [thema.Instance] to a []byte, then
//   - Return the resulting []byte, [thema.TranslationLacunas], and error
//
// The returned error may be from any of the above steps.
func NewByteMux(sch thema.Schema, codec Codec) ByteMux {
	f := NewUntypedMux(sch, codec)
	return func(b []byte) ([]byte, thema.TranslationLacunas, error) {
		ti, lac, err := f(b)
		if err != nil {
			return nil, lac, err
		}
		ob, err := codec.Encode(ti.Underlying())
		return ob, lac, err
	}
}
