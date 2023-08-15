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

	return []byte(buffer.String()), nil
}
