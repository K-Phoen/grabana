package jsonschema

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/ast"
)

type Config struct {
	// Package name used to generate code into.
	Package string
}

type newGenerator struct {
	schemaRootDefinition string

	file *ast.File
}

func GenerateAST(schemaReader io.Reader, c Config) (*ast.File, error) {
	g := &newGenerator{
		file: &ast.File{
			Package: c.Package,
		},
	}

	schema := Schema{}
	err := json.NewDecoder(schemaReader).Decode(&schema)
	if err != nil {
		return nil, err
	}

	g.schemaRootDefinition = schema.Ref

	for name, definition := range schema.Definitions {
		n, err := g.declareTopLevelType(name, definition)
		if err != nil {
			return nil, err
		}

		g.file.Types = append(g.file.Types, *n)
	}

	return g.file, nil
}

func (g *newGenerator) declareTopLevelType(name string, schema Schema) (*ast.Definition, error) {
	if schema.Enum != nil {
		return g.declareTopLevelEnum(name, schema)
	}

	if schema.Type.Exactly(TypeObject) {
		return g.declareTopLevelStruct(name, schema)
	}

	return nil, fmt.Errorf("unexpected top-level type '%s'", schema.Type)
}

func (g *newGenerator) declareTopLevelEnum(name string, schema Schema) (*ast.Definition, error) {
	if schema.Type.IsDisjunction() {
		return nil, fmt.Errorf("enums may only be generated from values of a single type: got '%s'", schema.Type)
	}

	if !schema.Type.Any(TypeString, TypeInteger, TypeNumber) {
		return nil, fmt.Errorf("enums may only be generated from strings, ints or numbers")
	}

	values, err := g.extractEnumValues(schema)
	if err != nil {
		return nil, err
	}

	typeDef := &ast.Definition{
		Type:         ast.TypeEnum,
		Values:       values,
		Name:         name,
		Comments:     schemaComments(schema),
		IsEntryPoint: "#/definitions/"+name == g.schemaRootDefinition,
	}

	return typeDef, nil
}

func (g *newGenerator) extractEnumValues(schema Schema) ([]ast.EnumValue, error) {
	fields := make([]ast.EnumValue, 0, len(schema.Enum))

	for _, value := range schema.Enum {
		fields = append(fields, ast.EnumValue{
			Type:  ast.TypeString,           // TODO
			Name:  fmt.Sprintf("%v", value), // TODO
			Value: value,
		})
	}

	return fields, nil
}

func (g *newGenerator) declareTopLevelStruct(name string, schema Schema) (*ast.Definition, error) {
	typeDef := &ast.Definition{
		Type:         ast.TypeStruct,
		Name:         name,
		Comments:     schemaComments(schema),
		IsEntryPoint: "#/definitions/"+name == g.schemaRootDefinition,
	}

	// explore struct fields
	for fieldName, property := range schema.Properties {
		node, err := g.declareNode(property)
		if err != nil {
			return nil, err
		}

		typeDef.Fields = append(typeDef.Fields, ast.FieldDefinition{
			Name:     fieldName,
			Comments: schemaComments(property),
			Required: stringInList(schema.Required, fieldName),
			Type:     *node,
		})
	}

	return typeDef, nil
}

func (g *newGenerator) declareNode(schema Schema) (*ast.Definition, error) {
	// This node is referring to another definition
	if schema.Ref != "" {
		parts := strings.Split(schema.Ref, "/")

		return &ast.Definition{
			Nullable: false,                           // TODO
			Type:     ast.TypeID(parts[len(parts)-1]), // this is definitely too naive
		}, nil
	}

	// Disjunctions
	if schema.Type.IsDisjunction() {
		return &ast.Definition{
			Type:     ast.TypeDisjunction,
			Branches: nil,   // TODO
			Nullable: false, // TODO
		}, nil
	}

	switch schema.Type[0] {
	case TypeNull:
		return &ast.Definition{Type: ast.TypeNull}, nil
	case TypeBoolean:
		return &ast.Definition{Type: ast.TypeBool}, nil
	case TypeString:
		return &ast.Definition{Type: ast.TypeString}, nil
	case TypeNumber, TypeInteger:
		return g.declareNumber(schema)
	case TypeArray:
		return g.declareList(schema)
	case TypeObject:
		return nil, fmt.Errorf("nested object definitions are not supported")
	default:
		return nil, fmt.Errorf("unexpected node with type '%s'", schema.Type.String())
	}
}

func (g *newGenerator) declareNumber(schema Schema) (*ast.Definition, error) {
	// TODO
	return &ast.Definition{
		Type:        ast.TypeInt64,
		Nullable:    false,
		Constraints: nil,
	}, nil
}

func (g *newGenerator) declareList(schema Schema) (*ast.Definition, error) {
	typeDef := &ast.Definition{
		Type:        ast.TypeArray,
		Nullable:    false,
		IndexType:   ast.TypeInt64,
		Constraints: nil,
	}

	expr, err := g.declareNode(*schema.Items)
	if err != nil {
		return nil, err
	}

	typeDef.ValueType = expr

	return typeDef, nil
}
