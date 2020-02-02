package graph

import (
	"testing"

	"github.com/K-Phoen/grabana/axis"

	"github.com/stretchr/testify/require"
)

func TestNewGraphPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Graph panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Graph panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelCanBeMadeEditable(t *testing.T) {
	req := require.New(t)

	panel := New("", Editable())

	req.True(panel.Builder.Editable)
}

func TestGraphPanelCanBeMadeReadOnly(t *testing.T) {
	req := require.New(t)

	panel := New("", ReadOnly())

	req.False(panel.Builder.Editable)
}

func TestGraphPanelCanHavePrometheusTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithPrometheusTarget(
		"rate(prometheus_http_requests_total[30s])",
	))

	req.Len(panel.Builder.GraphPanel.Targets, 1)
}

func TestGraphPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(6))

	req.Equal(float32(6), panel.Builder.Span)
}

func TestGraphPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *panel.Builder.Height)
}

func TestGraphPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("prometheus-default"))

	req.Equal("prometheus-default", *panel.Builder.Datasource)
}

func TestLeftYAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", LeftYAxis(axis.Hide()))

	req.False(panel.Builder.Yaxes[0].Show)
}

func TestRightYAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", RightYAxis())

	req.True(panel.Builder.Yaxes[0].Show)
}

func TestXAxisCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", XAxis(axis.Hide()))

	req.False(panel.Builder.Xaxis.Show)
}

func TestAlertsCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Alert("some alert"))

	req.NotNil(panel.Builder.Alert)
	req.Equal("some alert", panel.Builder.Alert.Name)
}

func TestDrawModeCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Draw(Lines, Points, Bars))

	req.True(panel.Builder.Lines)
	req.True(panel.Builder.Points)
	req.True(panel.Builder.Bars)
}
