package axis

import (
	"github.com/K-Phoen/sdk"
)

// PlacementMode represents the axis display placement mode.
type PlacementMode string

const (
	Hidden PlacementMode = "hidden"
	Auto   PlacementMode = "auto"
	Left   PlacementMode = "left"
	Right  PlacementMode = "right"
)

// ScaleMode represents the axis scale distribution.
type ScaleMode uint8

const (
	Linear ScaleMode = iota
	Log2
	Log10
)

// Option represents an option that can be used to configure an axis.
type Option func(axis *Axis)

// Axis represents a visualization axis.
type Axis struct {
	fieldConfig *sdk.FieldConfig
}

// New creates a new Axis configuration.
func New(fieldConfig *sdk.FieldConfig, options ...Option) *Axis {
	axis := &Axis{fieldConfig: fieldConfig}

	for _, opt := range options {
		opt(axis)
	}

	return axis
}

// Placement defines how the axis should be placed in the panel.
func Placement(placement PlacementMode) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Custom.AxisPlacement = string(placement)
	}
}

// SoftMin defines a soft minimum value for the axis.
func SoftMin(value int) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Custom.AxisSoftMin = &value
	}
}

// SoftMax defines a soft maximum value for the axis.
func SoftMax(value int) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Custom.AxisSoftMax = &value
	}
}

// Min defines a hard minimum value for the axis.
func Min(value int) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Min = &value
	}
}

// SoftMax defines a hard maximum value for the axis.
func Max(value int) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Max = &value
	}
}

// Unit sets the unit of the data displayed in this series.
func Unit(unit string) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Unit = unit
	}
}

// Scale sets the scale to use for the Y-axis values..
func Scale(mode ScaleMode) Option {
	return func(axis *Axis) {
		scaleConfig := struct {
			Type string `json:"type"`
			Log  int    `json:"log,omitempty"`
		}{
			Type: "linear",
		}

		switch mode {
		case Linear:
			scaleConfig.Type = "linear"
		case Log2:
			scaleConfig.Type = "log"
			scaleConfig.Log = 2
		case Log10:
			scaleConfig.Type = "log"
			scaleConfig.Log = 10
		}

		axis.fieldConfig.Defaults.Custom.ScaleDistribution = scaleConfig
	}
}

// Label sets a Y-axis text label.
func Label(label string) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Custom.AxisLabel = label
	}
}

// Decimals sets how many decimal points should be displayed.
func Decimals(decimals int) Option {
	return func(axis *Axis) {
		axis.fieldConfig.Defaults.Decimals = &decimals
	}
}
