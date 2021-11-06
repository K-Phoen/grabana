package timeseries

import (
	"testing"

	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewTimeSeriesPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("TimeSeries panel")

	req.False(panel.Builder.IsNew)
	req.Equal("TimeSeries panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestTimeSeriesPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithInfluxDBTarget("buckets()"))

	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestTimeSeriesPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestTimeSeriesPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestTimeSeriesPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestTimeSeriesPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestAlertsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Alert("some alert"))

	req.NotNil(panel.Builder.Alert)
	req.Equal("some alert", panel.Builder.Alert.Name)
}

func TestLineWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", LineWidth(3))

	req.Equal(3, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineWidth)
}

func TestRepeatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Repeat("ds"))

	req.NotNil(panel.Builder.Repeat)
	req.Equal("ds", *panel.Builder.Repeat)
}

func TestLegendCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Hide))

	req.Equal("hidden", panel.Builder.TimeseriesPanel.Options.Legend.DisplayMode)
}

func TestLegendCanBeDisplayedAsATable(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(AsTable))

	req.Equal("table", panel.Builder.TimeseriesPanel.Options.Legend.DisplayMode)
}

func TestLegendCanBeDisplayedAsAList(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(AsList))

	req.Equal("list", panel.Builder.TimeseriesPanel.Options.Legend.DisplayMode)
}

func TestLegendCanBeShownToTheRight(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(ToTheRight))

	req.Equal("right", panel.Builder.TimeseriesPanel.Options.Legend.Placement)
}

func TestLegendCanBeShownToTheBottom(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Bottom))

	req.Equal("bottom", panel.Builder.TimeseriesPanel.Options.Legend.Placement)
}

func TestLegendCanShowAvg(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Avg))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "mean")
}

func TestLegendCanShowMin(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Min))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "min")
}

func TestLegendCanShowMax(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Max))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "max")
}

func TestLegendCanShowTotal(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Total))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "sum")
}

func TestLegendCanShowCount(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Count))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "count")
}

func TestLegendCanShowRange(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Range))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "range")
}

func TestLegendCanShowFirst(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(First))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "first")
}

func TestLegendCanShowFirstNotNull(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(FirstNonNull))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "firstNotNull")
}

func TestLegendCanShowLast(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(Last))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "last")
}

func TestLegendCanShowLastNotNull(t *testing.T) {
	req := require.New(t)

	panel := New("", Legend(LastNonNull))

	req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, "lastNotNull")
}
