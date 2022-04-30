package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/alert"
)

var ErrNoAlertThresholdDefined = fmt.Errorf("no threshold defined")
var ErrInvalidAlertValueFunc = fmt.Errorf("invalid alert value function")

type Alert struct {
	Summary     string
	Description string            `yaml:",omitempty"`
	Runbook     string            `yaml:",omitempty"`
	Tags        map[string]string `yaml:",omitempty"`

	EvaluateEvery    string `yaml:"evaluate_every"`
	For              string
	OnNoData         string `yaml:"on_no_data"`
	OnExecutionError string `yaml:"on_execution_error"`

	If []AlertCondition
	// TODO queries definition
}

func (a Alert) toOptions() ([]alert.Option, error) {
	opts := []alert.Option{
		alert.EvaluateEvery(a.EvaluateEvery),
		alert.For(a.For),
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

type AlertThreshold struct {
	HasNoValue   bool       `yaml:"has_no_value,omitempty"`
	Above        *float64   `yaml:",omitempty"`
	Below        *float64   `yaml:",omitempty"`
	OutsideRange [2]float64 `yaml:"outside_range,omitempty,flow"`
	WithinRange  [2]float64 `yaml:"within_range,omitempty,flow"`
}

func (threshold AlertThreshold) toOption() (alert.ConditionEvaluator, error) {
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

type AlertCondition struct {
	Reducer   string `yaml:"func"`
	QueryRef  string `yaml:"ref"`
	Threshold AlertThreshold
}

func (c AlertCondition) toOption() (alert.Option, error) {
	threshold, err := c.Threshold.toOption()
	if err != nil {
		return nil, err
	}

	reducer, err := c.queryReducer()
	if err != nil {
		return nil, err
	}

	return alert.If(reducer, c.QueryRef, threshold), nil
}

func (c AlertCondition) queryReducer() (alert.QueryReducer, error) {
	var queryReducer alert.QueryReducer

	switch c.Reducer {
	case "avg":
		queryReducer = alert.Avg
	case "sum":
		queryReducer = alert.Sum
	case "count":
		queryReducer = alert.Count
	case "last":
		queryReducer = alert.Last
	case "min":
		queryReducer = alert.Min
	case "max":
		queryReducer = alert.Max
	case "median":
		queryReducer = alert.Median
	case "diff":
		queryReducer = alert.Diff
	case "percent_diff":
		queryReducer = alert.PercentDiff
	default:
		return "", ErrInvalidAlertValueFunc
	}

	return queryReducer, nil
}
