package errors

import (
	"github.com/cockroachdb/errors"
)

// ValidationCode represents different classes of validation errors that may
// occur vs. concrete data inputs.
type ValidationCode uint16

const (
	// KindConflict indicates a validation failure in which the schema and data
	// values are of differing, conflicting kinds - the schema value does not
	// subsume the data value. Example: data: "foo"; schema: int
	KindConflict ValidationCode = 1 << iota

	// OutOfBounds indicates a validation failure in which the data and schema have
	// the same (or subsuming) kinds, but the data is out of schema-defined bounds.
	// Example: data: 4; schema: int & <4
	OutOfBounds

	// MissingField indicates a validation failure in which the data lacks
	// a field that is required in the schema.
	MissingField

	// ExcessField indicates a validation failure in which the schema is treated as
	// closed, and the data contains a field not specified in the schema.
	ExcessField
)

// ValidationError is a subtype of
type ValidationError struct {
	msg string
}

// Unwrap implements standard Go error unwrapping, relied on by errors.Is.
//
// All ValidationErrors wrap the general ErrInvalidData sentinel error.
func (ve *ValidationError) Unwrap() error {
	return ErrInvalidData
}

// Validation error codes/types
var (
	// ErrInvalidData is the general error that indicates some data failed validation
	// against a Thema schema. Use it with errors.Is() to differentiate validation errors
	// from other classes of failure.
	ErrInvalidData = errors.New("data not a valid instance of schema")

	// ErrInvalidExcessField indicates a validation failure in which the schema is
	// treated as closed, and the data contains a field not specified in the schema.
	ErrInvalidExcessField = errors.New("data contains field not present in schema")

	// ErrInvalidMissingField indicates a validation failure in which the data lacks
	// a field that is required in the schema.
	ErrInvalidMissingField = errors.New("required field is absent in data")

	// ErrInvalidKindConflict indicates a validation failure in which the schema and
	// data values are of differing, conflicting kinds - the schema value does not
	// subsume the data value. Example: data: "foo"; schema: int
	ErrInvalidKindConflict = errors.New("schema and data are conflicting kinds")

	// ErrInvalidOutOfBounds indicates a validation failure in which the data and
	// schema have the same (or subsuming) kinds, but the data is out of
	// schema-defined bounds. Example: data: 4; schema: int & <3
	ErrInvalidOutOfBounds = errors.New("data is out of schema bounds")
)

// Translation errors. These all occur as a result of an invalid lens. Currently
// these may be returned from [thema.Instance.Translate]. Eventually, it is
// hoped that they will be caught statically in [thema.BindLineage] and cannot
// occur at runtime.
var (
	// ErrInvalidLens indicates that a lens is not correctly written. It is the parent
	// to all other lens and translation errors, and is a child of ErrInvalidLineage.
	ErrInvalidLens = errors.New("lens is invalid")

	// ErrLensIncomplete indicates that translating some valid data through
	// a lens produced a non-concrete result. This always indicates a problem with the
	// lens as it is written, and as such is a child of ErrInvalidLens.
	ErrLensIncomplete = errors.New("result of lens translation is not concrete")

	// ErrLensResultIsInvalidData indicates that translating some valid data through a
	// lens produced a result that was not an instance of the target schema. This
	// always indicates a problem with the lens as it is written, and as such is a
	// child of ErrInvalidLens.
	ErrLensResultIsInvalidData = errors.New("result of lens translation is not valid for target schema")
)

// Lower level general errors
var (
	// ErrValueNotExist indicates that a necessary CUE value did not exist.
	ErrValueNotExist = errors.New("cue value does not exist")

	// ErrValueNotALineage indicates that a provided CUE value is not a lineage.
	// This is almost always an end-user error - they oops'd and provided the
	// wrong path, file, etc.
	ErrValueNotALineage = errors.New("not a lineage")

	// ErrInvalidLineage indicates that a provided lineage does not fulfill one
	// or more of the Thema invariants.
	ErrInvalidLineage = errors.New("invalid lineage")

	// ErrInvalidSchemasOrder indicates that schemas in a lineage are not ordered
	// by version.
	ErrInvalidSchemasOrder = errors.New("schemas in lineage are not ordered by version")

	// ErrInvalidLensesOrder indicates that lenses are in the wrong order - they must be sorted by `to`, then `from`.
	ErrInvalidLensesOrder = errors.New("lenses in lineage are not ordered by version")

	// ErrDuplicateLenses indicates that a lens was defined declaratively in CUE, but the same lens
	// was also provided as a Go function to BindLineage.
	ErrDuplicateLenses = errors.New("lens is declared in both CUE and Go")

	// ErrMissingLenses indicates that the lenses provided to BindLineage in either
	// CUE or Go were missing at least one of the expected lenses determined by the
	// set of schemas in the lineage.
	ErrMissingLenses = errors.New("not all expected lenses were provided")

	// ErrErroneousLenses indicates that a lens was provided to BindLineage in either
	// CUE or Go that was not one of the expected lenses determined by the set of
	// schemas in the lineage.
	ErrErroneousLenses = errors.New("unexpected lenses were erroneously provided")

	// ErrVersionNotExist indicates that no schema exists in a lineage with a
	// given version.
	ErrVersionNotExist = errors.New("lineage does not contain schema with version") // ErrNoSchemaWithVersion

	// ErrMalformedSyntacticVersion indicates a string input of a syntactic
	// version was malformed.
	ErrMalformedSyntacticVersion = errors.New("not a valid syntactic version")
)
