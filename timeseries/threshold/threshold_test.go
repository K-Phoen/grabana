package threshold

import (
	"fmt"
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestThresholdModeCanBeConfigured(t *testing.T) {
	testCases := []struct {
		value    Mode
		expected string
	}{
		{value: Absolute, expected: "absolute"},
		{value: Percentage, expected: "percentage"},
	}

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("mode %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			cfg := &sdk.FieldConfig{}
			New(cfg, ValueMode(tc.value))

			req.Equal(tc.expected, cfg.Defaults.Thresholds.Mode)
		})
	}
}

func TestThresholdStyleCanBeConfigured(t *testing.T) {
	testCases := []struct {
		value    DisplayStyle
		expected string
	}{
		{value: Off, expected: "off"},
		{value: AsFilledRegions, expected: "area"},
		{value: AsLines, expected: "line"},
		{value: Both, expected: "line+area"},
	}

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("mode %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			cfg := &sdk.FieldConfig{}
			New(cfg, Style(tc.value))

			req.Equal(tc.expected, cfg.Defaults.Custom.ThresholdsStyle.Mode)
		})
	}
}

func TestBaseColorCanBeConfigured(t *testing.T) {
	req := require.New(t)

	threshold := New(&sdk.FieldConfig{}, BaseColor("red"))

	req.Equal("red", threshold.baseColor)
}

func TestStepsCanBeDefined(t *testing.T) {
	req := require.New(t)

	cfg := &sdk.FieldConfig{}
	New(cfg, BaseColor("red"), Steps(Step{
		Color: "green",
		Value: 10,
	}))

	steps := cfg.Defaults.Thresholds.Steps

	req.Len(steps, 2)

	// Base step
	req.Nil(steps[0].Value)
	req.Equal("red", steps[0].Color)

	// Threshold value
	req.Equal(10, *steps[1].Value)
	req.Equal("green", steps[1].Color)
}
