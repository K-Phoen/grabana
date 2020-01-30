package grabana

import (
	"github.com/grafana-tools/sdk"
)

type GraphOption func(graph *Graph)

type PrometheusTarget struct {
	RefID      string
	Datasource string

	Expr           string
	IntervalFactor int
	Interval       string
	Step           int
	LegendFormat   string
	Instant        bool
	Format         string
}

type Graph struct {
	builder *sdk.Panel
}

func GraphDefaults(graph *Graph) {
	graph.builder.AliasColors = make(map[string]interface{})
	graph.builder.IsNew = false
	graph.builder.Lines = true
	graph.builder.Linewidth = 1
	graph.builder.Fill = 1
	graph.builder.Tooltip.Sort = 2
	graph.builder.Tooltip.Shared = true
	graph.builder.GraphPanel.NullPointMode = "null as zero"
	graph.builder.GraphPanel.Lines = true
	graph.builder.Span = 6

	Editable()(graph)
	WithDefaultAxes()(graph)
	WithDefaultLegend()(graph)
}

func WithDefaultAxes() GraphOption {
	return func(graph *Graph) {
		graph.builder.GraphPanel.YAxis = true
		graph.builder.GraphPanel.XAxis = true
		graph.builder.GraphPanel.Yaxes = []sdk.Axis{
			{Format: "short", Show: true, LogBase: 1},
			{Format: "short", Show: false},
		}
		graph.builder.GraphPanel.Xaxis = sdk.Axis{
			Format:  "time",
			LogBase: 1,
			Show:    true,
		}
	}
}

func WithDefaultLegend() GraphOption {
	return func(graph *Graph) {
		graph.builder.Legend = struct {
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

func WithPrometheusTarget(target PrometheusTarget) GraphOption {
	return func(graph *Graph) {
		graph.builder.AddTarget(&sdk.Target{
			RefID:          target.RefID,
			Datasource:     target.Datasource,
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

func Editable() GraphOption {
	return func(graph *Graph) {
		graph.builder.Editable = true
	}
}

func ReadOnly() GraphOption {
	return func(graph *Graph) {
		graph.builder.Editable = false
	}
}

func WithDataSource(datasource string) GraphOption {
	return func(graph *Graph) {
		graph.builder.Datasource = &datasource
	}
}
