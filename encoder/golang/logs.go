package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) convertLogs(panel sdk.Panel) jen.Code {
	settings := []jen.Code{
		jen.Lit(panel.Title),
	}

	// TODO: span, height, links, targets, logs-specific options

	if panel.Description != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Description").Call(jen.Lit(*panel.Description)),
		)
	}
	if panel.Transparent {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Transparent").Call(jen.Lit(panel.Transparent)),
		)
	}
	if panel.Datasource != nil && panel.Datasource.LegacyName != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "DataSource").Call(jen.Lit(panel.Datasource.LegacyName)),
		)
	}

	for _, target := range panel.LogsPanel.Targets {
		settings = append(
			settings,
			encoder.convertLogsTarget(target),
		)
	}

	return jen.Qual(packageImportPath+"/row", "WithLogs").Call(
		settings...,
	)
}

func (encoder *Encoder) convertLogsTarget(target sdk.Target) jen.Code {
	settings := []jen.Code{
		jen.Lit(target.Expr),
	}

	if target.RefID != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/loki", "Ref").Call(jen.Lit(target.RefID)),
		)
	}
	if target.LegendFormat != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/loki", "Legend").Call(jen.Lit(target.LegendFormat)),
		)
	}
	if target.Hide {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/loki", "Hide").Call(),
		)
	}

	return jen.Qual(packageImportPath+"/logs", "WithLokiTarget").Call(
		settings...,
	)
}
