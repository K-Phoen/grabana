package series

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/sdk"
)

// OverrideOption represents an option that can be used alter a graph panel series.
type OverrideOption func(series *sdk.SeriesOverride) error

// Alias defines an alias/regex used to identify the series to override.
func Alias(alias string) OverrideOption {
	return func(series *sdk.SeriesOverride) error {
		series.Alias = alias

		return nil
	}
}

// Color overrides the color for the matched series.
func Color(color string) OverrideOption {
	return func(series *sdk.SeriesOverride) error {
		series.Color = &color

		return nil
	}
}

// Dashes enables/disables display of the series using dashes instead of lines.
func Dashes(enabled bool) OverrideOption {
	return func(series *sdk.SeriesOverride) error {
		series.Dashes = &enabled

		return nil
	}
}

// Lines enables/disables display of the series using dashes instead of dashes.
func Lines(enabled bool) OverrideOption {
	return func(series *sdk.SeriesOverride) error {
		series.Lines = &enabled

		return nil
	}
}

func Fill(opacity int) OverrideOption {
	return func(series *sdk.SeriesOverride) error {
		if opacity < 0 || opacity > 10 {
			return fmt.Errorf("fill must be between 0 and 10: %w", errors.ErrInvalidArgument)
		}

		series.Fill = &opacity

		return nil
	}
}

func LineWidth(width int) OverrideOption {
	return func(series *sdk.SeriesOverride) error {
		if width < 0 || width > 10 {
			return fmt.Errorf("line width must be between 0 and 10: %w", errors.ErrInvalidArgument)
		}

		series.LineWidth = &width

		return nil
	}
}
