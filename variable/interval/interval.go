package interval

import (
	"strings"

	"github.com/grafana-tools/sdk"
)

type Option func(interval *Interval)

// ValuesList represent a list of options for an interval variable.
type ValuesList []string

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

func Values(values ValuesList) Option {
	return func(interval *Interval) {
		interval.Builder.Query = strings.Join(values, ",")
	}
}

func Default(value string) Option {
	return func(interval *Interval) {
		interval.Builder.Current = sdk.Current{
			Text: value,
		}
	}
}

func Label(label string) Option {
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
