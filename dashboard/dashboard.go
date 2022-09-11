package dashboard

import (
	// We're not using it for security stuff, so it's fine.
	//nolint:gosec
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"

	"github.com/K-Phoen/grabana/alert"

	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/datasource"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
	"github.com/K-Phoen/sdk"
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
type Option func(dashboard *Builder) error

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
	board  *sdk.Board
	alerts []*alert.Alert
}

func NewFromBoard(board *sdk.Board) Builder {
	return Builder{
		board: board,
	}
}

// New creates a new dashboard builder.
func New(title string, options ...Option) (Builder, error) {
	board := sdk.NewBoard(title)
	board.ID = 0

	builder := &Builder{board: board}

	for _, opt := range append(defaults(), options...) {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
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
	return func(builder *Builder) error {
		builder.board.Timepicker = sdk.Timepicker{
			RefreshIntervals: []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"},
			TimeOptions:      []string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"},
		}

		return nil
	}
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// This method can be used to render the dashboard as JSON
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(builder.board)
}

// MarshalIndentJSON renders the dashboard as indented JSON
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *Builder) MarshalIndentJSON() ([]byte, error) {
	return json.MarshalIndent(builder.board, "", "  ")
}

// Alerts returns all the alerts defined in this dashboard.
func (builder *Builder) Alerts() []*alert.Alert {
	return builder.alerts
}

// Internal.
func (builder *Builder) Internal() *sdk.Board {
	return builder.board
}

// VariableAsConst adds a templated variable, defined as a set of constant
// values.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsConst(name string, options ...constant.Option) Option {
	return func(builder *Builder) error {
		templatedVar := constant.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)

		return nil
	}
}

// ID sets the ID used by the dashboard.
func ID(id uint) Option {
	return func(builder *Builder) error {
		builder.board.ID = id

		return nil
	}
}

// UID sets the UID used by the dashboard.
func UID(uid string) Option {
	return func(builder *Builder) error {
		validUID := uid

		if len(uid) > 40 {
			// We're not using it for security stuff, so it's fine.
			//nolint:gosec
			sha := sha1.Sum([]byte(uid))
			validUID = hex.EncodeToString(sha[:])
		}

		builder.board.UID = validUID

		return nil
	}
}

// VariableAsCustom adds a templated variable, defined as a set of custom
// values.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsCustom(name string, options ...custom.Option) Option {
	return func(builder *Builder) error {
		templatedVar := custom.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)

		return nil
	}
}

// VariableAsInterval adds a templated variable, defined as an interval.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsInterval(name string, options ...interval.Option) Option {
	return func(builder *Builder) error {
		templatedVar := interval.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)

		return nil
	}
}

// VariableAsQuery adds a templated variable, defined as a query.
// See https://grafana.com/docs/grafana/latest/reference/templating/#variable-types
func VariableAsQuery(name string, options ...query.Option) Option {
	return func(builder *Builder) error {
		templatedVar := query.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)

		return nil
	}
}

// VariableAsDatasource adds a templated variable, defined as a datasource.
// See https://grafana.com/docs/grafana/latest/variables/variable-types/add-data-source-variable/
func VariableAsDatasource(name string, options ...datasource.Option) Option {
	return func(builder *Builder) error {
		templatedVar := datasource.New(name, options...)

		builder.board.Templating.List = append(builder.board.Templating.List, templatedVar.Builder)

		return nil
	}
}

// ExternalLinks adds a dashboard-level link.
// See https://grafana.com/docs/grafana/latest/linking/dashboard-links/
func ExternalLinks(links ...ExternalLink) Option {
	return func(builder *Builder) error {
		for _, link := range links {
			builder.board.Links = append(builder.board.Links, link.asSdk())
		}

		return nil
	}
}

// Row adds a row to the dashboard.
func Row(title string, options ...row.Option) Option {
	return func(builder *Builder) error {
		r, err := row.New(builder.board, title, options...)
		if err != nil {
			return err
		}

		builder.alerts = append(builder.alerts, r.Alerts()...)

		return nil
	}
}

// TagsAnnotation adds a new source of annotation for the dashboard.
func TagsAnnotation(annotation TagAnnotation) Option {
	return func(builder *Builder) error {
		builder.board.Annotations.List = append(builder.board.Annotations.List, sdk.Annotation{
			Name:       annotation.Name,
			Datasource: &sdk.DatasourceRef{LegacyName: annotation.Datasource},
			IconColor:  annotation.IconColor,
			Enable:     true,
			Tags:       annotation.Tags,
			Type:       "tags",
		})

		return nil
	}
}

// Editable marks the dashboard as editable.
func Editable() Option {
	return func(builder *Builder) error {
		builder.board.Editable = true

		return nil
	}
}

// ReadOnly marks the dashboard as non-editable.
func ReadOnly() Option {
	return func(builder *Builder) error {
		builder.board.Editable = false

		return nil
	}
}

// SharedCrossHair configures the graph tooltip to be shared across panels.
func SharedCrossHair() Option {
	return func(builder *Builder) error {
		builder.board.SharedCrosshair = true

		return nil
	}
}

// DefaultTooltip configures the graph tooltip NOT to be shared across panels.
func DefaultTooltip() Option {
	return func(builder *Builder) error {
		builder.board.SharedCrosshair = false

		return nil
	}
}

// Tags adds the given set of tags to the dashboard.
func Tags(tags []string) Option {
	return func(builder *Builder) error {
		builder.board.Tags = tags

		return nil
	}
}

// AutoRefresh defines the auto-refresh interval for the dashboard.
func AutoRefresh(interval string) Option {
	return func(builder *Builder) error {
		builder.board.Refresh = &sdk.BoolString{Flag: true, Value: interval}

		return nil
	}
}

// Time defines the default time range for the dashboard, e.g. from "now-6h" to
// "now".
func Time(from, to string) Option {
	return func(builder *Builder) error {
		builder.board.Time = sdk.Time{From: from, To: to}

		return nil
	}
}

// Timezone defines the default timezone for the dashboard, e.g. "utc".
func Timezone(timezone TimezoneOption) Option {
	return func(builder *Builder) error {
		builder.board.Timezone = string(timezone)

		return nil
	}
}
