package grabana

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDashboardsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("My dashboard")

	req.Equal(uint(0), panel.board.ID)
	req.Equal("My dashboard", panel.board.Title)
	req.Empty(panel.board.Timezone)
	req.True(panel.board.SharedCrosshair)
	req.NotEmpty(panel.board.Timepicker.RefreshIntervals)
	req.NotEmpty(panel.board.Timepicker.TimeOptions)
	req.NotEmpty(panel.board.Time.From)
	req.NotEmpty(panel.board.Time.To)
}

func TestDashboardCanBeMadeEditable(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", Editable())

	req.True(panel.board.Editable)
}

func TestDashboardCanBeMadeReadOnly(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", ReadOnly())

	req.False(panel.board.Editable)
}

func TestDashboardCanBeAutoRefreshed(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", AutoRefresh("5s"))

	req.True(panel.board.Refresh.Flag)
	req.Equal("5s", panel.board.Refresh.Value)
}
