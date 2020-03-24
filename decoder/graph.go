package decoder

import (
	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
)

type dashboardGraph struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Targets    []target
	Axes       graphAxes
}

func (graphPanel dashboardGraph) toOption() (row.Option, error) {
	opts := []graph.Option{}

	if graphPanel.Span != 0 {
		opts = append(opts, graph.Span(graphPanel.Span))
	}
	if graphPanel.Height != "" {
		opts = append(opts, graph.Height(graphPanel.Height))
	}
	if graphPanel.Datasource != "" {
		opts = append(opts, graph.DataSource(graphPanel.Datasource))
	}
	if graphPanel.Axes.Left != nil {
		opts = append(opts, graph.LeftYAxis(graphPanel.Axes.Left.toOptions()...))
	}
	if graphPanel.Axes.Bottom != nil {
		opts = append(opts, graph.XAxis(graphPanel.Axes.Bottom.toOptions()...))
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

func (graphPanel *dashboardGraph) target(t target) (graph.Option, error) {
	if t.Prometheus != nil {
		return graph.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, ErrTargetNotConfigured
}

type graphAxis struct {
	Hidden  *bool
	Label   string
	Unit    *string
	Min     *float64
	Max     *float64
	LogBase int `yaml:"log_base"`
}

func (a graphAxis) toOptions() []axis.Option {
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

	return opts
}

type graphAxes struct {
	Left   *graphAxis
	Right  *graphAxis
	Bottom *graphAxis
}
