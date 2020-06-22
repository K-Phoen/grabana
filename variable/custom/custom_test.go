package custom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCustomVariablesCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("api_version")

	req.Equal("api_version", panel.Builder.Name)
	req.Equal("api_version", panel.Builder.Label)
	req.Equal("custom", panel.Builder.Type)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("custom var", Label("CustomVariable"))

	req.Equal("custom var", panel.Builder.Name)
	req.Equal("CustomVariable", panel.Builder.Label)
}

func TestValuesCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", Values(map[string]string{
		"v1": "v1-value",
		"v2": "v2-value",
	}))

	labels := make([]string, 0, 2)
	values := make([]string, 0, 2)

	for _, opt := range panel.Builder.Options {
		labels = append(labels, opt.Text)
		values = append(values, opt.Value)
	}

	req.Len(values, 2)
	req.ElementsMatch([]string{"v1", "v2"}, labels)
	req.ElementsMatch([]string{"v1-value", "v2-value"}, values)
}

func TestDefaultValueCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Default("99"))

	req.Equal("99", panel.Builder.Current.Text)
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

func TestAllValueCanBeOverriden(t *testing.T) {
	req := require.New(t)

	panel := New("", AllValue(".*"))

	req.Equal(".*", panel.Builder.AllValue)
}
