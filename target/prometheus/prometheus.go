package prometheus

type Option func(target *Prometheus)

type Prometheus struct {
	Ref            string
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
