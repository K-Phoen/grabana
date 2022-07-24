package discord

import (
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a "discord"
// contact point type.
type Option func(contactType *discordType)

type discordType struct {
	builder *sdk.ContactPointType
}

// With creates an Opsgenie contact point type with the given settings.
func With(webhookURL string, opts ...Option) alertmanager.ContactPointOption {
	discord := &discordType{
		builder: &sdk.ContactPointType{
			Type: "discord",
			Settings: map[string]interface{}{
				"url": webhookURL,
			},
		},
	}

	for _, opt := range opts {
		opt(discord)
	}

	return func(contact *alertmanager.Contact) {
		contact.Builder.GrafanaManagedReceivers = append(contact.Builder.GrafanaManagedReceivers, *discord.builder)
	}
}

// UseDiscordUsername uses the username configured in Discord's webhook settings.
// Otherwise, the username will be 'Grafana'.
func UseDiscordUsername() Option {
	return func(contactType *discordType) {
		contactType.builder.Settings["use_discord_username"] = true
	}
}
