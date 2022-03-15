package fields

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestUnit(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	Unit("short")(overrideCfg)

	req.Len(overrideCfg.Properties, 1)
	req.Equal("unit", overrideCfg.Properties[0].ID)
	req.Equal("short", overrideCfg.Properties[0].Value)
}

func TestFillOpacity(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	FillOpacity(0)(overrideCfg)

	req.Len(overrideCfg.Properties, 1)
	req.Equal("custom.fillOpacity", overrideCfg.Properties[0].ID)
	req.Equal(0, overrideCfg.Properties[0].Value)
}

func TestFixedColorScheme(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	FixedColorScheme("dark-blue")(overrideCfg)

	req.Len(overrideCfg.Properties, 1)
	req.Equal("color", overrideCfg.Properties[0].ID)

	values := overrideCfg.Properties[0].Value.(map[string]string)
	req.Equal("fixed", values["mode"])
	req.Equal("dark-blue", values["fixedColor"])
}
