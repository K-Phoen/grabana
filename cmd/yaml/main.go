package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/K-Phoen/grabana"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "Usage: go run -mod=vendor main.go file\n")
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(1)
	}

	builder, err := grabana.UnmarshalYAML(bytes.NewBuffer(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %s\n", err)
		os.Exit(1)
	}

	json, _ := builder.MarshalJSON()
	fmt.Printf("Dashboard JSON: %s\n", json)
}
