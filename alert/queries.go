package alert

import (
	"github.com/K-Phoen/grabana/alert/queries/graphite"
	"github.com/K-Phoen/grabana/alert/queries/influxdb"
	"github.com/K-Phoen/grabana/alert/queries/loki"
	"github.com/K-Phoen/grabana/alert/queries/prometheus"
	"github.com/K-Phoen/grabana/alert/queries/stackdriver"
)

// WithPrometheusQuery adds a prometheus query to the alert.
func WithPrometheusQuery(ref string, query string, options ...prometheus.Option) Option {
	return func(alert *Alert) {
		target := prometheus.New(ref, query, options...)
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, target.Builder)
	}
}

// WithGraphiteQuery adds a graphite query to the alert.
func WithGraphiteQuery(ref string, query string, options ...graphite.Option) Option {
	return func(alert *Alert) {
		target := graphite.New(ref, query, options...)
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, target.Builder)
	}
}

// WithLokiQuery adds a loki query to the alert.
func WithLokiQuery(ref string, query string, options ...loki.Option) Option {
	return func(alert *Alert) {
		target := loki.New(ref, query, options...)
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, target.Builder)
	}
}

// WithStackdriverQuery adds a Stackdriver query to the alert.
func WithStackdriverQuery(query *stackdriver.Stackdriver) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, query.Builder)
	}
}

// WithInfluxDBQuery adds an InfluxDB query to the alert.
func WithInfluxDBQuery(ref string, query string, options ...influxdb.Option) Option {
	return func(alert *Alert) {
		target := influxdb.New(ref, query, options...)
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, target.Builder)
	}
}
