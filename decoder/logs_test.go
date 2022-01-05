package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/stretchr/testify/require"
)

func TestLogsPanelsCanBeDecoded(t *testing.T) {
	req := require.New(t)

	panel := DashboardLogs{
		Title:       "awesome logs",
		Description: "awesome description",
		Span:        12,
		Height:      "300px",
		Transparent: true,
		Datasource:  "some-loki",
		Repeat:      "ds",
	}

	rowOption, err := panel.toOption()

	req.NoError(err)

	testBoard := dashboard.New("", dashboard.Row("", rowOption))
	req.Len(testBoard.Internal().Rows, 1)
	panels := testBoard.Internal().Rows[0].Panels
	req.Len(panels, 1)

	sdkPanel := panels[0]
	tsPanel := panels[0].LogsPanel

	req.NotNil(tsPanel)
	req.Equal(panel.Title, sdkPanel.Title)
	req.Equal(panel.Description, *sdkPanel.Description)
	req.Equal(panel.Span, sdkPanel.Span)
	req.True(sdkPanel.Transparent)
	req.Equal(panel.Datasource, *sdkPanel.Datasource)
	req.Equal(panel.Repeat, *sdkPanel.Repeat)
}
