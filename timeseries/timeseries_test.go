package timeseries

import (
	"fmt"
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/scheme"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/timeseries/fields"
	"github.com/K-Phoen/grabana/timeseries/threshold"
	"github.com/stretchr/testify/require"
)

func TestNewTimeSeriesPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("TimeSeries panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("TimeSeries panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestTimeSeriesPanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestTimeSeriesPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveLokiTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithLokiTarget(
		"rate({app=\"loki\"}[$__interval])",
	))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithInfluxDBTarget("buckets()"))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelCanHaveStackdriverTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithStackdriverTarget(stackdriver.Gauge("pubsub.googleapis.com/subscription/ack_message_count")))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.Targets, 1)
}

func TestTimeSeriesPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestInvalidSpanIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(20))

	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestTimeSeriesPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestTimeSeriesPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestTimeSeriesPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestTimeSeriesPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("prometheus-default"))

	req.NoError(err)
	req.Equal("prometheus-default", panel.Builder.Datasource.LegacyName)
}

func TestAlertsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("panel title", Alert("some alert"))

	req.NoError(err)
	req.NotNil(panel.Alert)
	req.Equal("panel title", panel.Alert.Builder.Name)
}

func TestLineWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", LineWidth(3))

	req.NoError(err)
	req.Equal(3, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineWidth)
}

func TestInvalidLineWidthIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", LineWidth(20))

	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestFillOpacityCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", FillOpacity(10))

	req.NoError(err)
	req.Equal(10, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.FillOpacity)
}

func TestInvalidFillOpacityIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", FillOpacity(-1))

	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestPointSizeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", PointSize(3))

	req.NoError(err)
	req.Equal(3, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize)
}

func TestInvalidPointSizeIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", PointSize(42))

	req.ErrorIs(err, errors.ErrInvalidArgument)
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
	req.Equal("hidden", panel.Builder.TimeseriesPanel.Options.Legend.DisplayMode)
}

func TestLegendCanBeDisplayedAsATable(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(AsTable))

	req.NoError(err)
	req.Equal("table", panel.Builder.TimeseriesPanel.Options.Legend.DisplayMode)
}

func TestLegendCanBeDisplayedAsAList(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(AsList))

	req.NoError(err)
	req.Equal("list", panel.Builder.TimeseriesPanel.Options.Legend.DisplayMode)
}

func TestLegendCanBeShownToTheRight(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(ToTheRight))

	req.NoError(err)
	req.Equal("right", panel.Builder.TimeseriesPanel.Options.Legend.Placement)
}

func TestLegendCanBeShownToTheBottom(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Legend(Bottom))

	req.NoError(err)
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

			panel, err := New("", Legend(tc.option))

			req.NoError(err)
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

			panel, err := New("", Lines(tc.mode))

			req.NoError(err)
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

			panel, err := New("", Bars(tc.mode))

			req.NoError(err)
			req.Equal(tc.expected, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.BarAlignment)
			req.Equal("bars", panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle)
		})
	}
}

func TestSeriesCanBeDisplayedAsPoints(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Points())

	req.NoError(err)
	req.Equal("points", panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle)
}

func TestAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Axis(axis.Decimals(2)))

	req.NoError(err)
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

			panel, err := New("", GradientMode(tc.mode))

			req.NoError(err)
			req.Equal(tc.expected, panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode)
		})
	}
}

func TestFieldOverridesCanBeDefined(t *testing.T) {
	req := require.New(t)

	panel, err := New("", FieldOverride(
		fields.ByQuery("A"),
		fields.Unit("short"),
	))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.FieldConfig.Overrides, 1)
}

func TestThresholdsCanBeDefined(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Thresholds(
		threshold.Steps(
			threshold.Step{
				Color: "green",
				Value: 10,
			},
		),
	))

	req.NoError(err)
	req.Len(panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Thresholds.Steps, 2)
}

func TestColorSchemeCanBeDefined(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ColorScheme(scheme.SingleColor("yellow")))

	req.NoError(err)
	req.Equal("fixed", panel.Builder.TimeseriesPanel.FieldConfig.Defaults.Color.Mode)
}
