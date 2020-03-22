package constant

import (
	"sort"
	"strings"

	"github.com/grafana-tools/sdk"
)

// Option represents an option that can be used to configure a constant.
type Option func(constant *Constant)

// ValuesMap represent a "label" to "value" map of options for a constant variable.
type ValuesMap map[string]string

func (values ValuesMap) asQuery() string {
	valuesList := make([]string, 0, len(values))

	for _, value := range values {
		valuesList = append(valuesList, value)
	}

	sort.Strings(valuesList)

	return strings.Join(valuesList, ",")
}

func (values ValuesMap) labelFor(value string) string {
	for label, val := range values {
		if val == value {
			return label
		}
	}

	return value
}

// Constant represents a "constant" templated variable.
type Constant struct {
	Builder sdk.TemplateVar
	values  ValuesMap
}

// New creates a new "constant" templated variable.
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

// Values sets the possible values for the variable.
func Values(values ValuesMap) Option {
	return func(constant *Constant) {
		for label, value := range values {
			constant.Builder.Options = append(constant.Builder.Options, sdk.Option{
				Text:  label,
				Value: value,
			})
		}

		constant.values = values
		constant.Builder.Query = values.asQuery()
	}
}

// Default sets the default value of the variable.
func Default(value string) Option {
	return func(constant *Constant) {
		constant.Builder.Current = sdk.Current{
			Text:  constant.values.labelFor(value),
			Value: value,
		}
	}
}

// Label sets the label of the variable.
func Label(label string) Option {
	return func(constant *Constant) {
		constant.Builder.Label = label
	}
}

// HideLabel ensures that this variable's label will not be displayed.
func HideLabel() Option {
	return func(constant *Constant) {
		constant.Builder.Hide = 1
	}
}

// Hide ensures that the variable will not be displayed.
func Hide() Option {
	return func(constant *Constant) {
		constant.Builder.Hide = 2
	}
}
