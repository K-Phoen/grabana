package axis

import (
	"fmt"
)

// Option represents an option that can be used to configure a Y axis.
type Option func(axis *YAxis)

type SDKAxis struct {
	Decimals    *int     `json:"decimals"`
	Format      string   `json:"format"`
	LogBase     int      `json:"logBase"`
	Show        bool     `json:"show"`
	Max         *string  `json:"max"`
	Min         *string  `json:"min"`
	SplitFactor *float64 `json:"splitFactor"`
}

// YAxis represents a the Y axis of a heatmap.
type YAxis struct {
	Builder *SDKAxis
}

// New creates a new YAxis configuration.
func New(options ...Option) *YAxis {
	axis := &YAxis{
		Builder: &SDKAxis{
			Format:  "short",
			LogBase: 1,
			Show:    true,
		},
	}

	for _, opt := range options {
		opt(axis)
	}

	return axis
}

// Unit sets the unit of the data displayed on this axis.
func Unit(unit string) Option {
	return func(axis *YAxis) {
		axis.Builder.Format = unit
	}
}

// Decimals set the number of decimals to be displayed on the axis.
func Decimals(decimals int) Option {
	return func(axis *YAxis) {
		axis.Builder.Decimals = &decimals
	}
}

// Min sets the minimum value expected on this axis.
func Min(min float64) Option {
	return func(axis *YAxis) {
		minStr := fmt.Sprintf("%f", min)
		axis.Builder.Min = &minStr
	}
}

// Max sets the maximum value expected on this axis.
func Max(max float64) Option {
	return func(axis *YAxis) {
		maxStr := fmt.Sprintf("%f", max)
		axis.Builder.Max = &maxStr
	}
}
