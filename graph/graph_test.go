package graph

import (
	"testing"

	"github.com/K-Phoen/grabana/graph/series"

	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewGraphPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Graph panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Graph panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithInfluxDBTarget("buckets()"))

	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestGraphPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestGraphPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestGraphPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestLeftYAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", LeftYAxis(axis.Hide()))

	req.False(panel.Builder.Yaxes[0].Show)
}

func TestRightYAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", RightYAxis())

	req.True(panel.Builder.Yaxes[0].Show)
}

func TestXAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", XAxis(axis.Hide()))

	req.False(panel.Builder.Xaxis.Show)
}

func TestAlertsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Alert("some alert"))

	req.NotNil(panel.Builder.Alert)
	req.Equal("some alert", panel.Builder.Alert.Name)
}

func TestDrawModeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Draw(Lines, Points, Bars))

	req.True(panel.Builder.Lines)
	req.True(panel.Builder.Points)
	req.True(panel.Builder.Bars)
}

func TestLineFillCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Fill(3))

	req.Equal(3, panel.Builder.Fill)
}

func TestLineWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", LineWidth(3))

	req.Equal(uint(3), panel.Builder.Linewidth)
}

func TestStaircaseModeCanBeEnabled(t *testing.T) {
	req := require.New(t)

	panel := New("", Staircase())

	req.True(panel.Builder.GraphPanel.SteppedLine)
}

func TestPointRadiusCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", PointRadius(3))

	req.Equal(float32(3), panel.Builder.Pointradius)
}

func TestNullValueModeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Null(AsNull))

	req.Equal("null", panel.Builder.GraphPanel.NullPointMode)
}

func TestLegendCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Hide))

	req.False(panel.Builder.GraphPanel.Legend.Show)
}

func TestLegendCanBeShownToTheRight(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(ToTheRight))

	req.True(panel.Builder.GraphPanel.Legend.RightSide)
}

func TestLegendCanBeDisplayedAsATable(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(AsTable))

	req.True(panel.Builder.GraphPanel.Legend.AlignAsTable)
}

func TestLegendCanShowAvg(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Avg))

	req.True(panel.Builder.GraphPanel.Legend.Avg)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowMin(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Min))

	req.True(panel.Builder.GraphPanel.Legend.Min)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowMax(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Max))

	req.True(panel.Builder.GraphPanel.Legend.Max)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowCurrent(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Current))

	req.True(panel.Builder.GraphPanel.Legend.Current)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowTotal(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Total))

	req.True(panel.Builder.GraphPanel.Legend.Total)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanHideZeroSeries(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(NoZeroSeries))

	req.True(panel.Builder.GraphPanel.Legend.HideZero)
}

func TestLegendCanHideNullSeries(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(NoNullSeries))

	req.True(panel.Builder.GraphPanel.Legend.HideEmpty)
}

func TestSeriesOverridesCanBeAdded(t *testing.T) {
	req := require.New(t)

	panel := New("", SeriesOverride(series.Alias("series"), series.Color("red")))

	req.Len(panel.Builder.GraphPanel.SeriesOverrides, 1)
}
