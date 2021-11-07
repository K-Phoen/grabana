package timeseries

import (
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a graph panel.
type Option func(timeseries *TimeSeries)

// TooltipMode configures which series will be displayed in the tooltip.
type TooltipMode string

const (
	// SingleSeries will only display the hovered series.
	SingleSeries TooltipMode = "single"
	// AllSeries will display all series.
	AllSeries TooltipMode = "multi"
	// NoSeries will hide the tooltip completely.
	NoSeries TooltipMode = "none"
)

// LineInterpolationMode defines how Grafana interpolates series lines when drawn as lines.
type LineInterpolationMode string

const (
	// Points are joined by straight lines.
	Linear LineInterpolationMode = "linear"
	// Points are joined by curved lines resulting in smooth transitions between points.
	Smooth LineInterpolationMode = "smooth"
	// The line is displayed as steps between points. Points are rendered at the end of the step.
	StepBefore LineInterpolationMode = "stepBefore"
	// Line is displayed as steps between points. Points are rendered at the beginning of the step.
	StepAfter LineInterpolationMode = "stepAfter"
)

// BarAlignment defines how Grafana aligns bars.
type BarAlignment int

const (
	// The bar is drawn around the point. The point is placed in the center of the bar.
	AlignCenter BarAlignment = 0
	// The bar is drawn before the point. The point is placed on the trailing corner of the bar.
	AlignBefore BarAlignment = -1
	// The bar is drawn after the point. The point is placed on the leading corner of the bar.
	AlignAfter BarAlignment = 1
)

// GradientType defines the mode of the gradient fill.
type GradientType string

const (
	// No gradient fill.
	NoGradient GradientType = "none"
	// Transparency of the gradient is calculated based on the values on the y-axis.
	// Opacity of the fill is increasing with the values on the Y-axis.
	Opacity GradientType = "opacity"
	// Gradient color is generated based on the hue of the line color.
	Hue GradientType = "hue"
	// In this mode the whole bar will use a color gradient defined by the color scheme.
	Scheme GradientType = "scheme"
)

// LegendOption allows to configure a legend.
type LegendOption uint16

const (
	// Hide keeps the legend from being displayed.
	Hide LegendOption = iota
	// AsTable displays the legend as a table.
	AsTable
	// AsList displays the legend as a list.
	AsList
	// Bottom displays the legend below the graph.
	Bottom
	// ToTheRight displays the legend on the right side of the graph.
	ToTheRight

	// Min displays the smallest value of the series.
	Min
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
		FillOpacity(25),
		PointSize(5),
		Tooltip(SingleSeries),
		Legend(Bottom, AsList),
		Lines(Linear),
		GradientMode(Opacity),
		Axis(
			axis.Placement(axis.Auto),
			axis.Scale(axis.Linear),
		),
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

// FillOpacity defines the opacity level of the series. The lower the value, the more transparent.
func FillOpacity(value int) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.FillOpacity = value
	}
}

// PointSize adjusts the size of points.
func PointSize(value int) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize = value
	}
}

// Lines displays the series as lines, with a given interpolation strategy.
func Lines(mode LineInterpolationMode) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineInterpolation = string(mode)
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle = "line"
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineStyle = struct {
			Fill string `json:"fill"`
		}{
			Fill: "solid",
		}
	}
}

// Bars displays the series as bars, with a given alignment strategy.
func Bars(alignment BarAlignment) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.BarAlignment = int(alignment)
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle = "bars"
	}
}

// Points displays the series as points.
func Points() Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle = "points"
	}
}

// GradientMode sets the mode of the gradient fill.
func GradientMode(mode GradientType) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode = string(mode)
	}
}

// Axis configures the axis for this time series.
func Axis(options ...axis.Option) Option {
	return func(timeseries *TimeSeries) {
		axis.New(&timeseries.Builder.TimeseriesPanel.FieldConfig, options...)
	}
}

// Legend defines what should be shown in the legend.
func Legend(opts ...LegendOption) Option {
	return func(timeseries *TimeSeries) {
		legend := sdk.TimeseriesLegendOptions{
			DisplayMode: "list",
			Placement:   "bottom",
			Calcs:       make([]string, 0),
		}

		for _, opt := range opts {
			switch opt {
			case Hide:
				legend.DisplayMode = "hidden"
			case AsList:
				legend.DisplayMode = "list"
			case AsTable:
				legend.DisplayMode = "table"
			case ToTheRight:
				legend.Placement = "right"
			case Bottom:
				legend.Placement = "bottom"

			case First:
				legend.Calcs = append(legend.Calcs, "first")
			case FirstNonNull:
				legend.Calcs = append(legend.Calcs, "firstNotNull")
			case Last:
				legend.Calcs = append(legend.Calcs, "last")
			case LastNonNull:
				legend.Calcs = append(legend.Calcs, "lastNotNull")

			case Min:
				legend.Calcs = append(legend.Calcs, "min")
			case Max:
				legend.Calcs = append(legend.Calcs, "max")
			case Avg:
				legend.Calcs = append(legend.Calcs, "mean")

			case Count:
				legend.Calcs = append(legend.Calcs, "count")
			case Total:
				legend.Calcs = append(legend.Calcs, "sum")
			case Range:
				legend.Calcs = append(legend.Calcs, "range")
			}
		}

		timeseries.Builder.TimeseriesPanel.Options.Legend = legend
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
