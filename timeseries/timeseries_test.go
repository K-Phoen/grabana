package timeseries

import (
	"fmt"
	"testing"

	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/grabana/timeseries/axis"
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

func TestFillOpacityCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", FillOpacity(10))

	req.Equal(10, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.FillOpacity)
}

func TestPointSizeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", PointSize(3))

	req.Equal(3, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize)
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

func TestLegendCanShowCalculatedData(t *testing.T) {
	testCases := []struct {
		option   LegendOption
		expected string
	}{
		{option: Min, expected: "min"},
		{option: Max, expected: "max"},
		{option: Avg, expected: "mean"},

		{option: Total, expected: "sum"},
		{option: Count, expected: "count"},
		{option: Range, expected: "range"},

		{option: First, expected: "first"},
		{option: FirstNonNull, expected: "firstNotNull"},
		{option: Last, expected: "last"},
		{option: LastNonNull, expected: "lastNotNull"},
	}

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("option %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			panel := New("", Legend(tc.option))

			req.Contains(panel.Builder.TimeseriesPanel.Options.Legend.Calcs, tc.expected)
		})
	}
}

func TestLineInterpolationCanBeConfigured(t *testing.T) {
	testCases := []struct {
		mode     LineInterpolationMode
		expected string
	}{
		{
			mode:     Linear,
			expected: string(Linear),
		},
		{
			mode:     Smooth,
			expected: string(Smooth),
		},
		{
			mode:     StepBefore,
			expected: string(StepBefore),
		},
		{
			mode:     StepAfter,
			expected: string(StepAfter),
		},
	}

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("interpolation mode %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			panel := New("", Lines(tc.mode))

			req.Equal(tc.expected, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineInterpolation)
			req.Equal("line", panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle)
		})
	}
}

func TestBarsAlignmentCanBeConfigured(t *testing.T) {
	testCases := []struct {
		mode     BarAlignment
		expected int
	}{
		{
			mode:     AlignCenter,
			expected: 0,
		},
		{
			mode:     AlignBefore,
			expected: -1,
		},
		{
			mode:     AlignAfter,
			expected: 1,
		},
	}

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("alignment %d", tc.expected), func(t *testing.T) {
			req := require.New(t)

			panel := New("", Bars(tc.mode))

			req.Equal(tc.expected, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.BarAlignment)
			req.Equal("bars", panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle)
		})
	}
}

func TestSeriesCanBeDisplayedAsPoints(t *testing.T) {
	req := require.New(t)

	panel := New("", Points())

	req.Equal("points", panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle)
}

func TestAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Axis(axis.Decimals(2)))

	req.Equal(2, *panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Decimals)
}

func TestGradientModeCanBeConfigured(t *testing.T) {
	testCases := []struct {
		mode     GradientType
		expected string
	}{
		{
			mode:     NoGradient,
			expected: "none",
		},
		{
			mode:     Opacity,
			expected: "opacity",
		},
		{
			mode:     Hue,
			expected: "hue",
		},
		{
			mode:     Scheme,
			expected: "scheme",
		},
	}

	for _, test := range testCases {
		tc := test
		t.Run(fmt.Sprintf("mode %s", tc.expected), func(t *testing.T) {
			req := require.New(t)

			panel := New("", GradientMode(tc.mode))

			req.Equal(tc.expected, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode)
		})
	}
}
