package table

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTablePanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Table panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Table panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestTablePanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestTablePanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestTablePanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestTablePanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithPrometheusTarget("go_threads"))

	req.Len(panel.Builder.TablePanel.Targets, 1)
}

func TestColumnsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideColumn("Time.*"))

	req.Len(panel.Builder.TablePanel.Styles, 2)
	req.Equal("Time.*", panel.Builder.TablePanel.Styles[1].Pattern)
	req.Equal("hidden", panel.Builder.TablePanel.Styles[1].Type)
}

func TestDataCanBeTransformedInTimeSeriesToRows(t *testing.T) {
	req := require.New(t)

	panel := New("", TimeSeriesToRows())

	req.Equal("timeseries_to_rows", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedInTimeSeriesToColumns(t *testing.T) {
	req := require.New(t)

	panel := New("", TimeSeriesToColumns())

	req.Equal("timeseries_to_columns", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsJSON(t *testing.T) {
	req := require.New(t)

	panel := New("", AsJSON())

	req.Equal("json", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsTable(t *testing.T) {
	req := require.New(t)

	panel := New("", AsTable())

	req.Equal("table", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsAnnotations(t *testing.T) {
	req := require.New(t)

	panel := New("", AsAnnotations())

	req.Equal("annotations", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsTimeSeriesAggregations(t *testing.T) {
	req := require.New(t)

	panel := New("", AsTimeSeriesAggregations([]Aggregation{
		{
			Label: "Average",
			Type:  AVG,
		},
	}))

	req.Equal("timeseries_aggregations", panel.Builder.TablePanel.Transform)
	req.Len(panel.Builder.TablePanel.Columns, 1)
	req.Equal("Average", panel.Builder.TablePanel.Columns[0].TextType)
	req.Equal(string(AVG), panel.Builder.TablePanel.Columns[0].Value)
}

func TestTextPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}
