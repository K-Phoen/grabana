package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeGraph(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "graph")

	settings = append(
		settings,
		encoder.encodeTargets(panel.GraphPanel.Targets, "graph")...,
	)

	settings = append(
		settings,
		encoder.encodeGraphLegend(panel.GraphPanel.Legend),
	)

	settings = append(
		settings,
		encoder.encodeGraphDraw(*panel.GraphPanel),
	)

	// TODO: SeriesOverride()

	settings = append(
		settings,
		encoder.encodeGraphAxis("XAxis", panel.GraphPanel.Xaxis),
	)

	if len(panel.GraphPanel.Yaxes) >= 1 {
		settings = append(
			settings,
			encoder.encodeGraphAxis("LeftYAxis", panel.GraphPanel.Yaxes[0]),
		)
	}
	if len(panel.GraphPanel.Yaxes) >= 2 {
		settings = append(
			settings,
			encoder.encodeGraphAxis("RightYAxis", panel.GraphPanel.Yaxes[1]),
		)
	}

	// Null
	if panel.GraphPanel.NullPointMode != "" {
		modes := map[string]string{
			"null as zero": "AsZero",
			"null":         "AsNull",
			"connected":    "Connected",
		}

		constName, ok := modes[panel.GraphPanel.NullPointMode]
		if !ok {
			encoder.logger.Warn("unknown null point mode in graph", zap.String("mode", panel.GraphPanel.NullPointMode))
		} else {
			settings = append(settings, graphQual("Null").Call(graphQual(constName)))
		}
	}

	// LineWidth
	if panel.GraphPanel.Linewidth != 0 {
		settings = append(settings, graphQual("LineWidth").Call(lit(panel.GraphPanel.Linewidth)))
	}
	// Fill
	if panel.GraphPanel.Fill != 0 {
		settings = append(settings, graphQual("Fill").Call(lit(panel.GraphPanel.Fill)))
	}
	// PointRadius
	if panel.GraphPanel.Pointradius != 0 {
		settings = append(settings, graphQual("PointRadius").Call(lit(panel.GraphPanel.Pointradius)))
	}
	// Staircase
	if panel.GraphPanel.SteppedLine {
		settings = append(settings, graphQual("Staircase").Call())
	}

	return qual("row", "WithGraph").MultiLineCall(
		settings...,
	)
}

func (encoder *Encoder) encodeGraphAxis(graphOptName string, axis sdk.Axis) jen.Code {
	var opts []jen.Code

	// Unit
	if axis.Format != "" {
		opts = append(opts, graphQual("Unit").Call(lit(axis.Format)))
	}
	// Hide
	if !axis.Show {
		opts = append(opts, graphQual("Hide").Call())
	}
	// LogBase
	if axis.LogBase != 0 {
		opts = append(opts, graphQual("LogBase").Call(lit(axis.LogBase)))
	}
	// Label
	if axis.Label != "" {
		opts = append(opts, graphQual("Label").Call(lit(axis.Label)))
	}
	// Min
	if axis.Min != nil {
		opts = append(opts, graphQual("Min").Call(lit(axis.Min.Value)))
	}
	// Max
	if axis.Max != nil {
		opts = append(opts, graphQual("Max").Call(lit(axis.Max.Value)))
	}

	if len(opts) == 0 {
		return nil
	}

	return graphQual(graphOptName).Call(opts...)
}

func (encoder *Encoder) encodeGraphLegend(legend sdk.Legend) jen.Code {
	var legendOpts []jen.Code

	// Hidden legend?
	if !legend.Show {
		legendOpts = append(legendOpts, graphQual("Hide"))
	}

	if legend.AlignAsTable {
		legendOpts = append(legendOpts, graphQual("AsTable"))
	}
	if legend.RightSide {
		legendOpts = append(legendOpts, graphQual("ToTheRight"))
	}
	if legend.Min {
		legendOpts = append(legendOpts, graphQual("Min"))
	}
	if legend.Max {
		legendOpts = append(legendOpts, graphQual("Max"))
	}
	if legend.Avg {
		legendOpts = append(legendOpts, graphQual("Avg"))
	}
	if legend.Current {
		legendOpts = append(legendOpts, graphQual("Current"))
	}
	if legend.Total {
		legendOpts = append(legendOpts, graphQual("Total"))
	}
	if legend.HideEmpty {
		legendOpts = append(legendOpts, graphQual("NoNullSeries"))
	}
	if legend.HideZero {
		legendOpts = append(legendOpts, graphQual("NoZeroSeries"))
	}

	if len(legendOpts) == 0 {
		return nil
	}

	return graphQual("Legend").Call(legendOpts...)
}

func (encoder *Encoder) encodeGraphDraw(panel sdk.GraphPanel) jen.Code {
	var opts []jen.Code

	if panel.Bars {
		opts = append(opts, graphQual("Bars"))
	}
	if panel.Lines {
		opts = append(opts, graphQual("Lines"))
	}
	if panel.Points {
		opts = append(opts, graphQual("Points"))
	}

	if len(opts) == 0 {
		return nil
	}

	return graphQual("Draw").Call(opts...)
}

func graphQual(name string) *jen.Statement {
	return qual("graph", name)
}
