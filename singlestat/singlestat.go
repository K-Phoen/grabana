package singlestat

import (
	"fmt"
	"strings"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a single stat panel.
type Option func(stat *SingleStat) error

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

	// Name will return the name value in the series.
	Name StatType = "name"
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
func New(title string, options ...Option) (*SingleStat, error) {
	panel := &SingleStat{Builder: sdk.NewSinglestat(title)}

	valueToText := "value to text"
	rangeToText := "range to text"

	panel.Builder.IsNew = false
	mappingType := uint(valueToTextMapping)
	panel.Builder.SinglestatPanel.MappingType = &mappingType
	panel.Builder.SinglestatPanel.MappingTypes = []*sdk.MapType{
		{
			Name:  &valueToText,
			Value: &valueToTextMapping,
		},
		{
			Name:  &rangeToText,
			Value: &rangeToTextMapping,
		},
	}
	panel.Builder.SinglestatPanel.SparkLine = sdk.SparkLine{}

	for _, opt := range append(defaults(), options...) {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

func defaults() []Option {
	return []Option{
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

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			singleStat.Builder.Links = append(singleStat.Builder.Links, link.Builder)
		}

		return nil
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(singleStat *SingleStat) error {
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

		return nil
	}
}

// WithGraphiteTarget adds a Graphite target to the graph.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(singleStat *SingleStat) error {
		singleStat.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(singleStat *SingleStat) error {
		singleStat.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.AddTarget(target.Builder)

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(singleStat *SingleStat) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		singleStat.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.Transparent = true

		return nil
	}
}

// Unit sets the unit of the data displayed on this axis.
func Unit(unit string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.Format = unit

		return nil
	}
}

// Decimals sets the number of decimals that should be displayed.
func Decimals(count int) Option {
	return func(singleStat *SingleStat) error {
		if count < 0 {
			return fmt.Errorf("decimals must be greater than 0: %w", errors.ErrInvalidArgument)
		}

		singleStat.Builder.SinglestatPanel.Decimals = count

		return nil
	}
}

// SparkLine displays the spark line summary of the series in addition to the
// single stat.
func SparkLine() Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.SparkLine.Show = true
		singleStat.Builder.SinglestatPanel.SparkLine.Full = false

		return nil
	}
}

// FullSparkLine displays a full height spark line summary of the series in
// addition to the single stat.
func FullSparkLine() Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.SparkLine.Show = true
		singleStat.Builder.SinglestatPanel.SparkLine.Full = true

		return nil
	}
}

// SparkLineColor sets the line color of the spark line.
func SparkLineColor(color string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.SparkLine.LineColor = &color

		return nil
	}
}

// SparkLineFillColor sets the color the spark line will be filled with.
func SparkLineFillColor(color string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.SparkLine.FillColor = &color

		return nil
	}
}

// SparkLineYMin defines the smallest value expected on the Y axis of the spark line.
func SparkLineYMin(value float64) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.SparkLine.YMin = &value

		return nil
	}
}

// SparkLineYMax defines the largest value expected on the Y axis of the spark line.
func SparkLineYMax(value float64) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.SparkLine.YMax = &value

		return nil
	}
}

// ValueType configures how the series will be reduced to a single value.
func ValueType(valueType StatType) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.ValueName = string(valueType)

		return nil
	}
}

// ValueFontSize sets the font size used to display the value (eg: "100%").
func ValueFontSize(size string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.ValueFontSize = size

		return nil
	}
}

// Prefix sets the text used as prefix of the value.
func Prefix(prefix string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.Prefix = &prefix

		return nil
	}
}

// PrefixFontSize sets the size used for the prefix text (eg: "110%").
func PrefixFontSize(size string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.PrefixFontSize = &size

		return nil
	}
}

// Postfix sets the text used as postfix of the value.
func Postfix(postfix string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.Postfix = &postfix

		return nil
	}
}

// PostfixFontSize sets the size used for the postfix text (eg: "110%")
func PostfixFontSize(size string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.PostfixFontSize = &size

		return nil
	}
}

// ColorValue will show the threshold's colors on the value itself.
func ColorValue() Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.ColorValue = true

		return nil
	}
}

// ColorBackground will show the threshold's colors in the background.
func ColorBackground() Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.ColorBackground = true

		return nil
	}
}

// Thresholds change the background and value colors dynamically within the
// panel, depending on the Singlestat value. The threshold is defined by 2
// values which represent 3 ranges that correspond to the three colors directly
// to the right.
func Thresholds(values [2]string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.Thresholds = strings.Join([]string{values[0], values[1]}, ",")

		return nil
	}
}

// Colors define which colors will be applied to the single value based on the
// threshold levels.
func Colors(values [3]string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.SinglestatPanel.Colors = []string{values[0], values[1], values[2]}

		return nil
	}
}

// ValuesToText allows to translate the value of the summary stat into explicit
// text.
func ValuesToText(mapping []ValueMap) Option {
	return func(singleStat *SingleStat) error {
		valueMap := make([]sdk.ValueMap, 0, len(mapping))

		for _, entry := range mapping {
			valueMap = append(valueMap, sdk.ValueMap{
				Op:       "=",
				TextType: entry.Text,
				Value:    entry.Value,
			})
		}

		mappingType := uint(valueToTextMapping)
		singleStat.Builder.SinglestatPanel.MappingType = &mappingType
		singleStat.Builder.SinglestatPanel.ValueMaps = valueMap

		return nil
	}
}

// RangesToText allows to translate the value of the summary stat into explicit
// text.
func RangesToText(mapping []RangeMap) Option {
	return func(singleStat *SingleStat) error {
		rangeMap := make([]*sdk.RangeMap, 0, len(mapping))

		for i := range mapping {
			rangeMap = append(rangeMap, &sdk.RangeMap{
				From: &mapping[i].From,
				To:   &mapping[i].To,
				Text: &mapping[i].Text,
			})
		}

		mappingType := uint(rangeToTextMapping)
		singleStat.Builder.SinglestatPanel.MappingType = &mappingType
		singleStat.Builder.SinglestatPanel.RangeMaps = rangeMap

		return nil
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(singleStat *SingleStat) error {
		singleStat.Builder.Repeat = &repeat

		return nil
	}
}
