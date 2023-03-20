package fields

import (
	"github.com/K-Phoen/sdk"
)

type Matcher func(field *sdk.FieldConfigOverride)

type FieldType string

const (
	FieldTypeTime FieldType = "time"
)

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

// ByRegex matches fields names using a regex.
func ByRegex(regex string) Matcher {
	return func(field *sdk.FieldConfigOverride) {
		field.Matcher.ID = "byRegexp"
		field.Matcher.Options = regex
	}
}

// ByType matches fields with a specific type.
func ByType(fieldType FieldType) Matcher {
	return func(field *sdk.FieldConfigOverride) {
		field.Matcher.ID = "byType"
		field.Matcher.Options = string(fieldType)
	}
}
