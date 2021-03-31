package interval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIntervalVariablesCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("interval")

	req.Equal("interval", panel.Builder.Name)
	req.Equal("interval", panel.Builder.Label)
	req.Equal("interval", panel.Builder.Type)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("interval var", Label("IntervalVariable"))

	req.Equal("interval var", panel.Builder.Name)
	req.Equal("IntervalVariable", panel.Builder.Label)
}

func TestValuesCanBeSet(t *testing.T) {
	req := require.New(t)
	values := []string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"}

	panel := New("", Values(values))

	req.Equal("30s,1m,5m,10m,30m,1h,6h,12h", panel.Builder.Query)
}

func TestValuesAreSortedByDuration(t *testing.T) {
	req := require.New(t)
	values := []string{"12h", "30m", "6h", "30d", "30s"}

	panel := New("", Values(values))

	req.Equal("30s,30m,6h,12h,30d", panel.Builder.Query)
}

func TestDefaultValueCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Default("99"))

	req.Equal([]string{"99"}, panel.Builder.Current.Text.Value)
}

func TestLabelCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideLabel())

	req.Equal(uint8(1), panel.Builder.Hide)
}

func TestVariableCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", Hide())

	req.Equal(uint8(2), panel.Builder.Hide)
}
