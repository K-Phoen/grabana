package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
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
	if panel.RepeatDirection != nil {
		directions := map[string]string{
			"v": "RepeatDirectionVertical",
			"h": "RepeatDirectionHorizontal",
		}

		constName, ok := directions[string(*panel.RepeatDirection)]
		if !ok {
			encoder.logger.Warn("unknown panel repeat direction", zap.String("direction", string(*panel.RepeatDirection)))
		} else {
			settings = append(settings,
				qual(grabanaPackage, "RepeatDirection").Call(jen.Qual(sdkImportPath, constName)),
			)
		}

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
