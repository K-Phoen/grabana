package dashboard

import "github.com/K-Phoen/sdk"

type LinkIcon string

const (
	IconExternal  LinkIcon = "external"
	IconDashboard LinkIcon = "dashboard"
	IconQuestion  LinkIcon = "question"
	IconInfo      LinkIcon = "info"
	IconBolt      LinkIcon = "bolt"
	IconDoc       LinkIcon = "doc"
	IconCloud     LinkIcon = "cloud"
)

// ExternalLink describes dashboard-level external link.
// See https://grafana.com/docs/grafana/latest/linking/dashboard-links/
type ExternalLink struct {
	Title                 string
	Description           string
	URL                   string
	Icon                  LinkIcon
	IncludeTimeRange      bool
	IncludeVariableValues bool
	OpenInNewTab          bool
}

func (link ExternalLink) asSdk() sdk.Link {
	falsePtr := false
	icon := string(IconExternal)
	if link.Icon != "" {
		icon = string(link.Icon)
	}

	return sdk.Link{
		Title:       link.Title,
		Tooltip:     &link.Description,
		URL:         &link.URL,
		AsDropdown:  &falsePtr,
		Icon:        &icon,
		IncludeVars: link.IncludeVariableValues,
		KeepTime:    &link.IncludeTimeRange,
		Tags:        make([]string, 0),
		TargetBlank: &link.OpenInNewTab,
		Type:        "link",
	}
}
