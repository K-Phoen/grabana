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

	return codejen.NewFile(file.Package+"_types_gen.go", output, jenny), nil
}

func (jenny GoRawTypes) generateFile(file *ast.File) ([]byte, error) {
	var buffer strings.Builder
	tr := newPreprocessor()

	tr.translateTypes(file.Types)

	buffer.WriteString(fmt.Sprintf("package %s\n\n", file.Package))

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

func (jenny GoRawTypes) formatTypeDef(def ast.TypeDefinition) ([]byte, error) {
	if def.Type == ast.DefinitionStruct {
		return jenny.formatStructDef(def)
	}

	return jenny.formatEnumDef(def)
}

func (jenny GoRawTypes) formatEnumDef(def ast.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("type %s %s\n", def.Name, def.SubType))

	buffer.WriteString("const (\n")
	for _, val := range def.Values {
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", strings.Title(val.Name), def.Name, val.Value))
	}
	buffer.WriteString(")\n")

	return []byte(buffer.String()), nil
}

func (jenny GoRawTypes) formatStructDef(def ast.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("type %s struct {\n", def.Name))

	for _, fieldDef := range def.Fields {
		fieldDefGen, err := jenny.formatField(fieldDef)
		if err != nil {
			return nil, err
		}

		buffer.WriteString("\t" + string(fieldDefGen))
	}

	buffer.WriteString("}\n")

	return []byte(buffer.String()), nil
}

func (jenny GoRawTypes) formatField(def ast.FieldDefinition) ([]byte, error) {
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

func formatType(def ast.FieldType) string {
	if def.Type == ast.TypeAny {
		return "any"
	}

	if def.Type == ast.TypeDisjunction {
		return formatDisjunction(def)
	}

	if def.Type == ast.TypeArray {
		return formatArray(def)
	}

	typeName := string(def.Type)
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

func formatArray(def ast.FieldType) string {
	var subTypeString string

	// we don't know what to do here (yet)
	if len(def.SubType) != 1 {
		subTypeString = formatDisjunction(ast.FieldType{
			SubType: def.SubType,
		})
	} else {
		subTypeString = formatType(def.SubType[0])
	}

	return fmt.Sprintf("[]%s", subTypeString)
}

func formatDisjunction(def ast.FieldType) string {
	typeName := string(def.Type)
	if def.SubType != nil {
		subTypes := make([]string, 0, len(def.SubType))
		for _, subType := range def.SubType {
			subTypes = append(subTypes, formatType(subType))
		}

		typeName = fmt.Sprintf("%s<%s>", typeName, strings.Join(subTypes, " | "))
	}

	return typeName
}
