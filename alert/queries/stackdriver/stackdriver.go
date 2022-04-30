package stackdriver

import (
	"time"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a stackdriver query.
type Option func(query *Stackdriver)

// Stackdriver represents a stackdriver query.
type Stackdriver struct {
	Builder sdk.AlertQuery
}

// New creates a new stackdriver query.
func New(ref string, query string, options ...Option) *Stackdriver {
	nope := false

	stackdriver := &Stackdriver{
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
					Type: "stackdriver",
				},
				Interval:   "",
				IntervalMs: 15000,
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(stackdriver)
	}

	return stackdriver
}

func defaults() []Option {
	return []Option{
		TimeRange(10*time.Minute, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.RelativeTimeRange.From = int(from.Seconds())
		stackdriver.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(stackdriver *Stackdriver) {
		stackdriver.Builder.Model.LegendFormat = legend
	}
}
