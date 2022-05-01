package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/alert"
	"github.com/stretchr/testify/require"
)

func TestDecodingSimpleAlert(t *testing.T) {
	req := require.New(t)

	threshold := float64(1)
	ref := "A"

	alertDef := Alert{
		Description:      "description",
		Runbook:          "runbook",
		Tags:             map[string]string{"service": "awesome"},
		EvaluateEvery:    "2m",
		For:              "5m",
		OnNoData:         "no_data",
		OnExecutionError: "error",
		If: []AlertCondition{
			{Avg: &ref, Above: &threshold},
		},
		Targets: []AlertTarget{
			{Prometheus: &AlertPrometheus{Ref: "A", Query: "some query"}},
		},
	}

	opts, err := alertDef.toOptions()
	req.NoError(err)

	alertBuilder := alert.New("", opts...)
	rule := alertBuilder.Builder.Rules[0]

	req.Equal("description", rule.Annotations["description"])
	req.Equal("runbook", rule.Annotations["runbook_url"])
	req.Len(rule.Labels, 1)
	req.Equal("awesome", rule.Labels["service"])
	req.Equal("5m", rule.For)
	req.Equal("2m", alertBuilder.Builder.Interval)
	req.Equal("NoData", rule.GrafanaAlert.NoDataState)
	req.Equal("Error", rule.GrafanaAlert.ExecutionErrorState)

	alertData := rule.GrafanaAlert.Data

	req.Len(alertData, 2)

	// The condition was added first by grabana
	req.Equal("_alert_condition_", alertData[0].RefID)
	req.Len(alertData[0].Model.Conditions, 1)
	req.Equal("avg", alertData[0].Model.Conditions[0].Reducer.Type)
	req.Equal("gt", alertData[0].Model.Conditions[0].Evaluator.Type)
	req.ElementsMatch([]float64{threshold}, alertData[0].Model.Conditions[0].Evaluator.Params)

	// The query was added last by grabana
	req.Equal("A", alertData[1].RefID)
	req.Equal("prometheus", alertData[1].Model.Datasource.Type)
}

func TestDecodingAlertFailsIfNoConditionIsDefined(t *testing.T) {
	req := require.New(t)

	alertDef := Alert{
		Targets: []AlertTarget{
			{Prometheus: &AlertPrometheus{Ref: "A", Query: "some query"}},
		},
	}

	_, err := alertDef.toOptions()
	req.ErrorIs(ErrNoConditionOnAlert, err)
}

func TestDecodingAlertFailsIfNoTargetIsDefined(t *testing.T) {
	req := require.New(t)

	threshold := float64(1)
	ref := "A"

	alertDef := Alert{
		If: []AlertCondition{
			{Avg: &ref, Above: &threshold},
		},
	}

	_, err := alertDef.toOptions()
	req.ErrorIs(ErrNoTargetOnAlert, err)
}
