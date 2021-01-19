package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeatmapCanBeDecoded(t *testing.T) {
	req := require.New(t)

	panel := DashboardHeatmap{
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
