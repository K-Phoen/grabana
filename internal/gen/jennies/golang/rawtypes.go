package golang

import (
	"fmt"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/ast"
	"github.com/grafana/codejen"
)

type GoRawTypes struct {
}

func (jenny GoRawTypes) JennyName() string {
	return "GoRawTypes"
}

func (jenny GoRawTypes) Generate(file *ast.File) (*codejen.File, error) {
	output, err := jenny.generateFile(file)
	if err != nil {
		return nil, err
	}

	return codejen.NewFile("types/"+file.Package+"_types_gen.go", output, jenny), nil
}

func (jenny GoRawTypes) generateFile(file *ast.File) ([]byte, error) {
	var buffer strings.Builder
	tr := newPreprocessor()

	tr.translateDefinitions(file.Types)

	buffer.WriteString("package types\n\n")

	for _, typeDef := range tr.sortedTypes() {
		typeDefGen, err := jenny.formatTypeDef(typeDef)
		if err != nil {
			return nil, err
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func (jenny GoRawTypes) formatTypeDef(def ast.Definition) ([]byte, error) {
	if def.Type == ast.TypeStruct {
		return jenny.formatStructDef(def)
	}

	return jenny.formatEnumDef(def)
}

func (jenny GoRawTypes) formatEnumDef(def ast.Definition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("type %s %s\n", def.Name, def.Values[0].Type))

	buffer.WriteString("const (\n")
	for _, val := range def.Values {
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", strings.Title(val.Name), def.Name, val.Value))
	}
	buffer.WriteString(")\n")

	return []byte(buffer.String()), nil
}

func (jenny GoRawTypes) formatStructDef(def ast.Definition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("type %s ", def.Name))
	buffer.WriteString(formatStructBody(def, ""))
	buffer.WriteString("\n")

	return []byte(buffer.String()), nil
}

func formatStructBody(def ast.Definition, typesPkg string) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	for _, fieldDef := range def.Fields {
		buffer.WriteString("\t" + formatField(fieldDef, typesPkg))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func formatField(def ast.FieldDefinition, typesPkg string) string {
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
		formatType(def.Type, def.Required, typesPkg),
		def.Name,
		jsonOmitEmpty,
	))

	return buffer.String()
}

func formatType(def ast.Definition, fieldIsRequired bool, typesPkg string) string {
	if def.Type == ast.TypeAny {
		return "any"
	}

	if def.Type == ast.TypeDisjunction {
		return formatDisjunction(def, typesPkg)
	}

	if def.Type == ast.TypeArray {
		return formatArray(def, typesPkg)
	}

	if def.Type == ast.TypeMap {
		return formatMap(def, typesPkg)
	}

	// anonymous struct
	if def.Type == ast.TypeStruct {
		return formatStructBody(def, typesPkg)
	}

	typeName := string(def.Type)

	if def.IsReference() && typesPkg != "" {
		typeName = typesPkg + "." + typeName
	}

	if def.Nullable || !fieldIsRequired {
		typeName = "*" + typeName
	}

	return typeName
}

func formatArray(def ast.Definition, typesPkg string) string {
	subTypeString := formatType(*def.ValueType, true, typesPkg)

	return fmt.Sprintf("[]%s", subTypeString)
}

func formatMap(def ast.Definition, typesPkg string) string {
	keyTypeString := def.IndexType
	valueTypeString := formatType(*def.ValueType, true, typesPkg)

	return fmt.Sprintf("map[%s]%s", keyTypeString, valueTypeString)
}

func formatDisjunction(def ast.Definition, typesPkg string) string {
	typeName := string(def.Type)
	subTypes := make([]string, 0, len(def.Branches))
	for _, subType := range def.Branches {
		subTypes = append(subTypes, formatType(subType, true, typesPkg))
	}

	typeName = fmt.Sprintf("%s<%s>", typeName, strings.Join(subTypes, " | "))

	return typeName
}
