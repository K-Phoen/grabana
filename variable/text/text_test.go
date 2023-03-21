package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTextVariablesCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("filter")

	req.Equal("filter", panel.Builder.Name)
	req.Equal("filter", panel.Builder.Label)
	req.Equal("textbox", panel.Builder.Type)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("filter", Label("Filter"))

	req.Equal("filter", panel.Builder.Name)
	req.Equal("Filter", panel.Builder.Label)
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
