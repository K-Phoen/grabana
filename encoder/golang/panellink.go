package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
)

func (encoder *Encoder) encodePanelLinks(links []sdk.Link) []jen.Code {
	linksStmts := make([]jen.Code, 0, len(links))

	for _, link := range links {
		linksStmts = append(linksStmts, encoder.encodePanelLink(link))
	}

	return linksStmts
}

func (encoder *Encoder) encodePanelLink(link sdk.Link) jen.Code {
	settings := []jen.Code{
		lit(link.Title),
		lit(*link.URL),
	}

	if link.TargetBlank != nil && *link.TargetBlank {
		settings = append(
			settings,
			qual("links", "OpenBlank").Call(),
		)
	}

	return qual("links", "New").Call(settings...)
}
