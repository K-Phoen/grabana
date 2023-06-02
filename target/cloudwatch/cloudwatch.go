package cloudwatch

import "github.com/K-Phoen/sdk"

// Option represents an option that can be used to configure a graphite query.
type Option func(target *Cloudwatch)

// Cloudwatch represents a cloudwatch query.
type Cloudwatch struct {
	Builder *sdk.Target
}

func New(dimensions map[string]string, statistics []string, namespace, metricName, period, region string, options ...Option) *Cloudwatch {
	cloudwatch := &Cloudwatch{
		Builder: &sdk.Target{
			Namespace:  namespace,
			MetricName: metricName,
			Dimensions: dimensions,
			Statistics: statistics,
			Period:     period,
			Region:     region,
		},
	}

	for _, opt := range options {
		opt(cloudwatch)
	}

	return cloudwatch
}

// Ref sets the reference ID for this query.
func Ref(ref string) Option {
	return func(cloudwatch *Cloudwatch) {
		cloudwatch.Builder.RefID = ref
	}
}

// Hide the query. Grafana does not send hidden queries to the data source,
// but they can still be referenced in alerts.
func Hide() Option {
	return func(cloudwatch *Cloudwatch) {
		cloudwatch.Builder.Hide = true
	}
}
