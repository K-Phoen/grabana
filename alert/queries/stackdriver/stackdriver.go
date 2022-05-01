package stackdriver

import (
	"time"

	"github.com/K-Phoen/sdk"
)

const AlignmentStackdriverAuto = "stackdriver-auto"
const AlignmentGrafanaAuto = "grafana-auto"

// PreprocessorMethod defines the available pre-processing options.
type PreprocessorMethod string

const PreprocessRate = "rate"
const PreprocessDelta = "delta"

// Aligner specifies the operation that will be applied to the data points in
// each alignment period in a time series.
// See https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.alertPolicies#Aligner
type Aligner string

const AlignNone Aligner = "ALIGN_NONE"
const AlignDelta Aligner = "ALIGN_DELTA"
const AlignRate Aligner = "ALIGN_RATE"
const AlignInterpolate Aligner = "ALIGN_INTERPOLATE"
const AlignNextOlder Aligner = "ALIGN_NEXT_OLDER"
const AlignMin Aligner = "ALIGN_MIN"
const AlignMax Aligner = "ALIGN_MAX"
const AlignMean Aligner = "ALIGN_MEAN"
const AlignCount Aligner = "ALIGN_COUNT"
const AlignSum Aligner = "ALIGN_SUM"
const AlignStdDev Aligner = "ALIGN_STDDEV"
const AlignCountTrue Aligner = "ALIGN_COUNT_TRUE"
const AlignCountFalse Aligner = "ALIGN_COUNT_FALSE"
const AlignFractionTrue Aligner = "ALIGN_FRACTION_TRUE"
const AlignPercentile99 Aligner = "ALIGN_PERCENTILE_99"
const AlignPercentile95 Aligner = "ALIGN_PERCENTILE_95"
const AlignPercentile50 Aligner = "ALIGN_PERCENTILE_50"
const AlignPercentile05 Aligner = "ALIGN_PERCENTILE_05"
const AlignPercentChange Aligner = "ALIGN_PERCENT_CHANGE"

// Reducer operation describes how to aggregate data points from multiple time
// series into a single time series, where the value of each data point in the
// resulting series is a function of all the already aligned values in the
// input time series.
// See https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.alertPolicies#reducer
type Reducer string

const ReduceNone Reducer = "REDUCE_NONE"
const ReduceMean Reducer = "REDUCE_MEAN"
const ReduceMin Reducer = "REDUCE_MIN"
const ReduceMax Reducer = "REDUCE_MAX"
const ReduceSum Reducer = "REDUCE_SUM"
const ReduceStdDev Reducer = "REDUCE_STDDEV"
const ReduceCount Reducer = "REDUCE_COUNT"
const ReduceCountTrue Reducer = "REDUCE_COUNT_TRUE"
const ReduceCountFalse Reducer = "REDUCE_COUNT_FALSE"
const ReduceCountFractionTrue Reducer = "REDUCE_FRACTION_TRUE"
const ReducePercentile99 Reducer = "REDUCE_PERCENTILE_99"
const ReducePercentile95 Reducer = "REDUCE_PERCENTILE_95"
const ReducePercentile50 Reducer = "REDUCE_PERCENTILE_50"
const ReducePercentile05 Reducer = "REDUCE_PERCENTILE_05"

// Option represents an option that can be used to configure a stackdriver query.
type Option func(query *Stackdriver)

// Stackdriver represents a stackdriver query.
type Stackdriver struct {
	Builder sdk.AlertQuery
}

// Delta represents the change in a value during a time interval.
func Delta(ref string, metricType string, options ...Option) *Stackdriver {
	return newMetric(ref, "DELTA", metricType, options...)
}

// Gauge represents an instantaneous measurement of a value.
func Gauge(ref string, metricType string, options ...Option) *Stackdriver {
	return newMetric(ref, "GAUGE", metricType, options...)
}

// Cumulative represents a value accumulated over a time interval. Cumulative
// measurements in a time series should have the same start time and
// increasing end times, until an event resets the cumulative value to zero
// and sets a new start time for the following points.
func Cumulative(ref string, metricType string, options ...Option) *Stackdriver {
	return newMetric(ref, "CUMULATIVE", metricType, options...)
}

// New creates a new stackdriver query.
func newMetric(ref string, metricKind string, metricType string, options ...Option) *Stackdriver {
	nope := false

	stackdriver := &Stackdriver{
		Builder: sdk.AlertQuery{
			RefID:             ref,
			QueryType:         "",
			DatasourceUID:     "__FILL_ME__",
			RelativeTimeRange: &sdk.AlertRelativeTimeRange{},
			Model: sdk.AlertModel{
				RefID:     ref,
				Format:    "time_series",
				QueryType: "metrics",
				Hide:      &nope,
				Datasource: sdk.AlertDatasourceRef{
					UID:  "__FILL_ME__",
					Type: "stackdriver",
				},
				Interval:   "",
				IntervalMs: 15000,

				MetricQuery: &sdk.StackdriverAlertQuery{
					MetricKind: metricKind,
					MetricType: metricType,
				},
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(stackdriver)
	}

	return stackdriver
}

func defaults() []Option {
	return []Option{
		TimeRange(10*time.Minute, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.RelativeTimeRange.From = int(from.Seconds())
		stackdriver.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.Model.MetricQuery.AliasBy = legend
	}
}

// Aggregation defines how the time series will be aggregated.
func Aggregation(reducer Reducer) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.Model.MetricQuery.CrossSeriesReducer = string(reducer)
	}
}

// Alignment defines how the time series will be aligned.
func Alignment(aligner Aligner, alignmentPeriod string) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.Model.MetricQuery.AlignmentPeriod = alignmentPeriod
		stackdriver.Builder.Model.MetricQuery.PerSeriesAligner = string(aligner)
	}
}

// Preprocessor defines how the time series should be pre-processed.
func Preprocessor(preprocessor PreprocessorMethod) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.Model.MetricQuery.Preprocessor = string(preprocessor)
	}
}

// Filter allows to specify which time series will be in the results.
func Filter(filters ...FilterOption) Option {
	return func(stackdriver *Stackdriver) {
		for i, filterOpt := range filters {
			f := &filter{}
			filterOpt(f)

			if i != 0 || len(stackdriver.Builder.Model.MetricQuery.Filters) != 0 {
				stackdriver.Builder.Model.MetricQuery.Filters = append(stackdriver.Builder.Model.MetricQuery.Filters, "AND")
			}

			stackdriver.Builder.Model.MetricQuery.Filters = append(stackdriver.Builder.Model.MetricQuery.Filters, f.leftOperand, f.operator, f.rightOperand)
		}
	}
}

// GroupBys defines a list of fields to group the query by.
func GroupBys(groupBys ...string) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.Model.MetricQuery.GroupBys = groupBys
	}
}
