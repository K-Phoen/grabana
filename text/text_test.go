package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTextPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Text panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Text panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
	req.Empty(panel.Builder.TextPanel.Content)
	req.Empty(panel.Builder.TextPanel.Mode)
}

func TestTextPanelsCanBeHTML(t *testing.T) {
	req := require.New(t)
	content := "<b>lala</b>"

	panel := New("", HTML(content))

	req.Equal(content, panel.Builder.TextPanel.Content)
	req.Equal("html", panel.Builder.TextPanel.Mode)
}

func TestTextPanelsCanBeMarkdown(t *testing.T) {
	req := require.New(t)
	content := "*lala*"

	panel := New("", Markdown(content))

	req.Equal(content, panel.Builder.TextPanel.Content)
	req.Equal("markdown", panel.Builder.TextPanel.Mode)
}

func TestTextPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestTextPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestTextPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestTextPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}
