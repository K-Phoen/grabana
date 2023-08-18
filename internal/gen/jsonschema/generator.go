package jsonschema

import (
	"github.com/K-Phoen/grabana/internal/gen/ast"
	"github.com/davecgh/go-spew/spew"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Config struct {
	// Package name used to generate code into.
	Package string
}

type newGenerator struct {
	file *ast.File
}

func GenerateAST(schemaURL string, c Config) (*ast.File, error) {
	g := &newGenerator{
		file: &ast.File{
			Package: c.Package,
		},
	}

	sch, err := jsonschema.Compile(schemaURL)
	if err != nil {
		return nil, err
	}

	spew.Dump(sch)

	return g.file, nil
}
