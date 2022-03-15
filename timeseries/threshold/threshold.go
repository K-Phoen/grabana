package threshold

import (
	"github.com/K-Phoen/sdk"
)

// Mode represents the type of value used as threshold.
type Mode string

const (
	Percentage Mode = "percentage"
	Absolute   Mode = "absolute"
)

// DisplayStyle represents how the threshold should be visualized.
type DisplayStyle string

const (
	Off             DisplayStyle = "off"
	AsFilledRegions DisplayStyle = "area"
	AsLines         DisplayStyle = "line"
	Both            DisplayStyle = "line+area"
)

// Option represents an option that can be used to configure an axis.
type Option func(threshold *Threshold)

type Step struct {
	Color string
	Value int
}

// Threshold represents a threshold visualization.
type Threshold struct {
	baseColor   string
	fieldConfig *sdk.FieldConfig
}

// New creates a new Threshold configuration.
func New(fieldConfig *sdk.FieldConfig, options ...Option) *Threshold {
	threshold := &Threshold{fieldConfig: fieldConfig}

	defaultOpts := []Option{
		Style(AsLines),
		ValueMode(Absolute),
		BaseColor("green"),
	}
	for _, opt := range append(defaultOpts, options...) {
		opt(threshold)
	}

	return threshold
}

// Style defines the thresholds display style.
func Style(style DisplayStyle) Option {
	return func(thresholds *Threshold) {
		thresholds.fieldConfig.Defaults.Custom.ThresholdsStyle.Mode = string(style)
	}
}

// BaseColor defines the color of the thresholds' base.
func BaseColor(color string) Option {
	return func(thresholds *Threshold) {
		thresholds.baseColor = color
	}
}

// ValueMode defines how to interpret the threshold values.
func ValueMode(mode Mode) Option {
	return func(thresholds *Threshold) {
		thresholds.fieldConfig.Defaults.Thresholds.Mode = string(mode)
	}
}

// Steps defines threshold steps.
func Steps(steps ...Step) Option {
	return func(thresholds *Threshold) {
		sdkSteps := make([]sdk.ThresholdStep, 0, len(steps))

		for i := range steps {
			sdkSteps = append(sdkSteps, sdk.ThresholdStep{
				Color: steps[i].Color,
				Value: &steps[i].Value,
			})
		}

		thresholds.fieldConfig.Defaults.Thresholds.Steps = append(
			// Base
			[]sdk.ThresholdStep{{Color: thresholds.baseColor}},
			// User-defined steps
			sdkSteps...,
		)
	}
}
