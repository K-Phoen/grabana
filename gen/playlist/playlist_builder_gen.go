package playlist

import "encoding/json"

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

func Interval(interval string) Option {
	return func(builder *Builder) error {

		builder.internal.Interval = interval

		return nil
	}
}

func Items(items []PlaylistItem) Option {
	return func(builder *Builder) error {

		builder.internal.Items = items

		return nil
	}
}

func Name(name string) Option {
	return func(builder *Builder) error {

		builder.internal.Name = name

		return nil
	}
}

func Xxx(xxx string) Option {
	return func(builder *Builder) error {

		builder.internal.Xxx = xxx

		return nil
	}
}
