package decoder

import (
	"github.com/K-Phoen/grabana/dashboard"
)

type DashboardExternalLink struct {
	Title                 string
	URL                   string `yaml:"url"`
	Description           string `yaml:",omitempty"`
	Icon                  string `yaml:"icon,omitempty"`
	IncludeTimeRange      bool   `yaml:"include_time_range,omitempty"`
	IncludeVariableValues bool   `yaml:"include_variable_values,omitempty"`
	OpenInNewTab          bool   `yaml:"open_in_new_tab,omitempty"`
}

func (l DashboardExternalLink) toModel() dashboard.ExternalLink {
	return dashboard.ExternalLink{
		Title:                 l.Title,
		Description:           l.Description,
		URL:                   l.URL,
		Icon:                  dashboard.LinkIcon(l.Icon),
		IncludeTimeRange:      l.IncludeTimeRange,
		IncludeVariableValues: l.IncludeVariableValues,
		OpenInNewTab:          l.OpenInNewTab,
	}
}
