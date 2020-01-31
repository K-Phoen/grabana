package constant

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConstantsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("percentile")

	req.Equal("percentile", panel.Builder.Name)
	req.Equal("percentile", panel.Builder.Label)
	req.Equal("constant", panel.Builder.Type)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", WithLabel("Constant"))

	req.Equal("const", panel.Builder.Name)
	req.Equal("Constant", panel.Builder.Label)
}

func TestValuesCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", WithValues([]Value{
		{Text: "90th", Value: "90"},
		{Text: "95th", Value: "95"},
		{Text: "99th", Value: "99"},
	}))

	req.Len(panel.Builder.Options, 3)
	req.Equal("90th", panel.Builder.Options[0].Text)
	req.Equal("90", panel.Builder.Options[0].Value)
}

func TestDefaultValueCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", WithDefault("99"))

	req.Equal("99", panel.Builder.Current.Text)
}
