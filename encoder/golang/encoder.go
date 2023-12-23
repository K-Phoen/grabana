package golang

import (
	"bytes"

	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

const packageImportPath = "github.com/K-Phoen/grabana"
const sdkImportPath = "github.com/K-Phoen/sdk"

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
		lit(dashboard.Title),
	}

	// TODO: links, annotations
	dashboardStatements = append(dashboardStatements, encoder.encodeGeneralSettings(dashboard)...)
	dashboardStatements = append(dashboardStatements, encoder.encodeVariables(dashboard.Templating.List)...)
	dashboardStatements = append(dashboardStatements, encoder.encodePanels(dashboard)...)

	file := jen.NewFile("main")
	file.Func().Id("main").Params().Block(
		jen.Id("builder, err").Op(":=").Qual(packageImportPath+"/dashboard", "New").MultiLineCall(
			dashboardStatements...,
		),
	)

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (encoder *Encoder) encodeGeneralSettings(dashboard sdk.Board) []jen.Code {
	var settings []jen.Code

	if dashboard.UID != "" {
		settings = append(settings, dashboardQual("UID").Call(lit(dashboard.UID)))
	}

	if dashboard.Slug != "" {
		settings = append(settings, dashboardQual("Slug").Call(lit(dashboard.Slug)))
	}

	if dashboard.SharedCrosshair {
		settings = append(settings, dashboardQual("SharedCrossHair").Call())
	}

	if !dashboard.Editable {
		settings = append(settings, dashboardQual("ReadOnly").Call())
	}

	if dashboard.Refresh != nil && dashboard.Refresh.Value != "" {
		settings = append(settings, dashboardQual("AutoRefresh").Call(lit(dashboard.Refresh.Value)))
	}

	if len(dashboard.Tags) != 0 {
		tagsLiterals := Map(dashboard.Tags, func(item string) jen.Code {
			return lit(item)
		})
		settings = append(settings, dashboardQual("Tags").Call(jen.List(tagsLiterals...)))
	}

	// TODO: timezone

	settings = append(
		settings,
		dashboardQual("Time").Call(
			lit(dashboard.Time.From),
			lit(dashboard.Time.To),
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
	case "graph":
		return encoder.encodeGraph(panel), true
	case "gauge":
		return encoder.encodeGauge(panel), true
	case "stat":
		return encoder.encodeStat(panel), true
	case "text":
		return encoder.encodeText(panel), true
	case "heatmap":
		return encoder.encodeHeatmap(panel), true
	/*
		case "singlestat":
			return encoder.encodeSingleStat(panel), true
		case "table":
			return encoder.encodeTable(panel), true
	*/
	default:
		encoder.logger.Warn("unhandled panel type: skipped", zap.String("type", panel.Type), zap.String("title", panel.Title))
	}

	return nil, false
}

func dashboardQual(name string) *jen.Statement {
	return qual("dashboard", name)
}
