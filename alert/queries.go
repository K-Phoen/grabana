package alert

import (
	"github.com/K-Phoen/grabana/alert/queries/prometheus"
)

// WithPrometheusQuery adds a prometheus query to the alert.
func WithPrometheusQuery(ref string, query string, options ...prometheus.Option) Option {
	return func(alert *Alert) {
		target := prometheus.New(ref, query, options...)
		alert.Builder.Rules[0].GrafanaAlert.Data = append(alert.Builder.Rules[0].GrafanaAlert.Data, target.Builder)
	}
}
