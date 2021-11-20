package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/timeseries"
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
			modeOpt, err := viz.gradientModeOption()

			req.NoError(err)

			tsPanel := timeseries.New("", modeOpt)

			req.Equal(string(tc.expectedMode), tsPanel.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode)
		})
	}
}

func TestTimeSeriesGradientModeRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	viz := TimeSeriesVisualization{
		GradientMode: "invalid",
	}
	_, err := viz.gradientModeOption()

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
			modeOpt, err := viz.tooltipOption()

			req.NoError(err)

			tsPanel := timeseries.New("", modeOpt)

			req.Equal(string(tc.expectedMode), tsPanel.Builder.TimeseriesPanel.Options.Tooltip.Mode)
		})
	}
}

func TestTimeSeriesTooltipRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	viz := TimeSeriesVisualization{
		Tooltip: "invalid",
	}
	_, err := viz.tooltipOption()

	req.Equal(ErrInvalidTooltipMode, err)
}
