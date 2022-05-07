package stat

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewStatPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("Stat panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("Stat panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestStatPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.NoError(err)
	req.Len(panel.Builder.StatPanel.Targets, 1)
}

func TestStatPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.NoError(err)
	req.Len(panel.Builder.StatPanel.Targets, 1)
}

func TestStatPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithInfluxDBTarget("buckets()"))

	req.NoError(err)
	req.Len(panel.Builder.StatPanel.Targets, 1)
}

func TestStatPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.NoError(err)
	req.Len(panel.Builder.StatPanel.Targets, 1)
}

func TestStatPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestStatRejectsInvalidSpans(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(16))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestStatPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestStatPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestStatPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestStatPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("prometheus-default"))

	req.NoError(err)
	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestRepeatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Repeat("ds"))

	req.NoError(err)
	req.NotNil(panel.Builder.Repeat)
	req.Equal("ds", *panel.Builder.Repeat)
}

func TestUnitCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Unit("bytes"))

	req.NoError(err)
	req.Equal("bytes", panel.Builder.StatPanel.FieldConfig.Defaults.Unit)
}

func TestDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Decimals(3))

	req.NoError(err)
	req.Equal(3, *panel.Builder.StatPanel.FieldConfig.Defaults.Decimals)
}

func TestInvalidDecimalsAreRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Decimals(-3))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestSparkLineCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SparkLine())

	req.NoError(err)
	req.Equal("area", panel.Builder.StatPanel.Options.GraphMode)
}

func TestSparkLineYMinCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SparkLineYMin(0))

	req.NoError(err)
	req.Equal(0, *panel.Builder.StatPanel.FieldConfig.Defaults.Min)
}

func TestSparkLineYMaxCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SparkLineYMax(0))

	req.NoError(err)
	req.Equal(0, *panel.Builder.StatPanel.FieldConfig.Defaults.Max)
}

func TestValueTypeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ValueType(LastNonNull))

	req.NoError(err)
	req.Len(panel.Builder.StatPanel.Options.ReduceOptions.Calcs, 1)
	req.Equal("lastNotNull", panel.Builder.StatPanel.Options.ReduceOptions.Calcs[0])
}

func TestValueFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ValueFontSize(120))

	req.NoError(err)
	req.Equal(120, panel.Builder.StatPanel.Options.Text.ValueSize)
}

func TestTitleFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", TitleFontSize(120))

	req.NoError(err)
	req.Equal(120, panel.Builder.StatPanel.Options.Text.TitleSize)
}

func TestValueCanBeColored(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ColorValue())

	req.NoError(err)
	req.Equal("value", panel.Builder.StatPanel.Options.ColorMode)
}

func TestBackgroundCanBeColored(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ColorBackground())

	req.NoError(err)
	req.Equal("background", panel.Builder.StatPanel.Options.ColorMode)
}

func TestRangeToTextMappingsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", RangesToText([]RangeMap{
		{
			From: "0",
			To:   "20",
			Text: "Low",
		}, {
			From: "20",
			To:   "30",
			Text: "Average",
		}, {
			From: "30",
			To:   "",
			Text: "High",
		},
	}))

	req.NoError(err)
	req.Len(panel.Builder.StatPanel.RangeMaps, 3)
}
