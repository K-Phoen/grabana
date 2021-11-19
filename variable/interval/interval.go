package interval

import (
	"sort"
	"strings"

	"github.com/K-Phoen/sdk"
	"github.com/prometheus/common/model"
)

// Option represents an option that can be used to configure an interval.
type Option func(interval *Interval)

// ValuesList represent a list of options for an interval variable.
type ValuesList []string

// Interval represents a "interval" templated variable.
type Interval struct {
	Builder sdk.TemplateVar
}

// New creates a new "interval" templated variable.
func New(name string, options ...Option) *Interval {
	interval := &Interval{Builder: sdk.TemplateVar{
		Name:    name,
		Label:   name,
		Type:    "interval",
		Options: []sdk.Option{},
	}}

	for _, opt := range options {
		opt(interval)
	}

	return interval
}

// Values sets the possible values for the variable.
func Values(values ValuesList) Option {
	return func(interval *Interval) {
		sort.SliceStable(values, func(i, j int) bool {
			iDuration, _ := model.ParseDuration(values[i])
			jDuration, _ := model.ParseDuration(values[j])
			return iDuration < jDuration
		})

		for _, value := range values {
			interval.Builder.Options = append(interval.Builder.Options, sdk.Option{
				Text:  value,
				Value: value,
			})
		}

		interval.Builder.Query = strings.Join(values, ",")
	}
}

// Default sets the default value of the variable.
func Default(value string) Option {
	return func(interval *Interval) {
		interval.Builder.Current = sdk.Current{
			Text:  &sdk.StringSliceString{Value: []string{value}, Valid: true},
			Value: value,
		}
	}
}

// Label sets the label of the variable.
func Label(label string) Option {
	return func(interval *Interval) {
		interval.Builder.Label = label
	}
}

// HideLabel ensures that this variable's label will not be displayed.
func HideLabel() Option {
	return func(interval *Interval) {
		interval.Builder.Hide = 1
	}
}

// Hide ensures that the variable will not be displayed.
func Hide() Option {
	return func(interval *Interval) {
		interval.Builder.Hide = 2
	}
}
