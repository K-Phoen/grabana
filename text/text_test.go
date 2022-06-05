package text

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/stretchr/testify/require"
)

func TestNewTextPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("Text panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("Text panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
	req.Empty(panel.Builder.TextPanel.Content)
	req.Empty(panel.Builder.TextPanel.Mode)
}

func TestTextPanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestTextPanelsCanBeHTML(t *testing.T) {
	req := require.New(t)
	content := "<b>lala</b>"

	panel, err := New("", HTML(content))

	req.NoError(err)
	req.Equal(content, panel.Builder.TextPanel.Content)
	req.Equal("html", panel.Builder.TextPanel.Mode)
}

func TestTextPanelsCanBeMarkdown(t *testing.T) {
	req := require.New(t)
	content := "*lala*"

	panel, err := New("", Markdown(content))

	req.NoError(err)
	req.Equal(content, panel.Builder.TextPanel.Content)
	req.Equal("markdown", panel.Builder.TextPanel.Mode)
}

func TestTextPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestInvalidTextPanelWidthIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(16))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestTextPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestTextPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestTextPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}
