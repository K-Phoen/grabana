package golang

import (
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) convertLogs(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "logs")

	for _, target := range panel.LogsPanel.Targets {
		settings = append(
			settings,
			encoder.encodeLogsTarget(target),
		)
	}

	return rowQual("WithLogs").MultiLineCall(
		append(settings, encoder.encodeLogsVizualizationSettings(panel)...)...,
	)
}

func (encoder *Encoder) encodeLogsVizualizationSettings(panel sdk.Panel) []jen.Code {
	var settings []jen.Code

	if panel.LogsPanel.Options.ShowTime {
		settings = append(settings, logsQual("Time").Call())
	}
	if panel.LogsPanel.Options.ShowLabels {
		settings = append(settings, logsQual("UniqueLabels").Call())
	}
	if panel.LogsPanel.Options.ShowCommonLabels {
		settings = append(settings, logsQual("CommonLabels").Call())
	}
	if panel.LogsPanel.Options.WrapLogMessage {
		settings = append(settings, logsQual("WrapLines").Call())
	}
	if panel.LogsPanel.Options.PrettifyLogMessage {
		settings = append(settings, logsQual("PrettifyJSON").Call())
	}
	if !panel.LogsPanel.Options.EnableLogDetails {
		settings = append(settings, logsQual("HideLogDetails").Call())
	}
	if panel.LogsPanel.Options.SortOrder != "" {
		settings = append(
			settings,
			logsQual("Order").Call(encoder.encodeLogsSortOrder(panel.LogsPanel.Options.SortOrder)),
		)
	}
	if panel.LogsPanel.Options.DedupStrategy != "" {
		settings = append(
			settings,
			logsQual("Deduplication").Call(encoder.encodeLogsDedupStrategy(panel.LogsPanel.Options.DedupStrategy)),
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

	return logsQual(constantName)
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

	return logsQual(constantName)
}

func (encoder *Encoder) encodeLogsTarget(target sdk.Target) jen.Code {
	settings := []jen.Code{
		lit(target.Expr),
	}

	if target.RefID != "" {
		settings = append(settings, qual("target/loki", "Ref").Call(lit(target.RefID)))
	}
	if target.LegendFormat != "" {
		settings = append(settings, qual("target/loki", "Legend").Call(lit(target.LegendFormat)))
	}
	if target.Hide {
		settings = append(settings, qual("target/loki", "Hide").Call())
	}

	return logsQual("WithLokiTarget").Call(settings...)
}

func logsQual(name string) *jen.Statement {
	return qual("logs", name)
}
