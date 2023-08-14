package main

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"github.com/K-Phoen/grabana/internal/gen/printer/golang"
	"github.com/K-Phoen/grabana/internal/gen/simplecue"
)

func main() {
	entrypoints := []string{"/home/kevin/sandbox/personal/grabana/schemas/core/dashboard/dashboard.cue"}

	// Load Cue files into Cue build.Instances slice
	// the second arg is a configuration object, we'll see this later
	bis := load.Instances(entrypoints, nil)

	runtimeInstances := cue.Build(bis)

	fmt.Printf("instances: %d\n", len(runtimeInstances))

	b, err := simplecue.GenerateAny(runtimeInstances[0].Value(), simplecue.Config{
		Export:  true,
		Package: "dashboard",
	}, golang.Printer)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("gen/dashboard.go", b, 0644)
	if err != nil {
		panic(err)
	}
}
