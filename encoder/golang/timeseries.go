package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
)

func (encoder *Encoder) encodeTimeseries(panel sdk.Panel) jen.Code {
	return jen.Qual(packageImportPath+"/row", "WithTimeSeries").Call(
		jen.Lit(panel.Title),
	)
}
