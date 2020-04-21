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
