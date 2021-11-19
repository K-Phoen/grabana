package datasource

import (
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a query.
type Option func(constant *Datasource)

const (
	// dashboardLoad will refresh the results every time the dashboard is loaded.
	dashboardLoad int64 = 1
)

// Datasource represents a "datasource" templated variable.
type Datasource struct {
	Builder sdk.TemplateVar
}

// New creates a new "query" templated variable.
func New(name string, options ...Option) *Datasource {
	refreshValue := dashboardLoad

	query := &Datasource{Builder: sdk.TemplateVar{
		Name:    name,
		Label:   name,
		Type:    "datasource",
		Options: []sdk.Option{},
		Refresh: sdk.BoolInt{
			Flag:  true,
			Value: &refreshValue,
		},
	}}

	for _, opt := range options {
		opt(query)
	}

	return query
}

// Type defines the datasource type. Example: "grafana", "stackdriver", "prometheus", ...
func Type(datasourceType string) Option {
	return func(query *Datasource) {
		query.Builder.Query = datasourceType
	}
}

// Regex defines a filter allowing to filter the values returned by the request/query.
func Regex(regex string) Option {
	return func(query *Datasource) {
		query.Builder.Regex = regex
	}
}

// Label sets the label of the variable.
func Label(label string) Option {
	return func(query *Datasource) {
		query.Builder.Label = label
	}
}

// HideLabel ensures that this variable's label will not be displayed.
func HideLabel() Option {
	return func(query *Datasource) {
		query.Builder.Hide = 1
	}
}

// Hide ensures that the variable will not be displayed.
func Hide() Option {
	return func(query *Datasource) {
		query.Builder.Hide = 2
	}
}

// Multi allows several values to be selected.
func Multi() Option {
	return func(query *Datasource) {
		query.Builder.Multi = true
	}
}

// IncludeAll adds an option to allow all values to be selected.
func IncludeAll() Option {
	return func(query *Datasource) {
		query.Builder.IncludeAll = true
		query.Builder.Options = append(query.Builder.Options, sdk.Option{
			Text:  "All",
			Value: "$__all",
		})
	}
}
