package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeStat(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "stat")

	// TODO: ColorScheme

	settings = append(
		settings,
		encoder.encodeTargets(panel.StatPanel.Targets, "stat")...,
	)

	settings = append(
		settings,
		encoder.encodeStatFieldConfigDefaults(panel.StatPanel.FieldConfig.Defaults)...,
	)

	settings = append(
		settings,
		encoder.encodeStatOptions(panel.StatPanel.Options)...,
	)

	return qual("row", "WithStat").MultiLineCall(
		settings...,
	)
}

func (encoder *Encoder) encodeStatOptions(options sdk.StatOptions) []jen.Code {
	var settings []jen.Code

	// SparkLine
	if options.GraphMode == "area" {
		settings = append(settings, statQual("SparkLine").Call())
	}

	// ValueFontSize
	if options.Text != nil && options.Text.ValueSize != 0 {
		settings = append(settings, statQual("ValueFontSize").Call(lit(options.Text.ValueSize)))
	}

	// TitleFontSize
	if options.Text != nil && options.Text.TitleSize != 0 {
		settings = append(settings, statQual("TitleFontSize").Call(lit(options.Text.TitleSize)))
	}

	// Text
	if options.TextMode != "" {
		modes := map[string]string{
			"auto":           "TextAuto",
			"value":          "TextValue",
			"name":           "TextName",
			"value_and_name": "TextValueAndName",
			"none":           "TextNone",
		}

		constName, ok := modes[options.TextMode]
		if !ok {
			encoder.logger.Warn("unknown text mode in stat", zap.String("text_mode", options.TextMode))
		} else {
			settings = append(settings, statQual("Text").Call(statQual(constName)))
		}
	}

	// Orientation
	if options.Orientation != "" {
		modes := map[string]string{
			"":           "OrientationAuto",
			"horizontal": "OrientationHorizontal",
			"vertical":   "OrientationVertical",
		}

		constName, ok := modes[options.Orientation]
		if !ok {
			encoder.logger.Warn("unknown orientation in stat", zap.String("orientation", options.Orientation))
		} else {
			settings = append(settings, statQual("Orientation").Call(statQual(constName)))
		}
	}

	if options.ColorMode != "" {
		switch options.ColorMode {
		case "none":
			settings = append(settings, statQual("ColorNone").Call())
		case "value":
			settings = append(settings, statQual("ColorValue").Call())
		case "background":
			settings = append(settings, statQual("ColorBackground").Call())
		}
	}

	// ValueType
	if len(options.ReduceOptions.Calcs) == 1 {
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

		for _, sdkCalc := range options.ReduceOptions.Calcs {
			constName, ok := calcs[sdkCalc]
			if !ok {
				encoder.logger.Warn("unknown calculation in timeseries legend", zap.String("calc", sdkCalc))
				continue
			}

			settings = append(settings, statQual("ValueType").Call(statQual(constName)))
		}
	}

	return settings
}

func (encoder *Encoder) encodeStatFieldConfigDefaults(defaults sdk.FieldConfigDefaults) []jen.Code {
	var settings []jen.Code

	// unit
	if defaults.Unit != "" {
		settings = append(settings, statQual("Unit").Call(lit(defaults.Unit)))
	}
	// decimals
	if defaults.Decimals != nil {
		settings = append(settings, statQual("Decimals").Call(lit(*defaults.Decimals)))
	}
	// sparkline Y min
	if defaults.Min != nil {
		settings = append(settings, statQual("SparkLineYMin").Call(lit(*defaults.Min)))
	}
	// sparkline Y max
	if defaults.Max != nil {
		settings = append(settings, statQual("SparkLineYMax").Call(lit(*defaults.Max)))
	}
	// NoValue
	if defaults.NoValue != "" {
		settings = append(settings, statQual("NoValue").Call(lit(defaults.NoValue)))
	}

	// RelativeThresholds/AbsoluteThresholds
	if defaults.Thresholds.Mode != "" && len(defaults.Thresholds.Steps) != 0 {
		// TODO: thresholds
	}

	return settings
}

func statQual(name string) *jen.Statement {
	return qual("stat", name)
}
