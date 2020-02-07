package grabana

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireJSON(t *testing.T, payload []byte) {
	var receiver map[string]interface{}
	if err := json.Unmarshal(payload, &receiver); err != nil {
		t.Fatalf("invalid json: %s", err)
	}
}

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

func TestDashboardCanBeMarshalledIntoJSON(t *testing.T) {
	req := require.New(t)

	builder := NewDashboardBuilder("Awesome dashboard")
	dashboardJSON, err := builder.MarshalJSON()

	req.NoError(err)
	requireJSON(t, dashboardJSON)
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

func TestDashboardCanHaveASharedCrossHair(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", SharedCrossHair())

	req.True(panel.board.SharedCrosshair)
}

func TestDashboardCanHaveADefaultTooltip(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", DefaultTooltip())

	req.False(panel.board.SharedCrosshair)
}

func TestDashboardCanBeAutoRefreshed(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", AutoRefresh("5s"))

	req.True(panel.board.Refresh.Flag)
	req.Equal("5s", panel.board.Refresh.Value)
}

func TestDashboardCanHaveTags(t *testing.T) {
	req := require.New(t)
	tags := []string{"generated", "grabana"}

	panel := NewDashboardBuilder("", Tags(tags))

	req.Len(panel.board.Tags, 2)
	req.ElementsMatch(tags, panel.board.Tags)
}

func TestDashboardCanHaveVariablesAsConstants(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", VariableAsConst("percentile"))

	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsCustom(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", VariableAsCustom("vX"))

	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsInterval(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", VariableAsInterval("interval"))

	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsQuery(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", VariableAsQuery("status"))

	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveRows(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", Row("Prometheus"))

	req.Len(panel.board.Rows, 1)
}

func TestDashboardCanHaveAnnotationsFromTags(t *testing.T) {
	req := require.New(t)

	panel := NewDashboardBuilder("", TagsAnnotation(TagAnnotation{}))

	req.Len(panel.board.Annotations.List, 1)
}
