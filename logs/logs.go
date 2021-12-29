package logs

import (
	"github.com/K-Phoen/grabana/target/loki"
	"github.com/K-Phoen/sdk"
)

// DedupStrategy represents a deduplication strategy.
type DedupStrategy string

const (
	None      DedupStrategy = "none"
	Exact     DedupStrategy = "exact"
	Numbers   DedupStrategy = "numbers"
	Signature DedupStrategy = "signature"
)

// SortOrder represents a sort order.
type SortOrder string

const (
	Asc  SortOrder = "Ascending"
	Desc SortOrder = "Descending"
)

// Option represents an option that can be used to configure a logs panel.
type Option func(logs *Logs)

// Logs represents a logs panel.
type Logs struct {
	Builder *sdk.Panel
}

// New creates a new heatmap panel.
func New(title string, options ...Option) *Logs {
	panel := &Logs{Builder: sdk.NewLogs(title)}
	panel.Builder.IsNew = false
	panel.Builder.LogsPanel.Options.EnableLogDetails = true

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Span(6),
		Order(Desc),
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(logs *Logs) {
		logs.Builder.Datasource = &source
	}
}

// WithLokiTarget adds a loki query to the graph.
func WithLokiTarget(query string, options ...loki.Option) Option {
	target := loki.New(query, options...)

	return func(logs *Logs) {
		logs.Builder.AddTarget(&sdk.Target{
			RefID:        target.Ref,
			Expr:         target.Expr,
			LegendFormat: target.LegendFormat,
		})
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(logs *Logs) {
		logs.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(logs *Logs) {
		logs.Builder.Height = &height
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(logs *Logs) {
		logs.Builder.Description = &content
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(logs *Logs) {
		logs.Builder.Transparent = true
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(logs *Logs) {
		logs.Builder.Repeat = &repeat
	}
}

// Time displays the "time" column. This is the timestamp associated with the
// log line as reported from the data source.
func Time() Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.ShowTime = true
	}
}

// UniqueLabels displays the "unique labels" column, which shows only non-common labels.
func UniqueLabels() Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.ShowLabels = true
	}
}

// CommonLabels displays the "common labels".
func CommonLabels() Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.ShowCommonLabels = true
	}
}

// WrapLines enables line wrapping.
func WrapLines() Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.WrapLogMessage = true
	}
}

// PrettifyJSON pretty prints all JSON logs. This setting does not affect logs
// in any format other than JSON.
func PrettifyJSON() Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.PrettifyLogMessage = true
	}
}

// HideLogDetails disables the log details view for each log row.
func HideLogDetails() Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.EnableLogDetails = false
	}
}

// Order display results in descending or ascending time order.
// The default is Descending, showing the newest logs first.
// Set to Ascending to show the oldest log lines first.
func Order(order SortOrder) Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.SortOrder = string(order)
	}
}

// Deduplication sets the deduplication strategy.
func Deduplication(dedup DedupStrategy) Option {
	return func(logs *Logs) {
		logs.Builder.LogsPanel.Options.DedupStrategy = string(dedup)
	}
}
