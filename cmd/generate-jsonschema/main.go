package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/K-Phoen/grabana/decoder"
	"github.com/invopop/jsonschema"
)

func main() {
	types := []struct {
		name  string
		input any
	}{
		{
			name:  "dashboard",
			input: &decoder.DashboardModel{},
		},
	}

	for _, t := range types {
		fmt.Printf("Generating schema for type '%s'\n", t.name)

		reflector := &jsonschema.Reflector{
			RequiredFromJSONSchemaTags: true, // we don't have good information as to what is required :/
			FieldNameTag:               "yaml",
			KeyNamer: func(key string) string {
				if strings.ToUpper(string(key[0])) == string(key[0]) {
					return strings.ToLower(key)
				}

				return key
			},
		}

		if err := reflector.AddGoComments("github.com/K-Phoen/grabana", "./decoder"); err != nil {
			panic(fmt.Errorf("could not add Go comments to reflector: %w", err))
		}

		schema := reflector.Reflect(t.input)
		schema.ID = jsonschema.ID(fmt.Sprintf("https://raw.githubusercontent.com/K-Phoen/grabana/master/schemas/%s.json", t.name))

		schemaJSON, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			panic(fmt.Errorf("could not marshal schema to JSON: %w", err))
		}

		if err := os.WriteFile(fmt.Sprintf("./schemas/%s.json", t.name), schemaJSON, 0600); err != nil {
			panic(fmt.Errorf("could not write schema: %w", err))
		}
	}
}
