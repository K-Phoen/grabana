package heatmap

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/heatmap/axis"
	"github.com/K-Phoen/grabana/links"
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
type Option func(stat *Heatmap) error

// Heatmap represents a heatmap panel.
type Heatmap struct {
	Builder *sdk.Panel
}

// New creates a new heatmap panel.
func New(title string, options ...Option) (*Heatmap, error) {
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
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
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
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.YAxis = *axis.New().Builder

		return nil
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			heatmap.Builder.Links = append(heatmap.Builder.Links, link.Builder)
		}

		return nil
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// DataFormat sets how the data should be interpreted.
func DataFormat(format DataFormatMode) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.DataFormat = string(format)

		return nil
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(heatmap *Heatmap) error {
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

		return nil
	}
}

// WithGraphiteTarget adds a Graphite target to the table.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(heatmap *Heatmap) error {
		heatmap.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(heatmap *Heatmap) error {
		heatmap.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.AddTarget(target.Builder)

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(heatmap *Heatmap) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		heatmap.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.Transparent = true

		return nil
	}
}

// Legend defines what should be shown in the legend.
func Legend(opts ...LegendOption) Option {
	return func(heatmap *Heatmap) error {
		for _, opt := range opts {
			if opt == Hide {
				heatmap.Builder.HeatmapPanel.Legend.Show = false
			}
		}

		return nil
	}
}

// ShowZeroBuckets forces the display of "zero" buckets.
func ShowZeroBuckets() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.HideZeroBuckets = false

		return nil
	}
}

// HideZeroBuckets hides "zero" buckets.
func HideZeroBuckets() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.HideZeroBuckets = true

		return nil
	}
}

// HighlightCards highlights bucket cards.
func HighlightCards() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.HighlightCards = true

		return nil
	}
}

// NoHighlightCards disables the highlighting of bucket cards.
func NoHighlightCards() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.HighlightCards = false

		return nil
	}
}

// ReverseYBuckets reverses the order of bucket on the Y-axis.
func ReverseYBuckets() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.ReverseYBuckets = true

		return nil
	}
}

// HideTooltip prevents the tooltip from being displayed.
func HideTooltip() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.Tooltip.Show = false

		return nil
	}
}

// HideTooltipHistogram prevents the histograms from being displayed in tooltips.
// Histogram represents the distribution of the bucket values for the specific timestamp.
func HideTooltipHistogram() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.Tooltip.ShowHistogram = false

		return nil
	}
}

// TooltipDecimals sets the number of decimals to be displayed in tooltips.
func TooltipDecimals(decimals int) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.TooltipDecimals = decimals

		return nil
	}
}

// HideXAxis prevents the X-axis from being displayed.
func HideXAxis() Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.XAxis.Show = false

		return nil
	}
}

// YAxis configures the Y axis.
func YAxis(opts ...axis.Option) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.HeatmapPanel.YAxis = *axis.New(opts...).Builder

		return nil
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(heatmap *Heatmap) error {
		heatmap.Builder.Repeat = &repeat

		return nil
	}
}
