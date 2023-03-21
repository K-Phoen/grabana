package text

import (
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a textbox variable.
type Option func(constant *Text)

// Text represents a "textbox" templated variable.
type Text struct {
	Builder sdk.TemplateVar
}

// New creates a new "query" templated variable.
func New(name string, options ...Option) *Text {
	query := &Text{Builder: sdk.TemplateVar{
		Name:    name,
		Label:   name,
		Type:    "textbox",
		Options: []sdk.Option{},
	}}

	for _, opt := range options {
		opt(query)
	}

	return query
}

// Label sets the label of the variable.
func Label(label string) Option {
	return func(query *Text) {
		query.Builder.Label = label
	}
}

// HideLabel ensures that this variable's label will not be displayed.
func HideLabel() Option {
	return func(query *Text) {
		query.Builder.Hide = 1
	}
}

// Hide ensures that the variable will not be displayed.
func Hide() Option {
	return func(query *Text) {
		query.Builder.Hide = 2
	}
}
