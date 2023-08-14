package ts_test

import (
	"testing"

	"github.com/K-Phoen/grabana/internal/gen/simplecue/ts"
	"github.com/matryer/is"
)

func TestCommentFromString(t *testing.T) {
	table := map[string]struct {
		input string
		plain string
		jsdoc string
	}{
		"basic": {
			input: "just some simple text",
			plain: "// just some simple text",
			jsdoc: `/**
 * just some simple text
 */
`,
		},
		"breaking": {
			input: "some more text, enough that it will be broken over multiple lines",
			plain: `// some more text, enough
// that it will be broken
// over multiple lines`,
			jsdoc: `/**
 * some more text, enough
 * that it will be broken
 * over multiple lines
 */
`,
		},
	}

	for name, tst := range table {
		tt := tst
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(ts.CommentFromString(tt.input, 25, false).String(), tt.plain)
			is.Equal(ts.CommentFromString(tt.input, 25, true).String(), tt.jsdoc)
		})
	}
}
