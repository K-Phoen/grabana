package dashboard

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
	IconColor  string   `yaml:"color"`
	Tags       []string `yaml:",flow"`
}

// Option represents an option that can be used to configure a
// dashboard.
type Option func(dashboard *Builder)

// TimezoneOption represents a possible value for the dashboard's timezone
// configuration.
type TimezoneOption string

// DefaultTimezone sets the dashboard's timezone to the default one used by
// Grafana.
const DefaultTimezone TimezoneOption = ""

// UTC sets the dashboard's timezone to UTC.
const UTC TimezoneOption = "utc"

// Browser sets the dashboard's timezone to the browser's one.
const Browser TimezoneOption = "browser"

// Builder is the main builder used to configure dashboards.
type Builder struct {
	board *sdk.Board
}

// New creates a new dashboard builder.
func New(title string, options ...Option) Builder {
	board := sdk.NewBoard(title)
	board.ID = 0

	builder := &Builder{board: board}

	for _, opt := range append(defaults(), options...) {
		opt(builder)
	}

	return *builder
}

func defaults() []Option {
	return []Option{
		defaultTimePicker(),
		Timezone(DefaultTimezone),
		Time("now-3h", "now"),
		SharedCrossHair(),
	}
}

func defaultTimePicker() Option {
	return func(builder *Builder) {
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
func (builder *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(builder.board)
}

// Internal.
func (builder *Builder) Internal() *sdk.Board {
	return builder.board
}

// VariableAsConst adds a templated variable, defined as a set of constant
// values.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsConst(name string, options ...constant.Option) Option {
	return func(builder *Builder) {
		templatedVar := constant.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// ID sets the ID used by the dashboard.
func ID(id uint) Option {
	return func(builder *Builder) {
		builder.board.ID = id
	}
}

// UID sets the UID used by the dashboard.
func UID(uid string) Option {
	return func(builder *Builder) {
		builder.board.UID = uid
	}
}

// VariableAsCustom adds a templated variable, defined as a set of custom
// values.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsCustom(name string, options ...custom.Option) Option {
	return func(builder *Builder) {
		templatedVar := custom.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// VariableAsInterval adds a templated variable, defined as an interval.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsInterval(name string, options ...interval.Option) Option {
	return func(builder *Builder) {
		templatedVar := interval.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// VariableAsQuery adds a templated variable, defined as a query.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsQuery(name string, options ...query.Option) Option {
	return func(builder *Builder) {
		templatedVar := query.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)
	}
}

// Row adds a row to the dashboard.
func Row(title string, options ...row.Option) Option {
	return func(builder *Builder) {
		row.New(builder.board, title, options...)
	}
}

// TagsAnnotation adds a new source of annotation for the dashboard.
func TagsAnnotation(annotation TagAnnotation) Option {
	return func(builder *Builder) {
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
func Editable() Option {
	return func(builder *Builder) {
		builder.board.Editable = true
	}
}

// ReadOnly marks the dashboard as non-editable.
func ReadOnly() Option {
	return func(builder *Builder) {
		builder.board.Editable = false
	}
}

// SharedCrossHair configures the graph tooltip to be shared across panels.
func SharedCrossHair() Option {
	return func(builder *Builder) {
		builder.board.SharedCrosshair = true
	}
}

// DefaultTooltip configures the graph tooltip NOT to be shared across panels.
func DefaultTooltip() Option {
	return func(builder *Builder) {
		builder.board.SharedCrosshair = false
	}
}

// Tags adds the given set of tags to the dashboard.
func Tags(tags []string) Option {
	return func(builder *Builder) {
		builder.board.Tags = tags
	}
}

// AutoRefresh defines the auto-refresh interval for the dashboard.
func AutoRefresh(interval string) Option {
	return func(builder *Builder) {
		builder.board.Refresh = &sdk.BoolString{Flag: true, Value: interval}
	}
}

// Time defines the default time range for the dashboard, e.g. from "now-6h" to
// "now".
func Time(from, to string) Option {
	return func(builder *Builder) {
		builder.board.Time = sdk.Time{From: from, To: to}
	}
}

// Timezone defines the default timezone for the dashboard, e.g. "utc".
func Timezone(timezone TimezoneOption) Option {
	return func(builder *Builder) {
		builder.board.Timezone = string(timezone)
	}
}
