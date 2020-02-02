package query

import (
	"github.com/grafana-tools/sdk"
)

type SortOption int
type RefreshOption int

// Option represents an option that can be used to configure a query.
type Option func(constant *Query)

const None SortOption = 0
const AlphabeticalAsc SortOption = 1
const AlphabeticalDesc SortOption = 2
const NumericalAsc SortOption = 3
const NumericalDesc SortOption = 4
const AlphabeticalNoCaseAsc SortOption = 5
const AlphabeticalNoCaseDesc SortOption = 6

const Never RefreshOption = 0
const DashboardLoad RefreshOption = 1
const TimeChange RefreshOption = 2

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
func Sort(order SortOption) Option {
	return func(query *Query) {
		query.Builder.Sort = int(order)
	}
}

// Refresh defines the interval in which the values will be refreshed.
func Refresh(refresh RefreshOption) Option {
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
		query.Builder.Current = sdk.Current{Text: "All", Value: "$_all"}
	}
}
