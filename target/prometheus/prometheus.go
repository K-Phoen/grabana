package prometheus

type Option func(target *Prometheus)

type Prometheus struct {
	Expr           string
	IntervalFactor int
	Interval       string
	Step           int
	LegendFormat   string
	Instant        bool
	Format         string
}

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

// WithLegend sets the legend format.
func WithLegend(legend string) Option {
	return func(prometheus *Prometheus) {
		prometheus.LegendFormat = legend
	}
}
