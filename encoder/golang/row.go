package golang

import (
	"github.com/dave/jennifer/jen"
)

type RowIR struct {
	Title     string
	RepeatFor *string
	Collapsed bool
	Panels    []jen.Code
}

func (encoder *Encoder) encodeRow(row RowIR) *jen.Statement {
	rowSettings := []jen.Code{
		jen.Lit(row.Title),
	}

	if row.RepeatFor != nil {
		rowSettings = append(
			rowSettings,
			jen.Qual(packageImportPath+"/row", "RepeatFor").Call(jen.Lit(*row.RepeatFor)),
		)
	}

	if row.Collapsed {
		rowSettings = append(
			rowSettings,
			jen.Qual(packageImportPath+"/row", "Collapse").Call(),
		)
	}

	for _, panel := range row.Panels {
		rowSettings = append(rowSettings, panel)
	}

	return jen.Qual(packageImportPath+"/dashboard", "Row").Call(
		rowSettings...,
	)
}
