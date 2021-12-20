package fields

import "github.com/K-Phoen/sdk"

type Matcher func(field *sdk.FieldConfigOverride)

func ByName(name string) Matcher {
	return func(field *sdk.FieldConfigOverride) {
		field.Matcher.ID = "byName"
		field.Matcher.Options = name
	}
}
