package webhook

import (
	"net/http"
	"testing"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/stretchr/testify/require"
)

func TestTo(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Call("webhook-url"))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("webhook", contactType.Type)
	req.Equal("webhook-url", contactType.Settings["url"])
}

func TestMethod(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Call("", Method(http.MethodPut)))

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("PUT", contactType.Settings["httpMethod"])
}

func TestCredentials(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Call("", Credentials("joe", "lafrite")))

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("joe", contactType.Settings["username"])
	req.Equal("lafrite", contactType.SecureSettings["password"])
}

func TestMaxAlerts(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", Call("", MaxAlerts(42)))

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("42", contactType.Settings["maxAlerts"])
}
