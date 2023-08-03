package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
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
		jen.Lit(variable.Name),
	}

	if variable.Label != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "Label").Call(jen.Lit(variable.Label)),
		)
	}
	if variable.IncludeAll {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "IncludeAll").Call(),
		)
	}
	if variable.Multi {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "Multiple").Call(),
		)
	}
	if variable.Hide == 1 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "HideLabel").Call(),
		)
	}
	if variable.Hide == 2 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "Hide").Call(),
		)
	}
	if variable.Regex != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "Regex").Call(jen.Lit(variable.Regex)),
		)
	}
	// TODO: eventually we should stop using legacy stuff... :|
	if variable.Datasource != nil && variable.Datasource.LegacyName != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "Datasource").Call(jen.Lit(variable.Datasource.LegacyName)),
		)
	}
	if variable.Current.Value == "$__all" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "DefaultAll").Call(),
		)
	}
	if variable.AllValue != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/query", "AllValue").Call(jen.Lit(variable.AllValue)),
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
			jen.Qual(packageImportPath+"/variable/query", "Refresh").Call(
				jen.Qual(packageImportPath+"/variable/query", refreshConstName),
			),
		)
	}
	if variable.Query != nil {
		if request, ok := variable.Query.(string); ok {
			settings = append(
				settings,
				jen.Qual(packageImportPath+"/variable/query", "Request").Call(jen.Lit(request)),
			)
		}
		if request, ok := variable.Query.(map[string]interface{}); ok {
			settings = append(
				settings,
				jen.Qual(packageImportPath+"/variable/query", "Request").Call(jen.Lit(request["query"].(string))),
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
			jen.Qual(packageImportPath+"/variable/query", "Sort").Call(
				jen.Qual(packageImportPath+"/variable/query", sortConstName),
			),
		)
	}

	return jen.Qual(packageImportPath+"/dashboard", "VariableAsQuery").Call(
		settings...,
	)
}

func (encoder *Encoder) encodeDatasourceVar(variable sdk.TemplateVar) jen.Code {
	settings := []jen.Code{
		jen.Lit(variable.Name),
	}

	if variable.Label != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "Label").Call(jen.Lit(variable.Label)),
		)
	}
	if variable.IncludeAll {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "IncludeAll").Call(),
		)
	}
	if variable.Multi {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "Multiple").Call(),
		)
	}
	if variable.Hide == 1 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "HideLabel").Call(),
		)
	}
	if variable.Hide == 2 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "Hide").Call(),
		)
	}
	if variable.Query != nil {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "Type").Call(jen.Lit(variable.Query.(string))),
		)
	}
	if variable.Regex != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/variable/datasource", "Regex").Call(jen.Lit(variable.Regex)),
		)
	}

	return jen.Qual(packageImportPath+"/dashboard", "VariableAsDatasource").Call(
		settings...,
	)
}
