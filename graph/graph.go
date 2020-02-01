package graph

import (
	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/grafana-tools/sdk"
)

type Option func(graph *Graph)

type Graph struct {
	Builder *sdk.Panel
}

func New(title string, options ...Option) *Graph {
	panel := &Graph{Builder: sdk.NewGraph(title)}

	panel.Builder.AliasColors = make(map[string]interface{})
	panel.Builder.IsNew = false
	panel.Builder.Lines = true
	panel.Builder.Linewidth = 1
	panel.Builder.Fill = 1
	panel.Builder.Tooltip.Sort = 2
	panel.Builder.Tooltip.Shared = true
	panel.Builder.GraphPanel.NullPointMode = "null as zero"
	panel.Builder.GraphPanel.Lines = true
	panel.Builder.Span = 6

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Editable(),
		defaultAxes(),
		defaultLegend(),
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

func defaultLegend() Option {
	return func(graph *Graph) {
		graph.Builder.Legend = struct {
			AlignAsTable bool  `json:"alignAsTable"`
			Avg          bool  `json:"avg"`
			Current      bool  `json:"current"`
			HideEmpty    bool  `json:"hideEmpty"`
			HideZero     bool  `json:"hideZero"`
			Max          bool  `json:"max"`
			Min          bool  `json:"min"`
			RightSide    bool  `json:"rightSide"`
			Show         bool  `json:"show"`
			SideWidth    *uint `json:"sideWidth,omitempty"`
			Total        bool  `json:"total"`
			Values       bool  `json:"values"`
		}{
			AlignAsTable: false,
			Avg:          false,
			Current:      false,
			HideEmpty:    true,
			HideZero:     true,
			Max:          false,
			Min:          false,
			RightSide:    false,
			Show:         true,
			SideWidth:    nil,
			Total:        false,
			Values:       false,
		}
	}
}

func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(graph *Graph) {
		graph.Builder.AddTarget(&sdk.Target{
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
