package sdk

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

type AlertManager struct {
	Config        AlertManagerConfig `json:"alertmanager_config"`
	TemplateFiles MessageTemplate    `json:"template_files"`
}

type ContactPoint struct {
	Name                    string             `json:"name"`
	GrafanaManagedReceivers []ContactPointType `json:"grafana_managed_receiver_configs,omitempty"`
}

type ContactPointType struct {
	UID                   string                 `json:"uid,omitempty"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	DisableResolveMessage bool                   `json:"disableResolveMessage"`
	Settings              map[string]interface{} `json:"settings"`
	SecureSettings        map[string]interface{} `json:"secureSettings,omitempty"`
}

type AlertManagerConfig struct {
	Receivers []ContactPoint       `json:"receivers"`
	Route     NotificationPolicies `json:"route"`
	Templates []string             `json:"templates"`
}

// Map of template name → template
type MessageTemplate map[string]string

type NotificationPolicies struct {
	// Default alert receiver
	Receiver string `json:"receiver"`

	// Default group bys
	GroupBy []string `json:"group_by,omitempty"`

	// Default timing settings
	GroupInterval  string `json:"group_interval,omitempty"`
	GroupWait      string `json:"group_wait,omitempty"`
	RepeatInterval string `json:"repeat_interval,omitempty"`

	// Routing policies
	Routes []NotificationRoutingPolicy `json:"routes,omitempty"`
}

type NotificationRoutingPolicy struct {
	// Alert receiver
	Receiver string `json:"receiver"`
	// Default timing settings overrides
	GroupInterval  string `json:"group_interval,omitempty"`
	GroupWait      string `json:"group_wait,omitempty"`
	RepeatInterval string `json:"repeat_interval,omitempty"`

	ObjectMatchers []AlertObjectMatcher `json:"object_matchers"`
}

type AlertObjectMatcher [3]string
