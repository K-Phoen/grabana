package custom

import (
	"github.com/grafana-tools/sdk"
)

// Option represents an option that can be used to configure a custom variable.
type Option func(constant *Custom)

// ValuesMap represent a "label" to "value" map of options for a custom variable.
type ValuesMap map[string]string

// Custom represents a "custom" templated variable.
type Custom struct {
	Builder sdk.TemplateVar
}

// New creates a new "custom" templated variable.
func New(name string, options ...Option) *Custom {
	constant := &Custom{Builder: sdk.TemplateVar{
		Name:  name,
		Label: name,
		Type:  "custom",
	}}

	for _, opt := range options {
		opt(constant)
	}

	return constant
}

// Values sets the possible values for the variable.
func Values(values ValuesMap) Option {
	return func(constant *Custom) {
		for label, value := range values {
			constant.Builder.Options = append(constant.Builder.Options, sdk.Option{
				Text:  label,
				Value: value,
			})
		}
	}
}

// Default sets the default value of the variable.
func Default(value string) Option {
	return func(constant *Custom) {
		constant.Builder.Current = sdk.Current{
			Text: value,
		}
	}
}

// Label sets the label of the variable.
func Label(label string) Option {
	return func(constant *Custom) {
		constant.Builder.Label = label
	}
}

// HideLabel ensures that this variable's label will not be displayed.
func HideLabel() Option {
	return func(constant *Custom) {
		constant.Builder.Hide = 1
	}
}

// Hide ensures that the variable will not be displayed.
func Hide() Option {
	return func(constant *Custom) {
		constant.Builder.Hide = 2
	}
}

// Multi allows several values to be selected.
func Multi() Option {
	return func(constant *Custom) {
		constant.Builder.Multi = true
	}
}

// IncludeAll adds an option to allow all values to be selected.
func IncludeAll() Option {
	return func(constant *Custom) {
		constant.Builder.IncludeAll = true
		constant.Builder.Options = append(constant.Builder.Options, sdk.Option{
			Text:  "All",
			Value: "$__all",
		})
	}
}
