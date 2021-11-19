package graphite

import "github.com/K-Phoen/sdk"

// Option represents an option that can be used to configure a graphite query.
type Option func(target *Graphite)

// Graphite represents a graphite query.
type Graphite struct {
	Builder *sdk.Target
}

func New(query string, options ...Option) *Graphite {
	graphite := &Graphite{
		Builder: &sdk.Target{
			Target: query,
		},
	}

	for _, opt := range options {
		opt(graphite)
	}

	return graphite
}

// Ref sets the reference ID for this query.
func Ref(ref string) Option {
	return func(graphite *Graphite) {
		graphite.Builder.RefID = ref
	}
}

// Hide the query. Grafana does not send hidden queries to the data source,
// but they can still be referenced in alerts.
func Hide() Option {
	return func(graphite *Graphite) {
		graphite.Builder.Hide = true
	}
}
