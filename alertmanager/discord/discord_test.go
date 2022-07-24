package discord

import (
	"testing"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/stretchr/testify/require"
)

func TestWith(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", With("url"))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("discord", contactType.Type)
	req.Equal("url", contactType.Settings["url"].(string))
}

func TestUseDiscordUsername(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", With("", UseDiscordUsername()))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.True(contactType.Settings["use_discord_username"].(bool))
}
