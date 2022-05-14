package loki

import (
	"time"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a loki query.
type Option func(query *Loki)

// Loki represents a loki query.
type Loki struct {
	Builder sdk.AlertQuery
}

// New creates a new loki query.
func New(ref string, query string, options ...Option) *Loki {
	nope := false

	loki := &Loki{
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
					Type: "loki",
				},
				Interval:   "",
				IntervalMs: 15000,
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(loki)
	}

	return loki
}

func defaults() []Option {
	return []Option{
		TimeRange(10*time.Minute, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(loki *Loki) {
		loki.Builder.RelativeTimeRange.From = int(from.Seconds())
		loki.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(loki *Loki) {
		loki.Builder.Model.LegendFormat = legend
	}
}
