package row

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestNewRowsCanBeCreated(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "Some row")

	req.Equal("Some row", panel.builder.Title)
	req.True(panel.builder.ShowTitle)
}

func TestRowsCanHaveHiddenTitle(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", HideTitle())

	req.False(panel.builder.ShowTitle)
}

func TestRowsCanHaveVisibleTitle(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", ShowTitle())

	req.True(panel.builder.ShowTitle)
}

func TestRowsCanHaveGraphs(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", WithGraph("HTTP Rate"))

	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveTimeSeries(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", WithTimeSeries("HTTP Rate"))

	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveTextPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", WithText("HTTP Rate"))

	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveTablePanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", WithTable("Some table"))

	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveSingleStatPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", WithSingleStat("Some stat"))

	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveHeatmapPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", WithHeatmap("Some heatmap"))

	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveRepeatedPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", RepeatFor("repeated"))

	req.Equal("repeated", *panel.builder.Repeat)
}

func TestRowsCanBeCollapsedByDefault(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel := New(board, "", Collapse())

	req.True(panel.builder.Collapse)
}
