package dashboard

import (
	"encoding/json"

	"github.com/K-Phoen/grabana/gen/dashboard/types"
)

type Option func(builder *Builder) error

type Builder struct {
	internal *types.Dashboard
}

func New(title string, options ...Option) (Builder, error) {
	dashboard := &types.Dashboard{
		Title: &title,
	}

	builder := &Builder{internal: dashboard}

	for _, opt := range options {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// This method can be used to render the dashboard as JSON
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(builder.internal)
}

// MarshalIndentJSON renders the dashboard as indented JSON
// which your configuration management tool of choice can then feed into
// Grafana's dashboard via its provisioning support.
// See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards
func (builder *Builder) MarshalIndentJSON() ([]byte, error) {
	return json.MarshalIndent(builder.internal, "", "  ")
}
func Id(id int64) Option {
	return func(builder *Builder) error {

		builder.internal.Id = &id

		return nil
	}
}

func Uid(uid string) Option {
	return func(builder *Builder) error {

		builder.internal.Uid = &uid

		return nil
	}
}

func Title(title string) Option {
	return func(builder *Builder) error {

		builder.internal.Title = &title

		return nil
	}
}

func Description(description string) Option {
	return func(builder *Builder) error {

		builder.internal.Description = &description

		return nil
	}
}

func Revision(revision int64) Option {
	return func(builder *Builder) error {

		builder.internal.Revision = &revision

		return nil
	}
}

func GnetId(gnetId string) Option {
	return func(builder *Builder) error {

		builder.internal.GnetId = &gnetId

		return nil
	}
}

func Tags(tags []string) Option {
	return func(builder *Builder) error {

		builder.internal.Tags = tags

		return nil
	}
}

func Style(style types.DashboardStyle) Option {
	return func(builder *Builder) error {

		builder.internal.Style = style

		return nil
	}
}

func Timezone(timezone string) Option {
	return func(builder *Builder) error {

		builder.internal.Timezone = &timezone

		return nil
	}
}

func Editable(editable bool) Option {
	return func(builder *Builder) error {

		builder.internal.Editable = editable

		return nil
	}
}

func GraphTooltip(graphTooltip types.DashboardCursorSync) Option {
	return func(builder *Builder) error {

		builder.internal.GraphTooltip = graphTooltip

		return nil
	}
}

func Time(time struct {
	From string `json:"from"`
	To   string `json:"to"`
}) Option {
	return func(builder *Builder) error {

		builder.internal.Time = time

		return nil
	}
}

func Timepicker(timepicker types.TimePicker) Option {
	return func(builder *Builder) error {

		builder.internal.Timepicker = &timepicker

		return nil
	}
}

func FiscalYearStartMonth(fiscalYearStartMonth uint8) Option {
	return func(builder *Builder) error {

		builder.internal.FiscalYearStartMonth = &fiscalYearStartMonth

		return nil
	}
}

func LiveNow(liveNow bool) Option {
	return func(builder *Builder) error {

		builder.internal.LiveNow = &liveNow

		return nil
	}
}

func WeekStart(weekStart string) Option {
	return func(builder *Builder) error {

		builder.internal.WeekStart = &weekStart

		return nil
	}
}

func Refresh(refresh types.StringOrBool) Option {
	return func(builder *Builder) error {

		builder.internal.Refresh = &refresh

		return nil
	}
}

func SchemaVersion(schemaVersion uint16) Option {
	return func(builder *Builder) error {

		builder.internal.SchemaVersion = schemaVersion

		return nil
	}
}

func Version(version uint32) Option {
	return func(builder *Builder) error {

		builder.internal.Version = &version

		return nil
	}
}

func Panels(panels []types.RowPanel) Option {
	return func(builder *Builder) error {

		builder.internal.Panels = panels

		return nil
	}
}

func Templating(templating types.DashboardTemplating) Option {
	return func(builder *Builder) error {

		builder.internal.Templating = &templating

		return nil
	}
}

func Annotations(annotations types.AnnotationContainer) Option {
	return func(builder *Builder) error {

		builder.internal.Annotations = &annotations

		return nil
	}
}

func Links(links []types.DashboardLink) Option {
	return func(builder *Builder) error {

		builder.internal.Links = links

		return nil
	}
}
