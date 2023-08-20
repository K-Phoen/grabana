package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/K-Phoen/grabana/internal/gen/ast"
	"github.com/K-Phoen/grabana/internal/gen/jennies"
	"github.com/K-Phoen/grabana/internal/gen/jennies/golang"
	"github.com/K-Phoen/grabana/internal/gen/jennies/typescript"
	"github.com/K-Phoen/grabana/internal/gen/simplecue"
	"github.com/grafana/codejen"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

func main() {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	entrypoints := []string{"/home/kevin/sandbox/personal/grabana/schemas/cue/core/dashboard_original"}
	pkg := "dashboard"

	overlayFS, err := dirToPrefixedFS(entrypoints[0], "")
	if err != nil {
		panic(err)
	}

	cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", "kind", overlayFS)
	if err != nil {
		panic(fmt.Errorf("could not load kindsys instance: %w", err))
	}

	props, err := kindsys.ToKindProps[kindsys.CoreProperties](cueInstance)
	if err != nil {
		panic(fmt.Errorf("could not convert cue value to kindsys props: %w", err))
	}

	kindDefinition := kindsys.Def[kindsys.CoreProperties]{
		V:          cueInstance,
		Properties: props,
	}

	boundKind, err := kindsys.BindCore(themaRuntime, kindDefinition)
	if err != nil {
		panic(fmt.Errorf("could not bind kind definition to kind: %w", err))
	}

	rawLatestSchemaAsCue := boundKind.Lineage().Latest().Underlying()
	latestSchemaAsCue := rawLatestSchemaAsCue.LookupPath(cue.MakePath(cue.Hid("_#schema", "github.com/grafana/thema")))

	schemaAst, err := simplecue.GenerateAST(latestSchemaAsCue, simplecue.Config{
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
		// Golang
		golang.GoRawTypes{},
		golang.GoBuilder{},

		// Typescript
		typescript.TypescriptRawTypes{},
	)
	generationTargets.AddPostprocessors(
		//golang.PostProcessFile,
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

func dirToPrefixedFS(directory string, prefix string) (fs.FS, error) {
	dirHandle, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	commonFS := fstest.MapFS{}
	for _, file := range dirHandle {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			return nil, err
		}

		commonFS[filepath.Join(prefix, file.Name())] = &fstest.MapFile{Data: content}
	}

	return commonFS, nil
}
