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

func TestSingleStatPanelCanBeMadeEditable(t *testing.T) {
	req := require.New(t)

	panel := New("", Editable())

	req.True(panel.Builder.Editable)
}

func TestSingleStatPanelCanBeMadeReadOnly(t *testing.T) {
	req := require.New(t)

	panel := New("", ReadOnly())

	req.False(panel.Builder.Editable)
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

func TestSingleStatPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestUnitCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Unit("bytes"))

	req.Equal("bytes", panel.Builder.Format)
}

func TestSparkLineCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", SparkLine())

	req.True(panel.Builder.SparkLine.Show)
	req.False(panel.Builder.SparkLine.Full)
}

func TestFullSparkLineCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", FullSparkLine())

	req.True(panel.Builder.SparkLine.Show)
	req.True(panel.Builder.SparkLine.Full)
}

func TestSparkLineColorCanBeSet(t *testing.T) {
	req := require.New(t)
	color := "rgb(31, 120, 193)"

	panel := New("", SparkLineColor(color))

	req.Equal(color, *panel.Builder.SparkLine.LineColor)
}

func TestSparkLineFillColorCanBeSet(t *testing.T) {
	req := require.New(t)
	color := "rgb(31, 120, 193)"

	panel := New("", SparkLineFillColor(color))

	req.Equal(color, *panel.Builder.SparkLine.FillColor)
}

func TestSparkLineYMinCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", SparkLineYMin(0))

	req.Equal(float64(0), *panel.Builder.SparkLine.YMin)
}

func TestSparkLineYMaxCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", SparkLineYMax(0))

	req.Equal(float64(0), *panel.Builder.SparkLine.YMax)
}

func TestValueTypeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", ValueType(Current))

	req.Equal(string(Current), panel.Builder.ValueName)
}

func TestValueFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", ValueFontSize("120%"))

	req.Equal("120%", panel.Builder.ValueFontSize)
}

func TestPrefixCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Prefix("joe"))

	req.Equal("joe", *panel.Builder.Prefix)
}

func TestPrefixFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", PrefixFontSize("120%"))

	req.Equal("120%", *panel.Builder.PrefixFontSize)
}

func TestPostfixCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Postfix("joe"))

	req.Equal("joe", *panel.Builder.Postfix)
}

func TestPostfixFontSizeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", PostfixFontSize("120%"))

	req.Equal("120%", *panel.Builder.PostfixFontSize)
}

func TestValueCanBeColored(t *testing.T) {
	req := require.New(t)

	panel := New("", ColorValue())

	req.True(panel.Builder.ColorValue)
}

func TestBackgroundCanBeColored(t *testing.T) {
	req := require.New(t)

	panel := New("", ColorBackground())

	req.True(panel.Builder.ColorBackground)
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

	req.Equal([]string{colors[0], colors[1], colors[2]}, panel.Builder.Colors)
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

	req.Len(panel.Builder.RangeMaps, 3)
}
