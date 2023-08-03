package golang

import (
	"bytes"

	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
	"go.uber.org/zap"
)

const packageImportPath = "github.com/K-Phoen/grabana"

type Encoder struct {
	logger *zap.Logger
}

func NewEncoder(logger *zap.Logger) *Encoder {
	return &Encoder{
		logger: logger,
	}
}

func (encoder *Encoder) EncodeDashboard(dashboard sdk.Board) (string, error) {
	dashboardStatements := []jen.Code{
		jen.Lit(dashboard.Title),
	}

	dashboardStatements = append(dashboardStatements, encoder.encodeGeneralSettings(dashboard)...)
	dashboardStatements = append(dashboardStatements, encoder.encodeVariables(dashboard.Templating.List)...)
	dashboardStatements = append(dashboardStatements, encoder.encodePanels(dashboard)...)

	file := jen.NewFile("main")
	file.Func().Id("main").Params().Block(
		jen.Id("builder, err").Op(":=").Qual(packageImportPath+"/dashboard", "New").Call(
			dashboardStatements...,
		),
	)

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return "", nil
	}

	return buffer.String(), nil
}

func (encoder *Encoder) encodeGeneralSettings(dashboard sdk.Board) []jen.Code {
	var settings []jen.Code

	if dashboard.UID != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/dashboard", "UID").Call(jen.Lit(dashboard.UID)),
		)
	}

	if dashboard.Slug != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/dashboard", "Slug").Call(jen.Lit(dashboard.Slug)),
		)
	}

	if dashboard.SharedCrosshair {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/dashboard", "SharedCrossHair").Call(),
		)
	}

	if !dashboard.Editable {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/dashboard", "ReadOnly").Call(),
		)
	}

	if dashboard.Refresh != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/dashboard", "AutoRefresh").Call(jen.Lit(dashboard.Refresh.Value)),
		)
	}

	// TODO: timezone, tags

	settings = append(
		settings,
		jen.Qual(packageImportPath+"/dashboard", "Time").Call(
			jen.Lit(dashboard.Time.From),
			jen.Lit(dashboard.Time.To),
		),
	)

	return settings
}

func (encoder *Encoder) encodePanels(dashboard sdk.Board) []jen.Code {
	var currentRow *RowIR
	var convertedRows []jen.Code

	for _, panel := range dashboard.Panels {
		if panel.Type == "row" {
			if currentRow != nil {
				convertedRows = append(convertedRows, encoder.encodeRow(*currentRow))
			}

			currentRow = &RowIR{
				Title:     panel.Title,
				RepeatFor: panel.Repeat,
			}

			if panel.RowPanel != nil && panel.RowPanel.Collapsed {
				currentRow.Collapsed = true
			}
			continue
		}

		if currentRow == nil {
			currentRow = &RowIR{
				Title: "Overview",
			}
		}

		convertedPanel, ok := encoder.encodeDataPanel(*panel)
		if ok {
			currentRow.Panels = append(currentRow.Panels, convertedPanel)
		}
	}

	if currentRow != nil {
		convertedRows = append(convertedRows, encoder.encodeRow(*currentRow))
	}

	return convertedRows
}

func (encoder *Encoder) encodeDataPanel(panel sdk.Panel) (jen.Code, bool) {
	switch panel.Type {
	case "logs":
		return encoder.convertLogs(panel), true
	case "timeseries":
		return encoder.encodeTimeseries(panel), true
	/*
		case "graph":
			return converter.convertGraph(panel), true
		case "heatmap":
			return converter.convertHeatmap(panel), true
		case "singlestat":
			return converter.convertSingleStat(panel), true
		case "stat":
			return converter.convertStat(panel), true
		case "table":
			return converter.convertTable(panel), true
		case "text":
			return converter.convertText(panel), true
		case "gauge":
			return converter.convertGauge(panel), true
	*/
	default:
		encoder.logger.Warn("unhandled panel type: skipped", zap.String("type", panel.Type), zap.String("title", panel.Title))
	}

	return nil, false

}
