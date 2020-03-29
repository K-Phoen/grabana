package decoder

import (
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/table"
)

type DashboardTable struct {
	Title                  string
	Span                   float32 `yaml:",omitempty"`
	Height                 string  `yaml:",omitempty"`
	Datasource             string  `yaml:",omitempty"`
	Targets                []Target
	HiddenColumns          []string            `yaml:"hidden_columns,flow"`
	TimeSeriesAggregations []table.Aggregation `yaml:"time_series_aggregations"`
}

func (tablePanel DashboardTable) toOption() (row.Option, error) {
	opts := []table.Option{}

	if tablePanel.Span != 0 {
		opts = append(opts, table.Span(tablePanel.Span))
	}
	if tablePanel.Height != "" {
		opts = append(opts, table.Height(tablePanel.Height))
	}
	if tablePanel.Datasource != "" {
		opts = append(opts, table.DataSource(tablePanel.Datasource))
	}

	for _, t := range tablePanel.Targets {
		opt, err := tablePanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, column := range tablePanel.HiddenColumns {
		opts = append(opts, table.HideColumn(column))
	}

	if len(tablePanel.TimeSeriesAggregations) != 0 {
		opts = append(opts, table.AsTimeSeriesAggregations(tablePanel.TimeSeriesAggregations))
	}

	return row.WithTable(tablePanel.Title, opts...), nil
}

func (tablePanel *DashboardTable) target(t Target) (table.Option, error) {
	if t.Prometheus != nil {
		return table.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, ErrTargetNotConfigured
}
