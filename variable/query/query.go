package query

import (
	"github.com/grafana-tools/sdk"
)

// Option represents an option that can be used to configure a query.
type Option func(constant *Query)

// SortOrder represents the ordering method applied to values.
type SortOrder int

const (
	// None will preserve the results ordering as returned by the query.
	None SortOrder = 0

	// AlphabeticalAsc will sort the results by ascending alphabetical order.
	AlphabeticalAsc SortOrder = 1

	// AlphabeticalDesc will sort the results by descending alphabetical order.
	AlphabeticalDesc SortOrder = 2

	// NumericalAsc will sort the results by ascending numerical order.
	NumericalAsc SortOrder = 3

	// NumericalDesc will sort the results by descending numerical order.
	NumericalDesc SortOrder = 4

	// AlphabeticalNoCaseAsc will sort the results by ascending alphabetical order, case-insensitive.
	AlphabeticalNoCaseAsc SortOrder = 5

	// AlphabeticalNoCaseDesc will sort the results by descending alphabetical order, case-insensitive.
	AlphabeticalNoCaseDesc SortOrder = 6
)

// RefreshInterval represents the interval at which the results of a query will
// be refreshed.
type RefreshInterval int

const (
	// Never will prevent the results from being refreshed.
	Never = 0

	// DashboardLoad will refresh the results every time the dashboard is loaded.
	DashboardLoad = 1

	// TimeChange will refresh the results every time the time interval changes.
	TimeChange = 2
)

// Query represents a "query" templated variable.
type Query struct {
	Builder sdk.TemplateVar
}

// New creates a new "query" templated variable.
func New(name string, options ...Option) *Query {
	query := &Query{Builder: sdk.TemplateVar{
		Name:  name,
		Label: name,
		Type:  "query",
	}}

	for _, opt := range append([]Option{Refresh(DashboardLoad)}, options...) {
		opt(query)
	}

	return query
}

// DataSource sets the data source to be used by the query.
func DataSource(source string) Option {
	return func(query *Query) {
		query.Builder.Datasource = &source
	}
}

// Request defines the query to be executed.
func Request(request string) Option {
	return func(query *Query) {
		query.Builder.Query = request
	}
}

// Sort defines the order in which the values will be sorted.
func Sort(order SortOrder) Option {
	return func(query *Query) {
		query.Builder.Sort = int(order)
	}
}

// Refresh defines the interval in which the values will be refreshed.
func Refresh(refresh RefreshInterval) Option {
	return func(query *Query) {
		value := int64(refresh)
		query.Builder.Refresh = sdk.BoolInt{Flag: true, Value: &value}
	}
}

// Regex defines a filter allowing to filter the values returned by the request/query.
func Regex(regex string) Option {
	return func(query *Query) {
		query.Builder.Regex = regex
	}
}

// Label sets the label of the variable.
func Label(label string) Option {
	return func(query *Query) {
		query.Builder.Label = label
	}
}

// HideLabel ensures that this variable's label will not be displayed.
func HideLabel() Option {
	return func(query *Query) {
		query.Builder.Hide = 1
	}
}

// Hide ensures that the variable will not be displayed.
func Hide() Option {
	return func(query *Query) {
		query.Builder.Hide = 2
	}
}

// Multi allows several values to be selected.
func Multi() Option {
	return func(query *Query) {
		query.Builder.Multi = true
	}
}

// IncludeAll adds an option to allow all values to be selected.
func IncludeAll() Option {
	return func(query *Query) {
		query.Builder.IncludeAll = true
		query.Builder.Options = append(query.Builder.Options, sdk.Option{
			Text:  "All",
			Value: "$__all",
		})
	}
}

// DefaultAll selects "All" values by default.
func DefaultAll() Option {
	return func(query *Query) {
		query.Builder.Current = sdk.Current{Text: "All", Value: "$__all"}
	}
}
