package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalLink(t *testing.T) {
	req := require.New(t)

	yamlLink := DashboardExternalLink{
		Title:                 "joe",
		URL:                   "http://foo",
		Description:           "bar",
		Icon:                  "cloud",
		IncludeTimeRange:      true,
		IncludeVariableValues: true,
		OpenInNewTab:          true,
	}

	model := yamlLink.toModel()

	req.Equal("joe", model.Title)
	req.Equal("http://foo", model.URL)
	req.Equal("bar", model.Description)
	req.Equal("cloud", string(model.Icon))
	req.True(model.IncludeTimeRange)
	req.True(model.IncludeVariableValues)
	req.True(model.OpenInNewTab)
}
