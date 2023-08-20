package kindsys

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/grafana/thema"
)

// SomeDef represents a single kind definition, having been loaded and
// validated by a func such as [LoadCoreKindDef].
//
// The underlying type of the Properties field indicates the category of kind.
type SomeDef struct {
	// V is the cue.Value containing the entire Kind definition.
	V cue.Value
	// Properties contains the kind's declarative non-schema properties.
	Properties SomeKindProperties
}

// BindKindLineage binds the lineage for the kind definition.
//
// For kinds with a corresponding Go type, it is left to the caller to associate
// that Go type with the lineage returned from this function by a call to
// [thema.BindType].
func (def SomeDef) BindKindLineage(rt *thema.Runtime, opts ...thema.BindOption) (thema.Lineage, error) {
	if rt == nil {
		return nil, fmt.Errorf("nil thema.Runtime")
	}
	return thema.BindLineage(def.V.LookupPath(cue.MakePath(cue.Str("lineage"))), rt, opts...)
}

// IsCore indicates whether the represented kind is a core kind.
func (def SomeDef) IsCore() bool {
	_, is := def.Properties.(CoreProperties)
	return is
}

// IsCustom indicates whether the represented kind is a custom kind.
func (def SomeDef) IsCustom() bool {
	_, is := def.Properties.(CustomProperties)
	return is
}

// IsComposable indicates whether the represented kind is a composable kind.
func (def SomeDef) IsComposable() bool {
	_, is := def.Properties.(ComposableProperties)
	return is
}

// Def represents a single kind definition, having been loaded and validated by
// a func such as [LoadCoreKindDef].
//
// Its type parameter indicates the category of kind.
//
// Thema lineages in the contained definition have not yet necessarily been
// validated.
type Def[T KindProperties] struct {
	// V is the cue.Value containing the entire Kind definition.
	V cue.Value
	// Properties contains the kind's declarative non-schema properties.
	Properties T
}

// Some converts the typed Def to the equivalent typeless SomeDef.
func (def Def[T]) Some() SomeDef {
	return SomeDef{
		V:          def.V,
		Properties: any(def.Properties).(SomeKindProperties),
	}
}
