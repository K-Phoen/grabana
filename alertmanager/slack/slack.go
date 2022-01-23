package slack

import (
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a "slack"
// contact point type.
type Option func(contactType *slackType)

type slackType struct {
	builder *sdk.ContactPointType
}

// Webhook creates a Slack contact point type that sends alerts to a Slack webhook.
// See https://api.slack.com/messaging/webhooks
func Webhook(webhookURL string, opts ...Option) alertmanager.ContactPointOption {
	slack := &slackType{
		builder: &sdk.ContactPointType{
			Type:     "slack",
			Settings: map[string]interface{}{},
			SecureSettings: map[string]interface{}{
				"url": webhookURL,
			},
		},
	}

	for _, opt := range opts {
		opt(slack)
	}

	return func(contact *alertmanager.Contact) {
		contact.Builder.GrafanaManagedReceivers = append(contact.Builder.GrafanaManagedReceivers, *slack.builder)
	}
}

// Title defines a templated title that will be sent in Slack messages.
func Title(templatedTitle string) Option {
	return func(contactType *slackType) {
		contactType.builder.Settings["title"] = templatedTitle
	}
}

// Body defines the body that will be sent in Slack messages.
func Body(body string) Option {
	return func(contactType *slackType) {
		contactType.builder.Settings["text"] = body
	}
}
