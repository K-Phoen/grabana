package jsonschema

import (
	"strings"
)

func stringInList(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

func schemaComments(schema Schema) []string {
	comment := schema.Comment
	if comment == "" {
		comment = schema.Description
	}

	lines := strings.Split(comment, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}
