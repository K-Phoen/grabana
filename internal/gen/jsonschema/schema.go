package jsonschema

import (
	"encoding/json"
	"strings"
)

// Naive types to represent a simplification of a draft7 jsonschema

type Type string

const (
	TypeNull    Type = "null"
	TypeBoolean Type = "boolean"
	TypeObject  Type = "object"
	TypeArray   Type = "array"
	TypeNumber  Type = "number"
	TypeString  Type = "string"
	TypeInteger Type = "integer"
)

type TypeSet []Type

func (typeSet TypeSet) IsDisjunction() bool {
	return len(typeSet) > 1
}

func (typeSet TypeSet) String() string {
	types := make([]string, 0, len(typeSet))
	for _, item := range typeSet {
		types = append(types, string(item))
	}

	return "[" + strings.Join(types, ", ") + "]"
}

func (typeSet TypeSet) Any(types ...Type) bool {
	for _, inputType := range types {
		for _, ts := range typeSet {
			if inputType == ts {
				return true
			}
		}
	}

	return false
}

func (typeSet TypeSet) Exactly(types ...Type) bool {
	if len(typeSet) != len(types) {
		return false
	}

	for _, inputType := range types {
		found := false
		for _, ts := range typeSet {
			if inputType == ts {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (ts *TypeSet) UnmarshalJSON(b []byte) error {
	if b[0] == '[' {
		type rawTypeSet TypeSet
		out := (*rawTypeSet)(ts)
		return json.Unmarshal(b, out)
	} else {
		var t Type
		err := json.Unmarshal(b, &t)
		if err != nil {
			*ts = nil
		} else {
			*ts = []Type{t}
		}
		return err
	}
}

type Schema struct {
	// Core
	Schema      string            `json:"$schema"`
	Vocabulary  map[string]bool   `json:"$vocabulary"`
	ID          string            `json:"$id"`
	Ref         string            `json:"$ref"`
	DynamicRef  string            `json:"$dynamicRef"`
	Definitions map[string]Schema `json:"definitions"`
	Comment     string            `json:"$comment"`

	// Applying subschemas with logic
	AllOf []Schema `json:"allOf"`
	AnyOf []Schema `json:"anyOf"`
	OneOf []Schema `json:"oneOf"`

	// Applying subschemas to arrays
	PrefixItems []Schema `json:"prefixItems"`
	Items       *Schema  `json:"items"`
	Contains    *Schema  `json:"contains"`

	// Applying subschemas to objects
	Properties           map[string]Schema `json:"properties"`
	PatternProperties    map[string]Schema `json:"patternProperties"`
	AdditionalProperties any               `json:"additionalProperties"` // nil or bool or *Schema.
	PropertyNames        *Schema           `json:"propertyNames"`

	// Validation
	Type  TypeSet       `json:"type"`
	Enum  []interface{} `json:"enum"`
	Const interface{}   `json:"const"`

	// Validation for numbers
	MultipleOf       json.Number `json:"multipleOf"`
	Maximum          json.Number `json:"maximum"`
	ExclusiveMaximum json.Number `json:"exclusiveMaximum"`
	Minimum          json.Number `json:"minimum"`
	ExclusiveMinimum json.Number `json:"exclusiveMinimum"`

	// Validation for strings
	MaxLength int    `json:"maxLength"`
	MinLength int    `json:"minLength"`
	Pattern   string `json:"pattern"`

	// Validation for arrays
	MaxItems    int  `json:"maxItems"`
	MinItems    int  `json:"minItems"`
	UniqueItems bool `json:"uniqueItems"`
	MaxContains int  `json:"maxContains"`
	MinContains int  `json:"minContains"`

	// Validation for objects
	MaxProperties     int                 `json:"maxProperties"`
	MinProperties     int                 `json:"minProperties"`
	Required          []string            `json:"required"`
	DependentRequired map[string][]string `json:"dependentRequired"`

	// Basic metadata annotations
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Default     interface{}   `json:"default"`
	Deprecated  bool          `json:"deprecated"`
	ReadOnly    bool          `json:"readOnly"`
	WriteOnly   bool          `json:"writeOnly"`
	Examples    []interface{} `json:"examples"`
}
