package alert

import (
	"github.com/K-Phoen/sdk"
)

// ErrorMode represents the behavior of an alert in case of execution error.
type ErrorMode string

// Alerting will set the alert state to "alerting".
const ErrorAlerting ErrorMode = "Alerting"

// LastState will set the alert state to its previous state.
const ErrorKO ErrorMode = "Error"

// LastState will set the alert state to its previous state.
const ErrorOK ErrorMode = "OK"

// NoDataMode represents the behavior of an alert when no data is returned by
// the query.
type NoDataMode string

// NoData will set the alert state to "no data".
const NoDataEmpty NoDataMode = "NoData"

// Error will set the alert state to "alerting".
const NoDataAlerting NoDataMode = "Alerting"

// OK will set the alert state to "ok".
const NoDataOK NoDataMode = "OK"

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

const alertConditionRef = "_alert_condition_"

// Alert represents an alert that can be triggered by a query.
type Alert struct {
	Builder *sdk.Alert

	// For internal use only
	Datasource   string
	DashboardUID string
	PanelID      string
}

// New creates a new alert.
func New(name string, options ...Option) *Alert {
	nope := false

	alert := &Alert{
		Builder: &sdk.Alert{
			Name: name,
			Rules: []sdk.AlertRule{
				{
					GrafanaAlert: &sdk.GrafanaAlert{
						Title:     name,
						Condition: alertConditionRef,
						Data: []sdk.AlertQuery{
							{
								RefID:         alertConditionRef,
								QueryType:     "",
								DatasourceUID: "-100",
								Model: sdk.AlertModel{
									RefID: alertConditionRef,
									Type:  "classic_conditions",
									Hide:  &nope,
									Datasource: sdk.AlertDatasourceRef{
										UID:  "-100",
										Type: "__expr__",
									},
									Conditions: []sdk.AlertCondition{},
								},
							},
						},
					},
					Annotations: map[string]string{},
					Labels:      map[string]string{},
				},
			},
		},
	}

	for _, opt := range append(defaults(), options...) {
		opt(alert)
	}

	return alert
}

func defaults() []Option {
	return []Option{
		EvaluateEvery("1m"),
		For("5m"),
		OnNoData(NoDataEmpty),
		OnExecutionError(ErrorAlerting),
	}
}

func (alert *Alert) HookDatasourceUID(uid string) {
	for _, rule := range alert.Builder.Rules {
		for i := range rule.GrafanaAlert.Data {
			query := &rule.GrafanaAlert.Data[i]

			if query.RefID == alertConditionRef {
				continue
			}

			query.DatasourceUID = uid
			query.Model.Datasource.UID = uid
		}
	}
}

func (alert *Alert) HookDashboardUID(uid string) {
	for _, rule := range alert.Builder.Rules {
		rule.Annotations["__dashboardUid__"] = uid
	}
}

func (alert *Alert) HookPanelID(id string) {
	for _, rule := range alert.Builder.Rules {
		rule.Annotations["__panelId__"] = id
	}
}

// Summary sets the summary associated to the alert.
func Summary(content string) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].Annotations["summary"] = content
	}
}

// Description sets the description associated to the alert.
func Description(content string) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].Annotations["description"] = content
	}
}

// Runbook sets the runbook URL associated to the alert.
func Runbook(url string) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].Annotations["runbook_url"] = url
	}
}

// For sets the time interval during which a query violating the threshold
// before the alert being actually triggered.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#for
func For(duration string) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].For = duration
	}
}

// EvaluateEvery defines the evaluation interval.
func EvaluateEvery(interval string) Option {
	return func(alert *Alert) {
		alert.Builder.Interval = interval
	}
}

// OnExecutionError defines the behavior on execution error.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#execution-errors-or-timeouts
func OnExecutionError(mode ErrorMode) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].GrafanaAlert.ExecutionErrorState = string(mode)
	}
}

// OnNoData defines the behavior when the query returns no data.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#no-data-null-values
func OnNoData(mode NoDataMode) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].GrafanaAlert.NoDataState = string(mode)
	}
}

// If defines a single condition that will trigger the alert.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#conditions
func If(reducer QueryReducer, queryRef string, evaluator ConditionEvaluator) Option {
	return ifOperand(And, reducer, queryRef, evaluator)
}

// IfOr defines a single condition that will trigger the alert.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#conditions
func IfOr(reducer QueryReducer, queryRef string, evaluator ConditionEvaluator) Option {
	return ifOperand(Or, reducer, queryRef, evaluator)
}

func ifOperand(operand Operator, reducer QueryReducer, queryRef string, evaluator ConditionEvaluator) Option {
	return func(alert *Alert) {
		cond := newCondition(reducer, queryRef, evaluator)
		cond.builder.Operator = sdk.AlertOperator{Type: string(operand)}

		alert.Builder.Rules[0].GrafanaAlert.Data[0].Model.Conditions = append(alert.Builder.Rules[0].GrafanaAlert.Data[0].Model.Conditions, *cond.builder)
	}
}

// Tags defines a set of tags that will be forwarded to the notifications
// channels when the alert will tbe triggered or used to route the alert.
func Tags(tags map[string]string) Option {
	return func(alert *Alert) {
		alert.Builder.Rules[0].Labels = tags
	}
}
