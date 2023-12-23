package golang

import (
	"strconv"

	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeHeatmap(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "heatmap")

	settings = append(
		settings,
		encoder.encodeTargets(panel.HeatmapPanel.Targets, "heatmap")...,
	)

	settings = append(settings, encoder.encodeHeatmapTooltip(panel.HeatmapPanel)...)
	settings = append(settings, encoder.encodeHeatmapYAxis(panel.HeatmapPanel))

	// DataFormat
	if panel.HeatmapPanel.DataFormat != "" {
		switch panel.HeatmapPanel.DataFormat {
		case "tsbuckets":
			settings = append(settings, heatmapQual("DataFormat").Call(heatmapQual("TimeSeriesBuckets")))
		case "timeseries":
			settings = append(settings, heatmapQual("DataFormat").Call(heatmapQual("TimeSeries")))
		default:
			encoder.logger.Warn("unknown heatmap data format", zap.String("data_format", panel.HeatmapPanel.DataFormat))
		}
	}

	// ShowZeroBuckets/HideZeroBuckets
	if !panel.HeatmapPanel.HideZeroBuckets {
		settings = append(settings, heatmapQual("ShowZeroBuckets").Call())
	} else {
		settings = append(settings, heatmapQual("HideZeroBuckets").Call())
	}

	// HighlightCards/NoHighlightCards
	if !panel.HeatmapPanel.HighlightCards {
		settings = append(settings, heatmapQual("NoHighlightCards").Call())
	} else {
		settings = append(settings, heatmapQual("HighlightCards").Call())
	}

	// ReverseYBuckets
	if panel.HeatmapPanel.ReverseYBuckets {
		settings = append(settings, heatmapQual("ReverseYBuckets").Call())
	}

	// HideXAxis
	if !panel.HeatmapPanel.XAxis.Show {
		settings = append(settings, heatmapQual("HideXAxis").Call())
	}

	// Legend()
	if !panel.HeatmapPanel.Legend.Show {
		settings = append(settings, heatmapQual("Legend").Call(heatmapQual("Hide")))
	}

	return qual("row", "WithHeatmap").MultiLineCall(
		settings...,
	)
}

func (encoder *Encoder) encodeHeatmapYAxis(panel *sdk.HeatmapPanel) jen.Code {
	var settings []jen.Code

	// Unit
	if panel.YAxis.Format != "" {
		settings = append(settings, heatmapAxisQual("Unit").Call(lit(panel.YAxis.Format)))
	}

	// Decimals
	if panel.YAxis.Decimals != nil {
		settings = append(settings, heatmapAxisQual("Decimals").Call(lit(*panel.YAxis.Decimals)))
	}

	// Min
	if panel.YAxis.Min != nil {
		asFloat, err := strconv.ParseFloat(*panel.YAxis.Min, 64)
		if err != nil {
			encoder.logger.Warn("could not parse heatmap YAxis min as float", zap.Error(err), zap.String("min", *panel.YAxis.Min))
		} else {
			settings = append(settings, heatmapAxisQual("Min").Call(lit(asFloat)))
		}
	}

	// Max
	if panel.YAxis.Max != nil {
		asFloat, err := strconv.ParseFloat(*panel.YAxis.Max, 64)
		if err != nil {
			encoder.logger.Warn("could not parse heatmap YAxis max as float", zap.Error(err), zap.String("min", *panel.YAxis.Max))
		} else {
			settings = append(settings, heatmapAxisQual("Max").Call(lit(asFloat)))
		}
	}

	if len(settings) == 0 {
		return nil
	}

	return heatmapQual("YAxis").Call(settings...)
}

func (encoder *Encoder) encodeHeatmapTooltip(panel *sdk.HeatmapPanel) []jen.Code {
	settings := []jen.Code{
		heatmapQual("HideTooltipHistogram").Call(lit(panel.TooltipDecimals)),
	}

	// HideTooltip
	if !panel.Tooltip.Show {
		settings = append(settings, heatmapQual("HideTooltip").Call())
	}

	// HideTooltipHistogram
	if !panel.Tooltip.ShowHistogram {
		settings = append(settings, heatmapQual("HideTooltipHistogram").Call())
	}

	return settings
}

func heatmapAxisQual(name string) *jen.Statement {
	return qual("heatmap/axis", name)
}

func heatmapQual(name string) *jen.Statement {
	return qual("heatmap", name)
}
