package kindsys

import (
	"github.com/grafana/thema"
)

// genericCore is a dynamically typed representation of a parsed and
// validated [Custom] kind, implemented with thema.
type genericCustom struct {
	def Def[CustomProperties]
	lin thema.Lineage
}

func (k genericCustom) FromBytes(b []byte, codec Decoder) (*UnstructuredResource, error) {
	inst, err := bytesToAnyInstance(k, b, codec)
	if err != nil {
		return nil, err
	}
	// we have a valid instance! decode into unstructured
	// TODO implement me
	_ = inst
	panic("implement me")
}

func (k genericCustom) Validate(b []byte, codec Decoder) error {
	// TODO implement me
	panic("implement me")
}

func (k genericCustom) CurrentVersion() thema.SyntacticVersion {
	return k.def.Properties.CurrentVersion
}

func (k genericCustom) Group() string {
	return k.def.Properties.CRD.Group
}

var _ Custom = genericCustom{}

// Props returns the generic SomeKindProperties
func (k genericCustom) Props() SomeKindProperties {
	return k.def.Properties
}

// Name returns the Name property
func (k genericCustom) Name() string {
	return k.def.Properties.Name
}

// MachineName returns the MachineName property
func (k genericCustom) MachineName() string {
	return k.def.Properties.MachineName
}

// Maturity returns the Maturity property
func (k genericCustom) Maturity() Maturity {
	return k.def.Properties.Maturity
}

// Def returns a Def with the type of ExtendedProperties, containing the bound ExtendedProperties
func (k genericCustom) Def() Def[CustomProperties] {
	return k.def
}

// Lineage returns the underlying bound Lineage
func (k genericCustom) Lineage() thema.Lineage {
	return k.lin
}

// BindCustom creates a Custom-implementing type from a def, runtime, and opts
//
//nolint:lll
func BindCustom(rt *thema.Runtime, def Def[CustomProperties], opts ...thema.BindOption) (Custom, error) {
	lin, err := def.Some().BindKindLineage(rt, opts...)
	if err != nil {
		return nil, err
	}

	return genericCustom{
		def: def,
		lin: lin,
	}, nil
}

// TODO docs
func BindCustomResource[R Resource](k Custom) (TypedCustom[R], error) {
	// TODO implement me
	panic("implement me")
}
