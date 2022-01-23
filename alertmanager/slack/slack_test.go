package slack

import (
	"testing"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/stretchr/testify/require"
)

func TestWebhook(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Webhook("webhook-url"))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("slack", contactType.Type)
	req.Equal("webhook-url", contactType.SecureSettings["url"].(string))
}

func TestTitle(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Webhook("", Title("{{ template \"slack.default.title\" . }}")))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("{{ template \"slack.default.title\" . }}", contactType.Settings["title"].(string))
}

func TestBody(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Webhook("", Body("some-body")))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("some-body", contactType.Settings["text"].(string))
}
