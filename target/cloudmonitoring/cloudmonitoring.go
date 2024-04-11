package cloudmonitoring

import "github.com/K-Phoen/sdk"

// AutoAlignmentPeriod lets Google Clound Monitoring decide what it the ideal alignment period.
const AutoAlignmentPeriod = "cloud-monitoring-auto"

// PreprocessorMethod defines the available pre-processing options.
type PreprocessorMethod string

const (
	PreprocessNone  PreprocessorMethod = "none"
	PreprocessRate  PreprocessorMethod = "rate"
	PreprocessDelta PreprocessorMethod = "delta"
)

// Aligner specifies the operation that will be applied to the data points in
// each alignment period in a time series.
// See https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.alertPolicies#Aligner
type Aligner string

const (
	AlignNone          Aligner = "ALIGN_NONE"
	AlignDelta         Aligner = "ALIGN_DELTA"
	AlignRate          Aligner = "ALIGN_RATE"
	AlignInterpolate   Aligner = "ALIGN_INTERPOLATE"
	AlignNextOlder     Aligner = "ALIGN_NEXT_OLDER"
	AlignMin           Aligner = "ALIGN_MIN"
	AlignMax           Aligner = "ALIGN_MAX"
	AlignMean          Aligner = "ALIGN_MEAN"
	AlignCount         Aligner = "ALIGN_COUNT"
	AlignSum           Aligner = "ALIGN_SUM"
	AlignStdDev        Aligner = "ALIGN_STDDEV"
	AlignCountTrue     Aligner = "ALIGN_COUNT_TRUE"
	AlignCountFalse    Aligner = "ALIGN_COUNT_FALSE"
	AlignFractionTrue  Aligner = "ALIGN_FRACTION_TRUE"
	AlignPercentile99  Aligner = "ALIGN_PERCENTILE_99"
	AlignPercentile95  Aligner = "ALIGN_PERCENTILE_95"
	AlignPercentile50  Aligner = "ALIGN_PERCENTILE_50"
	AlignPercentile05  Aligner = "ALIGN_PERCENTILE_05"
	AlignPercentChange Aligner = "ALIGN_PERCENT_CHANGE"
)

// Reducer operation describes how to aggregate data points from multiple time
// series into a single time series, where the value of each data point in the
// resulting series is a function of all the already aligned values in the
// input time series.
// See https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.alertPolicies#reducer
type Reducer string

const (
	ReduceNone              Reducer = "REDUCE_NONE"
	ReduceMean              Reducer = "REDUCE_MEAN"
	ReduceMin               Reducer = "REDUCE_MIN"
	ReduceMax               Reducer = "REDUCE_MAX"
	ReduceSum               Reducer = "REDUCE_SUM"
	ReduceStdDev            Reducer = "REDUCE_STDDEV"
	ReduceCount             Reducer = "REDUCE_COUNT"
	ReduceCountTrue         Reducer = "REDUCE_COUNT_TRUE"
	ReduceCountFalse        Reducer = "REDUCE_COUNT_FALSE"
	ReduceCountFractionTrue Reducer = "REDUCE_FRACTION_TRUE"
	ReducePercentile99      Reducer = "REDUCE_PERCENTILE_99"
	ReducePercentile95      Reducer = "REDUCE_PERCENTILE_95"
	ReducePercentile50      Reducer = "REDUCE_PERCENTILE_50"
	ReducePercentile05      Reducer = "REDUCE_PERCENTILE_05"
)

// FilterOperator describes the set of all possible operations applicable to filters.
type FilterOperator string

const (
	FilterOperatorEqual            FilterOperator = "="
	FilterOperatorNotEqual         FilterOperator = "!="
	FilterOperatorMatchesRegexp    FilterOperator = "=~"
	FilterOperatorNotMatchesRegexp FilterOperator = "!=~"
)

// Target is an interface regrouping all the different cloudmonitoring targets.
type Target interface {
	Target() *sdk.Target
}

// Target allows to return an alert model for a specific builder.
type AlertModel interface {
	AlertModel() sdk.AlertModel
}
