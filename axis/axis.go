package axis

import (
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure an axis.
type Option func(axis *Axis)

// Axis represents a visualization axis.
type Axis struct {
	Builder *sdk.Axis
}

// New creates a new Axis configuration.
func New(options ...Option) *Axis {
	axis := &Axis{Builder: &sdk.Axis{
		Format:  "short",
		Show:    true,
		LogBase: 1,
	}}

	for _, opt := range options {
		opt(axis)
	}

	return axis
}

// Unit sets the unit of the data displayed on this axis.
func Unit(unit string) Option {
	return func(axis *Axis) {
		axis.Builder.Format = unit
	}
}

// Hide makes the axis hidden.
func Hide() Option {
	return func(axis *Axis) {
		axis.Builder.Show = false
	}
}

// LogBase allows to change the logarithmic scale used to display the values.
func LogBase(base int) Option {
	return func(axis *Axis) {
		axis.Builder.LogBase = base
	}
}

// Label sets the label on this axis.
func Label(label string) Option {
	return func(axis *Axis) {
		axis.Builder.Label = label
	}
}

// Min sets the minimum value expected on this axis.
func Min(min float64) Option {
	return func(axis *Axis) {
		axis.Builder.Min = sdk.NewFloatString(min)
	}
}

// Max sets the maximum value expected on this axis.
func Max(max float64) Option {
	return func(axis *Axis) {
		axis.Builder.Max = sdk.NewFloatString(max)
	}
}
