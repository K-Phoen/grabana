package scheme

import (
	"github.com/K-Phoen/sdk"
)

type ColorMode string

const (
	Last ColorMode = "last"
	Min  ColorMode = "min"
	Max  ColorMode = "max"
)

// Option represents an option that can be used to configure an axis.
type Option func(scheme *Scheme)

type Step struct {
	Color string
	Value int
}

// Scheme represents a color scheme.
type Scheme struct {
	fieldConfig *sdk.FieldConfig
}

// New creates a new Scheme configuration.
func New(fieldConfig *sdk.FieldConfig, options ...Option) *Scheme {
	scheme := &Scheme{fieldConfig: fieldConfig}

	for _, opt := range options {
		opt(scheme)
	}

	return scheme
}

// SingleColor defines the color scheme with a single color.
func SingleColor(color string) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "fixed"
		thresholds.fieldConfig.Defaults.Color.FixedColor = color
	}
}

// ClassicPalette uses the classic palette color scheme.
func ClassicPalette() Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "palette-classic"
	}
}

// ThresholdsValue uses the thresholds colors.
func ThresholdsValue(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "thresholds"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}

// GreenYellowRed uses the green-yellow-red color scheme.
func GreenYellowRed(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "continuous-GrYlRd"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}

// YellowRed uses the yellow-red color scheme.
func YellowRed(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "continuous-YlRd"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}

// YellowBlue uses the yellow-blue color scheme.
func YellowBlue(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "continuous-YlBl"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}

// RedYellowGreen uses the red-yellow-green color scheme.
func RedYellowGreen(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "continuous-RdYlGr"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}

// BlueYellowRed uses the blue-yellow-red color scheme.
func BlueYellowRed(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "continuous-BlYlRd"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}

// BluePurple uses the blue-purple color scheme.
func BluePurple(colorBy ColorMode) Option {
	return func(thresholds *Scheme) {
		thresholds.fieldConfig.Defaults.Color.Mode = "continuous-BlPu"
		thresholds.fieldConfig.Defaults.Color.SeriesBy = string(colorBy)
	}
}
