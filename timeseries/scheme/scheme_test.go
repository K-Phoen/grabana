package scheme

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestSingleColor(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, SingleColor("red"))

	req.Equal("fixed", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal("red", threshold.fieldConfig.Defaults.Color.FixedColor)
}

func TestClassicPalette(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, ClassicPalette())

	req.Equal("palette-classic", threshold.fieldConfig.Defaults.Color.Mode)
}

func TestThresholdsValue(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, ThresholdsValue(Max))

	req.Equal("thresholds", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}

func TestGreenYellowRed(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, GreenYellowRed(Max))

	req.Equal("continuous-GrYlRd", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}

func TestYellowRed(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, YellowRed(Max))

	req.Equal("continuous-YlRd", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}

func TestYellowBlue(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, YellowBlue(Max))

	req.Equal("continuous-YlBl", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}

func TestRedYellowGreen(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, RedYellowGreen(Max))

	req.Equal("continuous-RdYlGr", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}

func TestBlueYellowRed(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, BlueYellowRed(Max))

	req.Equal("continuous-BlYlRd", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}

func TestBluePurple(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	threshold := New(cfg, BluePurple(Max))

	req.Equal("continuous-BlPu", threshold.fieldConfig.Defaults.Color.Mode)
	req.Equal(string(Max), threshold.fieldConfig.Defaults.Color.SeriesBy)
}
