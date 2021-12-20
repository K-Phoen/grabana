package fields

import "github.com/K-Phoen/sdk"

type OverrideOption func(field *sdk.FieldConfigOverride)

func Unit(unit string) OverrideOption {
	return func(field *sdk.FieldConfigOverride) {
		field.Properties = append(field.Properties,
			sdk.FieldConfigOverrideProperty{
				ID:    "unit",
				Value: unit,
			})
	}
}
