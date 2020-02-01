package query

import (
	"github.com/grafana-tools/sdk"
)

type SortOption int
type RefreshOption int
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

type Query struct {
	Builder sdk.TemplateVar
}

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

func DataSource(source string) Option {
	return func(query *Query) {
		query.Builder.Datasource = &source
	}
}

func Request(request string) Option {
	return func(query *Query) {
		query.Builder.Query = request
	}
}

func Sort(order SortOption) Option {
	return func(query *Query) {
		query.Builder.Sort = int(order)
	}
}

func Refresh(refresh RefreshOption) Option {
	return func(query *Query) {
		value := int64(refresh)
		query.Builder.Refresh = sdk.BoolInt{Flag: true, Value: &value}
	}
}

func Regex(regex string) Option {
	return func(query *Query) {
		query.Builder.Regex = regex
	}
}

func Label(label string) Option {
	return func(query *Query) {
		query.Builder.Label = label
	}
}

func HideLabel() Option {
	return func(query *Query) {
		query.Builder.Hide = 1
	}
}

func Hide() Option {
	return func(query *Query) {
		query.Builder.Hide = 2
	}
}

func Multi() Option {
	return func(query *Query) {
		query.Builder.Multi = true
	}
}

func IncludeAll() Option {
	return func(query *Query) {
		query.Builder.IncludeAll = true
		query.Builder.Options = append(query.Builder.Options, sdk.Option{
			Text:  "All",
			Value: "$__all",
		})
	}
}

func DefaultAll() Option {
	return func(query *Query) {
		query.Builder.Current = sdk.Current{Text: "All", Value: "$_all"}
	}
}
