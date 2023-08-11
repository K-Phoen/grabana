package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeTimeseries(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "timeseries")

	settings = append(
		settings,
		encoder.encodeTargets(panel.TimeseriesPanel.Targets, "timeseries")...,
	)

	settings = append(
		settings,
		encoder.encodeTimeseriesLegend(panel.TimeseriesPanel.Options.Legend),
	)

	settings = append(
		settings,
		encoder.encodeTimeseriesVizualization(panel)...,
	)

	return qual("row", "WithTimeSeries").MultiLineCall(
		settings...,
	)
}

func (encoder *Encoder) encodeTimeseriesLegend(legend sdk.TimeseriesLegendOptions) jen.Code {
	var legendOpts []jen.Code

	// Hidden legend?
	if legend.Show != nil && !*legend.Show {
		legendOpts = append(legendOpts, timeseriesQual("Hide"))
	} else {
		// Display mode
		switch legend.DisplayMode {
		case "list":
			legendOpts = append(legendOpts, timeseriesQual("AsList"))
		case "hidden":
			legendOpts = append(legendOpts, timeseriesQual("Hide"))
		default:
			legendOpts = append(legendOpts, timeseriesQual("AsTable"))
		}

		// Placement
		if legend.Placement == "right" {
			legendOpts = append(legendOpts, timeseriesQual("ToTheRight"))
		} else {
			legendOpts = append(legendOpts, timeseriesQual("Bottom"))
		}
	}

	// Automatic calculations
	calcs := map[string]string{
		"first":        "First",
		"firstNotNull": "FirstNonNull",
		"last":         "Last",
		"lastNotNull":  "LastNonNull",

		"min":  "Min",
		"max":  "Max",
		"mean": "Avg",

		"count": "Count",
		"sum":   "Total",
		"range": "Range",
	}

	for _, sdkCalc := range legend.Calcs {
		constName, ok := calcs[sdkCalc]
		if !ok {
			encoder.logger.Warn("unknown calculation in timeseries legend", zap.String("calc", sdkCalc))
			continue
		}

		legendOpts = append(legendOpts, timeseriesQual(constName))
	}

	return timeseriesQual("Legend").Call(legendOpts...)
}

func (encoder *Encoder) encodeTimeseriesVizualization(panel sdk.Panel) []jen.Code {
	var settings []jen.Code

	// Tooltip mode
	toolTipModeConst := "SingleSeries"
	switch panel.TimeseriesPanel.Options.Tooltip.Mode {
	case "none":
		toolTipModeConst = "NoSeries"
	case "multi":
		toolTipModeConst = "AllSeries"
	default:
		toolTipModeConst = "SingleSeries"
	}
	settings = append(
		settings,
		timeseriesQual("Tooltip").MultiLineCall(timeseriesQual(toolTipModeConst)),
	)

	return settings
}

func timeseriesQual(name string) *jen.Statement {
	return qual("timeseries", name)
}
