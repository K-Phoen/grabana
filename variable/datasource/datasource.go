package datasource

import (
	"github.com/grafana-tools/sdk"
)

// Option represents an option that can be used to configure a query.
type Option func(constant *Datasource)

// RefreshInterval represents the interval at which the results of a query will
// be refreshed.
type RefreshInterval int

const (
	// Never will prevent the results from being refreshed.
	Never RefreshInterval = 0

	// DashboardLoad will refresh the results every time the dashboard is loaded.
	DashboardLoad RefreshInterval = 1

	// TimeChange will refresh the results every time the time interval changes.
	TimeChange RefreshInterval = 2
)

// Datasource represents a "query" templated variable.
type Datasource struct {
	Builder sdk.TemplateVar
}

// New creates a new "query" templated variable.
func New(name string, options ...Option) *Datasource {
	query := &Datasource{Builder: sdk.TemplateVar{
		Name:  name,
		Label: name,
		Type:  "datasource",
	}}

	for _, opt := range append([]Option{Refresh(DashboardLoad)}, options...) {
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

// Refresh defines the interval in which the values will be refreshed.
func Refresh(refresh RefreshInterval) Option {
	return func(query *Datasource) {
		value := int64(refresh)
		query.Builder.Refresh = sdk.BoolInt{Flag: true, Value: &value}
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
