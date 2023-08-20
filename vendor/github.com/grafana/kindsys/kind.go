package kindsys

import (
	"fmt"

	"github.com/grafana/thema"

	"github.com/grafana/kindsys/encoding"
)

// TODO docs
type Maturity string

const (
	MaturityMerged       Maturity = "merged"
	MaturityExperimental Maturity = "experimental"
	MaturityStable       Maturity = "stable"
	MaturityMature       Maturity = "mature"
)

func maturityIdx(m Maturity) int {
	// icky to do this globally, this is effectively setting a default
	if string(m) == "" {
		m = MaturityMerged
	}

	for i, ms := range maturityOrder {
		if m == ms {
			return i
		}
	}
	panic(fmt.Sprintf("unknown maturity milestone %s", m))
}

var maturityOrder = []Maturity{
	MaturityMerged,
	MaturityExperimental,
	MaturityStable,
	MaturityMature,
}

func (m Maturity) Less(om Maturity) bool {
	return maturityIdx(m) < maturityIdx(om)
}

func (m Maturity) String() string {
	return string(m)
}

// Kind is a runtime representation of a Grafana kind definition.
//
// Kind definitions are canonically written in CUE. Loading and validating such
// CUE definitions produces instances of this type. Kind, and its
// sub-interfaces, are the expected canonical way of working with kinds in Go.
//
// Kind has six sub-interfaces, all of which provide:
//
// - Access to the kind's defined meta-properties, such as `name`, `pluralName`, or `maturity`
// - Access to the schemas defined in the kind
// - Methods for certain key operations on the kind and object instances of its schemas
//
// Kind definitions are written in CUE. The meta-schema specifying how to write
// kind definitions are also written in CUE. See the files at the root of
// [the kindsys repository].
//
// There are three categories of kinds, each having its own sub-interface:
// [Core], [Custom], and [Composable]. All kind definitions are in exactly one
// category (a kind can't be both Core and Composable). Correspondingly, all
// instances of Kind also implement exactly one of these sub-interfaces.
//
// Conceptually, kinds are similar to class definitions in object-oriented
// programming. They define a particular type of object, and how instances of
// that object should be created. The object defined by a [Core] or [Custom] kind
// is called a [Resource]. TODO name for the associated object for composable kinds
//
// [Core], [Custom] and [Composable] all provide methods for unmarshaling []byte
// into an unstructured Go type, [UnstructuredResource], similar to how
// json.Unmarshal can use map[string]any as a universal fallback. Relying on
// this untyped approach is recommended for use cases that need to work
// generically on any Kind. This is especially because untyped Kinds are
// portable, and can be loaded at runtime in Go: the original CUE definition is
// sufficient to create instances of [Core], [Custom] or [Composable].
//
// However, when working with a specific, finite set of kinds, it is usually
// preferable to use the typed interfaces:
//
// - [Core] -> [TypedCore]
// - [Custom] -> [TypedCustom]
// - [Composable] -> [TypedComposable] (TODO not yet implemented)
//
// Each embeds the corresponding untyped interface, and takes a generic type
// parameter. The provided struct is verified to be assignable to the latest
// schema defined in the kind. (See [thema.BindType]) Additional methods are
// provided on Typed* variants that do the same as their untyped counterparts,
// but using the type given in the generic type parameter.
//
// Directly implementing this interface is discouraged. Strongly prefer instead to
// rely on [BindCore], [BindCustom], or [BindComposable].
//
// [the kindsys repository]: https://github.com/grafana/kindsys
type Kind interface {
	// Name returns the kind's name, as specified in the name field of the kind definition.
	//
	// Note that this is the capitalized name of the kind. For other names and
	// common kind properties, see [Props.CommonProperties].
	Name() string

	// MachineName returns the kind's machine name, as specified in the machineName
	// field of the kind definition.
	MachineName() string

	// Maturity indicates the maturity of this kind, one of the enum of values we
	// accept in the maturity field of the kind definition.
	Maturity() Maturity

	// Props returns a [SomeKindProps], representing the properties
	// of the kind as declared in the .cue source. The underlying type is
	// determined by the category of kind.
	Props() SomeKindProperties

	// CurrentVersion returns the version number of the schema that is considered
	// the 'current' version, usually the latest version. When initializing object
	// instances of this Kind, the current version is used by default.
	CurrentVersion() thema.SyntacticVersion

	// Lineage returns the kind's [thema.Lineage]. The lineage contains the full
	// history of object schemas associated with the kind.
	//
	// TODO separate this onto an optional, additional interface
	Lineage() thema.Lineage
}

