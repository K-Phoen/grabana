package decoder

import (
	"github.com/K-Phoen/grabana/dashboard"
)

type DashboardExternalLink struct {
	Title                 string
	Type                  string   `yaml:"type,omitempty"`
	Tags                  []string `yaml:"tags,omitempty"`
	URL                   string   `yaml:"url,omitempty"`
	Description           string   `yaml:",omitempty"`
	Icon                  string   `yaml:"icon,omitempty"`
	IncludeTimeRange      bool     `yaml:"include_time_range,omitempty"`
	IncludeVariableValues bool     `yaml:"include_variable_values,omitempty"`
	OpenInNewTab          bool     `yaml:"open_in_new_tab,omitempty"`
}

func (l DashboardExternalLink) toModel() dashboard.ExternalLink {
	return dashboard.ExternalLink{
		Title:                 l.Title,
		Type:                  l.Type,
		Tags:                  l.Tags,
		Description:           l.Description,
		URL:                   l.URL,
		Icon:                  dashboard.LinkIcon(l.Icon),
		IncludeTimeRange:      l.IncludeTimeRange,
		IncludeVariableValues: l.IncludeVariableValues,
		OpenInNewTab:          l.OpenInNewTab,
	}
}
