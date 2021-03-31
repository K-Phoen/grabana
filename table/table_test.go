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

func TestTablePanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.Len(panel.Builder.TablePanel.Targets, 1)
}

func TestTablePanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithInfluxDBTarget("buckets()"))

	req.Len(panel.Builder.TablePanel.Targets, 1)
}

func TestColumnsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideColumn("Time.*"), HideColumn("Duration.*"))

	req.Len(panel.Builder.TablePanel.Styles, 3)
	req.Equal("Duration.*", panel.Builder.TablePanel.Styles[0].Pattern)
	req.Equal("hidden", panel.Builder.TablePanel.Styles[0].Type)
	req.Equal("Time.*", panel.Builder.TablePanel.Styles[1].Pattern)
	req.Equal("hidden", panel.Builder.TablePanel.Styles[1].Type)
	req.Equal("/.*/", panel.Builder.TablePanel.Styles[2].Pattern)
	req.Equal("string", panel.Builder.TablePanel.Styles[2].Type)
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

func TestTablePanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestTablePanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}
