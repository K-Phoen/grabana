package alert

import (
	"github.com/K-Phoen/grabana/alert/queries/prometheus"
)

type QueryOption func(alert *Alert)

func Queries(queries ...QueryOption) Option {
	return func(alert *Alert) {
		for _, opt := range queries {
			opt(alert)
		}
	}
}

// WithPrometheusQuery adds a prometheus query to the alert.
func WithPrometheusQuery(ref string, query string, options ...prometheus.Option) QueryOption {
	target := prometheus.New(ref, query, options...)

	return func(alert *Alert) {
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, target.Builder)
	}
}
