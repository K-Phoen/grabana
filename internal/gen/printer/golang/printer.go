package golang

import (
	"fmt"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/simplecue"
)

type translator struct {
	types map[string]simplecue.TypeDefinition
}

func newTranslator() *translator {
	return &translator{
		types: make(map[string]simplecue.TypeDefinition),
	}
}

func (translator *translator) emit(def simplecue.TypeDefinition) error {

	translator.types[def.Name] = def

	return nil
}

func Printer(file *simplecue.File) ([]byte, error) {
	var buffer strings.Builder
	tr := newTranslator()

	for _, typeDef := range file.Types {
		if err := tr.emit(typeDef); err != nil {
			return nil, err
		}
	}

	buffer.WriteString(fmt.Sprintf("package %s\n\n", file.Package))

	for _, typeDef := range file.Types {
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
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", val.Name, enumTypeName, val.Value))
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
	// we don't know what to do here (yet)
	if len(def.SubType) != 2 || !def.SubType.HasNullType() {
		return disjunctionPseudoCode(def)
	}

	finalType := def.SubType.NonNullTypes()[0]

	return fmt.Sprintf("*%s", formatType(finalType))
}

func disjunctionPseudoCode(def simplecue.FieldType) string {
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