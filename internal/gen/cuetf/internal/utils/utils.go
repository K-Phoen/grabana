package utils

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func Title(str string) string {
	var caser = cases.Title(language.English, cases.NoLower)
	return caser.String(str)
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToCamelCase(str string) string {
	words := strings.Split(str, "_")
	camelCase := ""
	for _, s := range words {
		camelCase += Title(s)
	}
	return camelCase
}

func CapitalizeFirstLetter(str string) string {
	sep := " "
	parts := strings.SplitN(str, sep, 2)
	if len(parts) != 2 {
		return Title(str)
	}
	return Title(parts[0]) + sep + parts[1]
}
