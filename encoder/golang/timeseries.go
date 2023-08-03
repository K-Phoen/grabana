package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) encodeTimeseries(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "timeseries")

	settings = append(
		settings,
		encoder.encodeTargets(panel.TimeseriesPanel.Targets, "timeseries")...,
	)

	return jen.Qual(packageImportPath+"/row", "WithTimeSeries").Call(
		settings...,
	)
}
