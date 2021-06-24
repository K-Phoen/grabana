package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/datasource"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
)

var ErrVariableNotConfigured = fmt.Errorf("variable not configured")
var ErrInvalidHideValue = fmt.Errorf("invalid hide value. Valid values are: 'label', 'variable', empty")

type DashboardVariable struct {
	Interval   *VariableInterval   `yaml:",omitempty"`
	Custom     *VariableCustom     `yaml:",omitempty"`
	Query      *VariableQuery      `yaml:",omitempty"`
	Const      *VariableConst      `yaml:",omitempty"`
	Datasource *VariableDatasource `yaml:",omitempty"`
}

func (variable *DashboardVariable) toOption() (dashboard.Option, error) {
	if variable.Query != nil {
		return variable.Query.toOption()
	}
	if variable.Interval != nil {
		return variable.Interval.toOption()
	}
	if variable.Const != nil {
		return variable.Const.toOption()
	}
	if variable.Custom != nil {
		return variable.Custom.toOption()
	}
	if variable.Datasource != nil {
		return variable.Datasource.toOption()
	}

	return nil, ErrVariableNotConfigured
}

type VariableInterval struct {
	Name    string
	Label   string
	Default string
	Values  []string `yaml:",flow"`
	Hide    string   `yaml:",omitempty"`
}

func (variable *VariableInterval) toOption() (dashboard.Option, error) {
	opts := []interval.Option{
		interval.Values(variable.Values),
	}

	if variable.Label != "" {
		opts = append(opts, interval.Label(variable.Label))
	}
	if variable.Default != "" {
		opts = append(opts, interval.Default(variable.Default))
	}

	switch variable.Hide {
	case "":
		// Nothing to do
		break
	case "label":
		opts = append(opts, interval.HideLabel())
	case "variable":
		opts = append(opts, interval.Hide())
	default:
		return dashboard.VariableAsInterval(variable.Name), ErrInvalidHideValue
	}

	return dashboard.VariableAsInterval(variable.Name, opts...), nil
}

type VariableCustom struct {
	Name       string
	Label      string
	Default    string
	ValuesMap  map[string]string `yaml:"values_map"`
	IncludeAll bool              `yaml:"include_all"`
	AllValue   string            `yaml:"all_value,omitempty"`
	Hide       string            `yaml:",omitempty"`
}

func (variable *VariableCustom) toOption() (dashboard.Option, error) {
	opts := []custom.Option{
		custom.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, custom.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, custom.Label(variable.Label))
	}
	if variable.AllValue != "" {
		opts = append(opts, custom.AllValue(variable.AllValue))
	}
	if variable.IncludeAll {
		opts = append(opts, custom.IncludeAll())
	}

	switch variable.Hide {
	case "":
		// Nothing to do
		break
	case "label":
		opts = append(opts, custom.HideLabel())
	case "variable":
		opts = append(opts, custom.Hide())
	default:
		return dashboard.VariableAsCustom(variable.Name), ErrInvalidHideValue
	}

	return dashboard.VariableAsCustom(variable.Name, opts...), nil
}

type VariableConst struct {
	Name      string
	Label     string
	Default   string
	ValuesMap map[string]string `yaml:"values_map"`
	Hide      string            `yaml:",omitempty"`
}

func (variable *VariableConst) toOption() (dashboard.Option, error) {
	opts := []constant.Option{
		constant.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, constant.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, constant.Label(variable.Label))
	}

	switch variable.Hide {
	case "":
		// Nothing to do
		break
	case "label":
		opts = append(opts, constant.HideLabel())
	case "variable":
		opts = append(opts, constant.Hide())
	default:
		return dashboard.VariableAsConst(variable.Name), ErrInvalidHideValue
	}

	return dashboard.VariableAsConst(variable.Name, opts...), nil
}

type VariableQuery struct {
	Name  string
	Label string

	Datasource string
	Request    string

	Regex      string
	IncludeAll bool   `yaml:"include_all"`
	DefaultAll bool   `yaml:"default_all"`
	AllValue   string `yaml:"all_value,omitempty"`
	Hide       string `yaml:",omitempty"`
}

func (variable *VariableQuery) toOption() (dashboard.Option, error) {
	opts := []query.Option{
		query.Request(variable.Request),
	}

	if variable.Datasource != "" {
		opts = append(opts, query.DataSource(variable.Datasource))
	}
	if variable.Label != "" {
		opts = append(opts, query.Label(variable.Label))
	}
	if variable.Regex != "" {
		opts = append(opts, query.Regex(variable.Regex))
	}
	if variable.AllValue != "" {
		opts = append(opts, query.AllValue(variable.AllValue))
	}
	if variable.IncludeAll {
		opts = append(opts, query.IncludeAll())
	}
	if variable.DefaultAll {
		opts = append(opts, query.DefaultAll())
	}

	switch variable.Hide {
	case "":
		// Nothing to do
		break
	case "label":
		opts = append(opts, query.HideLabel())
	case "variable":
		opts = append(opts, query.Hide())
	default:
		return dashboard.VariableAsQuery(variable.Name), ErrInvalidHideValue
	}

	return dashboard.VariableAsQuery(variable.Name, opts...), nil
}

type VariableDatasource struct {
	Name  string
	Label string

	Type string

	Regex      string
	IncludeAll bool   `yaml:"include_all"`
	Hide       string `yaml:",omitempty"`
}

func (variable *VariableDatasource) toOption() (dashboard.Option, error) {
	opts := []datasource.Option{
		datasource.Type(variable.Type),
	}

	if variable.Label != "" {
		opts = append(opts, datasource.Label(variable.Label))
	}
	if variable.Regex != "" {
		opts = append(opts, datasource.Regex(variable.Regex))
	}
	if variable.IncludeAll {
		opts = append(opts, datasource.IncludeAll())
	}

	switch variable.Hide {
	case "":
		// Nothing to do
		break
	case "label":
		opts = append(opts, datasource.HideLabel())
	case "variable":
		opts = append(opts, datasource.Hide())
	default:
		return dashboard.VariableAsDatasource(variable.Name), ErrInvalidHideValue
	}

	return dashboard.VariableAsDatasource(variable.Name, opts...), nil
}
