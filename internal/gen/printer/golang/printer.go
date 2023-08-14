package golang

import (
	"fmt"
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
func (translator *preprocessor) sortedTypes() []simplecue.TypeDefinition {
	typeNames := make([]string, 0, len(translator.types))
	for typeName := range translator.types {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	sorted := make([]simplecue.TypeDefinition, 0, len(translator.types))
	for _, k := range typeNames {
		sorted = append(sorted, translator.types[k])
	}

	return sorted
}

func (translator *preprocessor) translateTypes(definitions []simplecue.TypeDefinition) {
	for _, typeDef := range definitions {
		translator.translate(typeDef)
	}
}

func (translator *preprocessor) translate(def simplecue.TypeDefinition) {
	translator.types[def.Name] = translator.translateTypeDefinition(def)
}

func (translator *preprocessor) translateTypeDefinition(def simplecue.TypeDefinition) simplecue.TypeDefinition {
	newFields := make([]simplecue.FieldDefinition, 0, len(def.Fields))
	for _, fieldDef := range def.Fields {
		newFields = append(newFields, translator.translateFieldDefinition(fieldDef))
	}

	newDef := def
	newDef.Fields = newFields

	return newDef
}

func (translator *preprocessor) translateFieldDefinition(def simplecue.FieldDefinition) simplecue.FieldDefinition {
	if def.Type.Type != simplecue.TypeDisjunction {
		return def
	}

	newDef := def
	newDef.Type = translator.expandDisjunction(def.Type)

	return newDef
}

func (translator *preprocessor) expandDisjunction(def simplecue.FieldType) simplecue.FieldType {
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
	newTypeName := translator.disjunctionTypeName(def.SubType)

	if _, ok := translator.types[newTypeName]; !ok {
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

		translator.types[newTypeName] = newType
	}

	return simplecue.FieldType{
		Type:     simplecue.TypeID(newTypeName),
		Nullable: def.SubType.HasNullType(),
	}
}

func (translator *preprocessor) disjunctionTypeName(disjunctionTypes simplecue.FieldTypes) string {
	parts := make([]string, 0, len(disjunctionTypes))

	for _, subType := range disjunctionTypes {
		parts = append(parts, strings.Title(string(subType.Type)))
	}

	return strings.Title(strings.Join(parts, "Or"))
}

func Printer(file *simplecue.File) ([]byte, error) {
	var buffer strings.Builder
	tr := newPreprocessor()

	tr.translateTypes(file.Types)

	buffer.WriteString(fmt.Sprintf("package %s\n\n", file.Package))

	for _, typeDef := range tr.sortedTypes() {
		typeDefGen, err := formatTypeDef(typeDef)
		if err != nil {
			return nil, err
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func formatTypeDef(def simplecue.TypeDefinition) ([]byte, error) {
	if def.Type == simplecue.DefinitionStruct {
		return formatStructDef(def)
	}

	return formatEnumDef(def)
}

func formatEnumDef(def simplecue.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	enumTypeName := stripHashtag(def.Name)

	buffer.WriteString(fmt.Sprintf("type %s %s\n", enumTypeName, def.SubType))

	buffer.WriteString("const (\n")
	for _, val := range def.Values {
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", strings.Title(val.Name), enumTypeName, val.Value))
	}
	buffer.WriteString(")\n")

	return []byte(buffer.String()), nil
}

func formatStructDef(def simplecue.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("type %s struct {\n", stripHashtag(def.Name)))

	// TODO: fields
	for _, fieldDef := range def.Fields {
		fieldDefGen, err := formatField(fieldDef)
		if err != nil {
			return nil, err
		}

		buffer.WriteString("\t" + string(fieldDefGen))
	}

	buffer.WriteString("}\n")

	return []byte(buffer.String()), nil
}

func formatField(def simplecue.FieldDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	jsonOmitEmpty := ""
	if !def.Required {
		jsonOmitEmpty = ",omitempty"
	}

	buffer.WriteString(fmt.Sprintf(
		"%s %s `json:\"%s%s\"`\n",
		strings.Title(def.Name),
		formatType(def.Type),
		def.Name,
		jsonOmitEmpty,
	))

	return []byte(buffer.String()), nil
}

func formatType(def simplecue.FieldType) string {
	if def.Type == "unknown" {
		return "any"
	}

	if def.Type == "disjunction" {
		return formatDisjunction(def)
	}

	if def.Type == "array" {
		return formatArray(def)
	}

	typeName := stripHashtag(string(def.Type))
	if def.SubType != nil {
		subTypes := make([]string, 0, len(def.SubType))
		for _, subType := range def.SubType {
			subTypes = append(subTypes, formatType(subType))
		}

		typeName = fmt.Sprintf("%s<%s>", typeName, strings.Join(subTypes, " | "))
	}

	if def.Nullable {
		typeName = "*" + typeName
	}

	return typeName
}

func formatArray(def simplecue.FieldType) string {
	var subTypeString string

	// we don't know what to do here (yet)
	if len(def.SubType) != 1 {
		subTypeString = formatDisjunction(simplecue.FieldType{
			SubType: def.SubType,
		})
	} else {
		subTypeString = formatType(def.SubType[0])
	}

	return fmt.Sprintf("[]%s", subTypeString)
}

func formatDisjunction(def simplecue.FieldType) string {
	typeName := stripHashtag(string(def.Type))
	if def.SubType != nil {
		subTypes := make([]string, 0, len(def.SubType))
		for _, subType := range def.SubType {
			subTypes = append(subTypes, formatType(subType))
		}

		typeName = fmt.Sprintf("%s<%s>", typeName, strings.Join(subTypes, " | "))
	}

	return typeName
}

func stripHashtag(input string) string {
	return strings.TrimPrefix(input, "#")
}
