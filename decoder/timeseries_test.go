package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/stretchr/testify/require"
)

func TestTimeSeriesCanBeDecoded(t *testing.T) {
	req := require.New(t)

	panel := DashboardTimeSeries{
		Title:       "awesome timeseries",
		Description: "awesome description",
		Span:        12,
		Height:      "300px",
		Transparent: true,
		Datasource:  "some-prometheus",
		Repeat:      "ds",
		Legend:      []string{"hide"},
		Visualization: &TimeSeriesVisualization{
			GradientMode: "opacity",
			Tooltip:      "single_series",
			FillOpacity:  intPtr(5),
		},
		Axis: &TimeSeriesAxis{
			SoftMin:  intPtr(1),
			SoftMax:  intPtr(10),
			Min:      intPtr(0),
			Max:      intPtr(11),
			Decimals: intPtr(2),
			Display:  "auto",
			Scale:    "linear",
			Unit:     "short",
			Label:    "Some label",
		},
	}

	rowOption, err := panel.toOption()
	req.NoError(err)

	testBoard := dashboard.New("test-board", dashboard.Row("test row", rowOption))
	req.Len(testBoard.Internal().Rows, 1)
	panels := testBoard.Internal().Rows[0].Panels
	req.Len(panels, 1)

	sdkPanel := panels[0]
	tsPanel := panels[0].TimeseriesPanel

	req.NotNil(tsPanel)
	req.Equal(panel.Title, sdkPanel.Title)
	req.Equal(panel.Description, *sdkPanel.Description)
	req.Equal(panel.Datasource, *sdkPanel.Datasource)
	req.Equal(panel.Repeat, *sdkPanel.Repeat)
	req.Equal(panel.Span, sdkPanel.Span)
	req.True(sdkPanel.Transparent)
	req.Equal("hidden", tsPanel.Options.Legend.DisplayMode)

	// visualization
	req.Equal("opacity", tsPanel.FieldConfig.Defaults.Custom.GradientMode)
	req.Equal("single", tsPanel.Options.Tooltip.Mode)
	req.Equal(5, tsPanel.FieldConfig.Defaults.Custom.FillOpacity)

	// axis
	req.Equal("Some label", tsPanel.FieldConfig.Defaults.Custom.AxisLabel)
	req.Equal(2, *tsPanel.FieldConfig.Defaults.Decimals)
	req.Equal(0, *tsPanel.FieldConfig.Defaults.Min)
	req.Equal(11, *tsPanel.FieldConfig.Defaults.Max)
	req.Equal(1, *tsPanel.FieldConfig.Defaults.Custom.AxisSoftMin)
	req.Equal(10, *tsPanel.FieldConfig.Defaults.Custom.AxisSoftMax)
	req.Equal("short", tsPanel.FieldConfig.Defaults.Unit)
}

func TestTimeSeriesCanNotBeDecodedIfTargetIsInvalid(t *testing.T) {
	req := require.New(t)

	panel := DashboardTimeSeries{
		Span:        12,
		Height:      "300px",
		Transparent: true,
		Datasource:  "prometheus",
		Targets: []Target{
			{},
		},
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrTargetNotConfigured, err)
}

func TestTimeSeriesLegendRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	panel := DashboardTimeSeries{
		Legend: []string{"unknown"},
	}
	_, err := panel.legend()
	req.Error(err)
	req.Equal(ErrInvalidLegendAttribute, err)
}

func TestTimeSeriesLegendCanBeDecided(t *testing.T) {
	req := require.New(t)

	panel := DashboardTimeSeries{
		Legend: []string{
			"hide",
			"as_table",
			"as_list",
			"to_bottom",
			"to_the_right",
			"min",
			"max",
			"avg",
			"first",
			"first_non_null",
			"last",
			"last_non_null",
			"count",
			"total",
			"range",
		},
	}

	expectedOptions := []timeseries.LegendOption{
		timeseries.Hide,
		timeseries.AsTable,
		timeseries.AsList,
		timeseries.Bottom,
		timeseries.ToTheRight,
		timeseries.Min,
		timeseries.Max,
		timeseries.Avg,
		timeseries.First,
		timeseries.FirstNonNull,
		timeseries.Last,
		timeseries.LastNonNull,
		timeseries.Total,
		timeseries.Count,
		timeseries.Range,
	}

	legendOptions, err := panel.legend()
	req.NoError(err)
	req.ElementsMatch(expectedOptions, legendOptions)
}

