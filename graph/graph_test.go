package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGraphPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Graph panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Graph panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelCanBeMadeEditable(t *testing.T) {
	req := require.New(t)

	panel := New("", Editable())

	req.True(panel.Builder.Editable)
}

func TestGraphPanelCanBeMadeReadOnly(t *testing.T) {
	req := require.New(t)

	panel := New("", ReadOnly())

	req.False(panel.Builder.Editable)
}

func TestGraphPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestGraphPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}
