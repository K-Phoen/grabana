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
	fieldConfig := panel.TimeseriesPanel.FieldConfig

	settings := []jen.Code{
		timeseriesQual("FillOpacity").Call(lit(fieldConfig.Defaults.Custom.FillOpacity)),
		timeseriesQual("PointSize").Call(lit(fieldConfig.Defaults.Custom.PointSize)),
		timeseriesQual("LineWidth").Call(lit(fieldConfig.Defaults.Custom.LineWidth)),
	}

	// Line interpolation mode
	if fieldConfig.Defaults.Custom.DrawStyle == "line" {
		lineInterpolationConst := "Linear"

		switch fieldConfig.Defaults.Custom.LineInterpolation {
		case "smooth":
			lineInterpolationConst = "Smooth"
		case "linear":
			lineInterpolationConst = "Linear"
		case "stepBefore":
			lineInterpolationConst = "StepBefore"
		case "stepAfter":
			lineInterpolationConst = "StepAfter"
		default:
			encoder.logger.Warn("invalid line interpolation mode, defaulting to Linear", zap.String("interpolation_mode", fieldConfig.Defaults.Custom.LineInterpolation))
			lineInterpolationConst = "Linear"
		}

		// don't generate code for the default
		if lineInterpolationConst != "Linear" {
			settings = append(
				settings,
				timeseriesQual("Lines").Call(timeseriesQual(lineInterpolationConst)),
			)
		}
	}

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
	// don't generate code for the default
	if toolTipModeConst != "SingleSeries" {
		settings = append(
			settings,
			timeseriesQual("Tooltip").Call(timeseriesQual(toolTipModeConst)),
		)
	}

	// Gradient mode
	gradientModeConst := "Opacity"
	switch fieldConfig.Defaults.Custom.GradientMode {
	case "none":
		gradientModeConst = "NoGradient"
	case "hue":
		gradientModeConst = "Hue"
	case "scheme":
		gradientModeConst = "Scheme"
	default:
		gradientModeConst = "Opacity"
	}
	// don't generate code for the default
	if gradientModeConst != "Opacity" {
		settings = append(
			settings,
			timeseriesQual("GradientMode").Call(timeseriesQual(gradientModeConst)),
		)
	}

	// Stacking mode
	stackingModeConst := "Unstacked"
	switch fieldConfig.Defaults.Custom.Stacking.Mode {
	case "none":
		stackingModeConst = "Unstacked"
	case "normal":
		stackingModeConst = "NormalStack"
	case "percent":
		stackingModeConst = "PercentStack"
	default:
		stackingModeConst = "Unstacked"
	}
	// don't generate code for the default
	if stackingModeConst != "Unstacked" {
		settings = append(
			settings,
			timeseriesQual("Stack").Call(timeseriesQual(stackingModeConst)),
		)
	}

	return settings
}

func timeseriesQual(name string) *jen.Statement {
	return qual("timeseries", name)
}
