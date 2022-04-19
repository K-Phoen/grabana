package prometheus

import (
	"time"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a prometheus query.
type Option func(query *Prometheus)

// Prometheus represents a prometheus query.
type Prometheus struct {
	Builder sdk.AlertQuery
}

// New creates a new prometheus query.
func New(ref string, query string, options ...Option) *Prometheus {
	nope := false

	prometheus := &Prometheus{
		Builder: sdk.AlertQuery{
			RefID:             ref,
			QueryType:         "",
			DatasourceUID:     "__FILL_ME__",
			RelativeTimeRange: &sdk.AlertRelativeTimeRange{},
			Model: sdk.AlertModel{
				RefID:  ref,
				Expr:   query,
				Format: "time_series",
				Hide:   &nope,
				Datasource: sdk.AlertDatasourceRef{
					UID:  "__FILL_ME__",
					Type: "prometheus",
				},
				Interval:   "",
				IntervalMs: 15000,
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(prometheus)
	}

	return prometheus
}

func defaults() []Option {
	return []Option{
		TimeRange(1*time.Hour, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(prometheus *Prometheus) {
		prometheus.Builder.RelativeTimeRange.From = int(from.Seconds())
		prometheus.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(prometheus *Prometheus) {
		prometheus.Builder.Model.LegendFormat = legend
	}
}
