package text

import (
	"github.com/grafana-tools/sdk"
)

type Option func(text *Text)

type Text struct {
	Builder *sdk.Panel
}

func New(title string, options ...Option) *Text {
	panel := &Text{Builder: sdk.NewText(title)}

	panel.Builder.IsNew = false
	panel.Builder.Span = 6

	for _, opt := range options {
		opt(panel)
	}

	return panel
}

// HTML sets the content of the panel, to be rendered as HTML.
func HTML(content string) Option {
	return func(text *Text) {
		text.Builder.TextPanel.Mode = "html"
		text.Builder.TextPanel.Content = content
	}
}

// Markdown sets the content of the panel, to be rendered as markdown.
func Markdown(content string) Option {
	return func(text *Text) {
		text.Builder.TextPanel.Mode = "markdown"
		text.Builder.TextPanel.Content = content
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(text *Text) {
		text.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(text *Text) {
		text.Builder.Height = &height
	}
}
