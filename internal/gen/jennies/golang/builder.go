package golang

import (
	"fmt"
	"strings"

	"github.com/K-Phoen/grabana/internal/gen/simplecue"
	"github.com/grafana/codejen"
)

type GoBuilder struct {
}

func (jenny GoBuilder) JennyName() string {
	return "GoRawTypes"
}

func (jenny GoBuilder) Generate(file *simplecue.File) (*codejen.File, error) {
	output, err := jenny.generateFile(file)
	if err != nil {
		return nil, err
	}

	return codejen.NewFile(file.Package+"_builder_gen.go", output, jenny), nil
}

func (jenny GoBuilder) generateFile(file *simplecue.File) ([]byte, error) {
	var buffer strings.Builder
	tr := newPreprocessor()

	tr.translateTypes(file.Types)

	buffer.WriteString(fmt.Sprintf("package %s\n\n", file.Package))

	// Option type declaration
	buffer.WriteString("type Option func(builder *Builder) error\n\n")

	// Builder type declaration
	buffer.WriteString(`type Builder struct {
	internal *Dashboard
}
`)

	// Constructor type declaration
	buffer.WriteString(`func New(title string, options ...Option) (Builder, error) {
	dashboard := &Dashboard{
	Title: title,
}

	builder := &Builder{internal: dashboard}

	for _, opt := range options {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}
`)

	// (un)marshaling utilities
	buffer.WriteString(`
// MarshalJSON implements the encoding/json.Marshaler interface.
//
// This method can be used to render the dashboard as JSON
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(builder.internal)
}

// MarshalIndentJSON renders the dashboard as indented JSON
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *Builder) MarshalIndentJSON() ([]byte, error) {
	return json.MarshalIndent(builder.internal, "", "  ")
}
`)

	for _, typeDef := range tr.sortedTypes() {
		typeDefGen, err := jenny.formatTypeDef(typeDef)
		if err != nil {
			return nil, err
		}
		if typeDefGen == nil {
			continue
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func (jenny GoBuilder) formatTypeDef(def simplecue.TypeDefinition) ([]byte, error) {
	// nothing to do for enums & other non-struct types
	if def.Type != simplecue.DefinitionStruct {
		return nil, nil
	}

	// No options if not main/entrypoint type
	if !def.IsEntryPoint {
		return nil, nil
	}

	return jenny.formatMainTypeOptions(def)
}

func (jenny GoBuilder) formatMainTypeOptions(def simplecue.TypeDefinition) ([]byte, error) {
	var buffer strings.Builder

	for _, fieldDef := range def.Fields {
		buffer.WriteString(jenny.fieldToOption(fieldDef))
	}

	return []byte(buffer.String()), nil
}

func (jenny GoBuilder) fieldToOption(def simplecue.FieldDefinition) string {
	var buffer strings.Builder

	fieldName := strings.Title(def.Name)
	typeName := strings.TrimPrefix(formatType(def.Type), "*")
	generatedConstraints := strings.Join(jenny.constraints(def.Name, def.Type.Constraints), "\n")
	asPointer := ""
	if def.Type.Nullable {
		asPointer = "&"
	}

	buffer.WriteString(fmt.Sprintf(`
func %[1]s(%[2]s %[3]s) Option {
	return func(builder *Builder) error {
		%[4]s
		builder.internal.%[1]s = %[5]s%[2]s

		return nil
	}
}
`, fieldName, def.Name, typeName, generatedConstraints, asPointer))

	return buffer.String()
}

func (jenny GoBuilder) constraints(argumentName string, constraints []simplecue.TypeConstraint) []string {
	output := make([]string, 0, len(constraints))

	for _, constraint := range constraints {
		output = append(output, jenny.constraint(argumentName, constraint))
	}

	return output
}

func (jenny GoBuilder) constraint(argumentName string, constraint simplecue.TypeConstraint) string {
	var buffer strings.Builder

	// FIXME
	buffer.WriteString(fmt.Sprintf("if !(%[1]s %[2]s %[3]v) {\n", argumentName, constraint.Op, constraint.Args[0]))
	buffer.WriteString(fmt.Sprintf("return errors.New(\"%[1]s must be %[2]s %[3]v\")\n", argumentName, constraint.Op, constraint.Args[0]))
	buffer.WriteString("}\n")

	return buffer.String()
}
