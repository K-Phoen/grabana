package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/singlestat"
)

type dashboardSingleStat struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Unit       string
	Targets    []target
	Thresholds [2]string
	Colors     [3]string
	Color      []string
}

func (singleStatPanel dashboardSingleStat) toOption() (row.Option, error) {
	opts := []singlestat.Option{}

	if singleStatPanel.Span != 0 {
		opts = append(opts, singlestat.Span(singleStatPanel.Span))
	}
	if singleStatPanel.Height != "" {
		opts = append(opts, singlestat.Height(singleStatPanel.Height))
	}
	if singleStatPanel.Datasource != "" {
		opts = append(opts, singlestat.DataSource(singleStatPanel.Datasource))
	}
	if singleStatPanel.Unit != "" {
		opts = append(opts, singlestat.Unit(singleStatPanel.Unit))
	}
	if singleStatPanel.Thresholds[0] != "" {
		opts = append(opts, singlestat.Thresholds(singleStatPanel.Thresholds))
	}
	if singleStatPanel.Colors[0] != "" {
		opts = append(opts, singlestat.Colors(singleStatPanel.Colors))
	}

	for _, colorTarget := range singleStatPanel.Color {
		switch colorTarget {
		case "value":
			opts = append(opts, singlestat.ColorValue())
		case "background":
			opts = append(opts, singlestat.ColorBackground())
		default:
			return nil, fmt.Errorf("invalid coloring target '%s'", colorTarget)
		}
	}

	for _, t := range singleStatPanel.Targets {
		opt, err := singleStatPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithSingleStat(singleStatPanel.Title, opts...), nil
}

func (singleStatPanel dashboardSingleStat) target(t target) (singlestat.Option, error) {
	if t.Prometheus != nil {
		return singlestat.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, fmt.Errorf("target not configured")
}
