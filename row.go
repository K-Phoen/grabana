package grabana

import (
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/text"
	"github.com/grafana-tools/sdk"
)

type RowOption func(row *Row)

type Row struct {
	builder *sdk.Row
}

func rowDefaults() []RowOption {
	return []RowOption{
		ShowRowTitle(),
	}
}

func WithGraph(title string, options ...graph.Option) RowOption {
	return func(row *Row) {
		graphPanel := graph.New(title)

		for _, opt := range options {
			opt(graphPanel)
		}

		row.builder.Add(graphPanel.Builder)
	}
}

func WithText(title string, options ...text.Option) RowOption {
	return func(row *Row) {
		textPanel := text.New(title)

		for _, opt := range options {
			opt(textPanel)
		}

		row.builder.Add(textPanel.Builder)
	}
}

func ShowRowTitle() RowOption {
	return func(row *Row) {
		row.builder.ShowTitle = true
	}
}

func HideRowTitle() RowOption {
	return func(row *Row) {
		row.builder.ShowTitle = false
	}
}
