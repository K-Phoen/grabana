package graph

import (
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/grafana-tools/sdk"
)

// Option represents an option that can be used to configure a graph panel.
type Option func(graph *Graph)

// DrawMode represents a type of visualization that will be drawn in the graph
// (lines, bars, points)
type DrawMode uint8

const (
	// Bars will display bars.
	Bars DrawMode = iota
	// Lines will display lines.
	Lines
	// Points will display points.
	Points
)

// NullValue describes how null values are displayed.
type NullValue string

const (
	// AsZero treats null values as zero values.
	AsZero NullValue = "null as zero"

	// AsNull treats null values as null.
	AsNull NullValue = "null"

	// Connected connects null values.
	Connected NullValue = "connected"
)

// LegendOption allows to configure a legend.
type LegendOption uint16

const (
	// Hide keeps the legend from being displayed.
	Hide LegendOption = iota
	// AsTable displays the legend as a table.
	AsTable
	// ToTheRight displays the legend on the right side of the graph.
	ToTheRight
	// Min displays the smallest value of the series.
	Min
	// Max displays the largest value of the series.
	Max
	// Avg displays the average of the series.
	Avg
	// Current displays the current value of the series.
	Current
	// Total displays the total value of the series.
	Total
	// NoNullSeries hides series with only null values from the legend.
	NoNullSeries
	// NoZeroSeries hides series with only 0 values from the legend.
	NoZeroSeries
)

// Graph represents a graph panel.
type Graph struct {
	Builder *sdk.Panel
}

// New creates a new graph panel.
func New(title string, options ...Option) *Graph {
	panel := &Graph{Builder: sdk.NewGraph(title)}

	panel.Builder.AliasColors = make(map[string]interface{})
	panel.Builder.IsNew = false
	panel.Builder.Tooltip.Sort = 2
	panel.Builder.Tooltip.Shared = true

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Editable(),
		Draw(Lines),
		Span(6),
		Fill(1),
		Null(AsZero),
		LineWidth(1),
		Legend(NoZeroSeries, NoNullSeries),
		defaultAxes(),
	}
}

func defaultAxes() Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.YAxis = true
		graph.Builder.GraphPanel.XAxis = true
		graph.Builder.GraphPanel.Yaxes = []sdk.Axis{
			*axis.New().Builder,
			*axis.New(axis.Hide()).Builder,
		}
		graph.Builder.GraphPanel.Xaxis = *axis.New(axis.Unit("time")).Builder
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(graph *Graph) {
		graph.Builder.AddTarget(&sdk.Target{
			RefID:          target.Ref,
			Hide:           target.Hidden,
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
	return func(graph *Graph) {
		graph.Builder.AddTarget(target.Builder)
	}
}

// Editable marks the graph as editable.
func Editable() Option {
	return func(graph *Graph) {
		graph.Builder.Editable = true
	}
}

// ReadOnly marks the graph as non-editable.
func ReadOnly() Option {
	return func(graph *Graph) {
		graph.Builder.Editable = false
	}
}

// DataSource sets the data source to be used by the graph.
func DataSource(source string) Option {
	return func(graph *Graph) {
		graph.Builder.Datasource = &source
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(graph *Graph) {
		graph.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(graph *Graph) {
		graph.Builder.Height = &height
	}
}

// LeftYAxis configures the left Y axis.
func LeftYAxis(opts ...axis.Option) Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.Yaxes[0] = *axis.New(opts...).Builder
	}
}

// RightYAxis configures the right Y axis.
func RightYAxis(opts ...axis.Option) Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.Yaxes[1] = *axis.New(opts...).Builder
	}
}

// XAxis configures the X axis.
func XAxis(opts ...axis.Option) Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.Xaxis = *axis.New(opts...).Builder
	}
}

// Alert creates an alert for this graph.
func Alert(name string, opts ...alert.Option) Option {
	return func(graph *Graph) {
		graph.Builder.Alert = alert.New(name, opts...).Builder
	}
}

// Draw specifies how the graph will be drawn.
func Draw(modes ...DrawMode) Option {
	return func(graph *Graph) {
		graph.Builder.Bars = false
		graph.Builder.Lines = false
		graph.Builder.Points = false

		for _, mode := range modes {
			switch mode {
			case Bars:
				graph.Builder.Bars = true
			case Lines:
				graph.Builder.Lines = true
			case Points:
				graph.Builder.Points = true
			}
		}
	}
}

// Fill defines the amount of color fill for a series (default 1, max 10, 0 is none).
func Fill(value int) Option {
	return func(graph *Graph) {
		graph.Builder.Fill = value
	}
}

// LineWidth defines the width of the line for a series (default 1, max 10, 0 is none).
func LineWidth(value uint) Option {
	return func(graph *Graph) {
		graph.Builder.Linewidth = value
	}
}

// Staircase draws adjacent points as staircase.
func Staircase() Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.SteppedLine = true
	}
}

// PointRadius adjusts the size of points when Points are selected as Draw Mode.
func PointRadius(value float32) Option {
	return func(graph *Graph) {
		graph.Builder.Pointradius = value
	}
}

// Null configures how null values are displayed.
func Null(mode NullValue) Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.NullPointMode = string(mode)
	}
}

// Legend defines what should be shown in the legend.
func Legend(opts ...LegendOption) Option {
	return func(graph *Graph) {
		legend := sdk.Legend{Show: true}

		for _, opt := range opts {
			switch opt {
			case Hide:
				legend.Show = false
			case AsTable:
				legend.AlignAsTable = true
			case ToTheRight:
				legend.RightSide = true
			case Min:
				legend.Min = true
				legend.Values = true
			case Max:
				legend.Max = true
				legend.Values = true
			case Avg:
				legend.Avg = true
				legend.Values = true
			case Current:
				legend.Current = true
				legend.Values = true
			case Total:
				legend.Total = true
				legend.Values = true
			case NoNullSeries:
				legend.HideEmpty = true
			case NoZeroSeries:
				legend.HideZero = true
			}
		}

		graph.Builder.Legend = legend
	}
}
