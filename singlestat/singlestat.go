package singlestat

import (
	"strings"

	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
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

// ValueMap allows to map a value into explicit text.
type ValueMap struct {
	Value string
	Text  string
}

// RangeMap allows to map a range of values into explicit text.
type RangeMap struct {
	From string
	To   string
	Text string
}

// nolint: gochecknoglobals
var valueToTextMapping = 1

// nolint: gochecknoglobals
var rangeToTextMapping = 2

// SingleStat represents a single stat panel.
type SingleStat struct {
	Builder *sdk.Panel
}

// New creates a new single stat panel.
func New(title string, options ...Option) *SingleStat {
	panel := &SingleStat{Builder: sdk.NewSinglestat(title)}

	valueToText := "value to text"
	rangeToText := "range to text"

	panel.Builder.IsNew = false
	mappingType := uint(valueToTextMapping)
	panel.Builder.MappingType = &mappingType
	panel.Builder.MappingTypes = []*sdk.MapType{
		{
			Name:  &valueToText,
			Value: &valueToTextMapping,
		},
		{
			Name:  &rangeToText,
			Value: &rangeToTextMapping,
		},
	}
	panel.Builder.SparkLine = struct {
		FillColor *string  `json:"fillColor,omitempty"`
		Full      bool     `json:"full,omitempty"`
		LineColor *string  `json:"lineColor,omitempty"`
		Show      bool     `json:"show,omitempty"`
		YMin      *float64 `json:"ymin,omitempty"`
		YMax      *float64 `json:"ymax,omitempty"`
	}{}

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
		ValuesToText([]ValueMap{
			{
				Value: "null",
				Text:  "N/A",
			},
		}),
		SparkLineColor("rgb(31, 120, 193)"),
		SparkLineFillColor("rgba(31, 118, 189, 0.18)"),
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

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.AddTarget(target.Builder)
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

// SparkLine displays the spark line summary of the series in addition to the
// single stat.
func SparkLine() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SparkLine.Show = true
		singleStat.Builder.SparkLine.Full = false
	}
}

// FullSparkLine displays a full height spark line summary of the series in
// addition to the single stat.
func FullSparkLine() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SparkLine.Show = true
		singleStat.Builder.SparkLine.Full = true
	}
}

// SparkLineColor sets the line color of the spark line.
func SparkLineColor(color string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SparkLine.LineColor = &color
	}
}

// SparkLineFillColor sets the color the spark line will be filled with.
func SparkLineFillColor(color string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SparkLine.FillColor = &color
	}
}

// SparkLineYMin defines the smallest value expected on the Y axis of the spark line.
func SparkLineYMin(value float64) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SparkLine.YMin = &value
	}
}

// SparkLineYMax defines the largest value expected on the Y axis of the spark line.
func SparkLineYMax(value float64) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SparkLine.YMax = &value
	}
}

// ValueType configures how the series will be reduced to a single value.
func ValueType(valueType StatType) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ValueName = string(valueType)
	}
}

// ValueFontSize sets the font size used to display the value (eg: "100%").
func ValueFontSize(size string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ValueFontSize = size
	}
}

// Prefix sets the text used as prefix of the value.
func Prefix(prefix string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Prefix = &prefix
	}
}

// PrefixFontSize sets the size used for the prefix text (eg: "110%").
func PrefixFontSize(size string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.PrefixFontSize = &size
	}
}

// Postfix sets the text used as postfix of the value.
func Postfix(postfix string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.Postfix = &postfix
	}
}

// PostfixFontSize sets the size used for the postfix text (eg: "110%")
func PostfixFontSize(size string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.PostfixFontSize = &size
	}
}

// ColorValue will show the threshold's colors on the value itself.
func ColorValue() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ColorValue = true
	}
}

// ColorBackground will show the threshold's colors in the background.
func ColorBackground() Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.ColorBackground = true
	}
}

// Thresholds change the background and value colors dynamically within the
// panel, depending on the Singlestat value. The threshold is defined by 2
// values which represent 3 ranges that correspond to the three colors directly
// to the right.
func Thresholds(values [2]string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SinglestatPanel.Thresholds = strings.Join([]string{values[0], values[1]}, ",")
	}
}

// Colors define which colors will be applied to the single value based on the
// threshold levels.
func Colors(values [3]string) Option {
	return func(singleStat *SingleStat) {
		singleStat.Builder.SinglestatPanel.Colors = []string{values[0], values[1], values[2]}
	}
}

// ValuesToText allows to translate the value of the summary stat into explicit
// text.
func ValuesToText(mapping []ValueMap) Option {
	return func(singleStat *SingleStat) {
		valueMap := make([]sdk.ValueMap, 0, len(mapping))

		for _, entry := range mapping {
			valueMap = append(valueMap, sdk.ValueMap{
				Op:       "=",
				TextType: entry.Text,
				Value:    entry.Value,
			})
		}

		mappingType := uint(valueToTextMapping)
		singleStat.Builder.MappingType = &mappingType
		singleStat.Builder.ValueMaps = valueMap
	}
}

// RangesToText allows to translate the value of the summary stat into explicit
// text.
func RangesToText(mapping []RangeMap) Option {
	return func(singleStat *SingleStat) {
		rangeMap := make([]*sdk.RangeMap, 0, len(mapping))

		for i := range mapping {
			rangeMap = append(rangeMap, &sdk.RangeMap{
				From: &mapping[i].From,
				To:   &mapping[i].To,
				Text: &mapping[i].Text,
			})
		}

		mappingType := uint(rangeToTextMapping)
		singleStat.Builder.MappingType = &mappingType
		singleStat.Builder.RangeMaps = rangeMap
	}
}
