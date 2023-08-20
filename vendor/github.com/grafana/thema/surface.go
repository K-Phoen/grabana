package thema

import (
	"fmt"
	"strconv"
	"strings"

	"cuelang.org/go/cue"

	terrors "github.com/grafana/thema/errors"
	"github.com/grafana/thema/internal/envvars"
)

// A CUEWrapper wraps a cue.Value, and can return that value for inspection.
type CUEWrapper interface {
	// Underlying returns the underlying cue.Value wrapped by the object.
	Underlying() cue.Value
}

// A Lineage is the top-level container in Thema, holding the complete
// evolutionary history of a particular kind of object: every schema that has
// ever existed for that object, and the lenses that allow translating between
// those schema versions.
//
// Lineages may only be produced by calling [BindLineage].
type Lineage interface {
	CUEWrapper

	// Name returns the name of the object schematized by the lineage, as declared
	// in the lineage's `name` field.
	Name() string

	// ValidateAny checks that the provided data is valid with respect to at
	// least one of the schemas in the lineage. The oldest (smallest) schema against
	// which the data validates is chosen. A nil return indicates no validating
	// schema was found.
	//
	// While this method takes a cue.Value, this is only to avoid having to trigger
	// the translation internally; input values must be concrete. To use
	// incomplete CUE values with Thema schemas, prefer working directly in CUE,
	// or if you must, rely on Underlying().
	//
	// TODO should this instead be interface{} (ugh ugh wish Go had tagged unions) like FillPath?
	ValidateAny(data cue.Value) *Instance

	// Schema returns the schema identified by the provided version, if one exists.
	//
	// Only the [0, 0] schema is guaranteed to exist in all valid lineages.
	Schema(v SyntacticVersion) (Schema, error)

	// First returns the first Schema in the lineage (v0.0). Thema requires that all
	// valid lineages contain at least one schema, so this is guaranteed to exist.
	First() Schema

	// Latest returns the newest Schema in the lineage - largest minor version
	// within the largest major version.
	//
	// Thema requires that all valid lineages contain at least one schema, so schema
	// is is guaranteed to exist, even if it's the 0.0 version.
	//
	// EXERCISE CAUTION WITH THIS METHOD. Relying on Latest is appropriate and
	// necessary for some use cases, such as keeping Thema declarations and
	// generated code in sync within a single repository. But use in the wrong
	// context - usually cross repository, loosely coupled, dependency
	// management-like contexts - can completely undermine Thema's translatability
	// invariants.
	//
	// If you're not sure, ask yourself: when a breaking change to this lineage is
	// published, what would that break downstream, and will the users who experience
	// that breakage be expecting it to happen?
	//
	// If the user would be expecting the breakage, using Latest is probably appropriate.
	// Otherwise, it is probably preferable to pick an explicit version number.
	Latest() Schema

	// All returns all Schemas in the lineage. Thema requires that all valid lineages
	// contain at least one schema, so this is guaranteed to contain at least one element.
	All() []Schema

	// Runtime returns the thema.Runtime instance with which this lineage was built.
	Runtime() *Runtime

	// Lineage must be a private interface in order to ensure creation is only possible
	// through BindLineage().
	allVersions() versionList
}

// ImperativeLens is a lens transformation defined as a Go function, rather than
// in native CUE alongside the lineage.
//
// See [ImperativeLenses] for more information.
type ImperativeLens struct {
	To, From SyntacticVersion
	Mapper   func(inst *Instance, to Schema) (*Instance, error)
}

// SchemaP returns the schema identified by the provided version. If no schema
// exists in the lineage with the provided version, it panics.
//
// This is a simple convenience wrapper on the Lineage.Schema() method.
func SchemaP(lin Lineage, v SyntacticVersion) Schema {
	sch, err := lin.Schema(v)
	if err != nil {
		panic(err)
	}
	return sch
}

// LatestVersion returns the version number of the newest (largest) schema
// version in the provided lineage.
//
// Deprecated: call Lineage.Latest().Version().
func LatestVersion(lin Lineage) SyntacticVersion {
	return lin.Latest().Version()
}

// LatestVersionInSequence returns the version number of the newest (largest) schema
// version in the provided sequence number.
//
// An error indicates the number of the provided sequence does not exist.
//
// Deprecated: call Schema.LatestInMajor().Version() after loading a schema in the desired major version.
func LatestVersionInSequence(lin Lineage, seqv uint) (SyntacticVersion, error) {
	sch, err := lin.Schema(SV(seqv, 0))
	if err != nil {
		return SyntacticVersion{}, err
	}
	return sch.LatestInMajor().Version(), nil
}

