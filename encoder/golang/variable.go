package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeVariables(variables []sdk.TemplateVar) []jen.Code {
	varStmts := make([]jen.Code, 0, len(variables))

	for _, variable := range variables {
		varStmts = append(varStmts, encoder.encodeVariable(variable))
	}

	return varStmts
}

func (encoder *Encoder) encodeVariable(variable sdk.TemplateVar) jen.Code {
	switch variable.Type {
	/*
		case "interval":
			encoder.encodeIntervalVar(variable)
		case "custom":
			encoder.encodeCustomVar(variable)
		case "const":
			encoder.encodeConstVar(variable)

	*/
	case "query":
		return encoder.encodeQueryVar(variable)
	case "datasource":
		return encoder.encodeDatasourceVar(variable)
		/*
			case "textbox":
				encoder.encodeTextVar(variable)
		*/
	default:
		encoder.logger.Warn("unhandled variable type found: skipped", zap.String("type", variable.Type), zap.String("name", variable.Name))
	}

	return nil
}

func (encoder *Encoder) encodeQueryVar(variable sdk.TemplateVar) jen.Code {
	settings := []jen.Code{
		lit(variable.Name),
	}

	if variable.Label != "" {
		settings = append(
			settings,
			qual("variable/query", "Label").Call(lit(variable.Label)),
		)
	}
	if variable.IncludeAll {
		settings = append(settings, qual("variable/query", "IncludeAll").Call())
	}
	if variable.Multi {
		settings = append(settings, qual("variable/query", "Multiple").Call())
	}
	if variable.Hide == 1 {
		settings = append(settings, qual("variable/query", "HideLabel").Call())
	}
	if variable.Hide == 2 {
		settings = append(settings, qual("variable/query", "Hide").Call())
	}
	if variable.Regex != "" {
		settings = append(settings, qual("variable/query", "Regex").Call(lit(variable.Regex)))
	}
	// TODO: eventually we should stop using legacy stuff... :|
	if variable.Datasource != nil && variable.Datasource.LegacyName != "" {
		settings = append(
			settings,
			qual("variable/query", "Datasource").Call(lit(variable.Datasource.LegacyName)),
		)
	}
	if variable.Current.Value == "$__all" {
		settings = append(settings, qual("variable/query", "DefaultAll").Call())
	}
	if variable.AllValue != "" {
		settings = append(
			settings,
			qual("variable/query", "AllValue").Call(lit(variable.AllValue)),
		)
	}
	if variable.Refresh.Value != nil {
		refreshConstName := "DashboardLoad"
		switch *variable.Refresh.Value {
		case 1:
			refreshConstName = "DashboardLoad"
		case 2:
			refreshConstName = "TimeChange"
		default:
			encoder.logger.Warn("invalid refresh value for variable: using DashboardLoad", zap.Int64("refresh", *variable.Refresh.Value))
		}

		settings = append(
			settings,
			qual("variable/query", "Refresh").Call(
				qual("variable/query", refreshConstName),
			),
		)
	}
	if variable.Query != nil {
		if request, ok := variable.Query.(string); ok {
			settings = append(
				settings,
				qual("variable/query", "Request").Call(lit(request)),
			)
		}
		if request, ok := variable.Query.(map[string]interface{}); ok {
			settings = append(
				settings,
				qual("variable/query", "Request").Call(lit(request["query"].(string))),
			)
		}
	}
	if variable.Sort != 0 {
		sortConstName := "None"
		switch variable.Sort {
		case 1:
			sortConstName = "AlphabeticalAsc"
		case 2:
			sortConstName = "AlphabeticalDesc"
		case 3:
			sortConstName = "NumericalAsc"
		case 4:
			sortConstName = "NumericalDesc"
		case 5:
			sortConstName = "AlphabeticalNoCaseAsc"
		case 6:
			sortConstName = "AlphabeticalNoCaseDesc"
		default:
			encoder.logger.Warn("invalid sort value for variable: using None", zap.Int("sort", variable.Sort))
		}

		settings = append(
			settings,
			qual("variable/query", "Sort").Call(
				qual("variable/query", sortConstName),
			),
		)
	}

	return dashboardQual("VariableAsQuery").MultiLineCall(settings...)
}

func (encoder *Encoder) encodeDatasourceVar(variable sdk.TemplateVar) jen.Code {
	settings := []jen.Code{
		lit(variable.Name),
	}

	if variable.Label != "" {
		settings = append(
			settings,
			qual("variable/datasource", "Label").Call(lit(variable.Label)),
		)
	}
	if variable.IncludeAll {
		settings = append(settings, qual("variable/datasource", "IncludeAll").Call())
	}
	if variable.Multi {
		settings = append(settings, qual("variable/datasource", "Multiple").Call())
	}
	if variable.Hide == 1 {
		settings = append(settings, qual("variable/datasource", "HideLabel").Call())
	}
	if variable.Hide == 2 {
		settings = append(settings, qual("variable/datasource", "Hide").Call())
	}
	if variable.Query != nil {
		settings = append(settings, qual("variable/datasource", "Type").Call(lit(variable.Query.(string))))
	}
	if variable.Regex != "" {
		settings = append(settings, qual("variable/datasource", "Regex").Call(lit(variable.Regex)))
	}

	return dashboardQual("VariableAsDatasource").MultiLineCall(settings...)
}
