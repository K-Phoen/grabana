package golang

import (
	"sort"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/ast"
)

type preprocessor struct {
	types map[string]ast.TypeDefinition
}

func newPreprocessor() *preprocessor {
	return &preprocessor{
		types: make(map[string]ast.TypeDefinition),
	}
}

// inefficient, but I'm lazy. It's only used during code generation anyway.
func (preprocessor *preprocessor) sortedTypes() []ast.TypeDefinition {
	typeNames := make([]string, 0, len(preprocessor.types))
	for typeName := range preprocessor.types {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	sorted := make([]ast.TypeDefinition, 0, len(preprocessor.types))
	for _, k := range typeNames {
		sorted = append(sorted, preprocessor.types[k])
	}

	return sorted
}

func (preprocessor *preprocessor) translateTypes(definitions []ast.TypeDefinition) {
	for _, typeDef := range definitions {
		preprocessor.translate(typeDef)
	}
}

func (preprocessor *preprocessor) translate(def ast.TypeDefinition) {
	preprocessor.types[def.Name] = preprocessor.translateTypeDefinition(def)
}

func (preprocessor *preprocessor) translateTypeDefinition(def ast.TypeDefinition) ast.TypeDefinition {
	newFields := make([]ast.FieldDefinition, 0, len(def.Fields))
	for _, fieldDef := range def.Fields {
		newFields = append(newFields, preprocessor.translateFieldDefinition(fieldDef))
	}

	newDef := def
	newDef.Fields = newFields

	return newDef
}

func (preprocessor *preprocessor) translateFieldDefinition(def ast.FieldDefinition) ast.FieldDefinition {
	newDef := def
	newDef.Type = preprocessor.translateFieldType(def.Type)

	return newDef
}

// bool, string,..., [], disjunction
func (preprocessor *preprocessor) translateFieldType(def ast.FieldType) ast.FieldType {
	if def.Type == ast.TypeDisjunction || def.Type == ast.TypeArray {
		return preprocessor.expandDisjunction(def)
	}

	return def
}

// def is either a disjunction or a list of unknown sub-types
func (preprocessor *preprocessor) expandDisjunction(def ast.FieldType) ast.FieldType {
	if def.Type == ast.TypeArray {
		newSubTypes := make(ast.FieldTypes, 0, len(def.SubType))

		for _, subType := range def.SubType {
			newSubType := preprocessor.translateFieldType(subType)
			newSubTypes = append(newSubTypes, newSubType)
		}

		def.SubType = newSubTypes

		return def
	}

	// Ex: type | null
	if len(def.SubType) == 2 && def.SubType.HasNullType() {
		finalType := def.SubType.NonNullTypes()[0]

		return ast.FieldType{
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
		newType := ast.TypeDefinition{
			Type: ast.DefinitionStruct,
			Name: newTypeName,
		}

		for _, subType := range def.SubType {
			if subType.IsNull() {
				continue
			}

			newType.Fields = append(newType.Fields, ast.FieldDefinition{
				Name: "Val" + strings.Title(string(subType.Type)),
				Type: ast.FieldType{
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

	return ast.FieldType{
		Type:     ast.TypeID(newTypeName),
		Nullable: def.SubType.HasNullType(),
	}
}

func (preprocessor *preprocessor) disjunctionTypeName(disjunctionTypes ast.FieldTypes) string {
	parts := make([]string, 0, len(disjunctionTypes))

	for _, subType := range disjunctionTypes {
		parts = append(parts, strings.Title(string(subType.Type)))
	}

	return strings.Title(strings.Join(parts, "Or"))
}
