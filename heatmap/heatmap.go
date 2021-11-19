package heatmap

import (
	"github.com/K-Phoen/grabana/heatmap/axis"
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
)

// DataFormatMode represents the data format modes.
type DataFormatMode string

const (
	// Grafana does the bucketing by going through all time series values
	TimeSeriesBuckets DataFormatMode = "tsbuckets"

	// Each time series already represents a Y-Axis bucket.
	TimeSeries DataFormatMode = "timeseries"
)

// LegendOption allows to configure a legend.
type LegendOption uint16

const (
	// Hide keeps the legend from being displayed.
	Hide LegendOption = iota
)

// Option represents an option that can be used to configure a heatmap panel.
type Option func(stat *Heatmap)

// Heatmap represents a heatmap panel.
type Heatmap struct {
	Builder *sdk.Panel
}

// New creates a new heatmap panel.
func New(title string, options ...Option) *Heatmap {
	panel := &Heatmap{Builder: sdk.NewHeatmap(title)}
	panel.Builder.IsNew = false
	panel.Builder.HeatmapPanel.Cards = struct {
		CardPadding *float64 `json:"cardPadding"`
		CardRound   *float64 `json:"cardRound"`
	}{}
	panel.Builder.HeatmapPanel.Color = struct {
		CardColor   string   `json:"cardColor"`
		ColorScale  string   `json:"colorScale"`
		ColorScheme string   `json:"colorScheme"`
		Exponent    float64  `json:"exponent"`
		Min         *float64 `json:"min,omitempty"`
		Max         *float64 `json:"max,omitempty"`
		Mode        string   `json:"mode"`
	}{
		CardColor:   "#b4ff00",
		ColorScale:  "sqrt",
		ColorScheme: "interpolateSpectral",
		Exponent:    0.5,
		Mode:        "spectrum",
	}
	panel.Builder.HeatmapPanel.Legend = struct {
		Show bool `json:"show"`
	}{
		Show: true,
	}
	panel.Builder.HeatmapPanel.Tooltip = struct {
		Show          bool `json:"show"`
		ShowHistogram bool `json:"showHistogram"`
	}{
		Show:          true,
		ShowHistogram: true,
	}
	panel.Builder.HeatmapPanel.XAxis = struct {
		Show bool `json:"show"`
	}{
		Show: true,
	}
	panel.Builder.HeatmapPanel.YBucketBound = "auto"

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Span(6),
		DataFormat(TimeSeriesBuckets),
		HideZeroBuckets(),
		HighlightCards(),
		defaultYAxis(),
	}
}

func defaultYAxis() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.YAxis = *axis.New().Builder
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.Datasource = &source
	}
}

// DataFormat sets how the data should be interpreted.
func DataFormat(format DataFormatMode) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.DataFormat = string(format)
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(heatmap *Heatmap) {
		heatmap.Builder.AddTarget(&sdk.Target{
			RefID:          target.Ref,
			Expr:           target.Expr,
			IntervalFactor: target.IntervalFactor,
			Interval:       target.Interval,
			Step:           target.Step,
			LegendFormat:   target.LegendFormat,
			Instant:        target.Instant,
			Format:         target.Format,
		})
	}
}

// WithGraphiteTarget adds a Graphite target to the table.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(heatmap *Heatmap) {
		heatmap.Builder.AddTarget(target.Builder)
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(heatmap *Heatmap) {
		heatmap.Builder.AddTarget(target.Builder)
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.AddTarget(target.Builder)
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.Height = &height
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.Description = &content
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.Transparent = true
	}
}

// Legend defines what should be shown in the legend.
func Legend(opts ...LegendOption) Option {
	return func(heatmap *Heatmap) {
		for _, opt := range opts {
			if opt == Hide {
				heatmap.Builder.HeatmapPanel.Legend.Show = false
			}
		}
	}
}

// ShowZeroBuckets forces the display of "zero" buckets.
func ShowZeroBuckets() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.HideZeroBuckets = false
	}
}

// HideZeroBuckets hides "zero" buckets.
func HideZeroBuckets() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.HideZeroBuckets = true
	}
}

// HighlightCards highlights bucket cards.
func HighlightCards() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.HighlightCards = true
	}
}

// NoHighlightCards disables the highlighting of bucket cards.
func NoHighlightCards() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.HighlightCards = false
	}
}

// ReverseYBuckets reverses the order of bucket on the Y-axis.
func ReverseYBuckets() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.ReverseYBuckets = true
	}
}

// HideTooltip prevents the tooltip from being displayed.
func HideTooltip() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.Tooltip.Show = false
	}
}

// HideTooltipHistogram prevents the histograms from being displayed in tooltips.
// Histogram represents the distribution of the bucket values for the specific timestamp.
func HideTooltipHistogram() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.Tooltip.ShowHistogram = false
	}
}

// TooltipDecimals sets the number of decimals to be displayed in tooltips.
func TooltipDecimals(decimals int) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.TooltipDecimals = decimals
	}
}

// HideXAxis prevents the X-axis from being displayed.
func HideXAxis() Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.XAxis.Show = false
	}
}

// YAxis configures the Y axis.
func YAxis(opts ...axis.Option) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.HeatmapPanel.YAxis = *axis.New(opts...).Builder
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(heatmap *Heatmap) {
		heatmap.Builder.Repeat = &repeat
	}
}
