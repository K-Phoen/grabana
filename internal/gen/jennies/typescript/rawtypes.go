package typescript

import (
	"fmt"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/ast"
	"github.com/grafana/codejen"
)

type TypescriptRawTypes struct {
}

func (jenny TypescriptRawTypes) JennyName() string {
	return "TypescriptRawTypes"
}

func (jenny TypescriptRawTypes) Generate(file *ast.File) (*codejen.File, error) {
	output, err := jenny.generateFile(file)
	if err != nil {
		return nil, err
	}

	return codejen.NewFile(file.Package+"_types_gen.ts", output, jenny), nil
}

func (jenny TypescriptRawTypes) generateFile(file *ast.File) ([]byte, error) {
	var buffer strings.Builder

	for _, typeDef := range file.Types {
		typeDefGen, err := jenny.formatTypeDef(typeDef)
		if err != nil {
			return nil, err
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func (jenny TypescriptRawTypes) formatTypeDef(def ast.TypeDefinition) ([]byte, error) {
	if def.Type == ast.DefinitionStruct {
		return jenny.formatStructDef(def)
	}

	return jenny.formatEnumDef(def)
}

func (jenny TypescriptRawTypes) formatEnumDef(def ast.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("export enum %s {\n", def.Name))
	for _, val := range def.Values {
		buffer.WriteString(fmt.Sprintf("\t%s = %#v,\n", strings.Title(val.Name), val.Value))
	}
	buffer.WriteString("}\n")

	return []byte(buffer.String()), nil
}

func (jenny TypescriptRawTypes) formatStructDef(def ast.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("export interface %s {\n", def.Name))

	for i, fieldDef := range def.Fields {
		fieldDefGen, err := jenny.formatField(fieldDef)
		if err != nil {
			return nil, err
		}

		buffer.WriteString(
			strings.TrimSuffix(
				prefixLinesWith(string(fieldDefGen), "\t"),
				"\n\t",
			),
		)

		if i != len(def.Fields)-1 {
			buffer.WriteString("\n")
		}
	}

	buffer.WriteString("\n}\n")

	return []byte(buffer.String()), nil
}

func (jenny TypescriptRawTypes) formatField(def ast.FieldDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	required := ""
	if !def.Required {
		required = "?"
	}

	buffer.WriteString(fmt.Sprintf(
		"%s%s: %s;\n",
		def.Name,
		required,
		formatType(def.Type),
	))

	return []byte(buffer.String()), nil
}

func formatType(def ast.FieldType) string {
	// todo: handle nullable
	// maybe if nullable, append | null to the type?
	switch def.Type {
	case ast.TypeDisjunction:
		return formatDisjunction(def)
	case ast.TypeArray:
		return formatArray(def)

	case ast.TypeNull:
		return "null"
	case ast.TypeAny:
		return "any"

	case ast.TypeBytes, ast.TypeString:
		return "string"

	case ast.TypeFloat32, ast.TypeFloat64:
		return "number"
	case ast.TypeUint8, ast.TypeUint16, ast.TypeUint32, ast.TypeUint64:
		return "number"
	case ast.TypeInt8, ast.TypeInt16, ast.TypeInt32, ast.TypeInt64:
		return "number"

	case ast.TypeBool:
		return "boolean"

	default:
		return string(def.Type)
	}
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

	return fmt.Sprintf("%s[]", subTypeString)
}

func formatDisjunction(def ast.FieldType) string {
	typeName := string(def.Type)
	if def.SubType != nil {
		subTypes := make([]string, 0, len(def.SubType))
		for _, subType := range def.SubType {
			subTypes = append(subTypes, formatType(subType))
		}

		typeName = strings.Join(subTypes, " | ")
	}

	return typeName
}

func prefixLinesWith(input string, prefix string) string {
	lines := strings.Split(input, "\n")
	prefixed := make([]string, 0, len(lines))

	for _, line := range lines {
		prefixed = append(prefixed, prefix+line)
	}

	return strings.Join(prefixed, "\n")
}
