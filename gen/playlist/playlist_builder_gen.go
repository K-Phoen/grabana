package playlist

type Option func(builder *Builder) error

type Builder struct {
	internal *Playlist
}

func Name(name string) Option {
	return func(builder *Builder) error {

		builder.internal.Name = name

		return nil
	}
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

func Xxx(xxx string) Option {
	return func(builder *Builder) error {

		builder.internal.Xxx = xxx

		return nil
	}
}
