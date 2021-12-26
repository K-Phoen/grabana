package email

import (
	"strings"
	"testing"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/stretchr/testify/require"
)

func TestTo(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", To([]string{"foo@bar", "biz@localhost"}))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("email", contactType.Type)
	req.ElementsMatch([]string{"foo@bar", "biz@localhost"}, strings.Split(contactType.Settings["addresses"].(string), ","))
}

func TestSingle(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", To(nil, Single()))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.True(contactType.Settings["singleEmail"].(bool))
}

func TestMessage(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", To(nil, Message("test msg")))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("test msg", contactType.Settings["message"])
}
