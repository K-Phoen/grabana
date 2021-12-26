package opsgenie

import (
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure an "opsgenie"
// contact point type.
type Option func(contactType *opsgenieType)

// TagForwardMode describes how alert tags should be forwarded to Opsgenie.
type TagForwardMode string

const (
	Tags                   TagForwardMode = "tags"
	ExtraProperties        TagForwardMode = "details"
	TagsAndExtraProperties TagForwardMode = "both"
)

type opsgenieType struct {
	builder *sdk.ContactPointType
}

// With creates an Opsgenie contact point type with the given settings.
func With(apiURL string, apiKey string, opts ...Option) alertmanager.ContactPointOption {
	opsgenie := &opsgenieType{
		builder: &sdk.ContactPointType{
			Type: "opsgenie",
			Settings: map[string]interface{}{
				"apiUrl": apiURL,
			},
			SecureSettings: map[string]interface{}{
				"apiKey": apiKey,
			},
		},
	}

	defaultOpts := []Option{SentTagsAs(Tags)}
	for _, opt := range append(defaultOpts, opts...) {
		opt(opsgenie)
	}

	return func(contact *alertmanager.Contact) {
		contact.Builder.GrafanaManagedReceivers = append(contact.Builder.GrafanaManagedReceivers, *opsgenie.builder)
	}
}

// AutoClose automatically closes an alert in Opsgenie once it goes back to OK in Grafana.
func AutoClose() Option {
	return func(contactType *opsgenieType) {
		contactType.builder.Settings["autoClose"] = true
	}
}

// OverridePriority allows the alert priority to be set in Opsgenie based on
// the content of the `og_priority` annotation.
func OverridePriority() Option {
	return func(contactType *opsgenieType) {
		contactType.builder.Settings["overridePriority"] = true
	}
}

// SentTagsAs defines how alert tags should be forwarded to Opsgenie.
func SentTagsAs(mode TagForwardMode) Option {
	return func(contactType *opsgenieType) {
		contactType.builder.Settings["sendTagsAs"] = string(mode)
	}
}
