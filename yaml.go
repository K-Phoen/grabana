package grabana

import (
	"fmt"
	"io"

	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
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

	Rows []dashboardRow
}

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

	for _, row := range dashboard.Rows {
		opt, err := row.toOption()
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
	case "const":
		return variable.asConst(), nil
	case "custom":
		return variable.asCustom(), nil
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
	if variable.Default != "" {
		opts = append(opts, interval.Default(variable.Default))
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

func (variable *dashboardVariable) asConst() DashboardBuilderOption {
	opts := []constant.Option{
		constant.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, constant.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, constant.Label(variable.Label))
	}

	return VariableAsConst(variable.Name, opts...)
}

func (variable *dashboardVariable) asCustom() DashboardBuilderOption {
	opts := []custom.Option{
		custom.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, custom.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, custom.Label(variable.Label))
	}

	return VariableAsCustom(variable.Name, opts...)
}

type dashboardRow struct {
	Name   string
	Panels []dashboardPanel
}

func (r dashboardRow) toOption() (DashboardBuilderOption, error) {
	opts := []row.Option{}

	for _, panel := range r.Panels {
		opt, err := panel.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return Row(r.Name, opts...), nil
}

type dashboardPanel struct {
	Graph *dashboardGraph
	Table *dashboardTable
}

func (panel dashboardPanel) toOption() (row.Option, error) {
	if panel.Graph != nil {
		return panel.Graph.toOption()
	}
	if panel.Table != nil {
		return panel.Table.toOption()
	}

	return nil, fmt.Errorf("panel not configured")
}

type dashboardGraph struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Targets    []target
}

func (graphPanel dashboardGraph) toOption() (row.Option, error) {
	opts := []graph.Option{}

	if graphPanel.Span != 0 {
		opts = append(opts, graph.Span(graphPanel.Span))
	}
	if graphPanel.Height != "" {
		opts = append(opts, graph.Height(graphPanel.Height))
	}
	if graphPanel.Datasource != "" {
		opts = append(opts, graph.DataSource(graphPanel.Datasource))
	}

	for _, t := range graphPanel.Targets {
		opt, err := graphPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithGraph(graphPanel.Title, opts...), nil
}

func (graphPanel *dashboardGraph) target(t target) (graph.Option, error) {
	if t.Prometheus != nil {
		return graph.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, fmt.Errorf("target not configured")
}

type target struct {
	Prometheus *prometheusTarget
}

type prometheusTarget struct {
	Query  string
	Legend string
	Ref    string
}

func (t prometheusTarget) toOptions() []prometheus.Option {
	var opts []prometheus.Option

	if t.Legend != "" {
		opts = append(opts, prometheus.Legend(t.Legend))
	}
	if t.Ref != "" {
		opts = append(opts, prometheus.Legend(t.Ref))
	}

	return opts
}

type dashboardTable struct {
	Title                  string
	Span                   float32
	Height                 string
	Datasource             string
	Targets                []target
	HiddenColumns          []string            `yaml:"hidden_columns"`
	TimeSeriesAggregations []table.Aggregation `yaml:"time_series_aggregations"`
}

func (tablePanel dashboardTable) toOption() (row.Option, error) {
	opts := []table.Option{}

	if tablePanel.Span != 0 {
		opts = append(opts, table.Span(tablePanel.Span))
	}
	if tablePanel.Height != "" {
		opts = append(opts, table.Height(tablePanel.Height))
	}
	if tablePanel.Datasource != "" {
		opts = append(opts, table.DataSource(tablePanel.Datasource))
	}

	for _, t := range tablePanel.Targets {
		opt, err := tablePanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, column := range tablePanel.HiddenColumns {
		opts = append(opts, table.HideColumn(column))
	}

	if len(tablePanel.TimeSeriesAggregations) != 0 {
		opts = append(opts, table.AsTimeSeriesAggregations(tablePanel.TimeSeriesAggregations))
	}

	return row.WithTable(tablePanel.Title, opts...), nil
}

func (tablePanel *dashboardTable) target(t target) (table.Option, error) {
	if t.Prometheus != nil {
		return table.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, fmt.Errorf("target not configured")
}