// A LineageFactory returns a [Lineage], which is immutably bound to a single
// instance of #Lineage declared in CUE.
//
// LineageFactory funcs are intended to be the main Go entrypoint to all of the
// operations, guarantees, and capabilities of Thema lineages. Lineage authors
// should generally define and export one instance of LineageFactory per
// #Lineage instance.
//
// It is idiomatic to name LineageFactory funcs after the "name" field on the
// lineage they return:
//
//	func <name>Lineage ...
//
// If the Go package and lineage name are the same, the name should be omitted from
// the builder func to reduce stutter:
//
//	func Lineage ...
//
// Deprecated: having an explicit type for this adds little value.
type LineageFactory func(*Runtime, ...BindOption) (Lineage, error)

// A ConvergentLineageFactory is the same as a LineageFactory, but for a
// ConvergentLineage.
//
// There is no reason to provide both a ConvergentLineageFactory and a
// LineageFactory, as the latter is always reachable from the former. As such,
// idiomatic naming conventions are unchanged.
//
// Deprecated: having an explicit type for this adds little value.
type ConvergentLineageFactory[T Assignee] func(*Runtime, ...BindOption) (ConvergentLineage[T], error)

// A BindOption defines options that may be specified only at initial
// construction of a [Lineage] via [BindLineage].
type BindOption bindOption

// Internal representation of BindOption.
type bindOption func(c *bindConfig)

// Internal bind-time configuration options.
type bindConfig struct {
	skipbuggychecks bool
	implens         []ImperativeLens
}

// SkipBuggyChecks indicates that [BindLineage] should skip validation checks
// which have known bugs (e.g. panics) for certain should-be-valid CUE inputs.
//
// By default, BindLineage performs these checks anyway, as otherwise the
// default behavior of BindLineage is to not provide the guarantees it's
// supposed to provide.
//
// As Thema and CUE move towards maturity and the set of validations that are
// both a) necessary and b) buggy empties out, this will naturally become a
// no-op. At that point, this function will be marked deprecated.
//
// Ratcheting up verification checks in this way does mean that any code relying
// on this to bypass verification in BindLineage may begin failing in future
// versions of Thema if the underlying lineage being verified doesn't comply
// with a planned invariant.
func SkipBuggyChecks() BindOption {
	return func(c *bindConfig) {
		// We let the env var override this to make it easy to disable on tests.
		if !envvars.ForceVerify {
			c.skipbuggychecks = true
		}
	}
}

// ImperativeLenses takes a slice of [ImperativeLens]. These lenses will be
// executed on calls to [Instance.Translate].
//
// Currently, the entire lens set must be provided in either Go or CUE. This
// restriction may be relaxed in the future to allow a mix of Go and CUE lenses,
// or to allow Go funcs to supersede CUE lenses as a performance optimization.
//
// When providing lenses in this way, [BindLineage] will fail unless exactly the
// set of expected lenses is provided. The correctness of the function bodies
// cannot be pre-verified in this way, as Go is Turing-complete, but it is enforced
// at runtime that lenses return an [Instance] of the schema version they claim to
// in [ImperativeLens.To].
//
// Writing lenses in Go means that pure native CUE is no longer sufficient to
// produce a valid lineage. As a result, lineages are no longer portable outside
// of Go programs with compile-time access to the Go-defined lenses.
func ImperativeLenses(lenses ...ImperativeLens) BindOption {
	return func(c *bindConfig) {
		c.implens = append(c.implens, lenses...)
	}
}

// Schema represents a single, complete schema from a thema lineage. A Schema's
// Validate() method determines whether some data constitutes an Instance.
type Schema interface {
	CUEWrapper

	// Validate checks that the provided data is valid with respect to the
	// schema. If valid, the data is wrapped in an [Instance] and returned.
	// Otherwise, a nil Instance is returned along with an error detailing the
	// validation failure.
	//
	// While Validate takes a cue.Value, this is only to avoid having to trigger
	// the translation internally; input values must be concrete. Behavior of
	// this method is undefined for incomplete values.
	//
	// The concreteness requirement may be loosened in future versions of Thema. To
	// use incomplete CUE values with Thema schemas, prefer working directly in CUE,
	// or call [Schema.Underlying] to work directly with the underlying CUE API.
	//
	// TODO should this instead be interface{} (ugh ugh wish Go had tagged unions) like FillPath?
	Validate(data cue.Value) (*Instance, error)

	// Successor returns the next schema in the lineage, or nil if it is the last schema.
	Successor() Schema

	// Predecessor returns the previous schema in the lineage, or nil if it is the first schema.
	Predecessor() Schema

	// LatestInMajor returns the Schema with the newest (largest) minor version
	// within this Schema's major version. If the receiver Schema is the latest, it
	// will return itself.
	LatestInMajor() Schema

	// Version returns the schema's version number.
	Version() SyntacticVersion

	// Lineage returns the lineage that contains this schema.
	Lineage() Lineage

	// Examples returns the set of examples of this schema defined in the original
	// lineage. The string key is the name given to the example.
	Examples() map[string]*Instance

	// Schema must be a private interface in order to ensure all instances fully
	// conform to Thema invariants.
	_schema()
}

