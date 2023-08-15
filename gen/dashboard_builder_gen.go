package dashboard

import (
	"encoding/json"
	"errors"
)

type Option func(builder *Builder) error

type Builder struct {
	internal *Dashboard
}

func New(title string, options ...Option) (Builder, error) {
	dashboard := &Dashboard{
		Title: title,
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
		if !(id <= 9223372036854775807) {
			return errors.New("id must be <= 9223372036854775807")
		}

		builder.internal.Id = &id

		return nil
	}
}

func Uid(uid string) Option {
	return func(builder *Builder) error {

		builder.internal.Uid = uid

		return nil
	}
}

func Title(title string) Option {
	return func(builder *Builder) error {

		builder.internal.Title = title

		return nil
	}
}

func Description(description string) Option {
	return func(builder *Builder) error {

		builder.internal.Description = description

		return nil
	}
}

func Revision(revision int64) Option {
	return func(builder *Builder) error {
		if !(revision <= 9223372036854775807) {
			return errors.New("revision must be <= 9223372036854775807")
		}

		builder.internal.Revision = revision

		return nil
	}
}

func GnetId(gnetId string) Option {
	return func(builder *Builder) error {

		builder.internal.GnetId = gnetId

		return nil
	}
}

func Tags(tags []string) Option {
	return func(builder *Builder) error {

		builder.internal.Tags = tags

		return nil
	}
}

func Style(style DashboardStyle) Option {
	return func(builder *Builder) error {

		builder.internal.Style = style

		return nil
	}
}

func Timezone(timezone string) Option {
	return func(builder *Builder) error {

		builder.internal.Timezone = timezone

		return nil
	}
}

func Editable(editable bool) Option {
	return func(builder *Builder) error {

		builder.internal.Editable = editable

		return nil
	}
}

func GraphTooltip(graphTooltip DashboardCursorSync) Option {
	return func(builder *Builder) error {

		builder.internal.GraphTooltip = graphTooltip

		return nil
	}
}

func Time(time TimeInterval) Option {
	return func(builder *Builder) error {

		builder.internal.Time = time

		return nil
	}
}

func Timepicker(timepicker TimePicker) Option {
	return func(builder *Builder) error {

		builder.internal.Timepicker = timepicker

		return nil
	}
}

func FiscalYearStartMonth(fiscalYearStartMonth int64) Option {
	return func(builder *Builder) error {
		if !(fiscalYearStartMonth < 12) {
			return errors.New("fiscalYearStartMonth must be < 12")
		}

		builder.internal.FiscalYearStartMonth = fiscalYearStartMonth

		return nil
	}
}

func LiveNow(liveNow bool) Option {
	return func(builder *Builder) error {

		builder.internal.LiveNow = liveNow

		return nil
	}
}

func WeekStart(weekStart string) Option {
	return func(builder *Builder) error {

		builder.internal.WeekStart = weekStart

		return nil
	}
}

func Refresh(refresh StringOrBool) Option {
	return func(builder *Builder) error {

		builder.internal.Refresh = refresh

		return nil
	}
}

func SchemaVersion(schemaVersion int64) Option {
	return func(builder *Builder) error {
		if !(schemaVersion >= 0) {
			return errors.New("schemaVersion must be >= 0")
		}

		if !(schemaVersion <= 65535) {
			return errors.New("schemaVersion must be <= 65535")
		}

		builder.internal.SchemaVersion = schemaVersion

		return nil
	}
}

func Version(version int64) Option {
	return func(builder *Builder) error {
		if !(version >= 0) {
			return errors.New("version must be >= 0")
		}

		if !(version <= 4294967295) {
			return errors.New("version must be <= 4294967295")
		}

		builder.internal.Version = version

		return nil
	}
}

func Panels(panels []RowPanel) Option {
	return func(builder *Builder) error {

		builder.internal.Panels = panels

		return nil
	}
}

func Templating(templating DashboardTemplating) Option {
	return func(builder *Builder) error {

		builder.internal.Templating = templating

		return nil
	}
}

func Annotations(annotations AnnotationContainer) Option {
	return func(builder *Builder) error {

		builder.internal.Annotations = annotations

		return nil
	}
}

func Links(links []DashboardLink) Option {
	return func(builder *Builder) error {

		builder.internal.Links = links

		return nil
	}
}
