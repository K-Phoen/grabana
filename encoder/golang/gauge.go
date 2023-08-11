package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeGauge(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "gauge")

	settings = append(
		settings,
		encoder.encodeTargets(panel.GaugePanel.Targets, "gauge")...,
	)

	settings = append(settings, encoder.encodeGaugeSettings(panel)...)
	settings = append(settings, encoder.encodeGaugeOrientation(panel))
	settings = append(settings, encoder.encodeGaugeValueType(panel))

	// todo: thresholds

	return rowQual("WithGauge").MultiLineCall(settings...)
}

func (encoder *Encoder) encodeGaugeSettings(panel sdk.Panel) []jen.Code {
	settings := []jen.Code{
		gaugeQual("Unit").Call(lit(panel.GaugePanel.FieldConfig.Defaults.Unit)),
	}

	if panel.GaugePanel.FieldConfig.Defaults.Decimals != nil {
		settings = append(
			settings,
			gaugeQual("Decimals").Call(lit(*panel.GaugePanel.FieldConfig.Defaults.Decimals)),
		)
	}

	if panel.GaugePanel.Options.Text != nil {
		settings = append(
			settings,
			gaugeQual("ValueFontSize").Call(lit(panel.GaugePanel.Options.Text.ValueSize)),
		)
		settings = append(
			settings,
			gaugeQual("TitleFontSize").Call(lit(panel.GaugePanel.Options.Text.TitleSize)),
		)
	}

	return settings
}

func (encoder *Encoder) encodeGaugeOrientation(panel sdk.Panel) jen.Code {
	orientationConst := "OrientationAuto"
	switch panel.GaugePanel.Options.Orientation {
	case "":
		orientationConst = "OrientationAuto"
	case "auto":
		orientationConst = "OrientationAuto"
	case "horizontal":
		orientationConst = "OrientationHorizontal"
	case "vertical":
		orientationConst = "OrientationVertical"
	default:
		encoder.logger.Warn("unknown orientation, defaulting to auto", zap.String("orientation", panel.GaugePanel.Options.Orientation))
		orientationConst = "OrientationAuto"
	}

	return gaugeQual("Orientation").Call(gaugeQual(orientationConst))
}

func (encoder *Encoder) encodeGaugeValueType(panel sdk.Panel) jen.Code {
	valueTypeConst := "LastNonNull"

	if len(panel.GaugePanel.Options.ReduceOptions.Calcs) == 1 {
		valueType := panel.GaugePanel.Options.ReduceOptions.Calcs[0]

		switch valueType {
		case "first":
			valueTypeConst = "First"
		case "firstNotNull":
			valueTypeConst = "FirstNonNull"
		case "last":
			valueTypeConst = "Last"
		case "lastNotNull":
			valueTypeConst = "LastNonNull"

		case "min":
			valueTypeConst = "Min"
		case "max":
			valueTypeConst = "Max"
		case "mean":
			valueTypeConst = "Avg"

		case "count":
			valueTypeConst = "Count"
		case "sum":
			valueTypeConst = "Total"
		case "range":
			valueTypeConst = "Range"

		default:
			encoder.logger.Warn("unknown value type, defaulting to LastNonNull", zap.String("value type", valueType))
			valueTypeConst = "LastNonNull"
		}
	}

	return gaugeQual("ValueType").Call(gaugeQual(valueTypeConst))
}

func gaugeQual(name string) *jen.Statement {
	return jen.Qual(packageImportPath+"/gauge", name)
}
