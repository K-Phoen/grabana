package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/heatmap/axis"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestHeatmapCanBeDecoded(t *testing.T) {
	req := require.New(t)

	panel := DashboardHeatmap{
		Description:     "awesome description",
		Span:            12,
		Height:          "300px",
		Transparent:     true,
		Datasource:      "some-prometheus",
		DataFormat:      "time_series_buckets",
		Repeat:          "ds",
		RepeatDirection: "vertical",
		HideZeroBuckets: true,
		HighlightCards:  true,
		Targets:         nil,
		ReverseYBuckets: true,
		YAxis: &HeatmapYAxis{
			Unit: "none",
		},
	}

	_, err := panel.toOption()
	req.NoError(err)
}

func TestHeatmapTooltipCanBeHidden(t *testing.T) {
	req := require.New(t)

	decimals := 2

	panel := DashboardHeatmap{
		Tooltip: &HeatmapTooltip{
			Show:          false,
			ShowHistogram: false,
			Decimals:      &decimals,
		},
	}

	opts, err := panel.toOption()
	req.NoError(err)

	builder := sdk.NewBoard("")
	_, err = row.New(builder, "", opts)

	req.NoError(err)

	req.Len(builder.Rows, 1)
	req.Len(builder.Rows[0].Panels, 1)

	sdkPanel := builder.Rows[0].Panels[0]

	req.False(sdkPanel.HeatmapPanel.Tooltip.Show)
	req.False(sdkPanel.HeatmapPanel.Tooltip.ShowHistogram)
	req.Equal(decimals, sdkPanel.HeatmapPanel.TooltipDecimals)
}

func TestHeatmapYAxisCanBeDecoded(t *testing.T) {
	req := require.New(t)

	decimals := 2
	min := float64(1)
	max := float64(3)
	axisInput := HeatmapYAxis{
		Decimals: &decimals,
		Unit:     "none",
		Max:      &min,
		Min:      &max,
	}

	decoded := axis.New(axisInput.toOptions()...).Builder
	req.Equal("none", decoded.Format)
	req.Equal(2, *decoded.Decimals)
}

func TestHeatmapCanNotBeDecodedIfDataFormatIsInvalid(t *testing.T) {
	req := require.New(t)

	panel := DashboardHeatmap{
		Span:            12,
		Height:          "300px",
		Transparent:     true,
		Datasource:      "some-prometheus",
		DataFormat:      "invalid value here",
		HideZeroBuckets: true,
		HighlightCards:  true,
		Targets:         nil,
		ReverseYBuckets: true,
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrInvalidDataFormat, err)
}

func TestHeatmapCanNotBeDecodedIfTargetIsInvalid(t *testing.T) {
	req := require.New(t)

	panel := DashboardHeatmap{
		Span:            12,
		Height:          "300px",
		Transparent:     true,
		Datasource:      "prometheus",
		DataFormat:      "time_series_buckets",
		HideZeroBuckets: true,
		HighlightCards:  true,
		Targets: []Target{
			{},
		},
		ReverseYBuckets: true,
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrTargetNotConfigured, err)
}
