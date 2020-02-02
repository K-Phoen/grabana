package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testQuery() []string {
	return []string{"A", "1m", "now"}
}

func TestNewConditionsCanBeCreated(t *testing.T) {
	req := require.New(t)

	a := newCondition()

	req.Equal("query", a.builder.Type)
}

func TestAvgQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Avg(query[0], query[1], query[2]))

	req.Equal("avg", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestSumQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Sum(query[0], query[1], query[2]))

	req.Equal("sum", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestCountQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Count(query[0], query[1], query[2]))

	req.Equal("count", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestMinQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Min(query[0], query[1], query[2]))

	req.Equal("min", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestMaxQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Max(query[0], query[1], query[2]))

	req.Equal("max", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestLastQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Last(query[0], query[1], query[2]))

	req.Equal("last", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestMedianQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Median(query[0], query[1], query[2]))

	req.Equal("median", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestDiffQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(Diff(query[0], query[1], query[2]))

	req.Equal("diff", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestPercentDiffQueryCanBeExpressed(t *testing.T) {
	req := require.New(t)
	query := testQuery()

	a := newCondition(PercentDiff(query[0], query[1], query[2]))

	req.Equal("percent_diff", a.builder.Reducer.Type)
	req.ElementsMatch(query, a.builder.Query.Params)
}

func TestHasNoValueEvaluator(t *testing.T) {
	req := require.New(t)

	a := newCondition(HasNoValue())

	req.Equal("no_value", a.builder.Evaluator.Type)
	req.Empty(a.builder.Evaluator.Params)
}

func TestIsAboveEvaluator(t *testing.T) {
	req := require.New(t)

	a := newCondition(IsAbove(10))

	req.Equal("gt", a.builder.Evaluator.Type)
	req.Equal([]float64{10}, a.builder.Evaluator.Params)
}

func TestIsBelowEvaluator(t *testing.T) {
	req := require.New(t)

	a := newCondition(IsBelow(10))

	req.Equal("lt", a.builder.Evaluator.Type)
	req.Equal([]float64{10}, a.builder.Evaluator.Params)
}

func TestIsOutsideRangeEvaluator(t *testing.T) {
	req := require.New(t)

	a := newCondition(IsOutsideRange(10, 20))

	req.Equal("outside_range", a.builder.Evaluator.Type)
	req.Equal([]float64{10, 20}, a.builder.Evaluator.Params)
}

func TestIIsWithinRangeEvaluator(t *testing.T) {
	req := require.New(t)

	a := newCondition(IsWithinRange(10, 20))

	req.Equal("within_range", a.builder.Evaluator.Type)
	req.Equal([]float64{10, 20}, a.builder.Evaluator.Params)
}
