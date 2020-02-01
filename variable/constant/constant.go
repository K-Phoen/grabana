package constant

import (
	"github.com/grafana-tools/sdk"
)

type Option func(constant *Constant)

// ValuesMap represent a "label" to "value" map of options for a constant variable.
type ValuesMap map[string]string

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

func Values(values ValuesMap) Option {
	return func(constant *Constant) {
		for label, value := range values {
			constant.Builder.Options = append(constant.Builder.Options, sdk.Option{
				Text:  label,
				Value: value,
			})
		}
	}
}

func Default(value string) Option {
	return func(constant *Constant) {
		constant.Builder.Current = sdk.Current{
			Text: value,
		}
	}
}

func Label(label string) Option {
	return func(constant *Constant) {
		constant.Builder.Label = label
	}
}

func HideLabel() Option {
	return func(constant *Constant) {
		constant.Builder.Hide = 1
	}
}

func Hide() Option {
	return func(constant *Constant) {
		constant.Builder.Hide = 2
	}
}
