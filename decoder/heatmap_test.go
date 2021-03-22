package decoder

import (
	"testing"

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
		HideZeroBuckets: true,
		HightlightCards: true,
		Targets:         nil,
		ReverseYBuckets: true,
	}

	_, err := panel.toOption()
	req.NoError(err)
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
		HightlightCards: true,
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
		HightlightCards: true,
		Targets: []Target{
			{},
		},
		ReverseYBuckets: true,
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrTargetNotConfigured, err)
}
