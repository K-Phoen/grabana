package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) convertLogs(panel sdk.Panel) jen.Code {
	settings := []jen.Code{
		jen.Lit(panel.Title),
	}

	// TODO: links, logs-specific options

	span := panelSpan(panel)
	if span != 0 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Span").Call(jen.Lit(span)),
		)
	}

	if panel.Description != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Description").Call(jen.Lit(*panel.Description)),
		)
	}
	if panel.Height != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Height").Call(jen.Lit(*(panel.Height).(*string))),
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
		append(settings, encoder.encodeLogsVizualizationSettings(panel)...)...,
	)
}

func (encoder *Encoder) encodeLogsVizualizationSettings(panel sdk.Panel) []jen.Code {
	var settings []jen.Code

	// TODO: logs order, dedup strategy

	if panel.LogsPanel.Options.ShowTime {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Time").Call(),
		)
	}
	if panel.LogsPanel.Options.ShowLabels {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "UniqueLabels").Call(),
		)
	}
	if panel.LogsPanel.Options.ShowCommonLabels {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "CommonLabels").Call(),
		)
	}
	if panel.LogsPanel.Options.WrapLogMessage {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "WrapLines").Call(),
		)
	}
	if panel.LogsPanel.Options.PrettifyLogMessage {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "PrettifyJSON").Call(),
		)
	}
	if !panel.LogsPanel.Options.EnableLogDetails {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "HideLogDetails").Call(),
		)
	}

	return settings
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
