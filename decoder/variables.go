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

type DashboardVariable struct {
	Interval *VariableInterval `yaml:",omitempty"`
	Custom   *VariableCustom   `yaml:",omitempty"`
	Query    *VariableQuery    `yaml:",omitempty"`
	Const    *VariableConst    `yaml:",omitempty"`
}

func (variable *DashboardVariable) toOption() (dashboard.Option, error) {
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

type VariableInterval struct {
	Name    string
	Label   string
	Default string
	Values  []string `yaml:",flow"`
}

func (variable *VariableInterval) toOption() dashboard.Option {
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

type VariableCustom struct {
	Name      string
	Label     string
	Default   string
	ValuesMap map[string]string `yaml:"values_map"`
}

func (variable *VariableCustom) toOption() dashboard.Option {
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

type VariableConst struct {
	Name      string
	Label     string
	Default   string
	ValuesMap map[string]string `yaml:"values_map"`
}

func (variable *VariableConst) toOption() dashboard.Option {
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

type VariableQuery struct {
	Name  string
	Label string

	Datasource string
	Request    string

	IncludeAll bool `yaml:"include_all"`
	DefaultAll bool `yaml:"default_all"`
}

func (variable *VariableQuery) toOption() dashboard.Option {
	opts := []query.Option{
		query.Request(variable.Request),
	}

	if variable.Datasource != "" {
		opts = append(opts, query.DataSource(variable.Datasource))
	}
	if variable.Label != "" {
		opts = append(opts, query.Label(variable.Label))
	}
	if variable.IncludeAll {
		opts = append(opts, query.IncludeAll())
	}
	if variable.DefaultAll {
		opts = append(opts, query.DefaultAll())
	}

	return dashboard.VariableAsQuery(variable.Name, opts...)
}
