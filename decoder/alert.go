package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/alert"
)

var ErrNoAlertThresholdDefined = fmt.Errorf("no threshold defined")
var ErrInvalidAlertValueFunc = fmt.Errorf("invalid alert value function")
var ErrInvalidAlertOperand = fmt.Errorf("invalid alert operand")
var ErrNoConditionOnAlert = fmt.Errorf("no condition defined on alert")
var ErrNoTargetOnAlert = fmt.Errorf("no target defined on alert")

type Alert struct {
	Summary     string
	Description string            `yaml:",omitempty"`
	Runbook     string            `yaml:",omitempty"`
	Tags        map[string]string `yaml:",omitempty"`

	EvaluateEvery    string `yaml:"evaluate_every"`
	For              string
	OnNoData         string `yaml:"on_no_data"`
	OnExecutionError string `yaml:"on_execution_error"`

	If      []AlertCondition
	Targets []AlertTarget
}

func (a Alert) toOptions() ([]alert.Option, error) {
	opts := []alert.Option{}

	if len(a.If) == 0 {
		return nil, ErrNoConditionOnAlert
	}
	if len(a.Targets) == 0 {
		return nil, ErrNoTargetOnAlert
	}

	if a.EvaluateEvery != "" {
		opts = append(opts, alert.EvaluateEvery(a.EvaluateEvery))
	}
	if a.For != "" {
		opts = append(opts, alert.For(a.For))
	}

	if a.OnNoData != "" {
		noDataOpt, err := a.noDataOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, noDataOpt)
	}
	if a.OnExecutionError != "" {
		execErrorOpt, err := a.executionErrorOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, execErrorOpt)
	}
	if a.Description != "" {
		opts = append(opts, alert.Description(a.Description))
	}
	if a.Runbook != "" {
		opts = append(opts, alert.Runbook(a.Runbook))
	}
	if len(a.Tags) != 0 {
		opts = append(opts, alert.Tags(a.Tags))
	}

	for _, condition := range a.If {
		conditionOpt, err := condition.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, conditionOpt)
	}

	targetOpts, err := a.targetOptions()
	if err != nil {
		return nil, err
	}

	return append(opts, targetOpts...), nil
}

func (a Alert) targetOptions() ([]alert.Option, error) {
	opts := make([]alert.Option, 0, len(a.Targets))

	for _, alertTarget := range a.Targets {
		opt, err := alertTarget.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (a Alert) noDataOption() (alert.Option, error) {
	var mode alert.NoDataMode

	switch a.OnNoData {
	case "no_data":
		mode = alert.NoDataEmpty
	case "alerting":
		mode = alert.NoDataAlerting
	case "ok":
		mode = alert.NoDataOK
	default:
		return nil, fmt.Errorf("unknown on_no_data mode '%s'", a.OnNoData)
	}

	return alert.OnNoData(mode), nil
}

func (a Alert) executionErrorOption() (alert.Option, error) {
	var mode alert.ErrorMode

	switch a.OnExecutionError {
	case "alerting":
		mode = alert.ErrorAlerting
	case "error":
		mode = alert.ErrorKO
	case "ok":
		mode = alert.ErrorOK
	default:
		return nil, fmt.Errorf("unknown on_execution_error mode '%s'", a.OnExecutionError)
	}

	return alert.OnExecutionError(mode), nil
}

type AlertCondition struct {
	Operand *string `yaml:"operand,omitempty"`

	// Query reducers, only one should be used
	Avg         *string `yaml:"avg,omitempty"`
	Sum         *string `yaml:"sum,omitempty"`
	Count       *string `yaml:"count,omitempty"`
	Last        *string `yaml:"last,omitempty"`
	Min         *string `yaml:"min,omitempty"`
	Max         *string `yaml:"max,omitempty"`
	Median      *string `yaml:"median,omitempty"`
	Diff        *string `yaml:"diff,omitempty"`
	PercentDiff *string `yaml:"percent_diff,omitempty"`

	HasNoValue   bool       `yaml:"has_no_value,omitempty"`
	Above        *float64   `yaml:",omitempty"`
	Below        *float64   `yaml:",omitempty"`
	OutsideRange [2]float64 `yaml:"outside_range,omitempty,flow"`
	WithinRange  [2]float64 `yaml:"within_range,omitempty,flow"`
}

func (c AlertCondition) toOption() (alert.Option, error) {
	var err error
	alertOpt := alert.If

	if c.Operand != nil {
		switch *c.Operand {
		case string(alert.And):
			alertOpt = alert.If
		case string(alert.Or):
			alertOpt = alert.IfOr
		default:
			return nil, ErrInvalidAlertOperand
		}
	}

	reducer, queryRef, err := c.queryReducer()
	if err != nil {
		return nil, err
	}

	threshold, err := c.toThresholdOption()
	if err != nil {
		return nil, err
	}

	return alertOpt(reducer, queryRef, threshold), nil
}

func (c AlertCondition) queryReducer() (alert.QueryReducer, string, error) {
	if c.Avg != nil {
		return alert.Avg, *c.Avg, nil
	}
	if c.Sum != nil {
		return alert.Sum, *c.Sum, nil
	}
	if c.Count != nil {
		return alert.Count, *c.Count, nil
	}
	if c.Last != nil {
		return alert.Last, *c.Last, nil
	}
	if c.Min != nil {
		return alert.Min, *c.Min, nil
	}
	if c.Max != nil {
		return alert.Max, *c.Max, nil
	}
	if c.Median != nil {
		return alert.Median, *c.Median, nil
	}
	if c.Diff != nil {
		return alert.Diff, *c.Diff, nil
	}
	if c.PercentDiff != nil {
		return alert.PercentDiff, *c.PercentDiff, nil
	}

	return "", "", ErrInvalidAlertValueFunc
}

func (c AlertCondition) toThresholdOption() (alert.ConditionEvaluator, error) {
	if c.HasNoValue {
		return alert.HasNoValue(), nil
	}
	if c.Above != nil {
		return alert.IsAbove(*c.Above), nil
	}
	if c.Below != nil {
		return alert.IsBelow(*c.Below), nil
	}
	if c.OutsideRange[0] != 0 && c.OutsideRange[1] != 0 {
		return alert.IsOutsideRange(c.OutsideRange[0], c.OutsideRange[1]), nil
	}
	if c.WithinRange[0] != 0 && c.WithinRange[1] != 0 {
		return alert.IsWithinRange(c.WithinRange[0], c.WithinRange[1]), nil
	}

	return nil, ErrNoAlertThresholdDefined
}
