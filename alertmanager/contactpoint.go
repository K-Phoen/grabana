package alertmanager

import (
	"github.com/K-Phoen/sdk"
)

// ContactPointOption represents an option that can be used to configure a
// contact point.
type ContactPointOption func(contactPoint *Contact)

// Contact describes a contact point.
type Contact struct {
	Builder *sdk.ContactPoint
}

// ContactPoint defines a new contact point.
func ContactPoint(name string, opts ...ContactPointOption) Contact {
	contactPoint := &Contact{
		Builder: &sdk.ContactPoint{
			Name: name,
		},
	}

	for _, opt := range opts {
		opt(contactPoint)
	}

	return *contactPoint
}
