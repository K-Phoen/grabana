package singlestat

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewSingleStatPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("SingleStat panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("SingleStat panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestSingleStatPanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestSingleStatPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.NoError(err)
	req.Len(panel.Builder.SinglestatPanel.Targets, 1)
}

func TestSingleStatPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.NoError(err)
	req.Len(panel.Builder.SinglestatPanel.Targets, 1)
}

func TestSingleStatPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithInfluxDBTarget("buckets()"))

	req.NoError(err)
	req.Len(panel.Builder.SinglestatPanel.Targets, 1)
}

func TestSingleStatPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.NoError(err)
	req.Len(panel.Builder.SinglestatPanel.Targets, 1)
}

func TestSingleStatPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestSingleStatRejectsInvalidSpans(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(16))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestSingleStatPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestSingleStatPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestSingleStatPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestSingleStatPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("prometheus-default"))

	req.NoError(err)
	req.Equal("prometheus-default", panel.Builder.Datasource.LegacyName)
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
	req.Equal("bytes", panel.Builder.SinglestatPanel.Format)
}

func TestDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Decimals(3))

	req.NoError(err)
	req.Equal(3, panel.Builder.SinglestatPanel.Decimals)
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
	req.True(panel.Builder.SinglestatPanel.SparkLine.Show)
	req.False(panel.Builder.SinglestatPanel.SparkLine.Full)
}

func TestFullSparkLineCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", FullSparkLine())

	req.NoError(err)
	req.True(panel.Builder.SinglestatPanel.SparkLine.Show)
	req.True(panel.Builder.SinglestatPanel.SparkLine.Full)
}

func TestSparkLineColorCanBeSet(t *testing.T) {
	req := require.New(t)
	color := "rgb(31, 120, 193)"

	panel, err := New("", SparkLineColor(color))

	req.NoError(err)
	req.Equal(color, *panel.Builder.SinglestatPanel.SparkLine.LineColor)
}

func TestSparkLineFillColorCanBeSet(t *testing.T) {
	req := require.New(t)
	color := "rgb(31, 120, 193)"

	panel, err := New("", SparkLineFillColor(color))

	req.NoError(err)
	req.Equal(color, *panel.Builder.SinglestatPanel.SparkLine.FillColor)
}

func TestSparkLineYMinCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SparkLineYMin(0))

	req.NoError(err)
	req.Equal(float64(0), *panel.Builder.SinglestatPanel.SparkLine.YMin)
}

func TestSparkLineYMaxCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SparkLineYMax(0))

	req.NoError(err)
	req.Equal(float64(0), *panel.Builder.SinglestatPanel.SparkLine.YMax)
}

func TestValueTypeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ValueType(Current))

	req.NoError(err)
	req.Equal(string(Current), panel.Builder.SinglestatPanel.ValueName)
}

func TestValueFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ValueFontSize("120%"))

	req.NoError(err)
	req.Equal("120%", panel.Builder.SinglestatPanel.ValueFontSize)
}

func TestPrefixCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Prefix("joe"))

	req.NoError(err)
	req.Equal("joe", *panel.Builder.SinglestatPanel.Prefix)
}

func TestPrefixFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", PrefixFontSize("120%"))

	req.NoError(err)
	req.Equal("120%", *panel.Builder.SinglestatPanel.PrefixFontSize)
}

func TestPostfixCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Postfix("joe"))

	req.NoError(err)
	req.Equal("joe", *panel.Builder.SinglestatPanel.Postfix)
}

func TestPostfixFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", PostfixFontSize("120%"))

	req.NoError(err)
	req.Equal("120%", *panel.Builder.SinglestatPanel.PostfixFontSize)
}

func TestValueCanBeColored(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ColorValue())

	req.NoError(err)
	req.True(panel.Builder.SinglestatPanel.ColorValue)
}

func TestBackgroundCanBeColored(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ColorBackground())

	req.NoError(err)
	req.True(panel.Builder.SinglestatPanel.ColorBackground)
}

func TestThresholdsCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Thresholds([2]string{"20", "30"}))

	req.NoError(err)
	req.Equal("20,30", panel.Builder.SinglestatPanel.Thresholds)
}

func TestThresholdColorsCanBeSet(t *testing.T) {
	req := require.New(t)
	colors := [3]string{"#299c46", "rgba(237, 129, 40, 0.89)", "#d44a3a"}

	panel, err := New("", Colors(colors))

	req.NoError(err)
	req.Equal([]string{colors[0], colors[1], colors[2]}, panel.Builder.SinglestatPanel.Colors)
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
	req.Len(panel.Builder.SinglestatPanel.RangeMaps, 3)
}
