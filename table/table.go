package table

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a table panel.
type Option func(table *Table) error

// AggregationType represents an aggregation function used on values returned
// by the query.
type AggregationType string

const (
	// AVG aggregates results by computing the average.
	AVG AggregationType = "avg"

	// Count aggregates results by counting them.
	Count AggregationType = "count"

	// Current aggregates results by keeping only the current value.
	Current AggregationType = "current"

	// Min aggregates results by keeping only the smallest value.
	Min AggregationType = "min"

	// Max aggregates results by keeping only the largest value.
	Max AggregationType = "max"
)

// Aggregation configures how to display an aggregate in the table.
type Aggregation struct {
	Label string
	Type  AggregationType
}

// Table represents a table panel.
type Table struct {
	Builder *sdk.Panel
}

// New creates a new table panel.
func New(title string, options ...Option) (*Table, error) {
	panel := &Table{Builder: sdk.NewTable(title)}
	empty := ""

	panel.Builder.IsNew = false
	panel.Builder.TablePanel.Styles = []sdk.ColumnStyle{
		{
			Alias:   &empty,
			Pattern: "/.*/",
			Type:    "string",
		},
	}

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
		TimeSeriesToRows(),
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(table *Table) error {
		table.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			table.Builder.Links = append(table.Builder.Links, link.Builder)
		}

		return nil
	}
}

// WithPrometheusTarget adds a prometheus target to the table.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(table *Table) error {
		table.Builder.AddTarget(&sdk.Target{
			RefID:          target.Ref,
			Expr:           target.Expr,
			IntervalFactor: target.IntervalFactor,
			Interval:       target.Interval,
			Step:           target.Step,
			LegendFormat:   target.LegendFormat,
			Instant:        target.Instant,
			Format:         target.Format,
		})

		return nil
	}
}

// WithGraphiteTarget adds a Graphite target to the table.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(table *Table) error {
		table.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the table.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(table *Table) error {
		table.Builder.AddTarget(target.Builder)

		return nil
	}
}

// HideColumn hides the column having a label matching the given pattern.
func HideColumn(columnLabelPattern string) Option {
	return func(table *Table) error {
		table.Builder.TablePanel.Styles = append([]sdk.ColumnStyle{
			{
				Pattern: columnLabelPattern,
				Type:    "hidden",
			},
		}, table.Builder.TablePanel.Styles...)

		return nil
	}
}

// TimeSeriesToRows displays the data in rows.
func TimeSeriesToRows() Option {
	return func(table *Table) error {
		table.Builder.TablePanel.Transform = "timeseries_to_rows"

		return nil
	}
}

// TimeSeriesToColumns displays the data in columns.
func TimeSeriesToColumns() Option {
	return func(table *Table) error {
		table.Builder.TablePanel.Transform = "timeseries_to_columns"

		return nil
	}
}

// AsJSON displays the data as JSON.
func AsJSON() Option {
	return func(table *Table) error {
		table.Builder.TablePanel.Transform = "json"

		return nil
	}
}

// AsTable displays the data as a table.
func AsTable() Option {
	return func(table *Table) error {
		table.Builder.TablePanel.Transform = "table"

		return nil
	}
}

// AsAnnotations displays the data as annotations.
func AsAnnotations() Option {
	return func(table *Table) error {
		table.Builder.TablePanel.Transform = "annotations"

		return nil
	}
}

// AsTimeSeriesAggregations displays the data according to the given aggregation methods.
func AsTimeSeriesAggregations(aggregations []Aggregation) Option {
	return func(table *Table) error {
		columns := make([]sdk.Column, 0, len(aggregations))

		for _, aggregation := range aggregations {
			columns = append(columns, sdk.Column{
				TextType: aggregation.Label,
				Value:    string(aggregation.Type),
			})
		}

		table.Builder.TablePanel.Transform = "timeseries_aggregations"
		table.Builder.TablePanel.Columns = columns

		return nil
	}
}

// DataSource sets the data source to be used by the table.
func DataSource(source string) Option {
	return func(table *Table) error {
		table.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(table *Table) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		table.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(table *Table) error {
		table.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(table *Table) error {
		table.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(table *Table) error {
		table.Builder.Transparent = true

		return nil
	}
}
