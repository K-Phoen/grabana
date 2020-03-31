package alert

import (
	"github.com/grafana-tools/sdk"
)

// Operator represents a logical operator used to chain conditions.
type Operator string

// ConditionOption represents an option that can be used to configure a condition.
type ConditionOption func(condition *condition)

// And chains conditions with a logical AND
const And Operator = "and"

// Or chains conditions with a logical OR
const Or Operator = "or"

type condition struct {
	builder *sdk.AlertCondition
}

func newCondition(options ...ConditionOption) *condition {
	cond := &condition{builder: &sdk.AlertCondition{
		Type: "query",
	}}

	for _, opt := range options {
		opt(cond)
	}

	return cond
}

// Avg defines the query to execute and computes the average of the results.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Avg(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "avg", Params: []string{}}
	}
}

// Sum defines the query to execute and computes the sum of the results.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Sum(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "sum", Params: []string{}}
	}
}

// Count defines the query to execute and counts the results.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Count(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "count", Params: []string{}}
	}
}

// Last defines the query to execute and takes the last result.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Last(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "last", Params: []string{}}
	}
}

// Min defines the query to execute and takes the smallest result.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Min(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "min", Params: []string{}}
	}
}

// Max defines the query to execute and takes the largest result.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Max(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "max", Params: []string{}}
	}
}

// Median defines the query to execute and computes the mediam of the results.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Median(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "median", Params: []string{}}
	}
}

// Diff defines the query to execute.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func Diff(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "diff", Params: []string{}}
	}
}

// PercentDiff defines the query to execute.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
func PercentDiff(refID string, from string, to string) ConditionOption {
	return func(cond *condition) {
		cond.builder.Query = sdk.AlertQuery{Params: []string{refID, from, to}}
		cond.builder.Reducer = sdk.AlertReducer{Type: "percent_diff", Params: []string{}}
	}
}

// HasNoValue will match queries returning no value.
func HasNoValue() ConditionOption {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "no_value", Params: []float64{}}
	}
}

// IsAbove will match queries returning a value above the given threshold.
func IsAbove(value float64) ConditionOption {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "gt", Params: []float64{value}}
	}
}

// IsBelow will match queries returning a value below the given threshold.
func IsBelow(value float64) ConditionOption {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "lt", Params: []float64{value}}
	}
}

// IsOutsideRange will match queries returning a value outside the given range.
func IsOutsideRange(min float64, max float64) ConditionOption {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "outside_range", Params: []float64{min, max}}
	}
}

// IsWithinRange will match queries returning a value within the given range.
func IsWithinRange(min float64, max float64) ConditionOption {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "within_range", Params: []float64{min, max}}
	}
}
