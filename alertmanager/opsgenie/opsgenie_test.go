package opsgenie

import (
	"testing"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/stretchr/testify/require"
)

func TestWith(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", With("url", "key"))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.Equal("opsgenie", contactType.Type)
	req.Equal("url", contactType.Settings["apiUrl"].(string))
	req.Equal("key", contactType.SecureSettings["apiKey"].(string))
	req.Equal("tags", contactType.Settings["sendTagsAs"].(string))
}

func TestAutoClose(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", With("", "", AutoClose()))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.True(contactType.Settings["autoClose"].(bool))
}

func TestOverridePriority(t *testing.T) {
	req := require.New(t)

	contactPoint := alertmanager.ContactPoint("", With("", "", OverridePriority()))

	req.Len(contactPoint.Builder.GrafanaManagedReceivers, 1)

	contactType := contactPoint.Builder.GrafanaManagedReceivers[0]
	req.True(contactType.Settings["overridePriority"].(bool))
}
