package logs

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
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
type Option func(logs *Logs) error

// Logs represents a logs panel.
type Logs struct {
	Builder *sdk.Panel
}

// New creates a new logs panel.
func New(title string, options ...Option) (*Logs, error) {
	panel := &Logs{Builder: sdk.NewLogs(title)}
	panel.Builder.IsNew = false
	panel.Builder.LogsPanel.Options.EnableLogDetails = true

	for _, opt := range append(defaults(), options...) {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

func defaults() []Option {
	return []Option{
		Span(6),
		Order(Desc),
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(logs *Logs) error {
		logs.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			logs.Builder.Links = append(logs.Builder.Links, link.Builder)
		}

		return nil
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(logs *Logs) error {
		logs.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// WithLokiTarget adds a loki query to the graph.
func WithLokiTarget(query string, options ...loki.Option) Option {
	target := loki.New(query, options...)

	return func(logs *Logs) error {
		logs.Builder.AddTarget(&sdk.Target{
			RefID:        target.Ref,
			Expr:         target.Expr,
			LegendFormat: target.LegendFormat,
		})

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(logs *Logs) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		logs.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(logs *Logs) error {
		logs.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(logs *Logs) error {
		logs.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(logs *Logs) error {
		logs.Builder.Transparent = true

		return nil
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(logs *Logs) error {
		logs.Builder.Repeat = &repeat

		return nil
	}
}

// Time displays the "time" column. This is the timestamp associated with the
// log line as reported from the data source.
func Time() Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.ShowTime = true

		return nil
	}
}

// UniqueLabels displays the "unique labels" column, which shows only non-common labels.
func UniqueLabels() Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.ShowLabels = true

		return nil
	}
}

// CommonLabels displays the "common labels".
func CommonLabels() Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.ShowCommonLabels = true

		return nil
	}
}

// WrapLines enables line wrapping.
func WrapLines() Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.WrapLogMessage = true

		return nil
	}
}

// PrettifyJSON pretty prints all JSON logs. This setting does not affect logs
// in any format other than JSON.
func PrettifyJSON() Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.PrettifyLogMessage = true

		return nil
	}
}

// HideLogDetails disables the log details view for each log row.
func HideLogDetails() Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.EnableLogDetails = false

		return nil
	}
}

// Order display results in descending or ascending time order.
// The default is Descending, showing the newest logs first.
// Set to Ascending to show the oldest log lines first.
func Order(order SortOrder) Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.SortOrder = string(order)

		return nil
	}
}

// Deduplication sets the deduplication strategy.
func Deduplication(dedup DedupStrategy) Option {
	return func(logs *Logs) error {
		logs.Builder.LogsPanel.Options.DedupStrategy = string(dedup)

		return nil
	}
}
