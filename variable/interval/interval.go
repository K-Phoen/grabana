package interval

import (
	"strings"

	"github.com/grafana-tools/sdk"
)

type Option func(constant *Interval)

// Values represent a list of options for an interval variable.
type Values []string

type Interval struct {
	Builder sdk.TemplateVar
}

func New(name string, options ...Option) *Interval {
	interval := &Interval{Builder: sdk.TemplateVar{
		Name:  name,
		Label: name,
		Type:  "interval",
	}}

	for _, opt := range options {
		opt(interval)
	}

	return interval
}

func WithValues(values Values) Option {
	return func(interval *Interval) {
		interval.Builder.Query = strings.Join(values, ",")
	}
}

func WithDefault(value string) Option {
	return func(interval *Interval) {
		interval.Builder.Current = sdk.Current{
			Text: value,
		}
	}
}

func WithLabel(label string) Option {
	return func(interval *Interval) {
		interval.Builder.Label = label
	}
}

func HideLabel() Option {
	return func(interval *Interval) {
		interval.Builder.Hide = 1
	}
}

func Hide() Option {
	return func(interval *Interval) {
		interval.Builder.Hide = 2
	}
}
