package decoder

import (
	"github.com/K-Phoen/grabana/dashboard"
)

type DashboardInternalLink struct {
	Title                 string   `yaml:"title"`
	Tags                  []string `yaml:"tags"`
	AsDropdown            bool     `yaml:"as_dropdown,omitempty"`
	IncludeTimeRange      bool     `yaml:"include_time_range,omitempty"`
	IncludeVariableValues bool     `yaml:"include_variable_values,omitempty"`
	OpenInNewTab          bool     `yaml:"open_in_new_tab,omitempty"`
}

func (l DashboardInternalLink) toModel() dashboard.DashboardLink {
	return dashboard.DashboardLink{
		Title:                 l.Title,
		Tags:                  l.Tags,
		AsDropdown:            l.AsDropdown,
		IncludeTimeRange:      l.IncludeTimeRange,
		IncludeVariableValues: l.IncludeVariableValues,
		OpenInNewTab:          l.OpenInNewTab,
	}
}
