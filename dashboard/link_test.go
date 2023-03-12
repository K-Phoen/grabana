package dashboard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalLinkAsSdk(t *testing.T) {
	req := require.New(t)

	link := ExternalLink{
		Title:                 "my link",
		Description:           "super description",
		URL:                   "http://foo",
		Icon:                  IconBolt,
		IncludeTimeRange:      true,
		IncludeVariableValues: true,
		OpenInNewTab:          true,
	}
	sdkLink := link.asSdk()

	req.Equal("my link", sdkLink.Title)
	req.Equal("super description", *sdkLink.Tooltip)
	req.Equal("http://foo", *sdkLink.URL)
	req.Equal("bolt", *sdkLink.Icon)
	req.True(sdkLink.IncludeVars)
	req.True(*sdkLink.KeepTime)
	req.True(*sdkLink.TargetBlank)
	req.Equal("link", sdkLink.Type)
}

func TestExternalLinkIconHasDefault(t *testing.T) {
	req := require.New(t)

	link := ExternalLink{}
	sdkLink := link.asSdk()

	req.Equal("external", *sdkLink.Icon)
}

func TestDashboardLinkAsSdk(t *testing.T) {
	req := require.New(t)

	link := DashboardLink{
		Title:                 "my link",
		Tags:                  []string{"my-service"},
		AsDropdown:            true,
		IncludeTimeRange:      true,
		IncludeVariableValues: true,
		OpenInNewTab:          true,
	}
	sdkLink := link.asSdk()

	req.Equal("my link", sdkLink.Title)
	req.ElementsMatch([]string{"my-service"}, link.Tags)
	req.True(*sdkLink.AsDropdown)
	req.True(sdkLink.IncludeVars)
	req.True(*sdkLink.KeepTime)
	req.True(*sdkLink.TargetBlank)
	req.Equal("dashboards", sdkLink.Type)
}
