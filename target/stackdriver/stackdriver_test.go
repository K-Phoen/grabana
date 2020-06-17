package stackdriver_test

import (
	"testing"

	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestDeltaQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	target := stackdriver.Delta("some metric type for delta")

	req.Equal("DELTA", target.Builder.MetricKind)
	req.Equal("some metric type for delta", target.Builder.MetricType)
}

func TestGaugeQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	target := stackdriver.Gauge("some metric type for gauge")

	req.Equal("GAUGE", target.Builder.MetricKind)
	req.Equal("some metric type for gauge", target.Builder.MetricType)
}

func TestCumulativeQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	target := stackdriver.Cumulative("some metric type for cumulative")

	req.Equal("CUMULATIVE", target.Builder.MetricKind)
	req.Equal("some metric type for cumulative", target.Builder.MetricType)
}

func TestLegendCanBeConfigured(t *testing.T) {
	req := require.New(t)
	legend := "{{ code }} - {{ path }}"

	target := stackdriver.Delta("", stackdriver.Legend(legend))

	req.Equal(legend, target.Builder.AliasBy)
}

func TestRefCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := stackdriver.Delta("", stackdriver.Ref("A"))

	req.Equal("A", target.Builder.RefID)
}

func TestTargetCanBeHidden(t *testing.T) {
	req := require.New(t)

	target := stackdriver.Delta("", stackdriver.Hide())

	req.True(target.Builder.Hide)
}

func TestAggregationCanBeConfigured(t *testing.T) {
	req := require.New(t)
	reducers := []stackdriver.Reducer{
		stackdriver.ReduceNone,
		stackdriver.ReduceMean,
		stackdriver.ReduceMin,
		stackdriver.ReduceMax,
		stackdriver.ReduceSum,
		stackdriver.ReduceStdDev,
		stackdriver.ReduceCount,
		stackdriver.ReduceCountTrue,
		stackdriver.ReduceCountFalse,
		stackdriver.ReduceCountFractionTrue,
		stackdriver.ReducePercentile99,
		stackdriver.ReducePercentile95,
		stackdriver.ReducePercentile50,
		stackdriver.ReducePercentile05,
	}

	for _, reducer := range reducers {
		target := stackdriver.Delta("", stackdriver.Aggregation(reducer))

		req.Equal(string(reducer), target.Builder.CrossSeriesReducer)
	}
}

func TestAlignmentCanBeConfigured(t *testing.T) {
	req := require.New(t)
	aligners := []stackdriver.Aligner{
		stackdriver.AlignNone,
		stackdriver.AlignDelta,
		stackdriver.AlignRate,
		stackdriver.AlignNextOlder,
		stackdriver.AlignMin,
		stackdriver.AlignMax,
		stackdriver.AlignMean,
		stackdriver.AlignCount,
		stackdriver.AlignSum,
		stackdriver.AlignStdDev,
		stackdriver.AlignCountTrue,
		stackdriver.AlignCountFalse,
		stackdriver.AlignFractionTrue,
		stackdriver.AlignPercentile99,
		stackdriver.AlignPercentile95,
		stackdriver.AlignPercentile50,
		stackdriver.AlignPercentile05,
		stackdriver.AlignPercentChange,
	}

	for _, aligner := range aligners {
		target := stackdriver.Delta("", stackdriver.Alignment(aligner, "stackdriver-auto"))

		req.Equal(string(aligner), target.Builder.PerSeriesAligner)
		req.Equal("stackdriver-auto", target.Builder.AlignmentPeriod)
	}
}

func TestFiltersCanBeConfigured(t *testing.T) {
	testCases := []struct {
		desc     string
		opts     []stackdriver.FilterOption
		expected []string
	}{
		{
			desc: "simple eq",
			opts: []stackdriver.FilterOption{
				stackdriver.Eq("property", "value"),
			},
			expected: []string{"property", "=", "value"},
		},
		{
			desc: "simple neq",
			opts: []stackdriver.FilterOption{
				stackdriver.Neq("property", "value"),
			},
			expected: []string{"property", "!=", "value"},
		},
		{
			desc: "simple regex",
			opts: []stackdriver.FilterOption{
				stackdriver.Matches("property", "regex"),
			},
			expected: []string{"property", "=~", "regex"},
		},
		{
			desc: "simple NOT regex",
			opts: []stackdriver.FilterOption{
				stackdriver.NotMatches("property", "regex"),
			},
			expected: []string{"property", "!=~", "regex"},
		},

		{
			desc: "simple AND",
			opts: []stackdriver.FilterOption{
				stackdriver.Eq("property", "value"),
				stackdriver.Neq("other-property", "other-value"),
			},
			expected: []string{"property", "=", "value", "AND", "other-property", "!=", "other-value"},
		},
		{
			desc: "multiple AND",
			opts: []stackdriver.FilterOption{
				stackdriver.Eq("property", "value"),
				stackdriver.Neq("other-property", "other-value"),
				stackdriver.Matches("last-property", "last-value"),
			},
			expected: []string{"property", "=", "value", "AND", "other-property", "!=", "other-value", "AND", "last-property", "=~", "last-value"},
		},
	}

	//nolint: scopelint
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			req := require.New(t)

			target := stackdriver.Delta("", stackdriver.Filter(test.opts...))

			req.Equal(test.expected, target.Builder.Filters)
		})
	}
}

func TestTargetSupportsGroupBys(t *testing.T) {
	req := require.New(t)

	target := stackdriver.Delta("", stackdriver.GroupBys("field", "other"))

	req.ElementsMatch(target.Builder.GroupBys, []string{"field", "other"})
}
