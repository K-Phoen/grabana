package graph

import (
	"fmt"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/graph/series"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a graph panel.
type Option func(graph *Graph) error

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
	Alert   *alert.Alert
}

// New creates a new graph panel.
func New(title string, options ...Option) (*Graph, error) {
	panel := &Graph{Builder: sdk.NewGraph(title)}

	panel.Builder.AliasColors = make(map[string]interface{})
	panel.Builder.IsNew = false
	panel.Builder.GraphPanel.Tooltip.Sort = 2
	panel.Builder.GraphPanel.Tooltip.Shared = true

	for _, opt := range append(defaults(), options...) {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

func defaults() []Option {
	return []Option{
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
	return func(graph *Graph) error {
		graph.Builder.GraphPanel.YAxis = true
		graph.Builder.GraphPanel.XAxis = true
		graph.Builder.GraphPanel.Yaxes = []sdk.Axis{
			*axis.New().Builder,
			*axis.New(axis.Hide()).Builder,
		}
		graph.Builder.GraphPanel.Xaxis = *axis.New(axis.Unit("time")).Builder

		return nil
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(graph *Graph) error {
		graph.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			graph.Builder.Links = append(graph.Builder.Links, link.Builder)
		}

		return nil
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(graph *Graph) error {
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

		return nil
	}
}

// WithGraphiteTarget adds a Graphite target to the table.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(graph *Graph) error {
		graph.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(graph *Graph) error {
		graph.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(graph *Graph) error {
		graph.Builder.AddTarget(target.Builder)

		return nil
	}
}

// DataSource sets the data source to be used by the graph.
func DataSource(source string) Option {
	return func(graph *Graph) error {
		graph.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(graph *Graph) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		graph.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(graph *Graph) error {
		graph.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(graph *Graph) error {
		graph.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(graph *Graph) error {
		graph.Builder.Transparent = true

		return nil
	}
}

// LeftYAxis configures the left Y axis.
func LeftYAxis(opts ...axis.Option) Option {
	return func(graph *Graph) error {
		graph.Builder.GraphPanel.Yaxes[0] = *axis.New(opts...).Builder

		return nil
	}
}

// RightYAxis configures the right Y axis.
func RightYAxis(opts ...axis.Option) Option {
	return func(graph *Graph) error {
		graph.Builder.GraphPanel.Yaxes[1] = *axis.New(opts...).Builder

		return nil
	}
}

// XAxis configures the X axis.
func XAxis(opts ...axis.Option) Option {
	return func(graph *Graph) error {
		graph.Builder.GraphPanel.Xaxis = *axis.New(opts...).Builder

		return nil
	}
}

// Alert creates an alert for this graph.
func Alert(name string, opts ...alert.Option) Option {
	return func(graph *Graph) error {
		graph.Alert = alert.New(graph.Builder.Title, append(opts, alert.Summary(name))...)
		graph.Alert.Builder.Name = graph.Builder.Title

		return nil
	}
}

// Draw specifies how the graph will be drawn.
func Draw(modes ...DrawMode) Option {
	return func(graph *Graph) error {
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
			default:
				return errors.ErrInvalidArgument
			}
		}

		return nil
	}
}

// Fill defines the amount of color fill for a series (default 1, max 10, 0 is none).
func Fill(value int) Option {
	return func(graph *Graph) error {
		if value < 0 || value > 10 {
			return fmt.Errorf("fill must be between 0 and 10: %w", errors.ErrInvalidArgument)
		}

		graph.Builder.Fill = value

		return nil
	}
}

// LineWidth defines the width of the line for a series (default 1, max 10, 0 is none).
func LineWidth(value uint) Option {
	return func(graph *Graph) error {
		if value > 10 {
			return fmt.Errorf("line width must be between 0 and 10: %w", errors.ErrInvalidArgument)
		}

		graph.Builder.Linewidth = value

		return nil
	}
}

// Staircase draws adjacent points as staircase.
func Staircase() Option {
	return func(graph *Graph) error {
		graph.Builder.GraphPanel.SteppedLine = true

		return nil
	}
}

// PointRadius adjusts the size of points when Points are selected as Draw Mode.
func PointRadius(value float32) Option {
	return func(graph *Graph) error {
		if value < 0 || value > 10 {
			return fmt.Errorf("point radius must be between 0 and 10: %w", errors.ErrInvalidArgument)
		}

		graph.Builder.Pointradius = value

		return nil
	}
}

// Null configures how null values are displayed.
func Null(mode NullValue) Option {
	return func(graph *Graph) error {
		graph.Builder.GraphPanel.NullPointMode = string(mode)

		return nil
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(graph *Graph) error {
		graph.Builder.Repeat = &repeat

		return nil
	}
}

// SeriesOverride configures how null values are displayed.
// See https://grafana.com/docs/grafana/latest/panels/field-options/
func SeriesOverride(opts ...series.OverrideOption) Option {
	return func(graph *Graph) error {
		override := sdk.SeriesOverride{}

		for _, opt := range opts {
			if err := opt(&override); err != nil {
				return err
			}
		}

		graph.Builder.GraphPanel.SeriesOverrides = append(graph.Builder.GraphPanel.SeriesOverrides, override)

		return nil
	}
}

// Legend defines what should be shown in the legend.
func Legend(opts ...LegendOption) Option {
	return func(graph *Graph) error {
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
			default:
				return errors.ErrInvalidArgument
			}
		}

		graph.Builder.GraphPanel.Legend = legend

		return nil
	}
}
