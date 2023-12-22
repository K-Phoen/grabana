package decoder

import (
	"fmt"
	"time"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/alert/queries/graphite"
	"github.com/K-Phoen/grabana/alert/queries/loki"
	"github.com/K-Phoen/grabana/alert/queries/prometheus"
	"github.com/K-Phoen/grabana/alert/queries/stackdriver"
)

var ErrMissingRef = fmt.Errorf("target ref missing")
var ErrInvalidLookback = fmt.Errorf("invalid lookback")

type AlertTarget struct {
	Prometheus  *AlertPrometheus  `yaml:",omitempty"`
	Loki        *AlertLoki        `yaml:",omitempty"`
	Graphite    *AlertGraphite    `yaml:",omitempty"`
	Stackdriver *AlertStackdriver `yaml:",omitempty"`
}

func (t AlertTarget) toOption() (alert.Option, error) {
	if t.Prometheus != nil {
		return t.Prometheus.toOptions()
	}
	if t.Loki != nil {
		return t.Loki.toOptions()
	}
	if t.Graphite != nil {
		return t.Graphite.toOptions()
	}
	if t.Stackdriver != nil {
		return t.Stackdriver.toOptions()
	}

	return nil, ErrTargetNotConfigured
}

type AlertPrometheus struct {
	Ref      string `yaml:",omitempty"`
	Query    string
	Legend   string `yaml:",omitempty"`
	Lookback string `yaml:",omitempty"`
}

func (t AlertPrometheus) toOptions() (alert.Option, error) {
	opts := []prometheus.Option{
		prometheus.Legend(t.Legend),
	}

	if t.Ref == "" {
		return nil, ErrMissingRef
	}

	if t.Lookback != "" {
		from, err := time.ParseDuration(t.Lookback)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", err, ErrInvalidLookback)
		}

		opts = append(opts, prometheus.TimeRange(from, 0))
	}

	return alert.WithPrometheusQuery(t.Ref, t.Query, opts...), nil
}

type AlertLoki struct {
	Ref      string `yaml:",omitempty"`
	Query    string
	Legend   string `yaml:",omitempty"`
	Lookback string `yaml:",omitempty"`
}

func (t AlertLoki) toOptions() (alert.Option, error) {
	opts := []loki.Option{
		loki.Legend(t.Legend),
	}

	if t.Ref == "" {
		return nil, ErrMissingRef
	}

	if t.Lookback != "" {
		from, err := time.ParseDuration(t.Lookback)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", err, ErrInvalidLookback)
		}

		opts = append(opts, loki.TimeRange(from, 0))
	}

	return alert.WithLokiQuery(t.Ref, t.Query, opts...), nil
}

type AlertGraphite struct {
	Ref      string `yaml:",omitempty"`
	Query    string
	Lookback string `yaml:",omitempty"`
}

func (t AlertGraphite) toOptions() (alert.Option, error) {
	var opts []graphite.Option

	if t.Ref == "" {
		return nil, ErrMissingRef
	}

	if t.Lookback != "" {
		from, err := time.ParseDuration(t.Lookback)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", err, ErrInvalidLookback)
		}

		opts = append(opts, graphite.TimeRange(from, 0))
	}

	return alert.WithGraphiteQuery(t.Ref, t.Query, opts...), nil
}

type AlertStackdriver struct {
	Ref      string `yaml:",omitempty"`
	Lookback string `yaml:",omitempty"`

	Project      string `yaml:",omitempty"`
	Type         string
	Metric       string
	Filters      StackdriverAlertFilters    `yaml:",omitempty"`
	Aggregation  string                     `yaml:",omitempty"`
	Alignment    *StackdriverAlertAlignment `yaml:",omitempty"`
	Legend       string                     `yaml:",omitempty"`
	Preprocessor string                     `yaml:",omitempty"`
	Hidden       bool                       `yaml:",omitempty"`
	GroupBy      []string                   `yaml:"group_by,omitempty"`
}

type StackdriverAlertFilters struct {
	Eq         map[string]string `yaml:",omitempty"`
	Neq        map[string]string `yaml:",omitempty"`
	Matches    map[string]string `yaml:",omitempty"`
	NotMatches map[string]string `yaml:"not_matches,omitempty"`
}

type StackdriverAlertAlignment struct {
	Method string
	Period string
}

func (t AlertStackdriver) toOptions() (alert.Option, error) {
	if t.Ref == "" {
		return nil, ErrMissingRef
	}

	opts, err := t.targetOptions()
	if err != nil {
		return nil, err
	}

	var query *stackdriver.Stackdriver

	switch t.Type {
	case "delta":
		query = stackdriver.Delta(t.Ref, t.Metric, opts...)
	case "gauge":
		query = stackdriver.Gauge(t.Ref, t.Metric, opts...)
	case "cumulative":
		query = stackdriver.Cumulative(t.Ref, t.Metric, opts...)
	default:
		return nil, ErrInvalidStackdriverType
	}

	return alert.WithStackdriverQuery(query), nil
}

