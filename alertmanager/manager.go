package alertmanager

import (
	"encoding/json"

	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure an
// alert manager.
type Option func(manager *Manager)

// Manager represents an alert manager.
type Manager struct {
	builder *sdk.AlertManager
}

// New creates a new alert manager.
func New(opts ...Option) *Manager {
	manager := &Manager{
		builder: &sdk.AlertManager{},
	}

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

func ContactPoints(contactPoints ...Contact) Option {
	return func(manager *Manager) {
		config := &manager.builder.Config
		config.Receivers = nil

		for i, point := range contactPoints {
			config.Receivers = append(config.Receivers, *point.Builder)

			// we must have a default contact point, so we either use one that
			// explicitly was indicated as such, or the first one.
			if point.IsDefault || i == 0 {
				config.Route.Receiver = point.Builder.Name
			}
		}
	}
}

func Routing(policies ...RoutingPolicy) Option {
	return func(manager *Manager) {
		config := &manager.builder.Config
		config.Route.Routes = nil

		for _, policy := range policies {
			config.Route.Routes = append(config.Route.Routes, *policy.builder)
		}
	}
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (manager *Manager) MarshalJSON() ([]byte, error) {
	return json.Marshal(manager.builder)
}

// MarshalIndentJSON renders the manager as indented JSON.
func (manager *Manager) MarshalIndentJSON() ([]byte, error) {
	return json.MarshalIndent(manager.builder, "", "  ")
}
