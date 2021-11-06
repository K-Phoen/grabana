package timeseries

import (
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a graph panel.
type Option func(timeseries *TimeSeries)

// TooltipMode configures which series will be displayed in the tooltip
type TooltipMode string

const (
	// SingleSeries will only display the hovered series.
	SingleSeries TooltipMode = "single"
	// AllSeries will display all series.
	AllSeries = "multi"
	// NoSeries will hide the tooltip completely.
	NoSeries = "none"
)

// TimeSeries represents a time series panel.
type TimeSeries struct {
	Builder *sdk.Panel
}

// New creates a new time series panel.
func New(title string, options ...Option) *TimeSeries {
	panel := &TimeSeries{Builder: sdk.NewTimeseries(title)}
	panel.Builder.IsNew = false

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Span(6),
		LineWidth(1),
		Tooltip(SingleSeries),
	}
}

// DataSource sets the data source to be used by the graph.
func DataSource(source string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Datasource = &source
	}
}

// Tooltip configures the tooltip content.
func Tooltip(mode TooltipMode) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.Options.Tooltip.Mode = string(mode)
	}
}

// LineWidth defines the width of the line for a series (default 1, max 10, 0 is none).
func LineWidth(value int) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineWidth = value
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Height = &height
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Description = &content
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Transparent = true
	}
}

// Alert creates an alert for this graph.
func Alert(name string, opts ...alert.Option) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Alert = alert.New(name, opts...).Builder
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Repeat = &repeat
	}
}
