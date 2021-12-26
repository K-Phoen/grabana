package alertmanager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContactPoint(t *testing.T) {
	req := require.New(t)

	contact := ContactPoint("team-a")

	req.Equal("team-a", contact.Builder.Name)
	req.Empty(contact.Builder.GrafanaManagedReceivers)
}
