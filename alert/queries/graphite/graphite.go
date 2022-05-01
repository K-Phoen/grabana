package graphite

import (
	"time"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a graphite query.
type Option func(query *Graphite)

// Graphite represents a graphite query.
type Graphite struct {
	Builder sdk.AlertQuery
}

// New creates a new graphite query.
func New(ref string, query string, options ...Option) *Graphite {
	nope := false

	graphite := &Graphite{
		Builder: sdk.AlertQuery{
			RefID:             ref,
			QueryType:         "",
			DatasourceUID:     "__FILL_ME__",
			RelativeTimeRange: &sdk.AlertRelativeTimeRange{},
			Model: sdk.AlertModel{
				RefID:  ref,
				Expr:   query,
				Target: query,
				Format: "time_series",
				Hide:   &nope,
				Datasource: sdk.AlertDatasourceRef{
					UID:  "__FILL_ME__",
					Type: "graphite",
				},
				Interval:   "",
				IntervalMs: 15000,
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(graphite)
	}

	return graphite
}

func defaults() []Option {
	return []Option{
		TimeRange(10*time.Minute, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(graphite *Graphite) {
		graphite.Builder.RelativeTimeRange.From = int(from.Seconds())
		graphite.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(graphite *Graphite) {
		graphite.Builder.Model.LegendFormat = legend
	}
}
