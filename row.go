package grabana

import (
	"github.com/K-Phoen/grabana/text"
	"github.com/grafana-tools/sdk"
)

type RowOption func(row *Row)

type Row struct {
	builder *sdk.Row
}

func WithGraph(title string, options ...GraphOption) RowOption {
	return func(row *Row) {
		graph := &Graph{builder: sdk.NewGraph(title)}

		GraphDefaults(graph)

		for _, opt := range options {
			opt(graph)
		}

		row.builder.Add(graph.builder)
	}
}

func WithText(title string, options ...text.Option) RowOption {
	return func(row *Row) {
		textPanel := &text.Text{Builder: sdk.NewText(title)}

		text.Defaults(textPanel)

		for _, opt := range options {
			opt(textPanel)
		}

		row.builder.Add(textPanel.Builder)
	}
}

func rowDefaults() []RowOption {
	return []RowOption{
		ShowRowTitle(),
		EditableRow(),
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

func EditableRow() RowOption {
	return func(row *Row) {
		row.builder.Editable = true
	}
}

func ReadOnlyRow() RowOption {
	return func(row *Row) {
		row.builder.Editable = false
	}
}
