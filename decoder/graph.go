package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/graph/series"
	"github.com/K-Phoen/grabana/row"
)

var ErrInvalidLegendAttribute = fmt.Errorf("invalid legend attribute")

type DashboardGraph struct {
	Title         string
	Description   string  `yaml:",omitempty"`
	Span          float32 `yaml:",omitempty"`
	Height        string  `yaml:",omitempty"`
	Transparent   bool    `yaml:",omitempty"`
	Datasource    string  `yaml:",omitempty"`
	Repeat        string  `yaml:",omitempty"`
	Targets       []Target
	Links         DashboardPanelLinks `yaml:",omitempty"`
	Axes          *GraphAxes          `yaml:",omitempty"`
	Legend        []string            `yaml:",omitempty,flow"`
	Alert         *Alert              `yaml:",omitempty"`
	Visualization *GraphVisualization `yaml:",omitempty"`
}

func (graphPanel DashboardGraph) toOption() (row.Option, error) {
	opts := []graph.Option{}

	if graphPanel.Description != "" {
		opts = append(opts, graph.Description(graphPanel.Description))
	}
	if graphPanel.Span != 0 {
		opts = append(opts, graph.Span(graphPanel.Span))
	}
	if graphPanel.Height != "" {
		opts = append(opts, graph.Height(graphPanel.Height))
	}
	if graphPanel.Transparent {
		opts = append(opts, graph.Transparent())
	}
	if graphPanel.Datasource != "" {
		opts = append(opts, graph.DataSource(graphPanel.Datasource))
	}
	if graphPanel.Repeat != "" {
		opts = append(opts, graph.Repeat(graphPanel.Repeat))
	}
	if len(graphPanel.Links) != 0 {
		opts = append(opts, graph.Links(graphPanel.Links.toModel()...))
	}
	if graphPanel.Axes != nil && graphPanel.Axes.Right != nil {
		opts = append(opts, graph.RightYAxis(graphPanel.Axes.Right.toOptions()...))
	}
	if graphPanel.Axes != nil && graphPanel.Axes.Left != nil {
		opts = append(opts, graph.LeftYAxis(graphPanel.Axes.Left.toOptions()...))
	}
	if graphPanel.Axes != nil && graphPanel.Axes.Bottom != nil {
		opts = append(opts, graph.XAxis(graphPanel.Axes.Bottom.toOptions()...))
	}
	if len(graphPanel.Legend) != 0 {
		legendOpts, err := graphPanel.legend()
		if err != nil {
			return nil, err
		}

		opts = append(opts, graph.Legend(legendOpts...))
	}
	if graphPanel.Alert != nil {
		alertOpts, err := graphPanel.Alert.toOptions()
		if err != nil {
			return nil, err
		}

		opts = append(opts, graph.Alert(graphPanel.Alert.Summary, alertOpts...))
	}
	if graphPanel.Visualization != nil {
		opts = append(opts, graphPanel.Visualization.toOptions()...)
	}

	for _, t := range graphPanel.Targets {
		opt, err := graphPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithGraph(graphPanel.Title, opts...), nil
}

type GraphSeriesOverride struct {
	Alias     string
	Color     string `yaml:",omitempty"`
	Dashes    *bool  `yaml:",omitempty"`
	Lines     *bool  `yaml:",omitempty"`
	Fill      *int   `yaml:",omitempty"`
	LineWidth *int   `yaml:"line_width,omitempty"`
}

func (override *GraphSeriesOverride) toOption() graph.Option {
	overrideOpts := []series.OverrideOption{
		series.Alias(override.Alias),
	}

	if override.Color != "" {
		overrideOpts = append(overrideOpts, series.Color(override.Color))
	}
	if override.Dashes != nil {
		overrideOpts = append(overrideOpts, series.Dashes(*override.Dashes))
	}
	if override.Lines != nil {
		overrideOpts = append(overrideOpts, series.Dashes(*override.Lines))
	}
	if override.Fill != nil {
		overrideOpts = append(overrideOpts, series.Fill(*override.Fill))
	}
	if override.LineWidth != nil {
		overrideOpts = append(overrideOpts, series.LineWidth(*override.LineWidth))
	}

	return graph.SeriesOverride(overrideOpts...)
}

type GraphVisualization struct {
	NullValue string                `yaml:",omitempty"`
	Staircase bool                  `yaml:",omitempty"`
	Overrides []GraphSeriesOverride `yaml:"overrides,omitempty"`
}

func (graphViz *GraphVisualization) toOptions() []graph.Option {
	if graphViz == nil {
		return nil
	}

	opts := []graph.Option{}
	if graphViz.NullValue != "" {
		mode := graph.AsZero
		switch graphViz.NullValue {
		case "null as zero":
			mode = graph.AsZero
		case "null":
			mode = graph.AsNull
		case "connected":
			mode = graph.Connected
		}

		opts = append(opts, graph.Null(mode))
	}

	if graphViz.Staircase {
		opts = append(opts, graph.Staircase())
	}

	for _, override := range graphViz.Overrides {
		opts = append(opts, override.toOption())
	}

	return opts
}

func (graphPanel *DashboardGraph) legend() ([]graph.LegendOption, error) {
	opts := make([]graph.LegendOption, 0, len(graphPanel.Legend))

	for _, attribute := range graphPanel.Legend {
		var opt graph.LegendOption

		switch attribute {
		case "hide":
			opt = graph.Hide
		case "as_table":
			opt = graph.AsTable
		case "to_the_right":
			opt = graph.ToTheRight
		case "min":
			opt = graph.Min
		case "max":
			opt = graph.Max
		case "avg":
			opt = graph.Avg
		case "current":
			opt = graph.Current
		case "total":
			opt = graph.Total
		case "no_null_series":
			opt = graph.NoNullSeries
		case "no_zero_series":
			opt = graph.NoZeroSeries
		default:
			return nil, ErrInvalidLegendAttribute
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (graphPanel *DashboardGraph) target(t Target) (graph.Option, error) {
	if t.Prometheus != nil {
		return graph.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}
	if t.Graphite != nil {
		return graph.WithGraphiteTarget(t.Graphite.Query, t.Graphite.toOptions()...), nil
	}
	if t.InfluxDB != nil {
		return graph.WithInfluxDBTarget(t.InfluxDB.Query, t.InfluxDB.toOptions()...), nil
	}
	if t.Stackdriver != nil {
		stackdriverTarget, err := t.Stackdriver.toTarget()
		if err != nil {
			return nil, err
		}

		return graph.WithStackdriverTarget(stackdriverTarget), nil
	}

	return nil, ErrTargetNotConfigured
}

type GraphAxis struct {
	Hidden  *bool    `yaml:",omitempty"`
	Label   string   `yaml:",omitempty"`
	Unit    *string  `yaml:",omitempty"`
	Min     *float64 `yaml:",omitempty"`
	Max     *float64 `yaml:",omitempty"`
	LogBase int      `yaml:"log_base"`
}

func (a GraphAxis) toOptions() []axis.Option {
	opts := []axis.Option{}

	if a.Label != "" {
		opts = append(opts, axis.Label(a.Label))
	}
	if a.Unit != nil {
		opts = append(opts, axis.Unit(*a.Unit))
	}
	if a.Hidden != nil && *a.Hidden {
		opts = append(opts, axis.Hide())
	}
	if a.Min != nil {
		opts = append(opts, axis.Min(*a.Min))
	}
	if a.Max != nil {
		opts = append(opts, axis.Max(*a.Max))
	}
	if a.LogBase != 0 {
		opts = append(opts, axis.LogBase(a.LogBase))
	}

	return opts
}

type GraphAxes struct {
	Left   *GraphAxis `yaml:",omitempty"`
	Right  *GraphAxis `yaml:",omitempty"`
	Bottom *GraphAxis `yaml:",omitempty"`
}
