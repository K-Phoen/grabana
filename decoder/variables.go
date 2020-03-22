package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
)

type dashboardVariable struct {
	Type  string
	Name  string
	Label string

	// used for "interval", "const" and "custom"
	Default string

	// used for "interval"
	Values []string

	// used for "const" and "custom"
	ValuesMap map[string]string `yaml:"values_map"`

	// used for "query"
	Datasource string
	Request    string
}

func (variable *dashboardVariable) toOption() (dashboard.Option, error) {
	switch variable.Type {
	case "interval":
		return variable.asInterval(), nil
	case "query":
		return variable.asQuery(), nil
	case "const":
		return variable.asConst(), nil
	case "custom":
		return variable.asCustom(), nil
	}

	return nil, fmt.Errorf("unknown dashboard variable type '%s'", variable.Type)
}

func (variable *dashboardVariable) asInterval() dashboard.Option {
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

func (variable *dashboardVariable) asQuery() dashboard.Option {
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

func (variable *dashboardVariable) asConst() dashboard.Option {
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

func (variable *dashboardVariable) asCustom() dashboard.Option {
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
