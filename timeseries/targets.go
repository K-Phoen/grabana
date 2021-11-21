package timeseries

import (
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
)

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(graph *TimeSeries) {
		graph.Builder.AddTarget(&sdk.Target{
			RefID:          target.Ref,
			Hide:           target.Hidden,
			Expr:           target.Expr,
			IntervalFactor: target.IntervalFactor,
			Interval:       target.Interval,
			Step:           target.Step,
			LegendFormat:   target.LegendFormat,
			Instant:        target.Instant,
			Format:         target.Format,
		})
	}
}

// WithGraphiteTarget adds a Graphite target to the table.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(graph *TimeSeries) {
		graph.Builder.AddTarget(target.Builder)
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(graph *TimeSeries) {
		graph.Builder.AddTarget(target.Builder)
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(graph *TimeSeries) {
		graph.Builder.AddTarget(target.Builder)
	}
}
