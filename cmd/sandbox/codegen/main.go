package main

import (
	"context"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/K-Phoen/grabana/internal/gen/jennies/golang"
	"github.com/K-Phoen/grabana/internal/gen/simplecue"
	"github.com/grafana/codejen"
)

func main() {
	entrypoints := []string{"/home/kevin/sandbox/personal/grabana/schemas/core/dashboard/dashboard.cue"}

	// Load Cue files into Cue build.Instances slice
	// the second arg is a configuration object, we'll see this later
	bis := load.Instances(entrypoints, nil)

	values, err := cuecontext.New().BuildInstances(bis)
	if err != nil {
		panic(err)
	}

	ast, err := simplecue.GenerateAST(values[0].Value(), simplecue.Config{
		Package: "dashboard",
	})
	if err != nil {
		panic(err)
	}

	// Here begins the code generation setup
	generationTargets := codejen.JennyListWithNamer[*simplecue.File](func(f *simplecue.File) string {
		return f.Package
	})
	generationTargets.AppendOneToOne(
		golang.GoRawTypes{},
		golang.GoBuilder{},
	)
	generationTargets.AddPostprocessors(golang.PostProcessFile)

	rootCodeJenFS := codejen.NewFS()

	fs, err := generationTargets.GenerateFS(ast)
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
