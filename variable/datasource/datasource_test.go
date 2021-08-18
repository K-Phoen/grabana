package datasource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDatasourceVariablesCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("source")

	req.Equal("source", panel.Builder.Name)
	req.Equal("source", panel.Builder.Label)
	req.Equal("datasource", panel.Builder.Type)
	req.Equal(DashboardLoad, *panel.Builder.Refresh.Value)
	req.Equal(true, panel.Builder.Refresh.Flag)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("datasource var", Label("QueryVariable"))

	req.Equal("datasource var", panel.Builder.Name)
	req.Equal("QueryVariable", panel.Builder.Label)
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

func TestMultipleVariablesCanBeSelected(t *testing.T) {
	req := require.New(t)

	panel := New("", Multi())

	req.True(panel.Builder.Multi)
}

func TestAnOptionToIncludeAllCanBeAdded(t *testing.T) {
	req := require.New(t)

	panel := New("", IncludeAll())

	req.True(panel.Builder.IncludeAll)
}

func TestValuesCanBeFilteredByRegex(t *testing.T) {
	req := require.New(t)
	regex := "^4\\d+$"

	panel := New("", Regex(regex))

	req.Equal(regex, panel.Builder.Regex)
}

func TestDataSourceTypeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Type("prometheus"))

	req.Equal("prometheus", panel.Builder.Query)
}
