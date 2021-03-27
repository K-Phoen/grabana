package heatmap

import (
	"testing"

	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewHeatmapPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Heatmap panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Heatmap panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestHeatmapPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithInfluxDBTarget("buckets()"))

	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestHeatmapPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestHeatmapPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestHeatmapPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestHeatmapPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestDataFormatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataFormat(TimeSeriesBuckets))

	req.Equal("tsbuckets", panel.Builder.HeatmapPanel.DataFormat)
}

func TestLegendCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Hide))

	req.False(panel.Builder.HeatmapPanel.Legend.Show)
}

func TestZeroBucketsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideZeroBuckets())

	req.True(panel.Builder.HeatmapPanel.HideZeroBuckets)
}

func TestZeroBucketsCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", ShowZeroBuckets())

	req.False(panel.Builder.HeatmapPanel.HideZeroBuckets)
}

func TestCardsCanBeHighlighted(t *testing.T) {
	req := require.New(t)

	panel := New("", HightlightCards())

	req.True(panel.Builder.HeatmapPanel.HighlightCards)
}

func TestCardsCanBeNotHighlighted(t *testing.T) {
	req := require.New(t)

	panel := New("", NoHightlightCards())

	req.False(panel.Builder.HeatmapPanel.HighlightCards)
}

func TestYBucketsCanBeReversed(t *testing.T) {
	req := require.New(t)

	panel := New("", ReverseYBuckets())

	req.True(panel.Builder.HeatmapPanel.ReverseYBuckets)
}

func TestTooltipsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideTooltip())

	req.False(panel.Builder.HeatmapPanel.Tooltip.Show)
}

func TestTooltipHistogramsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideTooltipHistogram())

	req.False(panel.Builder.HeatmapPanel.Tooltip.ShowHistogram)
}

func TestTooltipDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", TooltipDecimals(3))

	req.Equal(3, panel.Builder.HeatmapPanel.TooltipDecimals)
}

func TestXAxisCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideXAxis())

	req.False(panel.Builder.HeatmapPanel.XAxis.Show)
}
