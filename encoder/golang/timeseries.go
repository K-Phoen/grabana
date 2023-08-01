package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) encodeTimeseries(panel sdk.Panel) jen.Code {
	settings := []jen.Code{
		jen.Lit(panel.Title),
	}

	settings = append(
		settings,
		encoder.encodeCommonPanelProperties(panel, "timeseries")...,
	)

	return jen.Qual(packageImportPath+"/row", "WithTimeSeries").Call(
		settings...,
	)
}
