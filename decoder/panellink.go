package decoder

import (
	"github.com/K-Phoen/grabana/links"
)

type DashboardPanelLinks []DashboardPanelLink

func (collection DashboardPanelLinks) toModel() []links.Link {
	models := make([]links.Link, 0, len(collection))

	for _, link := range collection {
		models = append(models, link.toModel())
	}

	return models
}

type DashboardPanelLink struct {
	Title        string
	URL          string `yaml:"url"`
	OpenInNewTab bool   `yaml:"open_in_new_tab,omitempty"`
}

func (l DashboardPanelLink) toModel() links.Link {
	if l.OpenInNewTab {
		return links.New(l.Title, l.URL, links.OpenBlank())
	}

	return links.New(l.Title, l.URL)
}
