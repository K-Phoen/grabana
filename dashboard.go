package grabana

import (
	"encoding/json"

	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
	"github.com/grafana-tools/sdk"
)

// TagAnnotation describes an annotation represented as a Tag.
// See https://grafana.com/docs/grafana/latest/reference/annotations/#query-by-tag
type TagAnnotation struct {
	Name       string
	Datasource string
	IconColor  string
	Tags       []string
}

// DashboardBuilderOption represents an option that can be used to configure a
// dashboard.
type DashboardBuilderOption func(dashboard *DashboardBuilder)

// DashboardBuilder is the main builder used to configure dashboards.
type DashboardBuilder struct {
	board *sdk.Board
}

// NewDashboardBuilder creates a new dashboard builder.
func NewDashboardBuilder(title string, options ...DashboardBuilderOption) DashboardBuilder {
	board := sdk.NewBoard(title)
	board.ID = 0
	board.Timezone = ""

	builder := &DashboardBuilder{board: board}

	for _, opt := range append(dashboardDefaults(), options...) {
		opt(builder)
	}

	return *builder
}

func dashboardDefaults() []DashboardBuilderOption {
	return []DashboardBuilderOption{
		defaultTimePicker(),
		defaultTime(),
		SharedCrossHair(),
	}
}

func defaultTime() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Time = sdk.Time{
			From: "now-3h",
			To:   "now",
		}
	}
}

func defaultTimePicker() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Timepicker = sdk.Timepicker{
			RefreshIntervals: []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"},
			TimeOptions:      []string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"},
		}
	}
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// This method can be used to render the dashboard into a JSON file
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *DashboardBuilder) MarshalJSON() ([]byte, error) {
	return json.Marshal(builder.board)
}

// VariableAsConst adds a templated variable, defined as a set of constant
// values.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsConst(name string, options ...constant.Option) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		templatedVar := constant.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// VariableAsCustom adds a templated variable, defined as a set of custom
// values.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsCustom(name string, options ...custom.Option) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		templatedVar := custom.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// VariableAsInterval adds a templated variable, defined as an interval.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsInterval(name string, options ...interval.Option) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		templatedVar := interval.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// VariableAsQuery adds a templated variable, defined as a query.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsQuery(name string, options ...query.Option) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		templatedVar := query.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// Row adds a row to the dashboard.
func Row(title string, options ...row.Option) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		row.New(builder.board, title, options...)
	}
}

// TagsAnnotation adds a new source of annotation for the dashboard.
func TagsAnnotation(annotation TagAnnotation) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Annotations.List = append(builder.board.Annotations.List, sdk.Annotation{
			Name:       annotation.Name,
			Datasource: &annotation.Datasource,
			IconColor:  annotation.IconColor,
			Enable:     true,
			Tags:       annotation.Tags,
			Type:       "tags",
		})
	}
}

// Editable marks the dashboard as editable.
func Editable() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Editable = true
	}
}

// ReadOnly marks the dashboard as non-editable.
func ReadOnly() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Editable = false
	}
}

// SharedCrossHair configures the graph tooltip to be shared across panels.
func SharedCrossHair() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.SharedCrosshair = true
	}
}

// DefaultTooltip configures the graph tooltip NOT to be shared across panels.
func DefaultTooltip() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.SharedCrosshair = false
	}
}

// Tags adds the given set of tags to the dashboard.
func Tags(tags []string) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Tags = tags
	}
}

// AutoRefresh defines the auto-refresh interval for the dashboard.
func AutoRefresh(interval string) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Refresh = &sdk.BoolString{Flag: true, Value: interval}
	}
}
