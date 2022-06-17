package decoder

import (
	"testing"
	"time"

	alertBuilder "github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/alert/queries/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestDecodingAlertTargetFailsIfNoTargetIsProvided(t *testing.T) {
	target := AlertTarget{}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrTargetNotConfigured, err)
}

func TestDecodingAPrometheusTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Prometheus: &AlertPrometheus{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingAPrometheusTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Prometheus: &AlertPrometheus{
			Ref:      "A",
			Query:    "prom-query",
			Legend:   "{{ code }}",
			Lookback: "15m",
		},
	}

	opt, err := target.toOption()

	req.NoError(err)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 2) // the query and the condition

	promQuery := alert.Builder.Rules[0].GrafanaAlert.Data[1]

	req.Equal("A", promQuery.RefID)
	req.Equal("A", promQuery.Model.RefID)
	req.Equal("prom-query", promQuery.Model.Expr)
	req.Equal("{{ code }}", promQuery.Model.LegendFormat)
	req.Equal(int((15 * time.Minute).Seconds()), promQuery.RelativeTimeRange.From)
}

func TestDecodingAPrometheusTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Prometheus: &AlertPrometheus{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingALokiTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Loki: &AlertLoki{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingALokiTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Loki: &AlertLoki{
			Ref:      "A",
			Query:    "loki-query",
			Legend:   "{{ status }}",
			Lookback: "15m",
		},
	}

	opt, err := target.toOption()

	req.NoError(err)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 2) // the query and the condition

	promQuery := alert.Builder.Rules[0].GrafanaAlert.Data[1]

	req.Equal("A", promQuery.RefID)
	req.Equal("A", promQuery.Model.RefID)
	req.Equal("loki-query", promQuery.Model.Expr)
	req.Equal("{{ status }}", promQuery.Model.LegendFormat)
	req.Equal(int((15 * time.Minute).Seconds()), promQuery.RelativeTimeRange.From)
}

func TestDecodingALokiTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Loki: &AlertLoki{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingAGraphiteTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Graphite: &AlertGraphite{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingAGraphiteTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Graphite: &AlertGraphite{
			Ref:      "A",
			Query:    "graphite-query",
			Lookback: "15m",
		},
	}

	opt, err := target.toOption()

	req.NoError(err)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 2) // the query and the condition

	promQuery := alert.Builder.Rules[0].GrafanaAlert.Data[1]

	req.Equal("A", promQuery.RefID)
	req.Equal("A", promQuery.Model.RefID)
	req.Equal("graphite-query", promQuery.Model.Expr)
	req.Equal(int((15 * time.Minute).Seconds()), promQuery.RelativeTimeRange.From)
}

func TestDecodingAGraphiteTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Graphite: &AlertGraphite{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingAStackdriverTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Stackdriver: &AlertStackdriver{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingAStackdriverTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Stackdriver: &AlertStackdriver{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingAStackdriverTargetWithInvalidType(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Stackdriver: &AlertStackdriver{
			Ref:  "A",
			Type: "invalid",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidStackdriverType)
}

func TestDecodingStackdriverTargetWithValidTypes(t *testing.T) {
	validTypes := []string{"delta", "gauge", "cumulative"}

	for _, inputType := range validTypes {
		t.Run(inputType, func(t *testing.T) {
			req := require.New(t)

			target := AlertTarget{
				Stackdriver: &AlertStackdriver{
					Ref:  "A",
					Type: inputType,
				},
			}

			_, err := target.toOption()

			req.NoError(err)
		})
	}
}

func TestDecodingStackdriverTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Stackdriver: &AlertStackdriver{
			Ref:         "A",
			Lookback:    "15m",
			Type:        "gauge",
			Metric:      "cloudsql.googleapis.com/database/cpu/utilization",
			GroupBy:     []string{"resource.label.database_id"},
			Aggregation: "mean",
			Alignment: &StackdriverAlertAlignment{
				Method: "mean",
				Period: "stackdriver-auto",
			},
		},
	}

	opt, err := target.toOption()

	req.NoError(err, ErrInvalidStackdriverType)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 2) // the query and the condition

	query := alert.Builder.Rules[0].GrafanaAlert.Data[1]
	stackdriverData := query.Model.MetricQuery

	req.Equal("A", query.RefID)
	req.Equal("A", query.Model.RefID)
	req.Equal("GAUGE", stackdriverData.MetricKind)
	req.Equal("cloudsql.googleapis.com/database/cpu/utilization", stackdriverData.MetricType)
	req.ElementsMatch([]string{"resource.label.database_id"}, stackdriverData.GroupBys)
	req.Equal("REDUCE_MEAN", stackdriverData.CrossSeriesReducer)
	req.Equal("ALIGN_MEAN", stackdriverData.PerSeriesAligner)
	req.Equal("stackdriver-auto", stackdriverData.AlignmentPeriod)
	req.Equal(int((15 * time.Minute).Seconds()), query.RelativeTimeRange.From)
}

