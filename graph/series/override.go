package series

import (
	"github.com/K-Phoen/sdk"
)

// OverrideOption represents an option that can be used alter a graph panel series.
type OverrideOption func(series *sdk.SeriesOverride)

// Alis defines an alias/regex used to identify the series to override.
func Alias(alias string) OverrideOption {
	return func(series *sdk.SeriesOverride) {
		series.Alias = alias
	}
}

// Color overrides the color for the matched series.
func Color(color string) OverrideOption {
	return func(series *sdk.SeriesOverride) {
		series.Color = &color
	}
}

// Dashes enables/disables display of the series using dashes instead of lines.
func Dashes(enabled bool) OverrideOption {
	return func(series *sdk.SeriesOverride) {
		series.Dashes = &enabled
	}
}

// Lines enables/disables display of the series using dashes instead of dashes.
func Lines(enabled bool) OverrideOption {
	return func(series *sdk.SeriesOverride) {
		series.Lines = &enabled
	}
}

func Fill(opacity int) OverrideOption {
	return func(series *sdk.SeriesOverride) {
		series.Fill = &opacity
	}
}

func LineWidth(width int) OverrideOption {
	return func(series *sdk.SeriesOverride) {
		series.LineWidth = &width
	}
}
