package row

import (
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/gauge"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/heatmap"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a row.
type Option func(row *Row) error

// Row represents a dashboard row.
type Row struct {
	builder *sdk.Row
	alerts  []*alert.Alert
}

// New creates a new row.
func New(board *sdk.Board, title string, options ...Option) (*Row, error) {
	panel := &Row{builder: board.AddRow(title)}

	for _, opt := range append(defaults(), options...) {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

func defaults() []Option {
	return []Option{
		ShowTitle(),
	}
}

// Alerts returns a list of alerts defined within this row.
func (row *Row) Alerts() []*alert.Alert {
	return row.alerts
}

// WithGraph adds a "graph" panel in the row.
// Deprecated: use WithTimeSeries() instead.
func WithGraph(title string, options ...graph.Option) Option {
	return func(row *Row) error {
		panel, err := graph.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		if panel.Alert == nil {
			return nil
		}

		if panel.Builder.Datasource != nil {
			panel.Alert.Datasource = panel.Builder.Datasource.LegacyName
		}

		row.alerts = append(row.alerts, panel.Alert)

		return nil
	}
}

// WithTimeSeries adds a "timeseries" panel in the row.
func WithTimeSeries(title string, options ...timeseries.Option) Option {
	return func(row *Row) error {
		panel, err := timeseries.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		if panel.Alert == nil {
			return nil
		}

		if panel.Builder.Datasource != nil {
			panel.Alert.Datasource = panel.Builder.Datasource.LegacyName
		}

		row.alerts = append(row.alerts, panel.Alert)

		return nil
	}
}

// WithGauge adds a "gauge" panel in the row.
func WithGauge(title string, options ...gauge.Option) Option {
	return func(row *Row) error {
		panel, err := gauge.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// WithLogs adds a "logs" panel in the row.
func WithLogs(title string, options ...logs.Option) Option {
	return func(row *Row) error {
		panel, err := logs.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// WithSingleStat adds a "single stat" panel in the row.
// Deprecated: use WithStat() instead
func WithSingleStat(title string, options ...singlestat.Option) Option {
	return func(row *Row) error {
		panel, err := singlestat.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// WithStat adds a "stat" panel in the row.
func WithStat(title string, options ...stat.Option) Option {
	return func(row *Row) error {
		panel, err := stat.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// WithTable adds a "table" panel in the row.
func WithTable(title string, options ...table.Option) Option {
	return func(row *Row) error {
		panel, err := table.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// WithText adds a "text" panel in the row.
func WithText(title string, options ...text.Option) Option {
	return func(row *Row) error {
		panel, err := text.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// WithHeatmap adds a "heatmap" panel in the row.
func WithHeatmap(title string, options ...heatmap.Option) Option {
	return func(row *Row) error {
		panel, err := heatmap.New(title, options...)
		if err != nil {
			return err
		}

		row.builder.Add(panel.Builder)

		return nil
	}
}

// ShowTitle ensures that the title of the row will be displayed.
func ShowTitle() Option {
	return func(row *Row) error {
		row.builder.ShowTitle = true

		return nil
	}
}

// HideTitle ensures that the title of the row will NOT be displayed.
func HideTitle() Option {
	return func(row *Row) error {
		row.builder.ShowTitle = false

		return nil
	}
}

// RepeatFor will repeat the row for all values of the given variable.
func RepeatFor(variable string) Option {
	return func(row *Row) error {
		row.builder.Repeat = &variable

		return nil
	}
}

// Collapse makes the row collapsed by default.
func Collapse() Option {
	return func(row *Row) error {
		row.builder.Collapse = true

		return nil
	}
}
