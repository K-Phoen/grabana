package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/variable/interval"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/query"
)

var ErrVariableNotConfigured = fmt.Errorf("variable not configured")

type dashboardVariable struct {
	Interval *variableInterval
	Custom   *variableCustom
	Query    *variableQuery
	Const    *variableConst
}

func (variable *dashboardVariable) toOption() (dashboard.Option, error) {
	if variable.Query != nil {
		return variable.Query.toOption(), nil
	}
	if variable.Interval != nil {
		return variable.Interval.toOption(), nil
	}
	if variable.Const != nil {
		return variable.Const.toOption(), nil
	}
	if variable.Custom != nil {
		return variable.Custom.toOption(), nil
	}

	return nil, ErrVariableNotConfigured
}

type variableInterval struct {
	Name    string
	Label   string
	Default string
	Values  []string
}

func (variable *variableInterval) toOption() dashboard.Option {
	opts := []interval.Option{
		interval.Values(variable.Values),
	}

	if variable.Label != "" {
		opts = append(opts, interval.Label(variable.Label))
	}
	if variable.Default != "" {
		opts = append(opts, interval.Default(variable.Default))
	}

	return dashboard.VariableAsInterval(variable.Name, opts...)
}

type variableCustom struct {
	Name      string
	Label     string
	Default   string
	ValuesMap map[string]string `yaml:"values_map"`
}

func (variable *variableCustom) toOption() dashboard.Option {
	opts := []custom.Option{
		custom.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, custom.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, custom.Label(variable.Label))
	}

	return dashboard.VariableAsCustom(variable.Name, opts...)
}

type variableConst struct {
	Name      string
	Label     string
	Default   string
	ValuesMap map[string]string `yaml:"values_map"`
}

func (variable *variableConst) toOption() dashboard.Option {
	opts := []constant.Option{
		constant.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, constant.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, constant.Label(variable.Label))
	}

	return dashboard.VariableAsConst(variable.Name, opts...)
}

type variableQuery struct {
	Name  string
	Label string

	Datasource string
	Request    string
}

func (variable *variableQuery) toOption() dashboard.Option {
	opts := []query.Option{
		query.Request(variable.Request),
	}

	if variable.Datasource != "" {
		opts = append(opts, query.DataSource(variable.Datasource))
	}
	if variable.Label != "" {
		opts = append(opts, query.Label(variable.Label))
	}

	return dashboard.VariableAsQuery(variable.Name, opts...)
}
