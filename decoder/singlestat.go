package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/singlestat"
)

var ErrInvalidColoringTarget = fmt.Errorf("invalid coloring target")
var ErrInvalidSparkLineMode = fmt.Errorf("invalid sparkline mode")
var ErrInvalidSingleStatValueType = fmt.Errorf("invalid single stat value type")

type DashboardSingleStat struct {
	Title           string
	Description     string              `yaml:",omitempty"`
	Span            float32             `yaml:",omitempty"`
	Height          string              `yaml:",omitempty"`
	Transparent     bool                `yaml:",omitempty"`
	Datasource      string              `yaml:",omitempty"`
	Repeat          string              `yaml:",omitempty"`
	Links           DashboardPanelLinks `yaml:",omitempty"`
	Unit            string
	Decimals        *int   `yaml:",omitempty"`
	ValueType       string `yaml:"value_type"`
	ValueFontSize   string `yaml:"value_font_size,omitempty"`
	PrefixFontSize  string `yaml:"prefix_font_size,omitempty"`
	PostfixFontSize string `yaml:"postfix_font_size,omitempty"`
	SparkLine       string `yaml:"sparkline"`
	Targets         []Target
	Thresholds      [2]string
	Colors          [3]string
	Color           []string              `yaml:",omitempty"`
	RangesToText    []singlestat.RangeMap `yaml:"ranges_to_text,omitempty"`
}

func (singleStatPanel DashboardSingleStat) toOption() (row.Option, error) {
	opts := []singlestat.Option{}

	if singleStatPanel.Description != "" {
		opts = append(opts, singlestat.Description(singleStatPanel.Description))
	}
	if singleStatPanel.Span != 0 {
		opts = append(opts, singlestat.Span(singleStatPanel.Span))
	}
	if singleStatPanel.Height != "" {
		opts = append(opts, singlestat.Height(singleStatPanel.Height))
	}
	if singleStatPanel.Transparent {
		opts = append(opts, singlestat.Transparent())
	}
	if singleStatPanel.Datasource != "" {
		opts = append(opts, singlestat.DataSource(singleStatPanel.Datasource))
	}
	if singleStatPanel.Repeat != "" {
		opts = append(opts, singlestat.Repeat(singleStatPanel.Repeat))
	}
	if len(singleStatPanel.Links) != 0 {
		opts = append(opts, singlestat.Links(singleStatPanel.Links.toModel()...))
	}
	if singleStatPanel.Unit != "" {
		opts = append(opts, singlestat.Unit(singleStatPanel.Unit))
	}
	if singleStatPanel.Decimals != nil {
		opts = append(opts, singlestat.Decimals(*singleStatPanel.Decimals))
	}
	if singleStatPanel.Thresholds[0] != "" {
		opts = append(opts, singlestat.Thresholds(singleStatPanel.Thresholds))
	}
	if singleStatPanel.Colors[0] != "" {
		opts = append(opts, singlestat.Colors(singleStatPanel.Colors))
	}
	if singleStatPanel.ValueFontSize != "" {
		opts = append(opts, singlestat.ValueFontSize(singleStatPanel.ValueFontSize))
	}
	if singleStatPanel.PrefixFontSize != "" {
		opts = append(opts, singlestat.PrefixFontSize(singleStatPanel.PrefixFontSize))
	}
	if singleStatPanel.PostfixFontSize != "" {
		opts = append(opts, singlestat.PostfixFontSize(singleStatPanel.PostfixFontSize))
	}
	if singleStatPanel.RangesToText != nil {
		opts = append(opts, singlestat.RangesToText(singleStatPanel.RangesToText))
	}

	switch singleStatPanel.SparkLine {
	case "bottom":
		opts = append(opts, singlestat.SparkLine())
	case "full":
		opts = append(opts, singlestat.FullSparkLine())
	case "":
	default:
		return nil, ErrInvalidSparkLineMode
	}

	if singleStatPanel.ValueType != "" {
		opt, err := singleStatPanel.valueType()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, colorTarget := range singleStatPanel.Color {
		switch colorTarget {
		case "value":
			opts = append(opts, singlestat.ColorValue())
		case "background":
			opts = append(opts, singlestat.ColorBackground())
		default:
			return nil, ErrInvalidColoringTarget
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

func (singleStatPanel DashboardSingleStat) valueType() (singlestat.Option, error) {
	switch singleStatPanel.ValueType {
	case "min":
		return singlestat.ValueType(singlestat.Min), nil
	case "max":
		return singlestat.ValueType(singlestat.Max), nil
	case "avg":
		return singlestat.ValueType(singlestat.Avg), nil
	case "current":
		return singlestat.ValueType(singlestat.Current), nil
	case "total":
		return singlestat.ValueType(singlestat.Total), nil
	case "first":
		return singlestat.ValueType(singlestat.First), nil
	case "delta":
		return singlestat.ValueType(singlestat.Delta), nil
	case "diff":
		return singlestat.ValueType(singlestat.Diff), nil
	case "range":
		return singlestat.ValueType(singlestat.Range), nil
	case "name":
		return singlestat.ValueType(singlestat.Name), nil
	default:
		return nil, ErrInvalidSingleStatValueType
	}
}

func (singleStatPanel DashboardSingleStat) target(t Target) (singlestat.Option, error) {
	if t.Prometheus != nil {
		return singlestat.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}
	if t.Graphite != nil {
		return singlestat.WithGraphiteTarget(t.Graphite.Query, t.Graphite.toOptions()...), nil
	}
	if t.InfluxDB != nil {
		return singlestat.WithInfluxDBTarget(t.InfluxDB.Query, t.InfluxDB.toOptions()...), nil
	}
	if t.Stackdriver != nil {
		stackdriverTarget, err := t.Stackdriver.toTarget()
		if err != nil {
			return nil, err
		}

		return singlestat.WithStackdriverTarget(stackdriverTarget), nil
	}

	return nil, ErrTargetNotConfigured
}
