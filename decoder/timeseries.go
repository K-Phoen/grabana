package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/timeseries/fields"
)

var ErrInvalidGradientMode = fmt.Errorf("invalid gradient mode")
var ErrInvalidLineInterpolationMode = fmt.Errorf("invalid line interpolation mode")
var ErrInvalidTooltipMode = fmt.Errorf("invalid tooltip mode")
var ErrInvalidStackMode = fmt.Errorf("invalid stack mode")
var ErrInvalidAxisDisplay = fmt.Errorf("invalid axis display")
var ErrInvalidAxisScale = fmt.Errorf("invalid axis scale")
var ErrInvalidOverrideMatcher = fmt.Errorf("invalid override matcher")

type DashboardTimeSeries struct {
	Title           string
	Description     string              `yaml:",omitempty"`
	Span            float32             `yaml:",omitempty"`
	Height          string              `yaml:",omitempty"`
	Transparent     bool                `yaml:",omitempty"`
	Datasource      string              `yaml:",omitempty"`
	Repeat          string              `yaml:",omitempty"`
	RepeatDirection string              `yaml:"repeat_direction,omitempty"`
	Links           DashboardPanelLinks `yaml:",omitempty"`
	Targets         []Target
	Legend          []string                 `yaml:",omitempty,flow"`
	Alert           *Alert                   `yaml:",omitempty"`
	Visualization   *TimeSeriesVisualization `yaml:",omitempty"`
	Axis            *TimeSeriesAxis          `yaml:",omitempty"`
	Overrides       []TimeSeriesOverride     `yaml:",omitempty"`
}

