package alert

import (
	"github.com/grafana-tools/sdk"
)

// ErrorMode represents the behavior of an alert in case of execution error.
type ErrorMode string

// Alerting will set the alert state to "alerting".
const Alerting ErrorMode = "alerting"

// LastState will set the alert state to its previous state.
const LastState ErrorMode = "keep_state"

// NoDataMode represents the behavior of an alert when no data is returned by
// the query.
type NoDataMode string

// NoData will set the alert state to "no data".
const NoData NoDataMode = "no_data"

// Error will set the alert state to "alerting".
const Error NoDataMode = "alerting"

// KeepLastState will set the alert state to its previous state.
const KeepLastState NoDataMode = "keep_state"

// OK will set the alert state to "ok".
const OK NoDataMode = "ok"

// Option represents an option that can be used to configure an alert.
type Option func(alert *Alert)

// Channel represents an alert notification channel.
// See https://grafana.com/docs/grafana/latest/alerting/notifications/#notification-channel-setup
type Channel struct {
	ID   uint   `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"Name"`
	Type string `json:"type"`
}

// Alert represents an alert that can be triggered by a query.
type Alert struct {
	Builder *sdk.Alert
}

// New creates a new alert.
func New(name string, options ...Option) *Alert {
	alert := &Alert{Builder: &sdk.Alert{
		Name:                name,
		Handler:             1, // TODO: what's that?
		ExecutionErrorState: string(LastState),
		NoDataState:         string(KeepLastState),
	}}

	for _, opt := range options {
		opt(alert)
	}

	return alert
}

// Notify adds a notification for this alert given a channel.
func Notify(channel *Channel) Option {
	return func(alert *Alert) {
		alert.Builder.Notifications = append(alert.Builder.Notifications, sdk.AlertNotification{
			UID: channel.UID,
		})
	}
}

// NotifyChannels appends the given notification channels to the list of
// channels for this alert.
func NotifyChannels(channels ...*Channel) Option {
	return func(alert *Alert) {
		for _, channel := range channels {
			alert.Builder.Notifications = append(alert.Builder.Notifications, sdk.AlertNotification{
				UID: channel.UID,
			})
		}
	}
}

// NotifyChannel adds a notification for this alert given a channel UID.
func NotifyChannel(channelUID string) Option {
	return func(alert *Alert) {
		alert.Builder.Notifications = append(alert.Builder.Notifications, sdk.AlertNotification{
			UID: channelUID,
		})
	}
}

// Message sets the message associated to the alert.
func Message(content string) Option {
	return func(alert *Alert) {
		alert.Builder.Message = content
	}
}

// For sets the time interval during which a query violating the threshold
// before the alert being actually triggered.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#for
func For(duration string) Option {
	return func(alert *Alert) {
		alert.Builder.For = duration
	}
}

// EvaluateEvery defines the evaluation interval.
func EvaluateEvery(interval string) Option {
	return func(alert *Alert) {
		alert.Builder.Frequency = interval
	}
}

// OnExecutionError defines the behavior on execution error.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#execution-errors-or-timeouts
func OnExecutionError(mode ErrorMode) Option {
	return func(alert *Alert) {
		alert.Builder.ExecutionErrorState = string(mode)
	}
}

// OnNoData defines the behavior when the query returns no data.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#no-data-null-values
func OnNoData(mode NoDataMode) Option {
	return func(alert *Alert) {
		alert.Builder.NoDataState = string(mode)
	}
}

// If adds a condition that could trigger the alert.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#conditions
func If(operator Operator, opts ...ConditionOption) Option {
	return func(alert *Alert) {
		cond := newCondition(opts...)
		cond.builder.Operator = sdk.AlertOperator{Type: string(operator)}

		alert.Builder.Conditions = append(alert.Builder.Conditions, *cond.builder)
	}
}
