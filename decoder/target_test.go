package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/target/stackdriver"

	"github.com/stretchr/testify/require"
)

func TestValidStackdriverAlignmentMethods(t *testing.T) {
	testCases := []struct {
		input    string
		expected stackdriver.Aligner
	}{
		{input: "none", expected: stackdriver.AlignNone},
		{input: "delta", expected: stackdriver.AlignDelta},
		{input: "rate", expected: stackdriver.AlignRate},
		{input: "interpolate", expected: stackdriver.AlignInterpolate},
		{input: "next_older", expected: stackdriver.AlignNextOlder},
		{input: "min", expected: stackdriver.AlignMin},
		{input: "max", expected: stackdriver.AlignMax},
		{input: "mean", expected: stackdriver.AlignMean},
		{input: "count", expected: stackdriver.AlignCount},
		{input: "sum", expected: stackdriver.AlignSum},
		{input: "stddev", expected: stackdriver.AlignStdDev},
		{input: "count_true", expected: stackdriver.AlignCountTrue},
		{input: "count_false", expected: stackdriver.AlignCountFalse},
		{input: "fraction_true", expected: stackdriver.AlignFractionTrue},
		{input: "percentile_99", expected: stackdriver.AlignPercentile99},
		{input: "percentile_95", expected: stackdriver.AlignPercentile95},
		{input: "percentile_50", expected: stackdriver.AlignPercentile50},
		{input: "percentile_05", expected: stackdriver.AlignPercentile05},
		{input: "percent_change", expected: stackdriver.AlignPercentChange},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := StackdriverAlignment{Method: tc.input}

			opt, err := panel.toOption()

			req.NoError(err)

			target := stackdriver.Delta("test")
			opt(target)

			req.Equal(string(tc.expected), target.Builder.PerSeriesAligner)
		})
	}
}

func TestValidStackdriverAggregations(t *testing.T) {
	testCases := []struct {
		input    string
		expected stackdriver.Reducer
	}{
		{input: "none", expected: stackdriver.ReduceNone},
		{input: "mean", expected: stackdriver.ReduceMean},
		{input: "min", expected: stackdriver.ReduceMin},
		{input: "max", expected: stackdriver.ReduceMax},
		{input: "sum", expected: stackdriver.ReduceSum},
		{input: "stddev", expected: stackdriver.ReduceStdDev},
		{input: "count", expected: stackdriver.ReduceCount},
		{input: "count_true", expected: stackdriver.ReduceCountTrue},
		{input: "count_false", expected: stackdriver.ReduceCountFalse},
		{input: "fraction_true", expected: stackdriver.ReduceCountFractionTrue},
		{input: "percentile_99", expected: stackdriver.ReducePercentile99},
		{input: "percentile_95", expected: stackdriver.ReducePercentile95},
		{input: "percentile_50", expected: stackdriver.ReducePercentile50},
		{input: "percentile_05", expected: stackdriver.ReducePercentile05},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := StackdriverTarget{Aggregation: tc.input}

			opt, err := panel.aggregation()

			req.NoError(err)

			target := stackdriver.Delta("test")
			opt(target)

			req.Equal(string(tc.expected), target.Builder.CrossSeriesReducer)
		})
	}
}

func TestValidStackdriverTargetTypes(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "delta", expected: "DELTA"},
		{input: "gauge", expected: "GAUGE"},
		{input: "cumulative", expected: "CUMULATIVE"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			targetPanel, err := StackdriverTarget{Type: tc.input}.toTarget()

			req.NoError(err)

			req.Equal(tc.expected, targetPanel.Builder.MetricKind)
		})
	}
}

func TestInvalidStackdriverTargetType(t *testing.T) {
	req := require.New(t)

	_, err := StackdriverTarget{Type: "invalid"}.toTarget()

	req.Error(err)
	req.Equal(ErrInvalidStackdriverType, err)
}

func TestStackdriverEqFilters(t *testing.T) {
	req := require.New(t)

	inputFilter := StackdriverFilters{Eq: map[string]string{
		"foo": "bar",
	}}

	options := inputFilter.toOptions()

	req.Len(options, 1)

	target := stackdriver.Delta("")
	stackdriver.Filter(options...)(target)

	req.Len(target.Builder.Filters, 3)
	req.Equal("foo", target.Builder.Filters[0])
	req.Equal("=", target.Builder.Filters[1])
	req.Equal("bar", target.Builder.Filters[2])
}

func TestStackdriverNeqFilters(t *testing.T) {
	req := require.New(t)

	inputFilter := StackdriverFilters{Neq: map[string]string{
		"neq": "val",
	}}

	options := inputFilter.toOptions()

	req.Len(options, 1)

	target := stackdriver.Delta("")
	stackdriver.Filter(options...)(target)

	req.Len(target.Builder.Filters, 3)
	req.Equal("neq", target.Builder.Filters[0])
	req.Equal("!=", target.Builder.Filters[1])
	req.Equal("val", target.Builder.Filters[2])
}

func TestStackdriverMatchesFilters(t *testing.T) {
	req := require.New(t)

	inputFilter := StackdriverFilters{Matches: map[string]string{
		"matches": "regex",
	}}

	options := inputFilter.toOptions()

	req.Len(options, 1)

	target := stackdriver.Delta("")
	stackdriver.Filter(options...)(target)

	req.Len(target.Builder.Filters, 3)
	req.Equal("matches", target.Builder.Filters[0])
	req.Equal("=~", target.Builder.Filters[1])
	req.Equal("regex", target.Builder.Filters[2])
}

func TestStackdriverNotMatchesFilters(t *testing.T) {
	req := require.New(t)

	inputFilter := StackdriverFilters{NotMatches: map[string]string{
		"notmatches": "regex",
	}}

	options := inputFilter.toOptions()

	req.Len(options, 1)

	target := stackdriver.Delta("")
	stackdriver.Filter(options...)(target)

	req.Len(target.Builder.Filters, 3)
	req.Equal("notmatches", target.Builder.Filters[0])
	req.Equal("!=~", target.Builder.Filters[1])
	req.Equal("regex", target.Builder.Filters[2])
}

func TestStackdriverHiddenTarget(t *testing.T) {
	req := require.New(t)

	target, err := StackdriverTarget{Type: "delta", Hidden: true}.toTarget()

	req.NoError(err)
	req.True(target.Builder.Hide)
}

func TestStackdriverGroupBy(t *testing.T) {
	req := require.New(t)

	target, err := StackdriverTarget{Type: "delta", GroupBy: []string{"field", "other"}}.toTarget()

	req.NoError(err)
	req.ElementsMatch(target.Builder.GroupBys, []string{"field", "other"})
}
