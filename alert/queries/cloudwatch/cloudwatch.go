package cloudwatch

import (
	"time"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a cloudwatch query.
type Option func(query *CloudWatch)

// CloudWatch represents a cloudwatch query.
type CloudWatch struct {
	Builder sdk.AlertQuery
}

// New creates a new cloudwatch query.
func New(ref string, query string, options ...Option) *CloudWatch {
	nope := false

	cloudwatch := &CloudWatch{
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
					Type: "cloudwatch",
				},
				Interval:   "",
				IntervalMs: 15000,
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(cloudwatch)
	}

	return cloudwatch
}

func defaults() []Option {
	return []Option{
		TimeRange(10*time.Minute, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(cloudwatch *CloudWatch) {
		cloudwatch.Builder.RelativeTimeRange.From = int(from.Seconds())
		cloudwatch.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(cloudwatch *CloudWatch) {
		cloudwatch.Builder.Model.LegendFormat = legend
	}
}
