package graph

import (
	"testing"

	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/graph/series"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewGraphPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("Graph panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("Graph panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestGraphPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.NoError(err)
	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.NoError(err)
	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithInfluxDBTarget("buckets()"))

	req.NoError(err)
	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.NoError(err)
	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestInvalidGraphPanelWidthAreRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(-6))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestGraphPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestGraphPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestGraphPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestGraphPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("prometheus-default"))

	req.NoError(err)
	req.Equal("prometheus-default", panel.Builder.Datasource.LegacyName)
}

func TestLeftYAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", LeftYAxis(axis.Hide()))

	req.NoError(err)
	req.False(panel.Builder.Yaxes[0].Show)
}

func TestRightYAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", RightYAxis())

	req.NoError(err)
	req.True(panel.Builder.Yaxes[0].Show)
}

func TestXAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", XAxis(axis.Hide()))

	req.NoError(err)
	req.False(panel.Builder.Xaxis.Show)
}

func TestAlertsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("panel title", Alert("some alert"))

	req.NoError(err)
	req.NotNil(panel.Alert)
	req.Equal("panel title", panel.Alert.Builder.Name)
}

func TestDrawModeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Draw(Lines, Points, Bars))

	req.NoError(err)
	req.True(panel.Builder.Lines)
	req.True(panel.Builder.Points)
	req.True(panel.Builder.Bars)
}

func TestLineFillCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Fill(3))

	req.NoError(err)
	req.Equal(3, panel.Builder.Fill)
}

func TestInvalidLineFillIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Fill(30))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestLineWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", LineWidth(3))

	req.NoError(err)
	req.Equal(uint(3), panel.Builder.Linewidth)
}

func TestInvalidLineWidthIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", LineWidth(22))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestStaircaseModeCanBeEnabled(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Staircase())

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.SteppedLine)
}

func TestPointRadiusCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", PointRadius(3))

	req.NoError(err)
	req.Equal(float32(3), panel.Builder.Pointradius)
}

func TestInvalidPointRadiusIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", PointRadius(-3))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestNullValueModeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Null(AsNull))

	req.NoError(err)
	req.Equal("null", panel.Builder.GraphPanel.NullPointMode)
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
	req.False(panel.Builder.GraphPanel.Legend.Show)
}

func TestLegendCanBeShownToTheRight(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(ToTheRight))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.RightSide)
}

func TestLegendCanBeDisplayedAsATable(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(AsTable))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.AlignAsTable)
}

func TestLegendCanShowAvg(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Avg))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.Avg)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowMin(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Min))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.Min)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowMax(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Max))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.Max)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowCurrent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Current))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.Current)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanShowTotal(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Total))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.Total)
	req.True(panel.Builder.GraphPanel.Legend.Values)
}

func TestLegendCanHideZeroSeries(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(NoZeroSeries))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.HideZero)
}

func TestLegendCanHideNullSeries(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(NoNullSeries))

	req.NoError(err)
	req.True(panel.Builder.GraphPanel.Legend.HideEmpty)
}

func TestSeriesOverridesCanBeAdded(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SeriesOverride(series.Alias("series"), series.Color("red")))

	req.NoError(err)
	req.Len(panel.Builder.GraphPanel.SeriesOverrides, 1)
}

func TestInvalidSeriesOverridesAreRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", SeriesOverride(series.Fill(-1)))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}
