package vmux

import "github.com/grafana/thema"

// ValueMux is a version multiplexer that maps a []byte containing data at any
// schematized version to a Go var of a type that a particular schematized
// version is [thema.AssignableTo].
type ValueMux[T thema.Assignee] func(b []byte) (T, thema.TranslationLacunas, error)

// NewValueMux creates a [ValueMux] func from the provided [thema.TypedSchema].
//
// When the returned mux func is called, it will:
//
//   - Decode the input []byte using the provided [Decoder], then
//   - Pass the result to [thema.TypedSchema.ValidateTyped], then
//   - Call [thema.Instance.Translate] on the result, to the version of the provided [thema.TypedSchema], then
//   - Populate an instance of T by calling [thema.TypedInstance.Value] on the result, then
//   - Return the resulting T, [thema.TranslationLacunas], and error
//
// The returned error may be from any of the above steps.
func NewValueMux[T thema.Assignee](sch thema.TypedSchema[T], dec Decoder) ValueMux[T] {
	f := NewTypedMux[T](sch, dec)
	return func(b []byte) (T, thema.TranslationLacunas, error) {
		ti, lac, err := f(b)
		if err != nil {
			return sch.NewT(), lac, err
		}
		t, err := ti.Value()
		return t, lac, err
	}
}
