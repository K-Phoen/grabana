package graph

import (
	"github.com/grafana-tools/sdk"
)

type Option func(graph *Graph)

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
	Builder *sdk.Panel
}

func New(title string) *Graph {
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

	Editable()(panel)
	WithDefaultAxes()(panel)
	WithDefaultLegend()(panel)

	return panel
}

func WithDefaultAxes() Option {
	return func(graph *Graph) {
		graph.Builder.GraphPanel.YAxis = true
		graph.Builder.GraphPanel.XAxis = true
		graph.Builder.GraphPanel.Yaxes = []sdk.Axis{
			{Format: "short", Show: true, LogBase: 1},
			{Format: "short", Show: false},
		}
		graph.Builder.GraphPanel.Xaxis = sdk.Axis{
			Format:  "time",
			LogBase: 1,
			Show:    true,
		}
	}
}

func WithDefaultLegend() Option {
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

func WithPrometheusTarget(target PrometheusTarget) Option {
	return func(graph *Graph) {
		graph.Builder.AddTarget(&sdk.Target{
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

func Editable() Option {
	return func(graph *Graph) {
		graph.Builder.Editable = true
	}
}

func ReadOnly() Option {
	return func(graph *Graph) {
		graph.Builder.Editable = false
	}
}

func DataSource(datasource string) Option {
	return func(graph *Graph) {
		graph.Builder.Datasource = &datasource
	}
}