func TestTimeSeriesVisualizationCanBeConfigured(t *testing.T) {
	req := require.New(t)

	tsViz := &TimeSeriesVisualization{
		GradientMode: "opacity",
		Tooltip:      "single_series",
		FillOpacity:  intPtr(30),
		PointSize:    intPtr(4),
	}

	opts, err := tsViz.toOptions()
	req.NoError(err)

	tsPanel := timeseries.New("", opts...)

	req.Equal("opacity", tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode)
	req.Equal(30, tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.FillOpacity)
	req.Equal(4, tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize)
	req.Equal("single", tsPanel.Builder.TimeseriesPanel.Options.Tooltip.Mode)
}

func TestTimeSeriesGradientModeCanBeDecided(t *testing.T) {
	testCases := []struct {
		mode         string
		expectedMode timeseries.GradientType
	}{
		{
			mode:         "no_gradient",
			expectedMode: timeseries.NoGradient,
		},
		{
			mode:         "opacity",
			expectedMode: timeseries.Opacity,
		},
		{
			mode:         "hue",
			expectedMode: timeseries.Hue,
		},
		{
			mode:         "scheme",
			expectedMode: timeseries.Scheme,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.mode, func(t *testing.T) {
			req := require.New(t)

			viz := TimeSeriesVisualization{
				GradientMode: tc.mode,
			}
			opts, err := viz.toOptions()

			req.NoError(err)

			tsPanel := timeseries.New("", opts...)

			req.Equal(string(tc.expectedMode), tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode)
		})
	}
}

func TestTimeSeriesGradientModeRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	viz := TimeSeriesVisualization{
		GradientMode: "invalid",
	}
	_, err := viz.toOptions()

	req.Equal(ErrInvalidGradientMode, err)
}

func TestTimeSeriesTooltipCanBeDecided(t *testing.T) {
	testCases := []struct {
		mode         string
		expectedMode timeseries.TooltipMode
	}{
		{
			mode:         "single_series",
			expectedMode: timeseries.SingleSeries,
		},
		{
			mode:         "all_series",
			expectedMode: timeseries.AllSeries,
		},
		{
			mode:         "no_series",
			expectedMode: timeseries.NoSeries,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.mode, func(t *testing.T) {
			req := require.New(t)

			viz := TimeSeriesVisualization{
				Tooltip: tc.mode,
			}
			opts, err := viz.toOptions()

			req.NoError(err)

			tsPanel := timeseries.New("", opts...)

			req.Equal(string(tc.expectedMode), tsPanel.Builder.TimeseriesPanel.Options.Tooltip.Mode)
		})
	}
}

func TestTimeSeriesTooltipRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	viz := TimeSeriesVisualization{
		Tooltip: "invalid",
	}
	_, err := viz.toOptions()

	req.Equal(ErrInvalidTooltipMode, err)
}

func TestTimeSeriesAxisSupportsDisplay(t *testing.T) {
	testCases := []struct {
		value    string
		expected axis.PlacementMode
	}{
		{
			value:    "none",
			expected: axis.Hidden,
		},
		{
			value:    "auto",
			expected: axis.Auto,
		},
		{
			value:    "left",
			expected: axis.Left,
		},
		{
			value:    "right",
			expected: axis.Right,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.value, func(t *testing.T) {
			req := require.New(t)

			tsAxis := TimeSeriesAxis{
				Display: tc.value,
			}
			opt, err := tsAxis.toOptions()

			req.NoError(err)

			tsPanel := timeseries.New("", timeseries.Axis(opt...))

			req.Equal(string(tc.expected), tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.AxisPlacement)
		})
	}
}

func TestTimeSeriesAxisRejectsInvalidDisplay(t *testing.T) {
	req := require.New(t)

	tsAxis := TimeSeriesAxis{
		Display: "invalid",
	}
	_, err := tsAxis.toOptions()

	req.Equal(ErrInvalidAxisDisplay, err)
}

func TestTimeSeriesAxisSupportsScale(t *testing.T) {
	testCases := []struct {
		value        string
		expectedType string
		expectedLog  int
	}{
		{
			value:        "linear",
			expectedType: "linear",
		},
		{
			value:        "log2",
			expectedType: "log",
			expectedLog:  2,
		},
		{
			value:        "log10",
			expectedType: "log",
			expectedLog:  10,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.value, func(t *testing.T) {
			req := require.New(t)

			tsAxis := TimeSeriesAxis{
				Scale: tc.value,
			}
			opt, err := tsAxis.toOptions()

			req.NoError(err)

			tsPanel := timeseries.New("", timeseries.Axis(opt...))

			req.Equal(tc.expectedType, tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Type)
			req.Equal(tc.expectedLog, tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Log)
		})
	}
}

func TestTimeSeriesAxisRejectsInvalidScale(t *testing.T) {
	req := require.New(t)

	tsAxis := TimeSeriesAxis{
		Scale: "invalid",
	}
	_, err := tsAxis.toOptions()

	req.Equal(ErrInvalidAxisScale, err)
}
