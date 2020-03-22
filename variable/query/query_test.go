package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQueryVariablesCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("query")

	req.Equal("query", panel.Builder.Name)
	req.Equal("query", panel.Builder.Label)
	req.Equal("query", panel.Builder.Type)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("query var", Label("QueryVariable"))

	req.Equal("query var", panel.Builder.Name)
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

func TestAnAllValuesCanBeTheDefault(t *testing.T) {
	req := require.New(t)

	panel := New("", DefaultAll())

	req.Equal("All", panel.Builder.Current.Text)
	req.Equal("$__all", panel.Builder.Current.Value)
}

func TestValuesCanBeFilteredByRegex(t *testing.T) {
	req := require.New(t)
	regex := "^4\\d+$"

	panel := New("", Regex(regex))

	req.Equal(regex, panel.Builder.Regex)
}

func TestValuesRefreshTimeCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Refresh(TimeChange))

	req.Equal(int64(TimeChange), *panel.Builder.Refresh.Value)
}

func TestValuesSortOrderCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Sort(AlphabeticalNoCaseDesc))

	req.Equal(int(AlphabeticalNoCaseDesc), panel.Builder.Sort)
}

func TestDataSourceCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestRequestCanBeSet(t *testing.T) {
	req := require.New(t)
	request := "label_values(prometheus_http_requests_total, code)"

	panel := New("", Request(request))

	req.Equal(request, panel.Builder.Query)
}
