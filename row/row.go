package row

import (
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/heatmap"
	"github.com/K-Phoen/grabana/logs"
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
	alerts  []*sdk.Alert
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

// Alerts returns a list of alerts defined within this row.
func (row *Row) Alerts() []*sdk.Alert {
	return row.alerts
}

// WithGraph adds a "graph" panel in the row.
func WithGraph(title string, options ...graph.Option) Option {
	return func(row *Row) {
		panel := graph.New(title, options...)
		row.builder.Add(panel.Builder)

		if panel.Alert != nil {
			row.alerts = append(row.alerts, panel.Alert)
		}
	}
}

// WithTimeSeries adds a "timeseries" panel in the row.
func WithTimeSeries(title string, options ...timeseries.Option) Option {
	return func(row *Row) {
		panel := timeseries.New(title, options...)
		row.builder.Add(panel.Builder)

		if panel.Alert != nil {
			row.alerts = append(row.alerts, panel.Alert)
		}
	}
}

// WithLogs adds a "logs" panel in the row.
func WithLogs(title string, options ...logs.Option) Option {
	return func(row *Row) {
		panel := logs.New(title, options...)
		row.builder.Add(panel.Builder)
	}
}

// WithSingleStat adds a "single stat" panel in the row.
func WithSingleStat(title string, options ...singlestat.Option) Option {
	return func(row *Row) {
		panel := singlestat.New(title, options...)
		row.builder.Add(panel.Builder)
	}
}

// WithTable adds a "table" panel in the row.
func WithTable(title string, options ...table.Option) Option {
	return func(row *Row) {
		panel := table.New(title, options...)
		row.builder.Add(panel.Builder)
	}
}

// WithText adds a "text" panel in the row.
func WithText(title string, options ...text.Option) Option {
	return func(row *Row) {
		panel := text.New(title, options...)
		row.builder.Add(panel.Builder)
	}
}

// WithHeatmap adds a "heatmap" panel in the row.
func WithHeatmap(title string, options ...heatmap.Option) Option {
	return func(row *Row) {
		panel := heatmap.New(title, options...)
		row.builder.Add(panel.Builder)
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
