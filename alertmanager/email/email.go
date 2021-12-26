package email

import (
	"strings"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure an "email"
// contact point type.
type Option func(contactType *emailType)

type emailType struct {
	builder *sdk.ContactPointType
}

// To creates an "email" contact point type.
func To(emails []string, opts ...Option) alertmanager.ContactPointOption {
	email := &emailType{
		builder: &sdk.ContactPointType{
			Type: "email",
			Settings: map[string]interface{}{
				"addresses": strings.Join(emails, ","),
			},
			SecureSettings: make(map[string]interface{}),
		},
	}

	for _, opt := range opts {
		opt(email)
	}

	return func(contact *alertmanager.Contact) {
		contact.Builder.GrafanaManagedReceivers = append(contact.Builder.GrafanaManagedReceivers, *email.builder)
	}
}

// Single send a single email to all recipients.
func Single() Option {
	return func(contactType *emailType) {
		contactType.builder.Settings["singleEmail"] = true
	}
}

// Message sets an optional message that will be included in the email.
// Variables are allowed.
func Message(content string) Option {
	return func(contactType *emailType) {
		contactType.builder.Settings["message"] = content
	}
}