func TestDecodingStackdriverPreprocessor(t *testing.T) {
	testCases := []struct {
		input         string
		expected      string
		expectedError error
	}{
		{
			input:         "delta",
			expected:      "delta",
			expectedError: nil,
		},
		{
			input:         "rate",
			expected:      "rate",
			expectedError: nil,
		},
		{
			input:         "invalid",
			expected:      "",
			expectedError: ErrInvalidStackdriverPreprocessor,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			target := AlertTarget{
				Stackdriver: &AlertStackdriver{
					Ref:          "A",
					Type:         "delta",
					Preprocessor: tc.input,
				},
			}

			opt, err := target.toOption()

			if tc.expectedError != nil {
				req.ErrorIs(tc.expectedError, err)
				return
			}

			req.NoError(err)

			alert := alertBuilder.New("", opt)
			query := alert.Builder.Rules[0].GrafanaAlert.Data[1]

			req.Equal(tc.expected, query.Model.MetricQuery.Preprocessor)
		})
	}
}

func TestDecodingStackdriverAggregation(t *testing.T) {
	testCases := []struct {
		input         string
		expected      stackdriver.Reducer
		expectedError error
	}{
		{
			input:         "none",
			expected:      stackdriver.ReduceNone,
			expectedError: nil,
		},
		{
			input:         "mean",
			expected:      stackdriver.ReduceMean,
			expectedError: nil,
		},
		{
			input:         "min",
			expected:      stackdriver.ReduceMin,
			expectedError: nil,
		},
		{
			input:         "max",
			expected:      stackdriver.ReduceMax,
			expectedError: nil,
		},
		{
			input:         "sum",
			expected:      stackdriver.ReduceSum,
			expectedError: nil,
		},
		{
			input:         "stddev",
			expected:      stackdriver.ReduceStdDev,
			expectedError: nil,
		},
		{
			input:         "count",
			expected:      stackdriver.ReduceCount,
			expectedError: nil,
		},
		{
			input:         "count_true",
			expected:      stackdriver.ReduceCountTrue,
			expectedError: nil,
		},
		{
			input:         "count_false",
			expected:      stackdriver.ReduceCountFalse,
			expectedError: nil,
		},
		{
			input:         "fraction_true",
			expected:      stackdriver.ReduceCountFractionTrue,
			expectedError: nil,
		},
		{
			input:         "percentile_99",
			expected:      stackdriver.ReducePercentile99,
			expectedError: nil,
		},
		{
			input:         "percentile_95",
			expected:      stackdriver.ReducePercentile95,
			expectedError: nil,
		},
		{
			input:         "percentile_50",
			expected:      stackdriver.ReducePercentile50,
			expectedError: nil,
		},
		{
			input:         "percentile_05",
			expected:      stackdriver.ReducePercentile05,
			expectedError: nil,
		},
		{
			input:         "invalid",
			expected:      "",
			expectedError: ErrInvalidStackdriverAggregation,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			target := AlertTarget{
				Stackdriver: &AlertStackdriver{
					Ref:         "A",
					Type:        "delta",
					Aggregation: tc.input,
				},
			}

			opt, err := target.toOption()

			if tc.expectedError != nil {
				req.ErrorIs(tc.expectedError, err)
				return
			}

			req.NoError(err)

			alert := alertBuilder.New("", opt)
			query := alert.Builder.Rules[0].GrafanaAlert.Data[1]

			req.Equal(string(tc.expected), query.Model.MetricQuery.CrossSeriesReducer)
		})
	}
}

