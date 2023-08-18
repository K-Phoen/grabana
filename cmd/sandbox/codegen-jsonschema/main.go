package main

import (
	"context"
	"os"

	"github.com/K-Phoen/grabana/internal/gen/ast"
	"github.com/K-Phoen/grabana/internal/gen/jennies"
	"github.com/K-Phoen/grabana/internal/gen/jennies/golang"
	"github.com/K-Phoen/grabana/internal/gen/jsonschema"
	"github.com/grafana/codejen"
)

func main() {
	entrypoint := "/home/kevin/sandbox/personal/grabana/schemas/jsonschema/core/playlist/playlist.json"
	pkg := "playlist"

	reader, err := os.Open(entrypoint)
	if err != nil {
		panic(err)
	}

	schemaAst, err := jsonschema.GenerateAST(reader, jsonschema.Config{
		Package: pkg, // TODO: extract from input schema/folder?
	})
	if err != nil {
		panic(err)
	}

	// Here begins the code generation setup
	generationTargets := codejen.JennyListWithNamer[*ast.File](func(f *ast.File) string {
		return f.Package
	})
	generationTargets.AppendOneToOne(
		golang.GoRawTypes{},
		golang.GoBuilder{},
	)
	generationTargets.AddPostprocessors(
		golang.PostProcessFile,
		jennies.Prefixer(pkg),
	)

	rootCodeJenFS := codejen.NewFS()

	fs, err := generationTargets.GenerateFS(schemaAst)
	if err != nil {
		panic(err)
	}

	err = rootCodeJenFS.Merge(fs)
	if err != nil {
		panic(err)
	}

	err = rootCodeJenFS.Write(context.Background(), "gen")
	if err != nil {
		panic(err)
	}
}
