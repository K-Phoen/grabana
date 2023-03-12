package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDashboardLink(t *testing.T) {
	req := require.New(t)

	yamlLink := DashboardInternalLink{
		Title:                 "joe",
		Tags:                  []string{"my-service"},
		AsDropdown:            true,
		IncludeTimeRange:      true,
		IncludeVariableValues: true,
		OpenInNewTab:          true,
	}

	model := yamlLink.toModel()

	req.Equal("joe", model.Title)
	req.ElementsMatch([]string{"my-service"}, model.Tags)
	req.True(model.AsDropdown)
	req.True(model.IncludeTimeRange)
	req.True(model.IncludeVariableValues)
	req.True(model.OpenInNewTab)
}
