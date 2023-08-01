package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) encodeTimeseries(panel sdk.Panel) jen.Code {
	settings := []jen.Code{
		jen.Lit(panel.Title),
	}

	if len(panel.Links) != 0 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/timeseries", "Links").Call(
				encoder.encodePanelLinks(panel.Links)...,
			),
		)
	}
	return jen.Qual(packageImportPath+"/row", "WithTimeSeries").Call(
		settings...,
	)
}
