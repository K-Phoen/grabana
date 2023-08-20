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
	if def.Kind == ast.KindStruct {
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
	switch def.Kind {
	case ast.KindDisjunction:
		return jenny.formatDisjunction(def)
	case ast.KindArray:
		return jenny.formatArray(def)
	case ast.KindStruct:
		return jenny.formatStructBody(def)
	case ast.KindMap:
		return jenny.formatMap(def)

	case ast.KindNull:
		return "null", nil
	case ast.KindAny:
		return "any", nil

	case ast.KindBytes, ast.KindString:
		return "string", nil

	case ast.KindFloat32, ast.KindFloat64:
		return "number", nil
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return "number", nil
	case ast.KindInt8, ast.KintInt16, ast.KindInt32, ast.KindInt64:
		return "number", nil

	case ast.KindBool:
		return "boolean", nil

	default:
		return string(def.Kind), nil
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
	typeName := string(def.Kind)
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
