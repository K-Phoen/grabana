package series

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestAliasCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := Alias("Error - .*")(series)

	req.NoError(err)
	req.Equal("Error - .*", series.Alias)
}

func TestColorCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := Color("#65c5db")(series)

	req.NoError(err)
	req.NotNil(series.Color)
	req.Equal("#65c5db", *series.Color)
}

func TestDashesCanBeEnabled(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := Dashes(true)(series)

	req.NoError(err)
	req.NotNil(series.Dashes)
	req.True(*series.Dashes)
}

func TestLinesCanBeEnabled(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := Lines(true)(series)

	req.NoError(err)
	req.NotNil(series.Lines)
	req.True(*series.Lines)
}

func TestFillOpacityCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := Fill(2)(series)

	req.NoError(err)
	req.NotNil(*series.Fill)
	req.Equal(2, *series.Fill)
}

func TestInvalidFillOpacityIsRejected(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := Fill(-2)(series)

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestLineWidthCanBeDefined(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := LineWidth(3)(series)

	req.NoError(err)
	req.NotNil(*series.LineWidth)
	req.Equal(3, *series.LineWidth)
}

func TestInvalidLineWidthIsRejected(t *testing.T) {
	req := require.New(t)
	series := &sdk.SeriesOverride{}

	err := LineWidth(-3)(series)

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}
