package alertmanager

import (
	"github.com/K-Phoen/sdk"
)

type ContactPointOption func(contactPoint *Contact)

type Contact struct {
	Builder   *sdk.ContactPoint
	IsDefault bool
}

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

func Default() ContactPointOption {
	return func(contactPoint *Contact) {
		contactPoint.IsDefault = true
	}
}
