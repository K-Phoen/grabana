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
