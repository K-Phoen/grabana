package fields

import (
	"testing"

	"github.com/K-Phoen/grabana/timeseries/axis"
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

func TestNegativeY(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	NegativeY()(overrideCfg)

	req.Len(overrideCfg.Properties, 1)
	req.Equal("custom.transform", overrideCfg.Properties[0].ID)
	req.Equal("negative-Y", overrideCfg.Properties[0].Value)
}

func TestAxisPlacement(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	AxisPlacement(axis.Hidden)(overrideCfg)

	req.Len(overrideCfg.Properties, 1)
	req.Equal("custom.axisPlacement", overrideCfg.Properties[0].ID)
	req.Equal("hidden", overrideCfg.Properties[0].Value)
}

func TestStack(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	Stack(PercentStack)(overrideCfg)

	req.Len(overrideCfg.Properties, 1)
	req.Equal("custom.stacking", overrideCfg.Properties[0].ID)

	values := overrideCfg.Properties[0].Value.(map[string]interface{})
	req.Equal("percent", values["mode"])
	req.Equal(false, values["group"])
}
