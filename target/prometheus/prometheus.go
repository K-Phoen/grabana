package prometheus

// Format switches between Table, Time series, or Heatmap. Table will only work
// in the Table panel. Heatmap is suitable for displaying metrics of the
// Histogram type on a Heatmap panel. Under the hood, it converts cumulative
// histograms to regular ones and sorts series by the bucket bound.
type FormatMode string

const (
	FormatTable      FormatMode = "table"
	FormatHeatmap    FormatMode = "heatmap"
	FormatTimeSeries FormatMode = "time_series"
)

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
		Format: string(FormatTimeSeries),
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

// Instant marks the query as "instant, which means Prometheus will only return the latest scrapped value.
func Instant() Option {
	return func(prometheus *Prometheus) {
		prometheus.Instant = true
	}
}

// Format indicates how the data should be returned.
func Format(format FormatMode) Option {
	return func(prometheus *Prometheus) {
		prometheus.Format = string(format)
	}
}

// IntervalFactor sets the resolution factor.
func IntervalFactor(factor int) Option {
	return func(prometheus *Prometheus) {
		prometheus.IntervalFactor = factor
	}
}
