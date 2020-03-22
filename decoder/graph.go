package decoder

import (
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
)

type dashboardGraph struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Targets    []target
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
