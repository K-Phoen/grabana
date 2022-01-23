package webhook

import (
	"strconv"
	"strings"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a "webhook"
// contact point type.
type Option func(contactType *contactType)

type contactType struct {
	builder *sdk.ContactPointType
}

// Call creates a "webhook" contact point type.
func Call(url string, opts ...Option) alertmanager.ContactPointOption {
	webhook := &contactType{
		builder: &sdk.ContactPointType{
			Type: "webhook",
			Settings: map[string]interface{}{
				"url": url,
			},
			SecureSettings: make(map[string]interface{}),
		},
	}

	for _, opt := range opts {
		opt(webhook)
	}

	return func(contact *alertmanager.Contact) {
		contact.Builder.GrafanaManagedReceivers = append(contact.Builder.GrafanaManagedReceivers, *webhook.builder)
	}
}

// Method defines the HTTP method used to call the webhook. Should be POST or PUT.
func Method(method string) Option {
	return func(contactType *contactType) {
		contactType.builder.Settings["httpMethod"] = strings.ToUpper(method)
	}
}

// Credentials sets the credentials used to call the webhook.
func Credentials(username string, password string) Option {
	return func(contactType *contactType) {
		contactType.builder.Settings["username"] = username
		contactType.builder.SecureSettings["password"] = password
	}
}

// MaxAlerts sets the maximum number of alerts to include in a single call.
// Remaining alerts in the same batch will be ignored above this number.
// 0 means no limit.
func MaxAlerts(max int) Option {
	return func(contactType *contactType) {
		contactType.builder.Settings["maxAlerts"] = strconv.Itoa(max)
	}
}
