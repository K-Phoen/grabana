package axis

import (
	"fmt"
	"testing"

	"github.com/K-Phoen/grabana/errors"
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

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("placement %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			cfg := &sdk.FieldConfig{}
			_, err := New(cfg, Placement(tc.value))

			req.NoError(err)
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

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("placement %s %d", tc.expectedType, tc.expectedLog), func(t *testing.T) {
			req := require.New(t)

			cfg := &sdk.FieldConfig{}
			_, err := New(cfg, Scale(tc.value))

			req.NoError(err)
			req.Equal(tc.expectedType, cfg.Defaults.Custom.ScaleDistribution.Type)
			req.Equal(tc.expectedLog, cfg.Defaults.Custom.ScaleDistribution.Log)
		})
	}
}

func TestAxisSoftMinCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, SoftMin(0))

	req.NoError(err)
	req.Equal(0, *cfg.Defaults.Custom.AxisSoftMin)
}

func TestAxisSoftMaxCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, SoftMax(0))

	req.NoError(err)
	req.Equal(0, *cfg.Defaults.Custom.AxisSoftMax)
}

func TestAxisMinCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, Min(1.1))

	req.NoError(err)
	req.Equal(1.1, *cfg.Defaults.Min)
}

func TestAxisMaxCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, Max(2.2))

	req.NoError(err)
	req.Equal(2.2, *cfg.Defaults.Max)
}

func TestLabelCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, Label("Foo"))

	req.NoError(err)
	req.Equal("Foo", cfg.Defaults.Custom.AxisLabel)
}

func TestDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, Decimals(2))

	req.NoError(err)
	req.Equal(2, *cfg.Defaults.Decimals)
}

func TestInvalidDecimalsAreRejected(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, Decimals(-2))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestUnitCanBeConfigured(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	_, err := New(cfg, Unit("reqps"))

	req.NoError(err)
	req.Equal("reqps", cfg.Defaults.Unit)
}
