package links

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPanelLinksCanBeCreated(t *testing.T) {
	req := require.New(t)

	link := New("link name", "https://github.com/grafana/grafana/")

	req.Equal("link name", link.Builder.Title)
	req.Equal("https://github.com/grafana/grafana/", *link.Builder.URL)
}

func TestLinksCanBeOpenedAsBlank(t *testing.T) {
	req := require.New(t)

	link := New("", "", OpenBlank())

	req.True(*link.Builder.TargetBlank)
}
