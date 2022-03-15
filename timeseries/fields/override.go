package fields

import "github.com/K-Phoen/sdk"

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
