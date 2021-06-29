package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/heatmap/axis"

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
