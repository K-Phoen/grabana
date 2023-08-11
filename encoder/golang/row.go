package golang

import (
	"github.com/K-Phoen/jennifer/jen"
)

type RowIR struct {
	Title     string
	RepeatFor *string
	Collapsed bool
	Panels    []jen.Code
}

func (encoder *Encoder) encodeRow(row RowIR) *jen.Statement {
	rowSettings := []jen.Code{
		lit(row.Title),
	}

	if row.RepeatFor != nil {
		rowSettings = append(rowSettings, rowQual("RepeatFor").Call(lit(*row.RepeatFor)))
	}

	if row.Collapsed {
		rowSettings = append(rowSettings, rowQual("Collapse").Call())
	}

	for _, panel := range row.Panels {
		rowSettings = append(rowSettings, panel)
	}

	return dashboardQual("Row").MultiLineCall(rowSettings...)
}

func rowQual(name string) *jen.Statement {
	return qual("row", name)
}
