package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestCondition(evaluator ConditionEvaluator) *condition {
	return newCondition(Avg, "A", evaluator)
}

func TestQueryReducers(t *testing.T) {
	testCases := []struct {
		reducer             QueryReducer
		expectedReducerType string
	}{
		{reducer: Avg, expectedReducerType: "avg"},
		{reducer: Sum, expectedReducerType: "sum"},
		{reducer: Count, expectedReducerType: "count"},
		{reducer: Last, expectedReducerType: "last"},
		{reducer: Min, expectedReducerType: "min"},
		{reducer: Max, expectedReducerType: "max"},
		{reducer: Median, expectedReducerType: "median"},
		{reducer: Diff, expectedReducerType: "diff"},
		{reducer: PercentDiff, expectedReducerType: "percent_diff"},
	}

	for _, test := range testCases {
		tc := test

		t.Run(tc.expectedReducerType, func(t *testing.T) {
			req := require.New(t)

			a := newCondition(tc.reducer, "A", IsAbove(1))

			req.Equal(tc.expectedReducerType, a.builder.Reducer.Type)
			req.Equal("A", a.builder.Query.Params[0])
		})
	}
}

func TestHasNoValueEvaluator(t *testing.T) {
	req := require.New(t)

	a := newTestCondition(HasNoValue())

	req.Equal("no_value", a.builder.Evaluator.Type)
	req.Empty(a.builder.Evaluator.Params)
}

func TestIsAboveEvaluator(t *testing.T) {
	req := require.New(t)

	a := newTestCondition(IsAbove(10))

	req.Equal("gt", a.builder.Evaluator.Type)
	req.Equal([]float64{10}, a.builder.Evaluator.Params)
}

func TestIsBelowEvaluator(t *testing.T) {
	req := require.New(t)

	a := newTestCondition(IsBelow(10))

	req.Equal("lt", a.builder.Evaluator.Type)
	req.Equal([]float64{10}, a.builder.Evaluator.Params)
}

func TestIsOutsideRangeEvaluator(t *testing.T) {
	req := require.New(t)

	a := newTestCondition(IsOutsideRange(10, 20))

	req.Equal("outside_range", a.builder.Evaluator.Type)
	req.Equal([]float64{10, 20}, a.builder.Evaluator.Params)
}

func TestIIsWithinRangeEvaluator(t *testing.T) {
	req := require.New(t)

	a := newTestCondition(IsWithinRange(10, 20))

	req.Equal("within_range", a.builder.Evaluator.Type)
	req.Equal([]float64{10, 20}, a.builder.Evaluator.Params)
}
