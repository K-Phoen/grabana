package row

import (
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/heatmap"
	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a row.
type Option func(row *Row)

// Row represents a dashboard row.
type Row struct {
	builder *sdk.Row
}

// New creates a new row.
func New(board *sdk.Board, title string, options ...Option) *Row {
	panel := &Row{builder: board.AddRow(title)}

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		ShowTitle(),
	}
}

// WithGraph adds a "graph" panel in the row.
func WithGraph(title string, options ...graph.Option) Option {
	return func(row *Row) {
		graphPanel := graph.New(title, options...)

		row.builder.Add(graphPanel.Builder)
	}
}

// WithTimeSeries adds a "timeseries" panel in the row.
func WithTimeSeries(title string, options ...timeseries.Option) Option {
	return func(row *Row) {
		panel := timeseries.New(title, options...)

		row.builder.Add(panel.Builder)
	}
}

// WithSingleStat adds a "single stat" panel in the row.
func WithSingleStat(title string, options ...singlestat.Option) Option {
	return func(row *Row) {
		singleStatPanel := singlestat.New(title, options...)

		row.builder.Add(singleStatPanel.Builder)
	}
}

// WithTable adds a "table" panel in the row.
func WithTable(title string, options ...table.Option) Option {
	return func(row *Row) {
		tablePanel := table.New(title, options...)

		row.builder.Add(tablePanel.Builder)
	}
}

// WithText adds a "text" panel in the row.
func WithText(title string, options ...text.Option) Option {
	return func(row *Row) {
		textPanel := text.New(title, options...)

		row.builder.Add(textPanel.Builder)
	}
}

// WithHeatmap adds a "heatmap" panel in the row.
func WithHeatmap(title string, options ...heatmap.Option) Option {
	return func(row *Row) {
		heatmapPanel := heatmap.New(title, options...)

		row.builder.Add(heatmapPanel.Builder)
	}
}

// ShowTitle ensures that the title of the row will be displayed.
func ShowTitle() Option {
	return func(row *Row) {
		row.builder.ShowTitle = true
	}
}

// HideTitle ensures that the title of the row will NOT be displayed.
func HideTitle() Option {
	return func(row *Row) {
		row.builder.ShowTitle = false
	}
}

// RepeatFor will repeat the row for all values of the given variable.
func RepeatFor(variable string) Option {
	return func(row *Row) {
		row.builder.Repeat = &variable
	}
}

// Collapse makes the row collapsed by default.
func Collapse() Option {
	return func(row *Row) {
		row.builder.Collapse = true
	}
}
