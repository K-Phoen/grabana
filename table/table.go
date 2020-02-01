package table

import (
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/grafana-tools/sdk"
)

type AggregationType string
type Option func(table *Table)

const AVG AggregationType = "avg"
const Count AggregationType = "count"
const Current AggregationType = "current"
const Min AggregationType = "min"
const Max AggregationType = "max"

type Aggregation struct {
	Label string
	Type  AggregationType
}

type Table struct {
	Builder *sdk.Panel
}

func New(title string, options ...Option) *Table {
	panel := &Table{Builder: sdk.NewTable(title)}
	empty := ""

	panel.Builder.IsNew = false
	panel.Builder.Span = 6
	panel.Builder.TablePanel.Styles = []sdk.ColumnStyle{
		{
			Alias:   &empty,
			Pattern: "/.*/",
			Type:    "string",
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Editable(),
		TimeSeriesToRows(),
	}
}

// WithPrometheusTarget adds a prometheus target to the table.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(table *Table) {
		table.Builder.AddTarget(&sdk.Target{
			Expr:           target.Expr,
			IntervalFactor: target.IntervalFactor,
			Interval:       target.Interval,
			Step:           target.Step,
			LegendFormat:   target.LegendFormat,
			Instant:        target.Instant,
			Format:         target.Format,
		})
	}
}

// HideColumn hides the column having a label matching the given pattern.
func HideColumn(columnLabelPattern string) Option {
	return func(table *Table) {
		table.Builder.TablePanel.Styles = append(table.Builder.TablePanel.Styles, sdk.ColumnStyle{
			Pattern: columnLabelPattern,
			Type:    "hidden",
		})
	}
}

// TimeSeriesToRows displays the data in rows.
func TimeSeriesToRows() Option {
	return func(table *Table) {
		table.Builder.TablePanel.Transform = "timeseries_to_rows"
	}
}

// TimeSeriesToColumns displays the data in columns.
func TimeSeriesToColumns() Option {
	return func(table *Table) {
		table.Builder.TablePanel.Transform = "timeseries_to_columns"
	}
}

// AsJSON displays the data as JSON.
func AsJSON() Option {
	return func(table *Table) {
		table.Builder.TablePanel.Transform = "json"
	}
}

// AsTable displays the data as a table.
func AsTable() Option {
	return func(table *Table) {
		table.Builder.TablePanel.Transform = "table"
	}
}

// AsTable displays the data as annotations.
func AsAnnotations() Option {
	return func(table *Table) {
		table.Builder.TablePanel.Transform = "annotations"
	}
}

// AsTimeSeriesAggregations displays the data according to the given aggregation methods.
func AsTimeSeriesAggregations(aggregations []Aggregation) Option {
	return func(table *Table) {
		columns := make([]sdk.Column, 0, len(aggregations))

		for _, aggregation := range aggregations {
			columns = append(columns, sdk.Column{
				TextType: aggregation.Label,
				Value:    string(aggregation.Type),
			})
		}

		table.Builder.TablePanel.Transform = "timeseries_aggregations"
		table.Builder.TablePanel.Columns = columns
	}
}

// Editable marks the graph as editable.
func Editable() Option {
	return func(table *Table) {
		table.Builder.Editable = true
	}
}

// ReadOnly marks the graph as non-editable.
func ReadOnly() Option {
	return func(table *Table) {
		table.Builder.Editable = false
	}
}

// DataSource sets the data source to be used by the table.
func DataSource(source string) Option {
	return func(table *Table) {
		table.Builder.Datasource = &source
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(table *Table) {
		table.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(table *Table) {
		table.Builder.Height = &height
	}
}
