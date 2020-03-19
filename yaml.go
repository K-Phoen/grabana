package grabana

import (
	"fmt"
	"io"

	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
	"gopkg.in/yaml.v2"
)

func UnmarshalYAML(input io.Reader) (DashboardBuilder, error) {
	decoder := yaml.NewDecoder(input)
	decoder.SetStrict(true)

	parsed := &dashboardYaml{}
	if err := decoder.Decode(parsed); err != nil {
		return DashboardBuilder{}, err
	}

	return parsed.toDashboardBuilder()
}

type dashboardYaml struct {
	Title           string
	Editable        bool
	SharedCrosshair bool `yaml:"shared_crosshair"`
	Tags            []string
	AutoRefresh     string `yaml:"auto_refresh"`

	TagsAnnotation []TagAnnotation `yaml:"tags_annotations"`
	Variables      []dashboardVariable
}

type dashboardVariable struct {
	Type  string
	Name  string
	Label string

	// used for "interval" and "const"
	Values []string

	// used for "query"
	Datasource string
	Request    string
}

func (dashboard *dashboardYaml) toDashboardBuilder() (DashboardBuilder, error) {
	emptyDashboard := DashboardBuilder{}
	opts := []DashboardBuilderOption{
		dashboard.editable(),
		dashboard.sharedCrossHair(),
	}

	if len(dashboard.Tags) != 0 {
		opts = append(opts, Tags(dashboard.Tags))
	}

	if dashboard.AutoRefresh != "" {
		opts = append(opts, AutoRefresh(dashboard.AutoRefresh))
	}

	for _, tagAnnotation := range dashboard.TagsAnnotation {
		opts = append(opts, TagsAnnotation(tagAnnotation))
	}

	for _, variable := range dashboard.Variables {
		opt, err := variable.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	return NewDashboardBuilder(dashboard.Title, opts...), nil
}

func (dashboard *dashboardYaml) sharedCrossHair() DashboardBuilderOption {
	if dashboard.SharedCrosshair {
		return SharedCrossHair()
	}

	return DefaultTooltip()
}

func (dashboard *dashboardYaml) editable() DashboardBuilderOption {
	if dashboard.Editable {
		return Editable()
	}

	return ReadOnly()
}

func (variable *dashboardVariable) toOption() (DashboardBuilderOption, error) {
	switch variable.Type {
	case "interval":
		return variable.asInterval(), nil
	case "query":
		return variable.asQuery(), nil
	}

	return nil, fmt.Errorf("unknown dashboard variable type '%s'", variable.Type)
}

func (variable *dashboardVariable) asInterval() DashboardBuilderOption {
	opts := []interval.Option{
		interval.Values(variable.Values),
	}

	if variable.Label != "" {
		opts = append(opts, interval.Label(variable.Label))
	}

	return VariableAsInterval(variable.Name, opts...)
}

func (variable *dashboardVariable) asQuery() DashboardBuilderOption {
	opts := []query.Option{
		query.Request(variable.Request),
	}

	if variable.Datasource != "" {
		opts = append(opts, query.DataSource(variable.Datasource))
	}
	if variable.Label != "" {
		opts = append(opts, query.Label(variable.Label))
	}

	return VariableAsQuery(variable.Name, opts...)
}
