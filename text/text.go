package text

import (
	"github.com/grafana-tools/sdk"
)

type Option func(text *Text)

type Text struct {
	Builder *sdk.Panel
}

func New(title string) *Text {
	panel := &Text{Builder: sdk.NewText(title)}

	panel.Builder.IsNew = false
	panel.Builder.Span = 6

	return panel
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
