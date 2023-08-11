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

	return jen.Qual(packageImportPath+"/row", "WithTimeSeries").MultiLineCall(
		settings...,
	)
}

func (encoder *Encoder) encodeTimeseriesLegend(legend sdk.TimeseriesLegendOptions) jen.Code {
	var legendOpts []jen.Code

	// Hidden legend?
	if legend.Show != nil && !*legend.Show {
		legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", "Hide"))
	} else {
		// Display mode
		switch legend.DisplayMode {
		case "list":
			legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", "AsList"))
		case "hidden":
			legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", "Hide"))
		default:
			legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", "AsTable"))
		}

		// Placement
		if legend.Placement == "right" {
			legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", "ToTheRight"))
		} else {
			legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", "Bottom"))
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

		legendOpts = append(legendOpts, jen.Qual(packageImportPath+"/timeseries", constName))
	}

	return jen.Qual(packageImportPath+"/timeseries", "Legend").Call(
		legendOpts...,
	)
}
