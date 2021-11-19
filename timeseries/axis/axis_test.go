package axis

import (
	"fmt"
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestAxisPlacementCanBeConfigured(t *testing.T) {
	testCases := []struct {
		value    PlacementMode
		expected string
	}{
		{value: Hidden, expected: "hidden"},
		{value: Auto, expected: "auto"},
		{value: Left, expected: "left"},
		{value: Right, expected: "right"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("placement %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			cfg := &sdk.FieldConfig{}
			New(cfg, Placement(tc.value))

			req.Equal(tc.expected, cfg.Defaults.Custom.AxisPlacement)
		})
	}
}

func TestAxisScaleCanBeConfigured(t *testing.T) {
	testCases := []struct {
		value        ScaleMode
		expectedType string
		expectedLog  int
	}{
		{value: Linear, expectedType: "linear"},
		{value: Log2, expectedType: "log", expectedLog: 2},
		{value: Log10, expectedType: "log", expectedLog: 10},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("placement %s %d", tc.expectedType, tc.expectedLog), func(t *testing.T) {
			req := require.New(t)

			cfg := &sdk.FieldConfig{}
			New(cfg, Scale(tc.value))

			req.Equal(tc.expectedType, cfg.Defaults.Custom.ScaleDistribution.Type)
			req.Equal(tc.expectedLog, cfg.Defaults.Custom.ScaleDistribution.Log)
		})
	}
}

func TestAxisSoftMinCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, SoftMin(0))

	req.Equal(0, *cfg.Defaults.Custom.AxisSoftMin)
}

func TestAxisSoftMaxCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, SoftMax(0))

	req.Equal(0, *cfg.Defaults.Custom.AxisSoftMax)
}

func TestAxisMinCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, Min(0))

	req.Equal(0, *cfg.Defaults.Min)
}

func TestAxisMaxCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, Max(0))

	req.Equal(0, *cfg.Defaults.Max)
}

func TestLabelCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, Label("Foo"))

	req.Equal("Foo", cfg.Defaults.Custom.AxisLabel)
}

func TestDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, Decimals(2))

	req.Equal(2, *cfg.Defaults.Decimals)
}

func TestUnitCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, Unit("reqps"))

	req.Equal("reqps", cfg.Defaults.Unit)
}
