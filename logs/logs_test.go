package logs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogsPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("Logs panel")

	req.False(panel.Builder.IsNew)
	req.Equal("Logs panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
	req.Equal("Descending", panel.Builder.LogsPanel.Options.SortOrder)
	req.True(panel.Builder.LogsPanel.Options.EnableLogDetails)
}

func TestLogsPanelCanHaveLokiTargets(t *testing.T) {
	req := require.New(t)

	panel := New("", WithLokiTarget("{app=\"loki\"}"))

	req.Len(panel.Builder.LogsPanel.Targets, 1)
}

func TestLogsPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Span(8))

	req.Equal(float32(8), panel.Builder.Span)
}

func TestLogsPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Height("400px"))

	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestLogsPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel := New("", Transparent())

	req.True(panel.Builder.Transparent)
}

func TestLogsPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("", Description("lala"))

	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestLogsPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", DataSource("loki-default"))

	req.Equal("loki-default", *panel.Builder.Datasource)
}

func TestRepeatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel := New("", Repeat("ds"))

	req.NotNil(panel.Builder.Repeat)
	req.Equal("ds", *panel.Builder.Repeat)
}

func TestTimeCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", Time())

	req.True(panel.Builder.LogsPanel.Options.ShowTime)
}

func TestUniqueLabelsCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", UniqueLabels())

	req.True(panel.Builder.LogsPanel.Options.ShowLabels)
}

func TestCommonLabelsCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel := New("", CommonLabels())

	req.True(panel.Builder.LogsPanel.Options.ShowCommonLabels)
}

func TestLinesCanBeWrapped(t *testing.T) {
	req := require.New(t)

	panel := New("", WrapLines())

	req.True(panel.Builder.LogsPanel.Options.WrapLogMessage)
}

func TestJSONCanBePrettyPrinted(t *testing.T) {
	req := require.New(t)

	panel := New("", PrettifyJSON())

	req.True(panel.Builder.LogsPanel.Options.PrettifyLogMessage)
}

func TestLogDetailsCanBeDisabled(t *testing.T) {
	req := require.New(t)

	panel := New("", HideLogDetails())

	req.False(panel.Builder.LogsPanel.Options.EnableLogDetails)
}

func TestSortOrderCanBeSet(t *testing.T) {
	testCases := []struct {
		order    SortOrder
		expected string
	}{
		{order: Asc, expected: "Ascending"},
		{order: Desc, expected: "Descending"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.expected, func(t *testing.T) {
			req := require.New(t)

			panel := New("", Order(tc.order))

			req.Equal(tc.expected, panel.Builder.LogsPanel.Options.SortOrder)
		})
	}
}

func TestDedupStrategyCanBeSet(t *testing.T) {
	testCases := []struct {
		strategy DedupStrategy
		expected string
	}{
		{strategy: None, expected: "none"},
		{strategy: Exact, expected: "exact"},
		{strategy: Numbers, expected: "numbers"},
		{strategy: Signature, expected: "signature"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.expected, func(t *testing.T) {
			req := require.New(t)

			panel := New("", Deduplication(tc.strategy))

			req.Equal(tc.expected, panel.Builder.LogsPanel.Options.DedupStrategy)
		})
	}
}
