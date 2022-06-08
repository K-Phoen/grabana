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

// ContactPoints defines the contact points that can receive alerts.
func ContactPoints(contactPoints ...Contact) Option {
	return func(manager *Manager) {
		config := &manager.builder.Config
		config.Receivers = nil

		for i, point := range contactPoints {
			config.Receivers = append(config.Receivers, *point.Builder)

			// we must have a default contact point, so we use the first contact point
			// if none is already set.
			if i == 0 && config.Route.Receiver == "" {
				config.Route.Receiver = point.Builder.Name
			}
		}
	}
}

// DefaultContactPoint sets the default contact point to be used when no
// specific routing policy applies.
func DefaultContactPoint(contactPoint string) Option {
	return func(manager *Manager) {
		manager.builder.Config.Route.Receiver = contactPoint
	}
}

// Templates defines templates that can be used when sending messages to
// contact points.
// See https://prometheus.io/blog/2016/03/03/custom-alertmanager-templates/
func Templates(templates map[string]string) Option {
	return func(manager *Manager) {
		manager.builder.TemplateFiles = templates
	}
}

// Routing configures the routing policies to apply on alerts.
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
