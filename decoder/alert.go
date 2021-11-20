package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/alert"
)

var ErrNoAlertThresholdDefined = fmt.Errorf("no threshold defined")
var ErrInvalidAlertValueFunc = fmt.Errorf("invalid alert value function")

type Alert struct {
	Title            string
	EvaluateEvery    string `yaml:"evaluate_every"`
	For              string
	If               []AlertCondition
	Notify           string            `yaml:",omitempty"`
	Notifications    []string          `yaml:",omitempty,flow"`
	Message          string            `yaml:",omitempty"`
	OnNoData         string            `yaml:"on_no_data"`
	OnExecutionError string            `yaml:"on_execution_error"`
	Tags             map[string]string `yaml:",omitempty"`
}

func (a Alert) toOptions() ([]alert.Option, error) {
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
	if len(a.Tags) != 0 {
		opts = append(opts, alert.Tags(a.Tags))
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
	HasNoValue   bool       `yaml:"has_no_value,omitempty"`
	Above        *float64   `yaml:",omitempty"`
	Below        *float64   `yaml:",omitempty"`
	OutsideRange [2]float64 `yaml:"outside_range,omitempty,flow"`
	WithinRange  [2]float64 `yaml:"within_range,omitempty,flow"`
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
	Value     AlertValue `yaml:",flow"`
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
