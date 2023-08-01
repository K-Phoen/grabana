package golang

import (
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
	"go.uber.org/zap"
)

func (encoder *Encoder) convertLogs(panel sdk.Panel) jen.Code {
	settings := []jen.Code{
		jen.Lit(panel.Title),
	}

	settings = append(
		settings,
		encoder.encodeCommonPanelProperties(panel, "logs")...,
	)

	for _, target := range panel.LogsPanel.Targets {
		settings = append(
			settings,
			encoder.encodeLogsTarget(target),
		)
	}

	return jen.Qual(packageImportPath+"/row", "WithLogs").Call(
		append(settings, encoder.encodeLogsVizualizationSettings(panel)...)...,
	)
}

func (encoder *Encoder) encodeLogsVizualizationSettings(panel sdk.Panel) []jen.Code {
	var settings []jen.Code

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
	if panel.LogsPanel.Options.SortOrder != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Order").Call(
				encoder.encodeLogsSortOrder(panel.LogsPanel.Options.SortOrder),
			),
		)
	}
	if panel.LogsPanel.Options.DedupStrategy != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/logs", "Deduplication").Call(
				encoder.encodeLogsDedupStrategy(panel.LogsPanel.Options.DedupStrategy),
			),
		)
	}

	return settings
}

func (encoder *Encoder) encodeLogsDedupStrategy(sdkDedupStrategy string) jen.Code {
	var constantName string

	switch sdkDedupStrategy {
	case "none":
		constantName = "None"
	case "exact":
		constantName = "Exact"
	case "numbers":
		constantName = "Numbers"
	case "signature":
		constantName = "Signature"
	default:
		encoder.logger.Warn("unhandled logs dedup strategy: using none as default", zap.String("strategy", sdkDedupStrategy))
		constantName = "none"
	}

	return jen.Qual(packageImportPath+"/logs", constantName)
}

func (encoder *Encoder) encodeLogsSortOrder(sdkSortOrder string) jen.Code {
	var constantName string

	switch sdkSortOrder {
	case string(logs.Asc):
		constantName = "Asc"
	case string(logs.Desc):
		constantName = "Desc"
	default:
		encoder.logger.Warn("unhandled sort order: using desc as default", zap.String("order", sdkSortOrder))
		constantName = "Desc"
	}

	return jen.Qual(packageImportPath+"/logs", constantName)
}

func (encoder *Encoder) encodeLogsTarget(target sdk.Target) jen.Code {
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
