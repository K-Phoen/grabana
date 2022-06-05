package logs

import (
	"testing"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/stretchr/testify/require"
)

func TestNewLogsPanelsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel, err := New("Logs panel")

	req.NoError(err)
	req.False(panel.Builder.IsNew)
	req.Equal("Logs panel", panel.Builder.Title)
	req.Equal(float32(6), panel.Builder.Span)
	req.Equal("Descending", panel.Builder.LogsPanel.Options.SortOrder)
	req.True(panel.Builder.LogsPanel.Options.EnableLogDetails)
}

func TestLogsPanelCanHaveLinks(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Links(links.New("", "")))

	req.NoError(err)
	req.Len(panel.Builder.Links, 1)
}

func TestLogsPanelCanHaveLokiTargets(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WithLokiTarget("{app=\"loki\"}"))

	req.NoError(err)
	req.Len(panel.Builder.LogsPanel.Targets, 1)
}

func TestLogsPanelWidthCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Span(8))

	req.NoError(err)
	req.Equal(float32(8), panel.Builder.Span)
}

func TestLogsPanelRejectIncorrectWidth(t *testing.T) {
	req := require.New(t)

	_, err := New("", Span(-8))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestLogsPanelHeightCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Height("400px"))

	req.NoError(err)
	req.Equal("400px", *(panel.Builder.Height).(*string))
}

func TestLogsPanelBackgroundCanBeTransparent(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Transparent())

	req.NoError(err)
	req.True(panel.Builder.Transparent)
}

func TestLogsPanelDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Description("lala"))

	req.NoError(err)
	req.NotNil(panel.Builder.Description)
	req.Equal("lala", *panel.Builder.Description)
}

func TestLogsPanelDataSourceCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", DataSource("loki-default"))

	req.NoError(err)
	req.Equal("loki-default", panel.Builder.Datasource.LegacyName)
}

func TestRepeatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Repeat("ds"))

	req.NoError(err)
	req.NotNil(panel.Builder.Repeat)
	req.Equal("ds", *panel.Builder.Repeat)
}

func TestTimeCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", Time())

	req.NoError(err)
	req.True(panel.Builder.LogsPanel.Options.ShowTime)
}

func TestUniqueLabelsCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", UniqueLabels())

	req.NoError(err)
	req.True(panel.Builder.LogsPanel.Options.ShowLabels)
}

func TestCommonLabelsCanBeDisplayed(t *testing.T) {
	req := require.New(t)

	panel, err := New("", CommonLabels())

	req.NoError(err)
	req.True(panel.Builder.LogsPanel.Options.ShowCommonLabels)
}

func TestLinesCanBeWrapped(t *testing.T) {
	req := require.New(t)

	panel, err := New("", WrapLines())

	req.NoError(err)
	req.True(panel.Builder.LogsPanel.Options.WrapLogMessage)
}

func TestJSONCanBePrettyPrinted(t *testing.T) {
	req := require.New(t)

	panel, err := New("", PrettifyJSON())

	req.NoError(err)
	req.True(panel.Builder.LogsPanel.Options.PrettifyLogMessage)
}

func TestLogDetailsCanBeDisabled(t *testing.T) {
	req := require.New(t)

	panel, err := New("", HideLogDetails())

	req.NoError(err)
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

			panel, err := New("", Order(tc.order))

			req.NoError(err)
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

			panel, err := New("", Deduplication(tc.strategy))

			req.NoError(err)
			req.Equal(tc.expected, panel.Builder.LogsPanel.Options.DedupStrategy)
		})
	}
}
