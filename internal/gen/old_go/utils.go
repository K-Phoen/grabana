package gocode

import (
	"strings"

	"cuelang.org/go/cue/ast"
)

func formatDoc(comments []*ast.CommentGroup) []string {
	var docLines []string
	for _, comment := range comments {
		for _, line := range strings.Split(strings.Trim(comment.Text(), "\n "), "\n") {
			docLines = append(docLines, line)
		}
	}

	return docLines
}