// ResourceKind represents a kind that defines a root object, or Kubernetish resource.
//
// TODO this is a temporary intermediate while we combine [Core] and [Custom]. This name will probably go away.
type ResourceKind interface {
	Kind

	// Validate takes a []byte representing an object instance of this kind and
	// checks that it is a valid instance of at least one schema defined in the
	// kind.
	//
	// A decoder must be provided that knows how to decode the []byte into an
	// intermediate form. At minimum, the right decoder must be chosen for the
	// format - for example, JSON vs YAML. For resource kinds, a decoder must
	// also know how to transform the input from a Kubernetes resource object
	// shape to Grafana's object shape. See [github.com/grafana/kindsys/encoding].
	Validate(b []byte, codec Decoder) error

	// FromBytes takes a []byte and a decoder, validates it against schema, and
	// if validation is successful, unmarshals it into an UnstructuredResource.
	FromBytes(b []byte, codec Decoder) (*UnstructuredResource, error)

	// Group returns the kind's group, as defined in the group field of the kind definition.
	//
	// This is equivalent to the group of a Kubernetes CRD.
	Group() string
}

// Core is the dynamically typed runtime representation of a Grafana core kind
// definition. It is one in a family of interfaces, see [Kind] for context.
//
// A Core kind provides interactions with its corresponding [Resource] using
// [UnstructuredResource].
type Core interface {
	ResourceKind

	// Def returns a wrapper around the underlying CUE value that represents the
	// loaded and validated kind definition.
	Def() Def[CoreProperties]

	// ToBytes takes a []byte and a decoder, validates it against schema, and
	// if validation is successful, unmarshals it into an UnstructuredResource.
	// ToBytes(UnstructuredResource, codec Encoder) ([]byte, error)
}

// Custom is the dynamically typed runtime representation of a Grafana custom kind
// definition. It is one in a family of interfaces, see [Kind] for context.
//
// A Custom kind provides interactions with its corresponding [Resource] using
// [UnstructuredResource].
//
// Custom kinds are declared in Grafana extensions, rather than in Grafana core. It
// is likely that this distinction will go away in the future, leaving only
// Custom kinds.
type Custom interface {
	ResourceKind

	// Def returns a wrapper around the underlying CUE value that represents the
	// loaded and validated kind definition.
	Def() Def[CustomProperties]
}

// Composable is the untyped runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// TODO sort out the Go type used for generic associated...objects? do we even need one?
type Composable interface {
	Kind

	// Def returns a wrapper around the underlying CUE value that represents the
	// loaded and validated kind definition.
	Def() Def[ComposableProperties]
}

// TypedCore is the statically typed runtime representation of a Grafana core
// kind definition. It is one in a family of interfaces, see [Kind] for context.
//
// A TypedCore provides typed interactions with the [Resource] type given as its
// generic type parameter. As it embeds [Core], untyped interaction is also available.
//
// A TypedCore is created by calling [BindCoreResource] on a [Core] with a
// Go type to which it is assignable (see [thema.BindType]).
type TypedCore[R Resource] interface {
	Core

	// TypeFromBytes is the same as [Core.FromBytes], but returns an instance of the
	// associated generic struct type instead of an [UnstructuredResource].
	TypeFromBytes(b []byte, codec Decoder) (R, error)
}

// TypedCustom is the statically typed runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// A TypedCustom provides typed interactions with the [Resource] type given as its
// generic type parameter. As it embeds [Custom], untyped interaction is also available.
//
// A TypedCustom is created by calling [BindCustomResource] on a [Custom] with a
// Go type to which it is assignable (see [thema.BindType]).
type TypedCustom[R Resource] interface {
	Custom

	// TypeFromBytes is the same as [Custom.FromBytes], but returns an instance of the
	// associated generic struct type instead of an [UnstructuredResource].
	TypeFromBytes(b []byte, codec Decoder) (R, error)
}

// Decoder takes a []byte representing a serialized resource and decodes it into
// the intermediate [encoding.GrafanaShapeBytes] form. Implementations should
// vary in the form of the []byte they expect to take - e.g. JSON vs. YAML;
// Kubernetes shape vs. Grafana shape.
//
// TODO do less with these by deciding on a single, consistent object format
type Decoder interface {
	Decode(b []byte) (encoding.GrafanaShapeBytes, error)
}

type Encoder interface {
	Encode(bytes encoding.GrafanaShapeBytes) ([]byte, error)
}
