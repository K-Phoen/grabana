package influxdb

import "github.com/K-Phoen/sdk"

// Option represents an option that can be used to configure a influxdb query.
type Option func(target *InfluxDB)

// InfluxDB represents a influxdb query.
type InfluxDB struct {
	Builder *sdk.Target
}

func New(query string, options ...Option) *InfluxDB {
	influxdb := &InfluxDB{
		Builder: &sdk.Target{
			Measurement: query,
		},
	}

	for _, opt := range options {
		opt(influxdb)
	}

	return influxdb
}

// Ref sets the reference ID for this query.
func Ref(ref string) Option {
	return func(influxdb *InfluxDB) {
		influxdb.Builder.RefID = ref
	}
}

// Hide the query. Grafana does not send hidden queries to the data source,
// but they can still be referenced in alerts.
func Hide() Option {
	return func(influxdb *InfluxDB) {
		influxdb.Builder.Hide = true
	}
}
