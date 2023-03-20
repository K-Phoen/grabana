package axis

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
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
type Option func(axis *Axis) error

// Axis represents a visualization axis.
type Axis struct {
	fieldConfig *sdk.FieldConfig
}

// New creates a new Axis configuration.
func New(fieldConfig *sdk.FieldConfig, options ...Option) (*Axis, error) {
	axis := &Axis{fieldConfig: fieldConfig}

	for _, opt := range options {
		if err := opt(axis); err != nil {
			return nil, err
		}
	}

	return axis, nil
}

// Placement defines how the axis should be placed in the panel.
func Placement(placement PlacementMode) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Custom.AxisPlacement = string(placement)

		return nil
	}
}

// SoftMin defines a soft minimum value for the axis.
func SoftMin(value int) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Custom.AxisSoftMin = &value

		return nil
	}
}

// SoftMax defines a soft maximum value for the axis.
func SoftMax(value int) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Custom.AxisSoftMax = &value

		return nil
	}
}

// Min defines a hard minimum value for the axis.
func Min(value float64) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Min = &value

		return nil
	}
}

// Max defines a hard maximum value for the axis.
func Max(value float64) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Max = &value

		return nil
	}
}

// Unit sets the unit of the data displayed in this series.
func Unit(unit string) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Unit = unit

		return nil
	}
}

// Scale sets the scale to use for the Y-axis values..
func Scale(mode ScaleMode) Option {
	return func(axis *Axis) error {
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

		return nil
	}
}

// Label sets a Y-axis text label.
func Label(label string) Option {
	return func(axis *Axis) error {
		axis.fieldConfig.Defaults.Custom.AxisLabel = label

		return nil
	}
}

// Decimals sets how many decimal points should be displayed.
func Decimals(decimals int) Option {
	return func(axis *Axis) error {
		if decimals < 0 {
			return fmt.Errorf("decimals must be greater than 0: %w", errors.ErrInvalidArgument)
		}

		axis.fieldConfig.Defaults.Decimals = &decimals

		return nil
	}
}
