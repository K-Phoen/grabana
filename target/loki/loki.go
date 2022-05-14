package loki

// Option represents an option that can be used to configure a loki query.
type Option func(target *Loki)

// Loki represents a loki query.
type Loki struct {
	Ref          string
	Hidden       bool
	Expr         string
	LegendFormat string
}

// New creates a new prometheus query.
func New(query string, options ...Option) *Loki {
	loki := &Loki{
		Expr: query,
	}

	for _, opt := range options {
		opt(loki)
	}

	return loki
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(loki *Loki) {
		loki.LegendFormat = legend
	}
}

// Hide the query. Grafana does not send hidden queries to the data source,
// but they can still be referenced in alerts.
func Hide() Option {
	return func(loki *Loki) {
		loki.Hidden = true
	}
}
