package constant

import (
	"github.com/grafana-tools/sdk"
)

type Option func(constant *Constant)

type Value struct {
	Text  string
	Value string
}

type Constant struct {
	Builder sdk.TemplateVar
}

func New(name string, options ...Option) *Constant {
	constant := &Constant{Builder: sdk.TemplateVar{
		Name:  name,
		Label: name,
		Type:  "constant",
	}}

	for _, opt := range options {
		opt(constant)
	}

	return constant
}

func WithValues(values []Value) Option {
	return func(constant *Constant) {
		for _, value := range values {
			constant.Builder.Options = append(constant.Builder.Options, sdk.Option{
				Text:  value.Text,
				Value: value.Value,
			})
		}
	}
}

func WithDefault(value string) Option {
	return func(constant *Constant) {
		constant.Builder.Current = sdk.Current{
			Text: value,
		}
	}
}

func WithLabel(label string) Option {
	return func(constant *Constant) {
		constant.Builder.Label = label
	}
}
