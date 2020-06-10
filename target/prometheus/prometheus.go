package prometheus

// Option represents an option that can be used to configure a prometheus query.
type Option func(target *Prometheus)

// Prometheus represents a prometheus query.
type Prometheus struct {
	Ref            string
	Hidden         bool
	Expr           string
	IntervalFactor int
	Interval       string
	Step           int
	LegendFormat   string
	Instant        bool
	Format         string
}

// New creates a new prometheus query.
func New(query string, options ...Option) *Prometheus {
	prometheus := &Prometheus{
		Expr:   query,
		Format: "time_series",
	}

	for _, opt := range options {
		opt(prometheus)
	}

	return prometheus
}

// Legend sets the legend format.
func Legend(legend string) Option {
	return func(prometheus *Prometheus) {
		prometheus.LegendFormat = legend
	}
}

// Ref sets the reference ID for this query.
func Ref(ref string) Option {
	return func(prometheus *Prometheus) {
		prometheus.Ref = ref
	}
}

// Hide the query. Grafana does not send hidden queries to the data source,
// but they can still be referenced in alerts.
func Hide() Option {
	return func(prometheus *Prometheus) {
		prometheus.Hidden = true
	}
}
