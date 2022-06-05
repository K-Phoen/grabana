package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/heatmap"
	"github.com/K-Phoen/grabana/heatmap/axis"
	"github.com/K-Phoen/grabana/row"
)

var ErrInvalidDataFormat = fmt.Errorf("invalid data format")

// DashboardHeatmap represents a heatmap panel.
type DashboardHeatmap struct {
	Title           string
	Description     string              `yaml:",omitempty"`
	Span            float32             `yaml:",omitempty"`
	Height          string              `yaml:",omitempty"`
	Transparent     bool                `yaml:",omitempty"`
	Datasource      string              `yaml:",omitempty"`
	Repeat          string              `yaml:",omitempty"`
	DataFormat      string              `yaml:"data_format,omitempty"`
	HideZeroBuckets bool                `yaml:"hide_zero_buckets"`
	HighlightCards  bool                `yaml:"highlight_cards"`
	Links           DashboardPanelLinks `yaml:",omitempty"`
	Targets         []Target
	ReverseYBuckets bool            `yaml:"reverse_y_buckets,omitempty"`
	Tooltip         *HeatmapTooltip `yaml:",omitempty"`
	YAxis           *HeatmapYAxis   `yaml:",omitempty"`
}

type HeatmapTooltip struct {
	Show          bool
	ShowHistogram bool
	Decimals      *int `yaml:",omitempty"`
}

func (tooltip *HeatmapTooltip) toOptions() []heatmap.Option {
	var opts []heatmap.Option

	if tooltip == nil {
		return nil
	}

	if !tooltip.Show {
		opts = append(opts, heatmap.HideTooltip())
	}
	if !tooltip.ShowHistogram {
		opts = append(opts, heatmap.HideTooltipHistogram())
	}
	if tooltip.Decimals != nil {
		opts = append(opts, heatmap.TooltipDecimals(*tooltip.Decimals))
	}

	return opts
}

type HeatmapYAxis struct {
	Decimals *int     `yaml:"decimals,omitempty"`
	Unit     string   `yaml:"unit,omitempty"`
	Max      *float64 `yaml:"max,omitempty"`
	Min      *float64 `yaml:"min,omitempty"`
}

func (yaxis *HeatmapYAxis) toOptions() []axis.Option {
	var opts []axis.Option

	if yaxis == nil {
		return nil
	}

	if yaxis.Decimals != nil {
		opts = append(opts, axis.Decimals(*yaxis.Decimals))
	}

	if yaxis.Unit != "" {
		opts = append(opts, axis.Unit(yaxis.Unit))
	}

	if yaxis.Min != nil {
		opts = append(opts, axis.Min(*yaxis.Min))
	}

	if yaxis.Max != nil {
		opts = append(opts, axis.Min(*yaxis.Max))
	}

	return opts
}

func (heatmapPanel DashboardHeatmap) toOption() (row.Option, error) {
	opts := []heatmap.Option{}

	if heatmapPanel.Description != "" {
		opts = append(opts, heatmap.Description(heatmapPanel.Description))
	}
	if heatmapPanel.Span != 0 {
		opts = append(opts, heatmap.Span(heatmapPanel.Span))
	}
	if heatmapPanel.Height != "" {
		opts = append(opts, heatmap.Height(heatmapPanel.Height))
	}
	if heatmapPanel.Transparent {
		opts = append(opts, heatmap.Transparent())
	}
	if heatmapPanel.Datasource != "" {
		opts = append(opts, heatmap.DataSource(heatmapPanel.Datasource))
	}
	if heatmapPanel.Repeat != "" {
		opts = append(opts, heatmap.Repeat(heatmapPanel.Repeat))
	}
	if len(heatmapPanel.Links) != 0 {
		opts = append(opts, heatmap.Links(heatmapPanel.Links.toModel()...))
	}
	if heatmapPanel.DataFormat != "" {
		switch heatmapPanel.DataFormat {
		case "time_series_buckets":
			opts = append(opts, heatmap.DataFormat(heatmap.TimeSeriesBuckets))
		case "time_series":
			opts = append(opts, heatmap.DataFormat(heatmap.TimeSeries))
		default:
			return nil, ErrInvalidDataFormat
		}
	}
	if heatmapPanel.HideZeroBuckets {
		opts = append(opts, heatmap.HideZeroBuckets())
	} else {
		opts = append(opts, heatmap.ShowZeroBuckets())
	}
	if heatmapPanel.HighlightCards {
		opts = append(opts, heatmap.HighlightCards())
	} else {
		opts = append(opts, heatmap.NoHighlightCards())
	}
	if heatmapPanel.ReverseYBuckets {
		opts = append(opts, heatmap.ReverseYBuckets())
	}
	if heatmapPanel.YAxis != nil {
		opts = append(opts, heatmap.YAxis(heatmapPanel.YAxis.toOptions()...))
	}
	opts = append(opts, heatmapPanel.Tooltip.toOptions()...)

	for _, t := range heatmapPanel.Targets {
		opt, err := heatmapPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithHeatmap(heatmapPanel.Title, opts...), nil
}

func (heatmapPanel DashboardHeatmap) target(t Target) (heatmap.Option, error) {
	if t.Prometheus != nil {
		return heatmap.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}
	if t.Graphite != nil {
		return heatmap.WithGraphiteTarget(t.Graphite.Query, t.Graphite.toOptions()...), nil
	}
	if t.InfluxDB != nil {
		return heatmap.WithInfluxDBTarget(t.InfluxDB.Query, t.InfluxDB.toOptions()...), nil
	}
	if t.Stackdriver != nil {
		stackdriverTarget, err := t.Stackdriver.toTarget()
		if err != nil {
			return nil, err
		}

		return heatmap.WithStackdriverTarget(stackdriverTarget), nil
	}

	return nil, ErrTargetNotConfigured
}
