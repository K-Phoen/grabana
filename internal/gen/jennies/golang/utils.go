package golang

import (
	"strings"
)

func stripHashtag(input string) string {
	return strings.TrimPrefix(input, "#")
}