func TestDecodingStackdriverAlignment(t *testing.T) {
	testCases := []struct {
		input         string
		expected      stackdriver.Aligner
		expectedError error
	}{
		{
			input:         "none",
			expected:      stackdriver.AlignNone,
			expectedError: nil,
		},
		{
			input:         "delta",
			expected:      stackdriver.AlignDelta,
			expectedError: nil,
		},
		{
			input:         "rate",
			expected:      stackdriver.AlignRate,
			expectedError: nil,
		},
		{
			input:         "interpolate",
			expected:      stackdriver.AlignInterpolate,
			expectedError: nil,
		},
		{
			input:         "next_older",
			expected:      stackdriver.AlignNextOlder,
			expectedError: nil,
		},
		{
			input:         "min",
			expected:      stackdriver.AlignMin,
			expectedError: nil,
		},
		{
			input:         "max",
			expected:      stackdriver.AlignMax,
			expectedError: nil,
		},
		{
			input:         "mean",
			expected:      stackdriver.AlignMean,
			expectedError: nil,
		},
		{
			input:         "count",
			expected:      stackdriver.AlignCount,
			expectedError: nil,
		},
		{
			input:         "sum",
			expected:      stackdriver.AlignSum,
			expectedError: nil,
		},
		{
			input:         "stddev",
			expected:      stackdriver.AlignStdDev,
			expectedError: nil,
		},
		{
			input:         "count_true",
			expected:      stackdriver.AlignCountTrue,
			expectedError: nil,
		},
		{
			input:         "count_false",
			expected:      stackdriver.AlignCountFalse,
			expectedError: nil,
		},
		{
			input:         "fraction_true",
			expected:      stackdriver.AlignFractionTrue,
			expectedError: nil,
		},
		{
			input:         "percentile_99",
			expected:      stackdriver.AlignPercentile99,
			expectedError: nil,
		},
		{
			input:         "percentile_95",
			expected:      stackdriver.AlignPercentile95,
			expectedError: nil,
		},
		{
			input:         "percentile_50",
			expected:      stackdriver.AlignPercentile50,
			expectedError: nil,
		},
		{
			input:         "percentile_05",
			expected:      stackdriver.AlignPercentile05,
			expectedError: nil,
		},
		{
			input:         "percent_change",
			expected:      stackdriver.AlignPercentChange,
			expectedError: nil,
		},
		{
			input:         "invalid",
			expected:      "",
			expectedError: ErrInvalidStackdriverAlignment,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			target := AlertTarget{
				Stackdriver: &AlertStackdriver{
					Ref:  "A",
					Type: "delta",
					Alignment: &StackdriverAlertAlignment{
						Method: tc.input,
						Period: "stackdriver-auto",
					},
				},
			}

			opt, err := target.toOption()

			if tc.expectedError != nil {
				req.ErrorIs(tc.expectedError, err)
				return
			}

			req.NoError(err)

			alert := alertBuilder.New("", opt)
			query := alert.Builder.Rules[0].GrafanaAlert.Data[1]

			req.Equal(string(tc.expected), query.Model.MetricQuery.PerSeriesAligner)
			req.Equal("stackdriver-auto", query.Model.MetricQuery.AlignmentPeriod)
		})
	}
}
