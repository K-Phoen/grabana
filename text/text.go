package text

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a text panel.
type Option func(text *Text) error

// Text represents a text panel.
type Text struct {
	Builder *sdk.Panel
}

// New creates a new text panel.
func New(title string, options ...Option) (*Text, error) {
	panel := &Text{Builder: sdk.NewText(title)}

	panel.Builder.IsNew = false
	panel.Builder.Span = 6

	for _, opt := range options {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(text *Text) error {
		text.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			text.Builder.Links = append(text.Builder.Links, link.Builder)
		}

		return nil
	}
}

// HTML sets the content of the panel, to be rendered as HTML.
func HTML(content string) Option {
	return func(text *Text) error {
		text.Builder.TextPanel.Mode = "html"
		text.Builder.TextPanel.Content = content

		return nil
	}
}

// Markdown sets the content of the panel, to be rendered as markdown.
func Markdown(content string) Option {
	return func(text *Text) error {
		text.Builder.TextPanel.Mode = "markdown"
		text.Builder.TextPanel.Content = content

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(text *Text) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		text.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(text *Text) error {
		text.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(text *Text) error {
		text.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(text *Text) error {
		text.Builder.Transparent = true

		return nil
	}
}
