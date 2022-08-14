package gauge

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/scheme"
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a stat panel.
type Option func(gauge *Gauge) error

type ThresholdStep struct {
	Color string
	Value *float64
}

// OrientationMode controls the layout.
type OrientationMode string

const (
	OrientationAuto       OrientationMode = ""
	OrientationHorizontal OrientationMode = "horizontal"
	OrientationVertical   OrientationMode = "vertical"
)

// ReductionType lets you set the function that your entire query is reduced into a
// single value with.
type ReductionType int

const (
	// Min displays the smallest value of the series.
	Min ReductionType = iota
	// Max displays the largest value of the series.
	Max
	// Avg displays the average of the series.
	Avg

	// First displays the first value of the series.
	First
	// FirstNonNull displays the first non-null value of the series.
	FirstNonNull
	// Last displays the last value of the series.
	Last
	// LastNonNull displays the last non-null value of the series.
	LastNonNull

	// Total displays the sum of values in the series.
	Total
	// Count displays the number of value in the series.
	Count
	// Range displays the difference between the minimum and maximum values.
	Range
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

// Gauge represents a stat panel.
type Gauge struct {
	Builder *sdk.Panel
}

// New creates a new gauge panel.
func New(title string, options ...Option) (*Gauge, error) {
	panel := &Gauge{Builder: sdk.NewGauge(title)}

	panel.Builder.IsNew = false

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
		ValueType(LastNonNull),
		Orientation(OrientationVertical),
		NoValue("N/A"),
		ColorScheme(scheme.ThresholdsValue(scheme.Last)),
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			gauge.Builder.Links = append(gauge.Builder.Links, link.Builder)
		}

		return nil
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(gauge *Gauge) error {
		gauge.Builder.AddTarget(&sdk.Target{
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

	return func(gauge *Gauge) error {
		gauge.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(gauge *Gauge) error {
		gauge.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.AddTarget(target.Builder)

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(gauge *Gauge) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		gauge.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(gauge *Gauge) error {
		gauge.Builder.Transparent = true

		return nil
	}
}

// Unit sets the unit of the data displayed on this axis.
func Unit(unit string) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.GaugePanel.FieldConfig.Defaults.Unit = unit

		return nil
	}
}

// Decimals sets the number of decimals that should be displayed.
func Decimals(count int) Option {
	return func(gauge *Gauge) error {
		if count < 0 {
			return fmt.Errorf("decimals must be greater than 0: %w", errors.ErrInvalidArgument)
		}

		gauge.Builder.GaugePanel.FieldConfig.Defaults.Decimals = &count

		return nil
	}
}

// ValueType configures how the series will be reduced to a single value.
func ValueType(valueType ReductionType) Option {
	return func(gauge *Gauge) error {
		var valType string

		switch valueType {
		case First:
			valType = "first"
		case FirstNonNull:
			valType = "firstNotNull"
		case Last:
			valType = "last"
		case LastNonNull:
			valType = "lastNotNull"

		case Min:
			valType = "min"
		case Max:
			valType = "max"
		case Avg:
			valType = "mean"

		case Count:
			valType = "count"
		case Total:
			valType = "sum"
		case Range:
			valType = "range"

		default:
			return fmt.Errorf("unknown value type: %w", errors.ErrInvalidArgument)
		}
		gauge.Builder.GaugePanel.Options.ReduceOptions.Calcs = []string{valType}

		return nil
	}
}

// ValueFontSize sets the font size used to display the value.
func ValueFontSize(size int) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.GaugePanel.Options.Text.ValueSize = size

		return nil
	}
}

// TitleFontSize sets the font size used to display the title.
func TitleFontSize(size int) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.GaugePanel.Options.Text.TitleSize = size

		return nil
	}
}

// AbsoluteThresholds changes the background and value colors dynamically within the
// panel, depending on the value. The threshold is defined by a series of steps
// values which, each having a value and an associated color.
func AbsoluteThresholds(steps []ThresholdStep) Option {
	return func(gauge *Gauge) error {
		sdkSteps := make([]sdk.ThresholdStep, 0, len(steps))
		for _, step := range steps {
			sdkSteps = append(sdkSteps, sdk.ThresholdStep{
				Color: step.Color,
				Value: step.Value,
			})
		}

		gauge.Builder.GaugePanel.FieldConfig.Defaults.Thresholds = sdk.Thresholds{
			Mode:  "absolute",
			Steps: sdkSteps,
		}

		return nil
	}
}

// RelativeThresholds changes the background and value colors dynamically within the
// panel, depending on the value. The threshold is defined by a series of steps
// values which, each having a value defined as a percentage and an associated color.
func RelativeThresholds(steps []ThresholdStep) Option {
	return func(gauge *Gauge) error {
		sdkSteps := make([]sdk.ThresholdStep, 0, len(steps))
		for _, step := range steps {
			sdkSteps = append(sdkSteps, sdk.ThresholdStep{
				Color: step.Color,
				Value: step.Value,
			})
		}

		gauge.Builder.GaugePanel.FieldConfig.Defaults.Thresholds = sdk.Thresholds{
			Mode:  "percentage",
			Steps: sdkSteps,
		}

		return nil
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.Repeat = &repeat

		return nil
	}
}

// ColorScheme configures the color scheme.
func ColorScheme(options ...scheme.Option) Option {
	return func(gauge *Gauge) error {
		scheme.New(&gauge.Builder.GaugePanel.FieldConfig, options...)

		return nil
	}
}

// NoValue defines what to show when there is no value.
func NoValue(text string) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.GaugePanel.FieldConfig.Defaults.NoValue = text

		return nil
	}
}

// Orientation changes the orientation of the layout.
func Orientation(mode OrientationMode) Option {
	return func(gauge *Gauge) error {
		gauge.Builder.GaugePanel.Options.Orientation = string(mode)

		return nil
	}
}
