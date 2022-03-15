package fields

import (
	"github.com/K-Phoen/sdk"
)

type Matcher func(field *sdk.FieldConfigOverride)

// ByName matches a specific field name.
func ByName(name string) Matcher {
	return func(field *sdk.FieldConfigOverride) {
		field.Matcher.ID = "byName"
		field.Matcher.Options = name
	}
}

// ByQuery matches all fields returned by the given query.
func ByQuery(ref string) Matcher {
	return func(field *sdk.FieldConfigOverride) {
		field.Matcher.ID = "byFrameRefID"
		field.Matcher.Options = ref
	}
}
