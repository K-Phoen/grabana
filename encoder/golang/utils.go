package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
)

func panelSpan(panel sdk.Panel) float32 {
	span := panel.Span
	if span == 0 && panel.GridPos.H != nil {
		span = float32(*panel.GridPos.W / 2) // 24 units per row to 12
	}

	return span
}

func qual(pkg string, name string) *jen.Statement {
	return jen.Qual(packageImportPath+"/"+pkg, name)
}

func lit(v interface{}) *jen.Statement {
	return jen.Lit(v)
}
