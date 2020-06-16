package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
)

var ErrNoAlertThresholdDefined = fmt.Errorf("no threshold defined")
var ErrInvalidAlertValueFunc = fmt.Errorf("invalid alert value function")
var ErrInvalidLegendAttribute = fmt.Errorf("invalid legend attribute")

type DashboardGraph struct {
	Title      string
	Span       float32 `yaml:",omitempty"`
	Height     string  `yaml:",omitempty"`
	Datasource string  `yaml:",omitempty"`
	Targets    []Target
	Axes       *GraphAxes  `yaml:",omitempty"`
	Legend     []string    `yaml:",omitempty,flow"`
	Alert      *GraphAlert `yaml:",omitempty"`
}

func (graphPanel DashboardGraph) toOption() (row.Option, error) {
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

		opts = append(opts, graph.Alert(graphPanel.Alert.Title, alertOpts...))
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

type GraphAlert struct {
	Title            string
	EvaluateEvery    string `yaml:"evaluate_every"`
	For              string
	If               []AlertCondition
	Notify           string
	Notifications    []string
	Message          string
	OnNoData         string `yaml:"on_no_data"`
	OnExecutionError string `yaml:"on_execution_error"`
}

func (a GraphAlert) toOptions() ([]alert.Option, error) {
	opts := []alert.Option{
		alert.EvaluateEvery(a.EvaluateEvery),
		alert.For(a.For),
	}

	if a.OnNoData != "" {
		var mode alert.NoDataMode

		switch a.OnNoData {
		case "no_data":
			mode = alert.NoData
		case "alerting":
			mode = alert.Error
		case "keep_state":
			mode = alert.KeepLastState
		case "ok":
			mode = alert.OK
		default:
			return nil, fmt.Errorf("unknown on_no_data mode '%s'", a.OnNoData)
		}

		opts = append(opts, alert.OnNoData(mode))
	}
	if a.OnExecutionError != "" {
		var mode alert.ErrorMode

		switch a.OnExecutionError {
		case "alerting":
			mode = alert.Alerting
		case "keep_state":
			mode = alert.LastState
		default:
			return nil, fmt.Errorf("unknown on_execution_error mode '%s'", a.OnExecutionError)
		}

		opts = append(opts, alert.OnExecutionError(mode))
	}
	if a.Notify != "" {
		opts = append(opts, alert.NotifyChannel(a.Notify))
	}
	if a.Message != "" {
		opts = append(opts, alert.Message(a.Message))
	}

	for _, channel := range a.Notifications {
		opts = append(opts, alert.NotifyChannel(channel))
	}

	for _, condition := range a.If {
		conditionOpt, err := condition.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, conditionOpt)
	}

	return opts, nil
}

type AlertThreshold struct {
	HasNoValue   bool `yaml:"has_no_value"`
	Above        *float64
	Below        *float64
	OutsideRange [2]float64 `yaml:"outside_range"`
	WithinRange  [2]float64 `yaml:"within_range"`
}

func (threshold AlertThreshold) toOption() (alert.ConditionOption, error) {
	if threshold.HasNoValue {
		return alert.HasNoValue(), nil
	}
	if threshold.Above != nil {
		return alert.IsAbove(*threshold.Above), nil
	}
	if threshold.Below != nil {
		return alert.IsBelow(*threshold.Below), nil
	}
	if threshold.OutsideRange[0] != 0 && threshold.OutsideRange[1] != 0 {
		return alert.IsOutsideRange(threshold.OutsideRange[0], threshold.OutsideRange[1]), nil
	}
	if threshold.WithinRange[0] != 0 && threshold.WithinRange[1] != 0 {
		return alert.IsWithinRange(threshold.WithinRange[0], threshold.WithinRange[1]), nil
	}

	return nil, ErrNoAlertThresholdDefined
}

type AlertValue struct {
	Func     string
	QueryRef string `yaml:"ref"`
	From     string
	To       string
}

func (v AlertValue) toOption() (alert.ConditionOption, error) {
	var alertFunc func(refID string, from string, to string) alert.ConditionOption

	switch v.Func {
	case "avg":
		alertFunc = alert.Avg
	case "sum":
		alertFunc = alert.Sum
	case "count":
		alertFunc = alert.Count
	case "last":
		alertFunc = alert.Last
	case "min":
		alertFunc = alert.Min
	case "max":
		alertFunc = alert.Max
	case "median":
		alertFunc = alert.Median
	case "diff":
		alertFunc = alert.Diff
	case "percent_diff":
		alertFunc = alert.PercentDiff
	default:
		return nil, ErrInvalidAlertValueFunc
	}

	return alertFunc(v.QueryRef, v.From, v.To), nil
}

type AlertCondition struct {
	Operand   string
	Value     AlertValue
	Threshold AlertThreshold
}

func (c AlertCondition) toOption() (alert.Option, error) {
	operand := alert.And
	if c.Operand == "or" {
		operand = alert.Or
	}

	threshold, err := c.Threshold.toOption()
	if err != nil {
		return nil, err
	}

	value, err := c.Value.toOption()
	if err != nil {
		return nil, err
	}

	return alert.If(operand, value, threshold), nil
}
