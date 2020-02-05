package singlestat

import (
	"strings"

	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/grafana-tools/sdk"
)

// Option represents an option that can be used to configure a single stat panel.
type Option func(stat *SingleStat)

// StatType let you set the function that your entire query is reduced into a
// single value with.
type StatType string

const (
	// Min will return the smallest value in the series.
	Min StatType = "min"

	// Max will return the largest value in the series.
	Max StatType = "max"

	// Avg will return the average of all the non-null values in the series.
	Avg StatType = "avg"

	// Current will return the last value in the series. If the series ends on
	// null the previous value will be used.
	Current StatType = "current"

	// Total will return the sum of all the non-null values in the series.
	Total StatType = "total"

	// First will return the first value in the series.
	First StatType = "first"

	// Delta will return the total incremental increase (of a counter) in the
	// series. An attempt is made to account for counter resets, but this will
	// only be accurate for single instance metrics. Used to show total
	// counter increase in time series.
	Delta StatType = "delta"

	// Diff will return difference between ‘current’ (last value) and ‘first’..
	Diff StatType = "diff"

	// Range will return the difference between ‘min’ and ‘max’. Useful to
	// show the range of change for a gauge..
	Range StatType = "range"
)

type ValueMap struct {
	Text  string
	Value string
}

// SingleStat represents a single stat panel.
type SingleStat struct {
	Builder *sdk.Panel
}

func strPtr(input string) *string {
	return &input
}

func intPtr(input int) *int {
	return &input
}

// New creates a new single stat panel.
func New(title string, options ...Option) *SingleStat {
	panel := &SingleStat{Builder: sdk.NewSinglestat(title)}

	panel.Builder.IsNew = false
	mappingType := uint(1)
	panel.Builder.MappingType = &mappingType
	panel.Builder.MappingTypes = []*sdk.MapType{
		{
			Name:  strPtr("value to text"),
			Value: intPtr(1),
		},
		{
			Name:  strPtr("range to text"),
			Value: intPtr(2),
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Editable(),
		Span(6),
		ValueFontSize("100%"),
		ValueType(Avg),
		Colors([3]string{"#299c46", "rgba(237, 129, 40, 0.89)", "#d44a3a"}),
		MapValuesToText([]ValueMap{
			{
				Text:  "N/A",
				Value: "null",
			},
		}),
	}
}

// Editable marks the graph as editable.
func Editable() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Editable = true
	}
}

// ReadOnly marks the graph as non-editable.
func ReadOnly() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Editable = false
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Datasource = &source
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(singleStat *SingleStat) {
		singleStat.Builder.AddTarget(&sdk.Target{
			RefID:          target.Ref,
			Expr:           target.Expr,
			IntervalFactor: target.IntervalFactor,
			Interval:       target.Interval,
			Step:           target.Step,
			LegendFormat:   target.LegendFormat,
			Instant:        target.Instant,
			Format:         target.Format,
		})
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Height = &height
	}
}

// Unit sets the unit of the data displayed on this axis.
func Unit(unit string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Format = unit
	}
}

func ValueType(valueType StatType) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ValueName = string(valueType)
	}
}

func ValueFontSize(size string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ValueFontSize = size
	}
}

func Prefix(prefix string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Prefix = &prefix
	}
}

func PrefixFontSize(size string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.PrefixFontSize = &size
	}
}

func Postfix(postfix string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Postfix = &postfix
	}
}

func PostfixFontSize(size string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.PostfixFontSize = &size
	}
}

func ColorValue() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ColorValue = true
	}
}

func ColorBackground() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ColorBackground = true
	}
}

func Thresholds(values []string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SinglestatPanel.Thresholds = strings.Join(values, ",")
	}
}

func Colors(values [3]string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SinglestatPanel.Colors = []string{values[0], values[1], values[2]}
	}
}

func MapValuesToText(mapping []ValueMap) Option {
	return func(singleStat *SingleStat) {
		valueMap := make([]sdk.ValueMap, 0, len(mapping))

		for _, entry := range mapping {
			valueMap = append(valueMap, sdk.ValueMap{
				Op:       "=",
				TextType: entry.Text,
				Value:    entry.Value,
			})
		}

		mappingType := uint(1)
		singleStat.Builder.MappingType = &mappingType
		singleStat.Builder.ValueMaps = valueMap
	}
}
