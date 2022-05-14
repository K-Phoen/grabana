package influxdb

import (
	"time"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a influxdb query.
type Option func(query *InfluxDB)

// InfluxDB represents a influxdb query.
type InfluxDB struct {
	Builder sdk.AlertQuery
}

// New creates a new influxdb query.
func New(ref string, query string, options ...Option) *InfluxDB {
	nope := false

	influxdb := &InfluxDB{
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
					Type: "influxdb",
				},
				Interval:   "",
				IntervalMs: 15000,
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(influxdb)
	}

	return influxdb
}

func defaults() []Option {
	return []Option{
		TimeRange(10*time.Minute, 0),
	}
}

// TimeRange sets the legend format.
func TimeRange(from time.Duration, to time.Duration) Option {
	return func(influxdb *InfluxDB) {
		influxdb.Builder.RelativeTimeRange.From = int(from.Seconds())
		influxdb.Builder.RelativeTimeRange.To = int(to.Seconds())
	}
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(influxdb *InfluxDB) {
		influxdb.Builder.Model.LegendFormat = legend
	}
}
