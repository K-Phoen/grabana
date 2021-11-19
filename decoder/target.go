package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
)

var ErrTargetNotConfigured = fmt.Errorf("target not configured")
var ErrInvalidStackdriverType = fmt.Errorf("invalid stackdriver target type")
var ErrInvalidStackdriverAggregation = fmt.Errorf("invalid stackdriver aggregation type")
var ErrInvalidStackdriverAlignment = fmt.Errorf("invalid stackdriver alignment method")

type Target struct {
	Prometheus  *PrometheusTarget  `yaml:",omitempty"`
	Graphite    *GraphiteTarget    `yaml:",omitempty"`
	InfluxDB    *InfluxDBTarget    `yaml:"influxdb,omitempty"`
	Stackdriver *StackdriverTarget `yaml:",omitempty"`
}

type PrometheusTarget struct {
	Query          string
	Legend         string `yaml:",omitempty"`
	Ref            string `yaml:",omitempty"`
	Hidden         bool   `yaml:",omitempty"`
	Format         string `yaml:",omitempty"`
	Instant        bool   `yaml:",omitempty"`
	IntervalFactor *int   `yaml:"interval_factor,omitempty"`
}

func (t PrometheusTarget) toOptions() []prometheus.Option {
	opts := []prometheus.Option{
		prometheus.Legend(t.Legend),
		prometheus.Ref(t.Ref),
	}

	if t.Hidden {
		opts = append(opts, prometheus.Hide())
	}
	if t.Instant {
		opts = append(opts, prometheus.Instant())
	}
	if t.IntervalFactor != nil {
		opts = append(opts, prometheus.IntervalFactor(*t.IntervalFactor))
	}
	if t.Format != "" {
		switch t.Format {
		case "heatmap":
			opts = append(opts, prometheus.Format(prometheus.FormatHeatmap))
		case "table":
			opts = append(opts, prometheus.Format(prometheus.FormatTable))
		case "time_series":
			opts = append(opts, prometheus.Format(prometheus.FormatTimeSeries))
		}
	}

	return opts
}

type GraphiteTarget struct {
	Query  string
	Ref    string `yaml:",omitempty"`
	Hidden bool   `yaml:",omitempty"`
}

func (t GraphiteTarget) toOptions() []graphite.Option {
	opts := []graphite.Option{
		graphite.Ref(t.Ref),
	}

	if t.Hidden {
		opts = append(opts, graphite.Hide())
	}

	return opts
}

type InfluxDBTarget struct {
	Query  string
	Ref    string `yaml:",omitempty"`
	Hidden bool   `yaml:",omitempty"`
}

func (t InfluxDBTarget) toOptions() []influxdb.Option {
	opts := []influxdb.Option{
		influxdb.Ref(t.Ref),
	}

	if t.Hidden {
		opts = append(opts, influxdb.Hide())
	}

	return opts
}

type StackdriverTarget struct {
	Project     string
	Type        string
	Metric      string
	Filters     StackdriverFilters    `yaml:",omitempty"`
	Aggregation string                `yaml:",omitempty"`
	Alignment   *StackdriverAlignment `yaml:",omitempty"`
	Legend      string                `yaml:",omitempty"`
	Ref         string                `yaml:",omitempty"`
	Hidden      bool                  `yaml:",omitempty"`
	GroupBy     []string              `yaml:"group_by,omitempty"`
}

type StackdriverFilters struct {
	Eq         map[string]string `yaml:",omitempty"`
	Neq        map[string]string `yaml:",omitempty"`
	Matches    map[string]string `yaml:",omitempty"`
	NotMatches map[string]string `yaml:"not_matches,omitempty"`
}

type StackdriverAlignment struct {
	Method string
	Period string
}

func (t StackdriverTarget) toTarget() (*stackdriver.Stackdriver, error) {
	opts, err := t.toOptions()
	if err != nil {
		return nil, err
	}

	switch t.Type {
	case "delta":
		return stackdriver.Delta(t.Metric, opts...), nil
	case "gauge":
		return stackdriver.Gauge(t.Metric, opts...), nil
	case "cumulative":
		return stackdriver.Cumulative(t.Metric, opts...), nil
	}

	return nil, ErrInvalidStackdriverType
}

func (t StackdriverTarget) toOptions() ([]stackdriver.Option, error) {
	opts := []stackdriver.Option{
		stackdriver.Legend(t.Legend),
		stackdriver.Ref(t.Ref),
	}

	if t.Hidden {
		opts = append(opts, stackdriver.Hide())
	}

	if t.Project != "" {
		opts = append(opts, stackdriver.Project(t.Project))
	}

	filters := t.Filters.toOptions()
	if len(filters) != 0 {
		opts = append(opts, stackdriver.Filter(filters...))
	}

	if len(t.GroupBy) != 0 {
		opts = append(opts, stackdriver.GroupBys(t.GroupBy...))
	}

	if t.Aggregation != "" {
		opt, err := t.aggregation()
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

func (t StackdriverTarget) aggregation() (stackdriver.Option, error) {
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

func (filters StackdriverFilters) toOptions() []stackdriver.FilterOption {
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

func (t StackdriverAlignment) toOption() (stackdriver.Option, error) {
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