func (timeseriesPanel DashboardTimeSeries) toOption() (row.Option, error) {
	opts := []timeseries.Option{}

	if timeseriesPanel.Description != "" {
		opts = append(opts, timeseries.Description(timeseriesPanel.Description))
	}
	if timeseriesPanel.Span != 0 {
		opts = append(opts, timeseries.Span(timeseriesPanel.Span))
	}
	if timeseriesPanel.Height != "" {
		opts = append(opts, timeseries.Height(timeseriesPanel.Height))
	}
	if timeseriesPanel.Transparent {
		opts = append(opts, timeseries.Transparent())
	}
	if timeseriesPanel.Datasource != "" {
		opts = append(opts, timeseries.DataSource(timeseriesPanel.Datasource))
	}
	if timeseriesPanel.Repeat != "" {
		opts = append(opts, timeseries.Repeat(timeseriesPanel.Repeat))
	}
	if timeseriesPanel.RepeatDirection != "" {
		direction, err := parsePanelRepeatDirection(timeseriesPanel.RepeatDirection)
		if err != nil {
			return nil, err
		}
		opts = append(opts, timeseries.RepeatDirection(direction))
	}
	if len(timeseriesPanel.Links) != 0 {
		opts = append(opts, timeseries.Links(timeseriesPanel.Links.toModel()...))
	}
	if len(timeseriesPanel.Legend) != 0 {
		legendOpts, err := timeseriesPanel.legend()
		if err != nil {
			return nil, err
		}

		opts = append(opts, timeseries.Legend(legendOpts...))
	}
	if timeseriesPanel.Alert != nil {
		alertOpts, err := timeseriesPanel.Alert.toOptions()
		if err != nil {
			return nil, err
		}

		opts = append(opts, timeseries.Alert(timeseriesPanel.Alert.Summary, alertOpts...))
	}
	if timeseriesPanel.Visualization != nil {
		vizOpts, err := timeseriesPanel.Visualization.toOptions()
		if err != nil {
			return nil, err
		}

		opts = append(opts, vizOpts...)
	}
	if timeseriesPanel.Axis != nil {
		axisOpts, err := timeseriesPanel.Axis.toOptions()
		if err != nil {
			return nil, err
		}

		opts = append(opts, timeseries.Axis(axisOpts...))
	}

	for _, override := range timeseriesPanel.Overrides {
		opt, err := override.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, t := range timeseriesPanel.Targets {
		opt, err := timeseriesPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithTimeSeries(timeseriesPanel.Title, opts...), nil
}

func (timeseriesPanel DashboardTimeSeries) legend() ([]timeseries.LegendOption, error) {
	opts := make([]timeseries.LegendOption, 0, len(timeseriesPanel.Legend))

	for _, attribute := range timeseriesPanel.Legend {
		var opt timeseries.LegendOption

		switch attribute {
		case "hide":
			opt = timeseries.Hide
		case "as_table":
			opt = timeseries.AsTable
		case "as_list":
			opt = timeseries.AsList
		case "to_bottom":
			opt = timeseries.Bottom
		case "to_the_right":
			opt = timeseries.ToTheRight

		case "min":
			opt = timeseries.Min
		case "max":
			opt = timeseries.Max
		case "avg":
			opt = timeseries.Avg

		case "first":
			opt = timeseries.First
		case "first_non_null":
			opt = timeseries.FirstNonNull
		case "last":
			opt = timeseries.Last
		case "last_non_null":
			opt = timeseries.LastNonNull

		case "count":
			opt = timeseries.Count
		case "total":
			opt = timeseries.Total
		case "range":
			opt = timeseries.Range
		default:
			return nil, ErrInvalidLegendAttribute
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (timeseriesPanel DashboardTimeSeries) target(t Target) (timeseries.Option, error) {
	if t.Prometheus != nil {
		return timeseries.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}
	if t.Graphite != nil {
		return timeseries.WithGraphiteTarget(t.Graphite.Query, t.Graphite.toOptions()...), nil
	}
	if t.InfluxDB != nil {
		return timeseries.WithInfluxDBTarget(t.InfluxDB.Query, t.InfluxDB.toOptions()...), nil
	}
	if t.Loki != nil {
		return timeseries.WithLokiTarget(t.Loki.Query, t.Loki.toOptions()...), nil
	}
	if t.Stackdriver != nil {
		stackdriverTarget, err := t.Stackdriver.toTarget()
		if err != nil {
			return nil, err
		}

		return timeseries.WithStackdriverTarget(stackdriverTarget), nil
	}

	return nil, ErrTargetNotConfigured
}

type TimeSeriesVisualization struct {
	GradientMode      string `yaml:"gradient_mode,omitempty"`
	Tooltip           string `yaml:"tooltip,omitempty"`
	Stack             string `yaml:"stack,omitempty"`
	FillOpacity       *int   `yaml:"fill_opacity,omitempty"`
	PointSize         *int   `yaml:"point_size,omitempty"`
	LineInterpolation string `yaml:"line_interpolation,omitempty"`
	LineWidth         *int   `yaml:"line_width,omitempty"`
	// TODO: draw: {bars: {}, lines: {}}
}

func (timeseriesViz *TimeSeriesVisualization) toOptions() ([]timeseries.Option, error) {
	if timeseriesViz == nil {
		return nil, nil
	}

	opts := []timeseries.Option{}

	if timeseriesViz.FillOpacity != nil {
		opts = append(opts, timeseries.FillOpacity(*timeseriesViz.FillOpacity))
	}
	if timeseriesViz.PointSize != nil {
		opts = append(opts, timeseries.PointSize(*timeseriesViz.PointSize))
	}
	if timeseriesViz.GradientMode != "" {
		gradient, err := timeseriesViz.gradientModeOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, gradient)
	}
	if timeseriesViz.Tooltip != "" {
		gradient, err := timeseriesViz.tooltipOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, gradient)
	}
	if timeseriesViz.Stack != "" {
		stack, err := timeseriesViz.stackOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, stack)
	}
	if timeseriesViz.LineInterpolation != "" {
		interpolationOpt, err := timeseriesViz.lineInterpolationOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, interpolationOpt)
	}
	if timeseriesViz.LineWidth != nil {
		opts = append(opts, timeseries.LineWidth(*timeseriesViz.LineWidth))
	}

	return opts, nil
}

func (timeseriesViz *TimeSeriesVisualization) lineInterpolationOption() (timeseries.Option, error) {
	var mode timeseries.LineInterpolationMode
	switch timeseriesViz.LineInterpolation {
	case "linear":
		mode = timeseries.Linear
	case "smooth":
		mode = timeseries.Smooth
	case "step_before":
		mode = timeseries.StepBefore
	case "step_after":
		mode = timeseries.StepAfter
	default:
		return nil, ErrInvalidLineInterpolationMode
	}

	return timeseries.Lines(mode), nil
}

func (timeseriesViz *TimeSeriesVisualization) gradientModeOption() (timeseries.Option, error) {
	var mode timeseries.GradientType
	switch timeseriesViz.GradientMode {
	case "none":
		mode = timeseries.NoGradient
	case "opacity":
		mode = timeseries.Opacity
	case "hue":
		mode = timeseries.Hue
	case "scheme":
		mode = timeseries.Scheme
	default:
		return nil, ErrInvalidGradientMode
	}

	return timeseries.GradientMode(mode), nil
}

func (timeseriesViz *TimeSeriesVisualization) tooltipOption() (timeseries.Option, error) {
	var mode timeseries.TooltipMode
	switch timeseriesViz.Tooltip {
	case "single_series":
		mode = timeseries.SingleSeries
	case "all_series":
		mode = timeseries.AllSeries
	case "none":
		mode = timeseries.NoSeries
	default:
		return timeseries.Tooltip(mode), ErrInvalidTooltipMode
	}

	return timeseries.Tooltip(mode), nil
}

func (timeseriesViz *TimeSeriesVisualization) stackOption() (timeseries.Option, error) {
	var mode timeseries.StackMode
	switch timeseriesViz.Stack {
	case "normal":
		mode = timeseries.NormalStack
	case "percent":
		mode = timeseries.PercentStack
	case "none":
		mode = timeseries.Unstacked
	default:
		return timeseries.Stack(mode), ErrInvalidStackMode
	}

	return timeseries.Stack(mode), nil
}

type TimeSeriesAxis struct {
	SoftMin *int     `yaml:"soft_min,omitempty"`
	SoftMax *int     `yaml:"soft_max,omitempty"`
	Min     *float64 `yaml:",omitempty"`
	Max     *float64 `yaml:",omitempty"`

	Decimals *int `yaml:",omitempty"`

	Display string `yaml:",omitempty"`
	Scale   string `yaml:",omitempty"`

	Unit  string `yaml:",omitempty"`
	Label string `yaml:",omitempty"`
}

func (tsAxis *TimeSeriesAxis) toOptions() ([]axis.Option, error) {
	opts := []axis.Option{}

	if tsAxis.SoftMin != nil {
		opts = append(opts, axis.SoftMin(*tsAxis.SoftMin))
	}
	if tsAxis.SoftMax != nil {
		opts = append(opts, axis.SoftMax(*tsAxis.SoftMax))
	}
	if tsAxis.Min != nil {
		opts = append(opts, axis.Min(*tsAxis.Min))
	}
	if tsAxis.Max != nil {
		opts = append(opts, axis.Max(*tsAxis.Max))
	}

	if tsAxis.Decimals != nil {
		opts = append(opts, axis.Decimals(*tsAxis.Decimals))
	}

	if tsAxis.Display != "" {
		opt, err := tsAxis.placementOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	if tsAxis.Scale != "" {
		opt, err := tsAxis.scaleOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}

	if tsAxis.Unit != "" {
		opts = append(opts, axis.Unit(tsAxis.Unit))
	}
	if tsAxis.Label != "" {
		opts = append(opts, axis.Label(tsAxis.Label))
	}

	return opts, nil
}

func (tsAxis *TimeSeriesAxis) placementOption() (axis.Option, error) {
	placementMode, err := axisPlacementFromString(tsAxis.Display)
	if err != nil {
		return nil, err
	}

	return axis.Placement(placementMode), nil
}

func axisPlacementFromString(input string) (axis.PlacementMode, error) {
	switch input {
	case "none":
		return axis.Hidden, nil
	case "hidden":
		return axis.Hidden, nil
	case "auto":
		return axis.Auto, nil
	case "left":
		return axis.Left, nil
	case "right":
		return axis.Right, nil
	default:
		return axis.Auto, ErrInvalidAxisDisplay
	}
}

func (tsAxis *TimeSeriesAxis) scaleOption() (axis.Option, error) {
	var scaleMode axis.ScaleMode

	switch tsAxis.Scale {
	case "linear":
		scaleMode = axis.Linear
	case "log2":
		scaleMode = axis.Log2
	case "log10":
		scaleMode = axis.Log10
	default:
		return nil, ErrInvalidAxisScale
	}

	return axis.Scale(scaleMode), nil
}

type TimeSeriesOverride struct {
	Matcher    TimeSeriesOverrideMatcher `yaml:"match,flow"`
	Properties TimeSeriesOverrideProperties
}

func (override TimeSeriesOverride) toOption() (timeseries.Option, error) {
	matcher, err := override.Matcher.toOption()
	if err != nil {
		return nil, err
	}

	overrideOpts, err := override.Properties.toOptions()
	if err != nil {
		return nil, err
	}

	return timeseries.FieldOverride(matcher, overrideOpts...), nil
}

type TimeSeriesOverrideMatcher struct {
	FieldName *string `yaml:"field_name,omitempty"`
	QueryRef  *string `yaml:"query_ref,omitempty"`
	Regex     *string `yaml:"regex,omitempty"`
	Type      *string `yaml:"field_type,omitempty"`
}

func (matcher TimeSeriesOverrideMatcher) toOption() (fields.Matcher, error) {
	if matcher.FieldName != nil {
		return fields.ByName(*matcher.FieldName), nil
	}
	if matcher.QueryRef != nil {
		return fields.ByQuery(*matcher.QueryRef), nil
	}
	if matcher.Regex != nil {
		return fields.ByRegex(*matcher.Regex), nil
	}
	if matcher.Type != nil {
		return fields.ByType(fields.FieldType(*matcher.Type)), nil
	}

	return nil, ErrInvalidOverrideMatcher
}

type TimeSeriesOverrideProperties struct {
	Unit        *string `yaml:",omitempty"`
	Color       *string `yaml:"color,omitempty"`
	FillOpacity *int    `yaml:"fill_opacity,omitempty"`
	NegativeY   *bool   `yaml:"negative_Y,omitempty"`
	AxisDisplay *string `yaml:"axis_display,omitempty"`
	Stack       *string `yaml:",omitempty"`
}

func (properties TimeSeriesOverrideProperties) toOptions() ([]fields.OverrideOption, error) {
	var opts []fields.OverrideOption

	if properties.Unit != nil {
		opts = append(opts, fields.Unit(*properties.Unit))
	}
	if properties.Color != nil {
		opts = append(opts, fields.FixedColorScheme(*properties.Color))
	}
	if properties.FillOpacity != nil {
		opts = append(opts, fields.FillOpacity(*properties.FillOpacity))
	}
	if properties.NegativeY != nil && *properties.NegativeY {
		opts = append(opts, fields.NegativeY())
	}
	if properties.AxisDisplay != nil {
		axisPlacement, err := axisPlacementFromString(*properties.AxisDisplay)
		if err != nil {
			return nil, err
		}

		opts = append(opts, fields.AxisPlacement(axisPlacement))
	}
	if properties.Stack != nil {
		opts = append(opts, fields.Stack(fields.StackMode(*properties.Stack)))
	}

	return opts, nil
}
