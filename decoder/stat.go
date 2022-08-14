package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
)

var ErrInvalidStatOrientation = fmt.Errorf("invalid orientation")
var ErrInvalidStatColorMode = fmt.Errorf("invalid stat color mode")
var ErrInvalidStatThresholdMode = fmt.Errorf("invalid stat threshold mode")
var ErrInvalidStatTextMode = fmt.Errorf("invalid text mode")
var ErrInvalidStatValueType = fmt.Errorf("invalid stat value type")

type StatThresholdStep struct {
	Color string
	Value *float64 `yaml:",omitempty"`
}

type DashboardStat struct {
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

	SparkLine     bool   `yaml:"sparkline,omitempty"`
	Orientation   string `yaml:",omitempty"`
	Text          string `yaml:",omitempty"`
	ValueType     string `yaml:"value_type,omitempty"`
	ColorMode     string `yaml:"color_mode,omitempty"`
	TitleFontSize int    `yaml:"title_font_size,omitempty"`
	ValueFontSize int    `yaml:"value_font_size,omitempty"`

	ThresholdMode string              `yaml:"threshold_mode,omitempty"`
	Thresholds    []StatThresholdStep `yaml:",omitempty"`
}

func (statPanel DashboardStat) toOption() (row.Option, error) {
	opts := []stat.Option{}

	if statPanel.Description != "" {
		opts = append(opts, stat.Description(statPanel.Description))
	}
	if statPanel.Span != 0 {
		opts = append(opts, stat.Span(statPanel.Span))
	}
	if statPanel.Height != "" {
		opts = append(opts, stat.Height(statPanel.Height))
	}
	if statPanel.Transparent {
		opts = append(opts, stat.Transparent())
	}
	if statPanel.Datasource != "" {
		opts = append(opts, stat.DataSource(statPanel.Datasource))
	}
	if statPanel.Repeat != "" {
		opts = append(opts, stat.Repeat(statPanel.Repeat))
	}
	if len(statPanel.Links) != 0 {
		opts = append(opts, stat.Links(statPanel.Links.toModel()...))
	}
	if statPanel.Unit != "" {
		opts = append(opts, stat.Unit(statPanel.Unit))
	}
	if statPanel.Decimals != nil {
		opts = append(opts, stat.Decimals(*statPanel.Decimals))
	}
	if statPanel.SparkLine {
		opts = append(opts, stat.SparkLine())
	}
	if statPanel.TitleFontSize != 0 {
		opts = append(opts, stat.TitleFontSize(statPanel.TitleFontSize))
	}
	if statPanel.ValueFontSize != 0 {
		opts = append(opts, stat.ValueFontSize(statPanel.ValueFontSize))
	}

	if statPanel.Orientation != "" {
		opt, err := statPanel.orientationOpt()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	if statPanel.Text != "" {
		opt, err := statPanel.textOpt()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	if statPanel.ValueType != "" {
		opt, err := statPanel.valueType()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}
	if statPanel.ColorMode != "" {
		opt, err := statPanel.colorMode()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}
	if len(statPanel.Thresholds) != 0 {
		opt, err := statPanel.thresholds()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, t := range statPanel.Targets {
		opt, err := statPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithStat(statPanel.Title, opts...), nil
}

func (statPanel DashboardStat) thresholds() (stat.Option, error) {
	thresholds := make([]stat.ThresholdStep, 0, len(statPanel.Thresholds))
	for _, threshold := range statPanel.Thresholds {
		thresholds = append(thresholds, stat.ThresholdStep{
			Color: threshold.Color,
			Value: threshold.Value,
		})
	}

	switch statPanel.ThresholdMode {
	case "absolute":
		return stat.AbsoluteThresholds(thresholds), nil
	case "":
		return stat.AbsoluteThresholds(thresholds), nil
	case "relative":
		return stat.RelativeThresholds(thresholds), nil
	}

	return nil, fmt.Errorf("got mode '%s': %w", statPanel.ThresholdMode, ErrInvalidStatThresholdMode)
}

func (statPanel DashboardStat) colorMode() (stat.Option, error) {
	switch statPanel.ColorMode {
	case "background":
		return stat.ColorBackground(), nil
	case "value":
		return stat.ColorValue(), nil
	case "none":
		return stat.ColorNone(), nil
	default:
		return nil, ErrInvalidStatColorMode
	}
}

func (statPanel DashboardStat) valueType() (stat.Option, error) {
	switch statPanel.ValueType {
	case "min":
		return stat.ValueType(stat.Min), nil
	case "max":
		return stat.ValueType(stat.Max), nil
	case "avg":
		return stat.ValueType(stat.Avg), nil

	case "count":
		return stat.ValueType(stat.Count), nil
	case "total":
		return stat.ValueType(stat.Total), nil
	case "range":
		return stat.ValueType(stat.Range), nil

	case "first":
		return stat.ValueType(stat.First), nil
	case "first_non_null":
		return stat.ValueType(stat.FirstNonNull), nil
	case "last":
		return stat.ValueType(stat.Last), nil
	case "last_non_null":
		return stat.ValueType(stat.LastNonNull), nil
	default:
		return nil, ErrInvalidStatValueType
	}
}

func (statPanel DashboardStat) orientationOpt() (stat.Option, error) {
	switch statPanel.Orientation {
	case "horizontal":
		return stat.Orientation(stat.OrientationHorizontal), nil
	case "vertical":
		return stat.Orientation(stat.OrientationVertical), nil
	case "auto":
		return stat.Orientation(stat.OrientationAuto), nil
	default:
		return nil, ErrInvalidStatOrientation
	}
}

func (statPanel DashboardStat) textOpt() (stat.Option, error) {
	switch statPanel.Text {
	case "value":
		return stat.Text(stat.TextValue), nil
	case "name":
		return stat.Text(stat.TextName), nil
	case "value_and_name":
		return stat.Text(stat.TextValueAndName), nil
	case "none":
		return stat.Text(stat.TextNone), nil
	case "auto":
		return stat.Text(stat.TextAuto), nil
	default:
		return nil, ErrInvalidStatTextMode
	}
}

func (statPanel DashboardStat) target(t Target) (stat.Option, error) {
	if t.Prometheus != nil {
		return stat.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}
	if t.Graphite != nil {
		return stat.WithGraphiteTarget(t.Graphite.Query, t.Graphite.toOptions()...), nil
	}
	if t.InfluxDB != nil {
		return stat.WithInfluxDBTarget(t.InfluxDB.Query, t.InfluxDB.toOptions()...), nil
	}
	if t.Stackdriver != nil {
		stackdriverTarget, err := t.Stackdriver.toTarget()
		if err != nil {
			return nil, err
		}

		return stat.WithStackdriverTarget(stackdriverTarget), nil
	}

	return nil, ErrTargetNotConfigured
}