func (t AlertStackdriver) targetOptions() ([]stackdriver.Option, error) {
	opts := []stackdriver.Option{
		stackdriver.Legend(t.Legend),
	}

	if t.Lookback != "" {
		from, err := time.ParseDuration(t.Lookback)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", err, ErrInvalidLookback)
		}

		opts = append(opts, stackdriver.TimeRange(from, 0))
	}

	filters := t.Filters.toOptions()
	if len(filters) != 0 {
		opts = append(opts, stackdriver.Filter(filters...))
	}

	if len(t.GroupBy) != 0 {
		opts = append(opts, stackdriver.GroupBys(t.GroupBy...))
	}

	if t.Project != "" {
		opts = append(opts, stackdriver.Project(t.Project))
	}

	if t.Aggregation != "" {
		opt, err := t.aggregation()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	if t.Preprocessor != "" {
		opt, err := t.preprocessor()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	if t.Alignment != nil {
		opt, err := t.Alignment.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (t AlertStackdriver) aggregation() (stackdriver.Option, error) {
	switch t.Aggregation {
	case "none":
		return stackdriver.Aggregation(stackdriver.ReduceNone), nil
	case "mean":
		return stackdriver.Aggregation(stackdriver.ReduceMean), nil
	case "min":
		return stackdriver.Aggregation(stackdriver.ReduceMin), nil
	case "max":
		return stackdriver.Aggregation(stackdriver.ReduceMax), nil
	case "sum":
		return stackdriver.Aggregation(stackdriver.ReduceSum), nil
	case "stddev":
		return stackdriver.Aggregation(stackdriver.ReduceStdDev), nil
	case "count":
		return stackdriver.Aggregation(stackdriver.ReduceCount), nil
	case "count_true":
		return stackdriver.Aggregation(stackdriver.ReduceCountTrue), nil
	case "count_false":
		return stackdriver.Aggregation(stackdriver.ReduceCountFalse), nil
	case "fraction_true":
		return stackdriver.Aggregation(stackdriver.ReduceCountFractionTrue), nil
	case "percentile_99":
		return stackdriver.Aggregation(stackdriver.ReducePercentile99), nil
	case "percentile_95":
		return stackdriver.Aggregation(stackdriver.ReducePercentile95), nil
	case "percentile_50":
		return stackdriver.Aggregation(stackdriver.ReducePercentile50), nil
	case "percentile_05":
		return stackdriver.Aggregation(stackdriver.ReducePercentile05), nil
	default:
		return nil, ErrInvalidStackdriverAggregation
	}
}

func (t AlertStackdriver) preprocessor() (stackdriver.Option, error) {
	switch t.Preprocessor {
	case "delta":
		return stackdriver.Preprocessor(stackdriver.PreprocessDelta), nil
	case "rate":
		return stackdriver.Preprocessor(stackdriver.PreprocessRate), nil
	default:
		return nil, ErrInvalidStackdriverPreprocessor
	}
}

func (filters StackdriverAlertFilters) toOptions() []stackdriver.FilterOption {
	opts := []stackdriver.FilterOption{}

	for key, value := range filters.Eq {
		opts = append(opts, stackdriver.Eq(key, value))
	}
	for key, value := range filters.Neq {
		opts = append(opts, stackdriver.Neq(key, value))
	}
	for key, value := range filters.Matches {
		opts = append(opts, stackdriver.Matches(key, value))
	}
	for key, value := range filters.NotMatches {
		opts = append(opts, stackdriver.NotMatches(key, value))
	}

	return opts
}

func (t StackdriverAlertAlignment) toOption() (stackdriver.Option, error) {
	switch t.Method {
	case "none":
		return stackdriver.Alignment(stackdriver.AlignNone, t.Period), nil
	case "delta":
		return stackdriver.Alignment(stackdriver.AlignDelta, t.Period), nil
	case "rate":
		return stackdriver.Alignment(stackdriver.AlignRate, t.Period), nil
	case "interpolate":
		return stackdriver.Alignment(stackdriver.AlignInterpolate, t.Period), nil
	case "next_older":
		return stackdriver.Alignment(stackdriver.AlignNextOlder, t.Period), nil
	case "min":
		return stackdriver.Alignment(stackdriver.AlignMin, t.Period), nil
	case "max":
		return stackdriver.Alignment(stackdriver.AlignMax, t.Period), nil
	case "mean":
		return stackdriver.Alignment(stackdriver.AlignMean, t.Period), nil
	case "count":
		return stackdriver.Alignment(stackdriver.AlignCount, t.Period), nil
	case "sum":
		return stackdriver.Alignment(stackdriver.AlignSum, t.Period), nil
	case "stddev":
		return stackdriver.Alignment(stackdriver.AlignStdDev, t.Period), nil
	case "count_true":
		return stackdriver.Alignment(stackdriver.AlignCountTrue, t.Period), nil
	case "count_false":
		return stackdriver.Alignment(stackdriver.AlignCountFalse, t.Period), nil
	case "fraction_true":
		return stackdriver.Alignment(stackdriver.AlignFractionTrue, t.Period), nil
	case "percentile_99":
		return stackdriver.Alignment(stackdriver.AlignPercentile99, t.Period), nil
	case "percentile_95":
		return stackdriver.Alignment(stackdriver.AlignPercentile95, t.Period), nil
	case "percentile_50":
		return stackdriver.Alignment(stackdriver.AlignPercentile50, t.Period), nil
	case "percentile_05":
		return stackdriver.Alignment(stackdriver.AlignPercentile05, t.Period), nil
	case "percent_change":
		return stackdriver.Alignment(stackdriver.AlignPercentChange, t.Period), nil
	default:
		return nil, ErrInvalidStackdriverAlignment
	}
}