// ConvergentLineage is a lineage where exactly one of its contained schemas
// is associated with a Go type - a TypedSchema[Assignee], as returned from
// [BindType].
//
// This variant of lineage is intended to directly support the primary
// anticipated use pattern for Thema within a Go program: accepting all
// historical forms of an object's schema as input to the program, but writing
// the program against just one version.
//
// This process is known as version multiplexing. See
// [github.com/grafana/thema/vmux].
type ConvergentLineage[T Assignee] interface {
	Lineage

	TypedSchema() TypedSchema[T]
}

// TypedSchema is a Thema schema that has been bound to a particular Go type, per
// Thema's assignability rules.
type TypedSchema[T Assignee] interface {
	Schema

	// NewT returns a new instance of T, but with schema-specified defaults for its
	// field values instead of Go zero values. Fields without a schema-specified default
	// are populated with standard Go zero values.
	NewT() T

	// ValidateTyped performs validation identically to [Schema.Validate], but
	// returns a TypedInstance on success.
	ValidateTyped(data cue.Value) (*TypedInstance[T], error)

	// ConvergentLineage returns the ConvergentLineage that contains this schema.
	ConvergentLineage() ConvergentLineage[T]
}

// Assignee is a type constraint used by Thema generics for type parameters
// where there exists a particular Schema that is [AssignableTo] the type.
//
// This property is not representable in Go's static type system, as Thema types
// are dynamic, and AssignableTo() is a runtime check. Thus, the only actual
// type constraint Go's type system can be made aware of is any.
//
// Instead, Thema's implementation guarantees that it is only possible to
// instantiate a generic type with an Assignee type parameter if the relevant
// AssignableTo() relation has already been verified, and there is an
// unambiguous relationship between the generic type and the relevant [Schema].
//
// For example: for TypedSchema[T Assignee], it is the related Schema. With
// TypedInstance[T Assignee], the related schema is returned from its
// TypedSchema() method.
//
// As this type constraint is simply any, it exists solely as a signal to the
// human reader that the relation to a Schema exists, and that the relation
// has been verified in any properly instantiated type carrying this generic
// type constraint. (Improperly instantiated generic Thema types panic upon
// calls to any of their methods)
type Assignee any

// SyntacticVersion is a two-tuple of uints describing the position of a schema
// within a lineage. Syntactic versions are Thema's canonical version numbering
// system.
//
// The first element is the index of the sequence containing the schema within
// the lineage, and the second element is the index of the schema within that
// sequence.
type SyntacticVersion [2]uint

// SV creates a [SyntacticVersion].
//
// A trivial helper to avoid repetitive Go-stress disorder from countless
// instances of typing:
//
//	SyntacticVersion{0, 0}
func SV(seqv, schv uint) SyntacticVersion {
	return [2]uint{seqv, schv}
}

// Less reports whether the receiver [SyntacticVersion] is less than the
// provided one, consistent with the expectations of the stdlib sort package.
func (sv SyntacticVersion) Less(osv SyntacticVersion) bool {
	return sv[0] < osv[0] || (sv[0] == osv[0] && sv[1] < osv[1])
}

func (sv SyntacticVersion) String() string {
	return fmt.Sprintf("%v.%v", sv[0], sv[1])
}

// ParseSyntacticVersion parses a canonical representation of a
// [SyntacticVersion] (e.g. "0.0") from a string.
func ParseSyntacticVersion(s string) (SyntacticVersion, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return synv(), fmt.Errorf("%w: %q", terrors.ErrMalformedSyntacticVersion, s)
	}

	// i mean 4 billion is probably enough version numbers
	seqv, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return synv(), fmt.Errorf("%w: %q has invalid sequence number %q", terrors.ErrMalformedSyntacticVersion, s, parts[0])
	}

	// especially when squared
	schv, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return synv(), fmt.Errorf("%w: %q has invalid schema number %q", terrors.ErrMalformedSyntacticVersion, s, parts[1])
	}
	return synv(uint(seqv), uint(schv)), nil
}

type versionList []SyntacticVersion

func (vl versionList) String() string {
	vstrings := make([]string, 0, len(vl))
	for _, v := range vl {
		vstrings = append(vstrings, v.String())
	}
	return strings.Join(vstrings, ", ")
}
