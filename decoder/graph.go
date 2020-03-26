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

type dashboardGraph struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Targets    []target
	Axes       graphAxes
	Alert      *graphAlert
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

type graphAlert struct {
	Title            string
	EvaluateEvery    string `yaml:"evaluate_every"`
	For              string
	If               []alertCondition
	Notify           *int64
	Message          string
	OnNoData         string `yaml:"on_no_data"`
	OnExecutionError string `yaml:"on_execution_error"`
}

func (a graphAlert) toOptions() ([]alert.Option, error) {
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
	if a.Notify != nil {
		opts = append(opts, alert.NotifyChannel(*a.Notify))
	}
	if a.Message != "" {
		opts = append(opts, alert.Message(a.Message))
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

type alertThreshold struct {
	HasNoValue   bool `yaml:"has_no_value"`
	Above        *float64
	Below        *float64
	OutsideRange [2]float64 `yaml:"outside_range"`
	WithinRange  [2]float64 `yaml:"within_range"`
}

func (threshold alertThreshold) toOption() (alert.ConditionOption, error) {
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

type alertValue struct {
	Func     string
	QueryRef string `yaml:"ref"`
	From     string
	To       string
}

func (v alertValue) toOption() (alert.ConditionOption, error) {
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

type alertCondition struct {
	Operand   string
	Value     alertValue
	Threshold alertThreshold
}

func (c alertCondition) toOption() (alert.Option, error) {
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
