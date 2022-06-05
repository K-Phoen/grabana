package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPanelLink(t *testing.T) {
	req := require.New(t)

	yamlLink := DashboardPanelLink{
		Title:        "joe",
		URL:          "http://foo",
		OpenInNewTab: true,
	}

	model := yamlLink.toModel()

	req.Equal("joe", model.Builder.Title)
	req.Equal("http://foo", *model.Builder.URL)
	req.True(*model.Builder.TargetBlank)
}

func TestPanelLinks(t *testing.T) {
	req := require.New(t)

	yamlLinks := DashboardPanelLinks{
		{Title: "joe", URL: "http://foo"},
		{Title: "bar", URL: "http://baz"},
	}

	model := yamlLinks.toModel()

	req.Len(model, 2)
}
