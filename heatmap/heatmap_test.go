package heatmap

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/heatmap/axis"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewHeatmapPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("Heatmap panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("Heatmap panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestHeatmapPanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestHeatmapPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.NoError(err)
	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.NoError(err)
	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithInfluxDBTarget("buckets()"))

	req.NoError(err)
	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.NoError(err)
	req.Len(panel.Builder.HeatmapPanel.Targets, 1)
}

func TestHeatmapPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestInvalidHeatmapPanelWidthIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(-6))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestHeatmapPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestHeatmapPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestHeatmapPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestHeatmapPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("prometheus-default"))

	req.NoError(err)
	req.Equal("prometheus-default", panel.Builder.Datasource.LegacyName)
}

func TestDataFormatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataFormat(TimeSeriesBuckets))

	req.NoError(err)
	req.Equal("tsbuckets", panel.Builder.HeatmapPanel.DataFormat)
}

func TestRepeatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Repeat("ds"))

	req.NoError(err)
	req.NotNil(panel.Builder.Repeat)
	req.Equal("ds", *panel.Builder.Repeat)
}

func TestLegendCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Hide))

	req.NoError(err)
	req.False(panel.Builder.HeatmapPanel.Legend.Show)
}

func TestZeroBucketsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HideZeroBuckets())

	req.NoError(err)
	req.True(panel.Builder.HeatmapPanel.HideZeroBuckets)
}

func TestZeroBucketsCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ShowZeroBuckets())

	req.NoError(err)
	req.False(panel.Builder.HeatmapPanel.HideZeroBuckets)
}

func TestCardsCanBeHighlighted(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HighlightCards())

	req.NoError(err)
	req.True(panel.Builder.HeatmapPanel.HighlightCards)
}

func TestCardsCanBeNotHighlighted(t *testing.T) {
	req := require.New(t)

	panel, err := New("", NoHighlightCards())

	req.NoError(err)
	req.False(panel.Builder.HeatmapPanel.HighlightCards)
}

func TestYBucketsCanBeReversed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ReverseYBuckets())

	req.NoError(err)
	req.True(panel.Builder.HeatmapPanel.ReverseYBuckets)
}

func TestTooltipsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HideTooltip())

	req.NoError(err)
	req.False(panel.Builder.HeatmapPanel.Tooltip.Show)
}

func TestTooltipHistogramsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HideTooltipHistogram())

	req.NoError(err)
	req.False(panel.Builder.HeatmapPanel.Tooltip.ShowHistogram)
}

func TestTooltipDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", TooltipDecimals(3))

	req.NoError(err)
	req.Equal(3, panel.Builder.HeatmapPanel.TooltipDecimals)
}

func TestXAxisCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HideXAxis())

	req.NoError(err)
	req.False(panel.Builder.HeatmapPanel.XAxis.Show)
}

func TestYAxisCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", YAxis(axis.Unit("none")))

	req.NoError(err)
	req.Equal("none", panel.Builder.HeatmapPanel.YAxis.Format)
}
