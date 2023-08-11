package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
)

func (encoder *Encoder) encodeCommonPanelProperties(panel sdk.Panel, grabanaPackage string) []jen.Code {
	settings := []jen.Code{
		lit(panel.Title),
	}

	if len(panel.Links) != 0 {
		settings = append(
			settings,
			qual(grabanaPackage, "Links").Call(
				encoder.encodePanelLinks(panel.Links)...,
			),
		)
	}

	span := panelSpan(panel)
	if span != 0 {
		settings = append(
			settings,
			qual(grabanaPackage, "Span").Call(lit(span)),
		)
	}

	if panel.Description != nil {
		settings = append(
			settings,
			qual(grabanaPackage, "Description").Call(lit(*panel.Description)),
		)
	}
	if panel.Height != nil {
		settings = append(
			settings,
			qual(grabanaPackage, "Height").Call(lit(*(panel.Height).(*string))),
		)
	}
	if panel.Transparent {
		settings = append(
			settings,
			qual(grabanaPackage, "Transparent").Call(lit(panel.Transparent)),
		)
	}
	if panel.Repeat != nil {
		settings = append(
			settings,
			qual(grabanaPackage, "Repeat").Call(lit(*panel.Repeat)),
		)
	}
	if panel.Datasource != nil && panel.Datasource.LegacyName != "" {
		settings = append(
			settings,
			qual(grabanaPackage, "DataSource").Call(lit(panel.Datasource.LegacyName)),
		)
	}

	return settings
}
