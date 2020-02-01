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

func TestItCanBeHidden(t *testing.T) {
	req := require.New(t)

	a := New(Hide())

	req.False(a.Builder.Show)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Label("memory"))

	req.Equal("memory", a.Builder.Label)
}

func TestLogBaseCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(LogBase(2))

	req.Equal(2, a.Builder.LogBase)
}

func TestMinCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Min(1))

	req.Equal(float64(1), a.Builder.Min.Value)
}

func TestMaxCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New(Max(1))

	req.Equal(float64(1), a.Builder.Max.Value)
}
