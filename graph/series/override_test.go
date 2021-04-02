package series

import (
	"testing"

	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/require"
)

func TestAliasCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	Alias("Error - .*")(series)

	req.Equal("Error - .*", series.Alias)
}

func TestColorCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	Color("#65c5db")(series)

	req.NotNil(series.Color)
	req.Equal("#65c5db", *series.Color)
}

func TestDashesCanBeEnabled(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	Dashes(true)(series)

	req.NotNil(series.Dashes)
	req.True(*series.Dashes)
}

func TestLinesCanBeEnabled(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	Lines(true)(series)

	req.NotNil(series.Lines)
	req.True(*series.Lines)
}

func TestFillOpacityCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	Fill(2)(series)

	req.NotNil(*series.Fill)
	req.Equal(2, *series.Fill)
}

func TestLineWidthCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	LineWidth(3)(series)

	req.NotNil(*series.LineWidth)
	req.Equal(3, *series.LineWidth)
}
