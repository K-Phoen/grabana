package singlestat

import (
	"testing"

	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/stretchr/testify/require"
)

func TestNewSingleStatPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("SingleStat panel")

	req.False(panel.Builder.IsNew)
	req.Equal("SingleStat panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestSingleStatPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.Len(panel.Builder.SinglestatPanel.Targets, 1)
}

func TestSingleStatPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.Len(panel.Builder.SinglestatPanel.Targets, 1)
}

func TestSingleStatPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestSingleStatPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestSingleStatPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestSingleStatPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestSingleStatPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestUnitCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Unit("bytes"))

	req.Equal("bytes", panel.Builder.SinglestatPanel.Format)
}

func TestDecimalsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Decimals(3))

	req.Equal(3, panel.Builder.SinglestatPanel.Decimals)
}

func TestSparkLineCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", SparkLine())

	req.True(panel.Builder.SinglestatPanel.SparkLine.Show)
	req.False(panel.Builder.SinglestatPanel.SparkLine.Full)
}

func TestFullSparkLineCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", FullSparkLine())

	req.True(panel.Builder.SinglestatPanel.SparkLine.Show)
	req.True(panel.Builder.SinglestatPanel.SparkLine.Full)
}

func TestSparkLineColorCanBeSet(t *testing.T) {
	req := require.New(t)
	color := "rgb(31, 120, 193)"

	panel := New("", SparkLineColor(color))

	req.Equal(color, *panel.Builder.SinglestatPanel.SparkLine.LineColor)
}

func TestSparkLineFillColorCanBeSet(t *testing.T) {
	req := require.New(t)
	color := "rgb(31, 120, 193)"

	panel := New("", SparkLineFillColor(color))

	req.Equal(color, *panel.Builder.SinglestatPanel.SparkLine.FillColor)
}

func TestSparkLineYMinCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", SparkLineYMin(0))

	req.Equal(float64(0), *panel.Builder.SinglestatPanel.SparkLine.YMin)
}

func TestSparkLineYMaxCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", SparkLineYMax(0))

	req.Equal(float64(0), *panel.Builder.SinglestatPanel.SparkLine.YMax)
}

func TestValueTypeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", ValueType(Current))

	req.Equal(string(Current), panel.Builder.SinglestatPanel.ValueName)
}

func TestValueFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", ValueFontSize("120%"))

	req.Equal("120%", panel.Builder.SinglestatPanel.ValueFontSize)
}

func TestPrefixCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Prefix("joe"))

	req.Equal("joe", *panel.Builder.SinglestatPanel.Prefix)
}

func TestPrefixFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", PrefixFontSize("120%"))

	req.Equal("120%", *panel.Builder.SinglestatPanel.PrefixFontSize)
}

func TestPostfixCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Postfix("joe"))

	req.Equal("joe", *panel.Builder.SinglestatPanel.Postfix)
}

func TestPostfixFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", PostfixFontSize("120%"))

	req.Equal("120%", *panel.Builder.SinglestatPanel.PostfixFontSize)
}

func TestValueCanBeColored(t *testing.T) {
	req := require.New(t)

	panel := New("", ColorValue())

	req.True(panel.Builder.SinglestatPanel.ColorValue)
}

func TestBackgroundCanBeColored(t *testing.T) {
	req := require.New(t)

	panel := New("", ColorBackground())

	req.True(panel.Builder.SinglestatPanel.ColorBackground)
}

func TestThresholdsCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Thresholds([2]string{"20", "30"}))

	req.Equal("20,30", panel.Builder.SinglestatPanel.Thresholds)
}

func TestThresholdColorsCanBeSet(t *testing.T) {
	req := require.New(t)
	colors := [3]string{"#299c46", "rgba(237, 129, 40, 0.89)", "#d44a3a"}

	panel := New("", Colors(colors))

	req.Equal([]string{colors[0], colors[1], colors[2]}, panel.Builder.SinglestatPanel.Colors)
}

func TestRangeToTextMappingsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", RangesToText([]RangeMap{
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

	req.Len(panel.Builder.SinglestatPanel.RangeMaps, 3)
}
