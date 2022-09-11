package dashboard

import (
	"encoding/json"
	"testing"

	"github.com/K-Phoen/grabana/variable/datasource"
	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func requireJSON(t *testing.T, payload []byte) {
	var receiver map[string]interface{}
	if err := json.Unmarshal(payload, &receiver); err != nil {
		t.Fatalf("invalid json: %s", err)
	}
}

func TestNewDashboardsCanBeCreatedFromSdkBoard(t *testing.T) {
	req := require.New(t)

	testBoard := sdk.NewBoard("My dashboard")
	testBoard.Timezone = "UTC"

	builder := NewFromBoard(testBoard)

	req.Equal("My dashboard", builder.board.Title)
	req.Equal("UTC", builder.board.Timezone)
}

func TestNewDashboardsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("My dashboard")

	req.NoError(err)
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

	builder, err := New("Awesome dashboard")
	req.NoError(err)

	dashboardJSON, err := builder.MarshalJSON()

	req.NoError(err)
	requireJSON(t, dashboardJSON)
}

func TestDashboardCanBeMarshalledIntoIndentedJSON(t *testing.T) {
	req := require.New(t)

	builder, err := New("Awesome dashboard")
	req.NoError(err)

	dashboardJSON, err := builder.MarshalIndentJSON()

	req.NoError(err)
	requireJSON(t, dashboardJSON)
}

func TestDashboardCanBeMadeEditable(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Editable())

	req.NoError(err)
	req.True(panel.board.Editable)
}

func TestDashboardIDCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ID(42))

	req.NoError(err)
	req.Equal(uint(42), panel.board.ID)
}

func TestDashboardUIDCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", UID("foo"))

	req.NoError(err)
	req.Equal("foo", panel.board.UID)
}

func TestDashboardUIDWillBeHashedWhenTooLongForGrafana(t *testing.T) {
	req := require.New(t)

	originalUID := "this-uid-is-more-than-forty-characters-and-grafana-does-not-like-it"
	panel, err := New("", UID(originalUID))

	req.NoError(err)
	req.NotEqual(originalUID, panel.board.UID)
	req.Len(panel.board.UID, 40)
}

func TestDashboardCanBeMadeReadOnly(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ReadOnly())

	req.NoError(err)
	req.False(panel.board.Editable)
}

func TestDashboardCanHaveASharedCrossHair(t *testing.T) {
	req := require.New(t)

	panel, err := New("", SharedCrossHair())

	req.NoError(err)
	req.True(panel.board.SharedCrosshair)
}

func TestDashboardCanHaveADefaultTooltip(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DefaultTooltip())

	req.NoError(err)
	req.False(panel.board.SharedCrosshair)
}

func TestDashboardCanBeAutoRefreshed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", AutoRefresh("5s"))

	req.NoError(err)
	req.True(panel.board.Refresh.Flag)
	req.Equal("5s", panel.board.Refresh.Value)
}

func TestDashboardCanHaveTime(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Time("now-6h", "now"))

	req.NoError(err)
	req.Equal("now-6h", panel.board.Time.From)
	req.Equal("now", panel.board.Time.To)
}

func TestDashboardCanHaveTimezone(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Timezone(UTC))

	req.NoError(err)
	req.Equal("utc", panel.board.Timezone)
}

func TestDashboardCanHaveTags(t *testing.T) {
	req := require.New(t)
	tags := []string{"generated", "grabana"}

	panel, err := New("", Tags(tags))

	req.NoError(err)
	req.Len(panel.board.Tags, 2)
	req.ElementsMatch(tags, panel.board.Tags)
}

func TestDashboardCanHaveVariablesAsConstants(t *testing.T) {
	req := require.New(t)

	panel, err := New("", VariableAsConst("percentile"))

	req.NoError(err)
	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsCustom(t *testing.T) {
	req := require.New(t)

	panel, err := New("", VariableAsCustom("vX"))

	req.NoError(err)
	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsInterval(t *testing.T) {
	req := require.New(t)

	panel, err := New("", VariableAsInterval("interval"))

	req.NoError(err)
	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsQuery(t *testing.T) {
	req := require.New(t)

	panel, err := New("", VariableAsQuery("status"))

	req.NoError(err)
	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveVariablesAsDatasource(t *testing.T) {
	req := require.New(t)

	panel, err := New("", VariableAsDatasource("source", datasource.Type("prometheus")))

	req.NoError(err)
	req.Len(panel.board.Templating.List, 1)
}

func TestDashboardCanHaveRows(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Row("Prometheus"))

	req.NoError(err)
	req.Len(panel.board.Rows, 1)
}

func TestDashboardCanHaveAnnotationsFromTags(t *testing.T) {
	req := require.New(t)

	panel, err := New("", TagsAnnotation(TagAnnotation{}))

	req.NoError(err)
	req.Len(panel.board.Annotations.List, 1)
}

func TestDashboardCanHaveExternalLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", ExternalLinks(ExternalLink{}, ExternalLink{}))

	req.NoError(err)
	req.Len(panel.board.Links, 2)
}
