package alert

import (
	"github.com/K-Phoen/sdk"
)

// ConditionEvaluator represents an option that can be used to configure a condition.
type ConditionEvaluator func(condition *condition)

// QueryReducer represents a function used to reduce a query to a single value
// that can then be fed to the evaluator to determine if the alert will be
// triggered or not.
type QueryReducer string

const (
	// Avg defines the query to execute and computes the average of the results.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Avg QueryReducer = "avg"

	// Sum defines the query to execute and computes the sum of the results.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Sum QueryReducer = "sum"

	// Count defines the query to execute and counts the results.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Count QueryReducer = "count"

	// Last defines the query to execute and takes the last result.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Last QueryReducer = "last"

	// Min defines the query to execute and takes the smallest result.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Min QueryReducer = "min"

	// Max defines the query to execute and takes the largest result.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Max QueryReducer = "max"

	// Median defines the query to execute and computes the mediam of the results.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Median QueryReducer = "median"

	// Diff defines the query to execute.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	Diff QueryReducer = "diff"

	// PercentDiff defines the query to execute.
	// See https://grafana.com/docs/grafana/latest/alerting/rules/#query-condition-example
	PercentDiff QueryReducer = "percent_diff"
)

// Operator represents a logical operator used to chain conditions.
type Operator string

// And chains conditions with a logical AND
const And Operator = "and"

// Or chains conditions with a logical OR
const Or Operator = "or"

type condition struct {
	builder *sdk.AlertCondition
}

func newCondition(reducer QueryReducer, queryRef string, evaluator ConditionEvaluator) *condition {
	cond := &condition{
		builder: &sdk.AlertCondition{
			Type:    "query",
			Query:   sdk.AlertConditionQueryRef{Params: []string{queryRef}},
			Reducer: sdk.AlertReducer{Type: string(reducer), Params: []string{}},
		},
	}

	evaluator(cond)

	return cond
}

// HasNoValue will match queries returning no value.
func HasNoValue() ConditionEvaluator {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "no_value", Params: []float64{}}
	}
}

// IsAbove will match queries returning a value above the given threshold.
func IsAbove(value float64) ConditionEvaluator {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "gt", Params: []float64{value}}
	}
}

// IsBelow will match queries returning a value below the given threshold.
func IsBelow(value float64) ConditionEvaluator {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "lt", Params: []float64{value}}
	}
}

// IsOutsideRange will match queries returning a value outside the given range.
func IsOutsideRange(min float64, max float64) ConditionEvaluator {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "outside_range", Params: []float64{min, max}}
	}
}

// IsWithinRange will match queries returning a value within the given range.
func IsWithinRange(min float64, max float64) ConditionEvaluator {
	return func(cond *condition) {
		cond.builder.Evaluator = sdk.AlertEvaluator{Type: "within_range", Params: []float64{min, max}}
	}
}
