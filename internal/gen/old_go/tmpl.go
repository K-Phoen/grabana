package gocode

import (
	"embed"
	"strings"
	"text/template"
)

// All the parsed templates in the tmpl subdirectory.
var tmpls *template.Template

func init() {
	base := template.New("gocode").Funcs(template.FuncMap{
		"lowerCase":  strings.ToLower,
		"startsWith": strings.HasPrefix,
		"join":       strings.Join,
	})
	tmpls = template.Must(base.ParseFS(tmplFS, "templates/*.tmpl"))
}

//go:embed templates/*.tmpl
var tmplFS embed.FS
