package dashboard

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
