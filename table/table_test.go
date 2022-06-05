package table

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/stretchr/testify/require"
)

func TestNewTablePanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("Table panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("Table panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestTablePanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestTablePanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(6))

	req.NoError(err)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestInvalidTablePanelWidthIsRejected(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(-6))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestTablePanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestTablePanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("prometheus-default"))

	req.NoError(err)
	req.Equal("prometheus-default", panel.Builder.Datasource.LegacyName)
}

func TestTablePanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithPrometheusTarget("go_threads"))

	req.NoError(err)
	req.Len(panel.Builder.TablePanel.Targets, 1)
}

func TestTablePanelCanHaveGraphiteTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithGraphiteTarget("stats_counts.statsd.packets_received"))

	req.NoError(err)
	req.Len(panel.Builder.TablePanel.Targets, 1)
}

func TestTablePanelCanHaveInfluxDBTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithInfluxDBTarget("buckets()"))

	req.NoError(err)
	req.Len(panel.Builder.TablePanel.Targets, 1)
}

func TestColumnsCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HideColumn("Time.*"), HideColumn("Duration.*"))

	req.NoError(err)
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

	panel, err := New("", TimeSeriesToRows())

	req.NoError(err)
	req.Equal("timeseries_to_rows", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedInTimeSeriesToColumns(t *testing.T) {
	req := require.New(t)

	panel, err := New("", TimeSeriesToColumns())

	req.NoError(err)
	req.Equal("timeseries_to_columns", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsJSON(t *testing.T) {
	req := require.New(t)

	panel, err := New("", AsJSON())

	req.NoError(err)
	req.Equal("json", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsTable(t *testing.T) {
	req := require.New(t)

	panel, err := New("", AsTable())

	req.NoError(err)
	req.Equal("table", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsAnnotations(t *testing.T) {
	req := require.New(t)

	panel, err := New("", AsAnnotations())

	req.NoError(err)
	req.Equal("annotations", panel.Builder.TablePanel.Transform)
}

func TestDataCanBeTransformedAsTimeSeriesAggregations(t *testing.T) {
	req := require.New(t)

	panel, err := New("", AsTimeSeriesAggregations([]Aggregation{
		{
			Label: "Average",
			Type:  AVG,
		},
	}))

	req.NoError(err)
	req.Equal("timeseries_aggregations", panel.Builder.TablePanel.Transform)
	req.Len(panel.Builder.TablePanel.Columns, 1)
	req.Equal("Average", panel.Builder.TablePanel.Columns[0].TextType)
	req.Equal(string(AVG), panel.Builder.TablePanel.Columns[0].Value)
}

func TestTablePanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestTablePanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}
