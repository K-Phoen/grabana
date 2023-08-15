package golang

import (
	"sort"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/simplecue"
)

type preprocessor struct {
	types map[string]simplecue.TypeDefinition
}

func newPreprocessor() *preprocessor {
	return &preprocessor{
		types: make(map[string]simplecue.TypeDefinition),
	}
}

// inefficient, but I'm lazy. It's only used during code generation anyway.
func (preprocessor *preprocessor) sortedTypes() []simplecue.TypeDefinition {
	typeNames := make([]string, 0, len(preprocessor.types))
	for typeName := range preprocessor.types {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	sorted := make([]simplecue.TypeDefinition, 0, len(preprocessor.types))
	for _, k := range typeNames {
		sorted = append(sorted, preprocessor.types[k])
	}

	return sorted
}

func (preprocessor *preprocessor) translateTypes(definitions []simplecue.TypeDefinition) {
	for _, typeDef := range definitions {
		preprocessor.translate(typeDef)
	}
}

func (preprocessor *preprocessor) translate(def simplecue.TypeDefinition) {
	preprocessor.types[def.Name] = preprocessor.translateTypeDefinition(def)
}

func (preprocessor *preprocessor) translateTypeDefinition(def simplecue.TypeDefinition) simplecue.TypeDefinition {
	newFields := make([]simplecue.FieldDefinition, 0, len(def.Fields))
	for _, fieldDef := range def.Fields {
		newFields = append(newFields, preprocessor.translateFieldDefinition(fieldDef))
	}

	newDef := def
	newDef.Fields = newFields

	return newDef
}

func (preprocessor *preprocessor) translateFieldDefinition(def simplecue.FieldDefinition) simplecue.FieldDefinition {
	if def.Type.Type != simplecue.TypeDisjunction {
		return def
	}

	newDef := def
	newDef.Type = preprocessor.expandDisjunction(def.Type)

	return newDef
}

func (preprocessor *preprocessor) expandDisjunction(def simplecue.FieldType) simplecue.FieldType {
	// Ex: type | null
	if len(def.SubType) == 2 && def.SubType.HasNullType() {
		finalType := def.SubType.NonNullTypes()[0]

		return simplecue.FieldType{
			Type:        finalType.Type,
			Nullable:    true,
			Constraints: finalType.Constraints,
		}
	}

	// type | otherType | something (| null)?
	// generate a type with a nullable field for every branch of the disjunction,
	// add it to preprocessor.types, and use it instead.
	newTypeName := preprocessor.disjunctionTypeName(def.SubType)

	if _, ok := preprocessor.types[newTypeName]; !ok {
		newType := simplecue.TypeDefinition{
			Type: simplecue.DefinitionStruct,
			Name: newTypeName,
		}

		for _, subType := range def.SubType {
			if subType.IsNull() {
				continue
			}

			newType.Fields = append(newType.Fields, simplecue.FieldDefinition{
				Name: "Val" + strings.Title(string(subType.Type)),
				Type: simplecue.FieldType{
					Nullable:    true,
					Type:        subType.Type,
					SubType:     subType.SubType,
					Constraints: subType.Constraints,
				},
				Required: false,
			})
		}

		preprocessor.types[newTypeName] = newType
	}

	return simplecue.FieldType{
		Type:     simplecue.TypeID(newTypeName),
		Nullable: def.SubType.HasNullType(),
	}
}

func (preprocessor *preprocessor) disjunctionTypeName(disjunctionTypes simplecue.FieldTypes) string {
	parts := make([]string, 0, len(disjunctionTypes))

	for _, subType := range disjunctionTypes {
		parts = append(parts, strings.Title(string(subType.Type)))
	}

	return strings.Title(strings.Join(parts, "Or"))
}
