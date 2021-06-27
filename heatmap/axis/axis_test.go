package axis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAxisCanBeCreated(t *testing.T) {
	req := require.New(t)

	a := New()

	req.Equal("short", a.Builder.Format)
	req.Equal(1, a.Builder.LogBase)
	req.True(a.Builder.Show)
}

func TestUnitCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Unit("time"))

	req.Equal("time", a.Builder.Format)
}

func TestDecimalsCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Decimals(2))

	req.Equal(2, *a.Builder.Decimals)
}

func TestMinCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Min(1))

	req.Equal("1.000000", *a.Builder.Min)
}

func TestMaxCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Max(1))

	req.Equal("1.000000", *a.Builder.Max)
}
