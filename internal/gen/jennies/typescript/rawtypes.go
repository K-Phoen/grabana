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

func (jenny TypescriptRawTypes) formatTypeDef(def ast.Definition) ([]byte, error) {
	if def.Type == ast.TypeStruct {
		return jenny.formatStructDef(def)
	}

	return jenny.formatEnumDef(def)
}

func (jenny TypescriptRawTypes) formatEnumDef(def ast.Definition) ([]byte, error) {
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

func (jenny TypescriptRawTypes) formatStructDef(def ast.Definition) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(fmt.Sprintf("export interface %s ", def.Name))

	body, err := jenny.formatStructBody(def)
	if err != nil {
		return nil, nil
	}

	buffer.WriteString(body + "\n")

	return []byte(buffer.String()), nil
}

func (jenny TypescriptRawTypes) formatStructBody(def ast.Definition) (string, error) {
	var buffer strings.Builder

	buffer.WriteString("{\n")

	for i, fieldDef := range def.Fields {
		fieldDefGen, err := jenny.formatField(fieldDef)
		if err != nil {
			return "", err
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

	buffer.WriteString("\n}")

	return buffer.String(), nil
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

	formattedType, err := jenny.formatType(def.Type)
	if err != nil {
		return nil, err
	}

	buffer.WriteString(fmt.Sprintf(
		"%s%s: %s;\n",
		def.Name,
		required,
		formattedType,
	))

	return []byte(buffer.String()), nil
}

func (jenny TypescriptRawTypes) formatType(def ast.Definition) (string, error) {
	// todo: handle nullable
	// maybe if nullable, append | null to the type?
	switch def.Type {
	case ast.TypeDisjunction:
		return jenny.formatDisjunction(def)
	case ast.TypeArray:
		return jenny.formatArray(def)
	case ast.TypeStruct:
		return jenny.formatStructBody(def)
	case ast.TypeMap:
		return jenny.formatMap(def)

	case ast.TypeNull:
		return "null", nil
	case ast.TypeAny:
		return "any", nil

	case ast.TypeBytes, ast.TypeString:
		return "string", nil

	case ast.TypeFloat32, ast.TypeFloat64:
		return "number", nil
	case ast.TypeUint8, ast.TypeUint16, ast.TypeUint32, ast.TypeUint64:
		return "number", nil
	case ast.TypeInt8, ast.TypeInt16, ast.TypeInt32, ast.TypeInt64:
		return "number", nil

	case ast.TypeBool:
		return "boolean", nil

	default:
		return string(def.Type), nil
	}
}

func (jenny TypescriptRawTypes) formatArray(def ast.Definition) (string, error) {
	// we don't know what to do here (yet)
	subTypeString, err := jenny.formatType(*def.ValueType)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s[]", subTypeString), nil
}

func (jenny TypescriptRawTypes) formatDisjunction(def ast.Definition) (string, error) {
	typeName := string(def.Type)
	subTypes := make([]string, 0, len(def.Branches))
	for _, subType := range def.Branches {
		formatted, err := jenny.formatType(subType)
		if err != nil {
			return "", err
		}

		subTypes = append(subTypes, formatted)
	}

	typeName = strings.Join(subTypes, " | ")

	return typeName, nil
}

func (jenny TypescriptRawTypes) formatMap(def ast.Definition) (string, error) {
	keyTypeString := def.IndexType
	valueTypeString, err := jenny.formatType(*def.ValueType)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Record<%s, %s>", keyTypeString, valueTypeString), nil
}

func prefixLinesWith(input string, prefix string) string {
	lines := strings.Split(input, "\n")
	prefixed := make([]string, 0, len(lines))

	for _, line := range lines {
		prefixed = append(prefixed, prefix+line)
	}

	return strings.Join(prefixed, "\n")
}
