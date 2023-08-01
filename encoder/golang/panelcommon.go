package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) encodeCommonPanelProperties(panel sdk.Panel, grabanaPackage string) []jen.Code {
	settings := []jen.Code{
		jen.Lit(panel.Title),
	}

	if len(panel.Links) != 0 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "Links").Call(
				encoder.encodePanelLinks(panel.Links)...,
			),
		)
	}

	span := panelSpan(panel)
	if span != 0 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "Span").Call(jen.Lit(span)),
		)
	}

	if panel.Description != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "Description").Call(jen.Lit(*panel.Description)),
		)
	}
	if panel.Height != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "Height").Call(jen.Lit(*(panel.Height).(*string))),
		)
	}
	if panel.Transparent {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "Transparent").Call(jen.Lit(panel.Transparent)),
		)
	}
	if panel.Repeat != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "Repeat").Call(jen.Lit(*panel.Repeat)),
		)
	}
	if panel.Datasource != nil && panel.Datasource.LegacyName != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/"+grabanaPackage, "DataSource").Call(jen.Lit(panel.Datasource.LegacyName)),
		)
	}

	return settings
}
