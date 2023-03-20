package fields

import (
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/sdk"
)

// StackMode configures mode of series stacking.
// FIXME: copied here to avoid circular imports with parent package
type StackMode string

const (
	// Unstacked will not stack series
	Unstacked StackMode = "none"
	// NormalStack will stack series as absolute numbers
	NormalStack StackMode = "normal"
	// PercentStack will stack series as percents
	PercentStack StackMode = "percent"
)

type OverrideOption func(field *sdk.FieldConfigOverride)

// Unit overrides the unit.
func Unit(unit string) OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID:    "unit",
				Value: unit,
			})
	}
}

// FillOpacity overrides the opacity.
func FillOpacity(opacity int) OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID:    "custom.fillOpacity",
				Value: opacity,
			})
	}
}

// FixedColorScheme forces the use of a fixed color scheme.
func FixedColorScheme(color string) OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID: "color",
				Value: map[string]string{
					"mode":       "fixed",
					"fixedColor": color,
				},
			})
	}
}

// NegativeY flips the results to negative values on the Y axis.
func NegativeY() OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID:    "custom.transform",
				Value: "negative-Y",
			})
	}
}

// AxisPlacement overrides how the axis should be placed in the panel.
func AxisPlacement(placement axis.PlacementMode) OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID:    "custom.axisPlacement",
				Value: string(placement),
			})
	}
}

// Stack overrides if the series should be stacked and using which mode (default not stacked).
func Stack(mode StackMode) OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID: "custom.stacking",
				Value: map[string]interface{}{
					"group": false,
					"mode":  string(mode),
				},
			})
	}
}
