package text

import (
	"github.com/grafana-tools/sdk"
)

type Option func(text *Text)

type Text struct {
	Builder *sdk.Panel
}

func Defaults(text *Text) {
	text.Builder.IsNew = false
	text.Builder.Span = 6

	Editable()(text)
}

func HTML(content string) Option {
	return func(text *Text) {
		text.Builder.TextPanel.Mode = "html"
		text.Builder.TextPanel.Content = content
	}
}

func Markdown(content string) Option {
	return func(text *Text) {
		text.Builder.TextPanel.Mode = "markdown"
		text.Builder.TextPanel.Content = content
	}
}

func Editable() Option {
	return func(text *Text) {
		text.Builder.Editable = true
	}
}

func ReadOnly() Option {
	return func(text *Text) {
		text.Builder.Editable = false
	}
}
