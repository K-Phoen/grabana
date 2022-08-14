package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/gauge"
	"github.com/K-Phoen/grabana/row"
)

var ErrInvalidGaugeThresholdMode = fmt.Errorf("invalid gauge threshold mode")
var ErrInvalidGaugeValueType = fmt.Errorf("invalid gauge value type")
var ErrInvalidGaugeOrientation = fmt.Errorf("invalid gauge orientation")

type GaugeThresholdStep struct {
	Color string
	Value *float64 `yaml:",omitempty"`
}

type DashboardGauge struct {
	Title       string
	Description string              `yaml:",omitempty"`
	Span        float32             `yaml:",omitempty"`
	Height      string              `yaml:",omitempty"`
	Transparent bool                `yaml:",omitempty"`
	Datasource  string              `yaml:",omitempty"`
	Repeat      string              `yaml:",omitempty"`
	Links       DashboardPanelLinks `yaml:",omitempty"`
	Targets     []Target

	Unit     string `yaml:",omitempty"`
	Decimals *int   `yaml:",omitempty"`

	Orientation   string `yaml:",omitempty"`
	ValueType     string `yaml:"value_type,omitempty"`
	TitleFontSize int    `yaml:"title_font_size,omitempty"`
	ValueFontSize int    `yaml:"value_font_size,omitempty"`

	ThresholdMode string               `yaml:"threshold_mode,omitempty"`
	Thresholds    []GaugeThresholdStep `yaml:",omitempty"`
}

func (gaugePanel DashboardGauge) toOption() (row.Option, error) {
	opts := []gauge.Option{}

	if gaugePanel.Description != "" {
		opts = append(opts, gauge.Description(gaugePanel.Description))
	}
	if gaugePanel.Span != 0 {
		opts = append(opts, gauge.Span(gaugePanel.Span))
	}
	if gaugePanel.Height != "" {
		opts = append(opts, gauge.Height(gaugePanel.Height))
	}
	if gaugePanel.Transparent {
		opts = append(opts, gauge.Transparent())
	}
	if gaugePanel.Datasource != "" {
		opts = append(opts, gauge.DataSource(gaugePanel.Datasource))
	}
	if gaugePanel.Repeat != "" {
		opts = append(opts, gauge.Repeat(gaugePanel.Repeat))
	}
	if len(gaugePanel.Links) != 0 {
		opts = append(opts, gauge.Links(gaugePanel.Links.toModel()...))
	}
	if gaugePanel.Unit != "" {
		opts = append(opts, gauge.Unit(gaugePanel.Unit))
	}
	if gaugePanel.Decimals != nil {
		opts = append(opts, gauge.Decimals(*gaugePanel.Decimals))
	}
	if gaugePanel.TitleFontSize != 0 {
		opts = append(opts, gauge.TitleFontSize(gaugePanel.TitleFontSize))
	}
	if gaugePanel.ValueFontSize != 0 {
		opts = append(opts, gauge.ValueFontSize(gaugePanel.ValueFontSize))
	}

	if gaugePanel.Orientation != "" {
		opt, err := gaugePanel.orientationOpt()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	if gaugePanel.ValueType != "" {
		opt, err := gaugePanel.valueType()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	if len(gaugePanel.Thresholds) != 0 {
		opt, err := gaugePanel.thresholds()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, t := range gaugePanel.Targets {
		opt, err := gaugePanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithGauge(gaugePanel.Title, opts...), nil
}

func (gaugePanel DashboardGauge) thresholds() (gauge.Option, error) {
	thresholds := make([]gauge.ThresholdStep, 0, len(gaugePanel.Thresholds))
	for _, threshold := range gaugePanel.Thresholds {
		thresholds = append(thresholds, gauge.ThresholdStep{
			Color: threshold.Color,
			Value: threshold.Value,
		})
	}

	switch gaugePanel.ThresholdMode {
	case "absolute":
		return gauge.AbsoluteThresholds(thresholds), nil
	case "":
		return gauge.AbsoluteThresholds(thresholds), nil
	case "relative":
		return gauge.RelativeThresholds(thresholds), nil
	}

	return nil, fmt.Errorf("got mode '%s': %w", gaugePanel.ThresholdMode, ErrInvalidGaugeThresholdMode)
}

func (gaugePanel DashboardGauge) valueType() (gauge.Option, error) {
	switch gaugePanel.ValueType {
	case "min":
		return gauge.ValueType(gauge.Min), nil
	case "max":
		return gauge.ValueType(gauge.Max), nil
	case "avg":
		return gauge.ValueType(gauge.Avg), nil

	case "count":
		return gauge.ValueType(gauge.Count), nil
	case "total":
		return gauge.ValueType(gauge.Total), nil
	case "range":
		return gauge.ValueType(gauge.Range), nil

	case "first":
		return gauge.ValueType(gauge.First), nil
	case "first_non_null":
		return gauge.ValueType(gauge.FirstNonNull), nil
	case "last":
		return gauge.ValueType(gauge.Last), nil
	case "last_non_null":
		return gauge.ValueType(gauge.LastNonNull), nil
	default:
		return nil, ErrInvalidGaugeValueType
	}
}

func (gaugePanel DashboardGauge) orientationOpt() (gauge.Option, error) {
	switch gaugePanel.Orientation {
	case "horizontal":
		return gauge.Orientation(gauge.OrientationHorizontal), nil
	case "vertical":
		return gauge.Orientation(gauge.OrientationVertical), nil
	case "auto":
		return gauge.Orientation(gauge.OrientationAuto), nil
	default:
		return nil, ErrInvalidGaugeOrientation
	}
}

func (gaugePanel DashboardGauge) target(t Target) (gauge.Option, error) {
	if t.Prometheus != nil {
		return gauge.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}
	if t.Graphite != nil {
		return gauge.WithGraphiteTarget(t.Graphite.Query, t.Graphite.toOptions()...), nil
	}
	if t.InfluxDB != nil {
		return gauge.WithInfluxDBTarget(t.InfluxDB.Query, t.InfluxDB.toOptions()...), nil
	}
	if t.Stackdriver != nil {
		stackdriverTarget, err := t.Stackdriver.toTarget()
		if err != nil {
			return nil, err
		}

		return gauge.WithStackdriverTarget(stackdriverTarget), nil
	}

	return nil, ErrTargetNotConfigured
}
