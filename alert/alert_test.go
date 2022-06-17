package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAlertCanBeCreated(t *testing.T) {
	req := require.New(t)
	alertTitle := "some alert"

	a := New(alertTitle)

	req.Len(a.Builder.Rules, 1)

	req.Equal(alertTitle, a.Builder.Name)
	req.Equal(alertTitle, a.Builder.Rules[0].GrafanaAlert.Title)

	req.Equal(string(NoDataEmpty), a.Builder.Rules[0].GrafanaAlert.NoDataState)
	req.Equal(string(ErrorAlerting), a.Builder.Rules[0].GrafanaAlert.ExecutionErrorState)
}

func TestConditionsCanBeCombined(t *testing.T) {
	req := require.New(t)

	a := New(
		"",
		IfOr(Avg, "A", IsBelow(10)),
		IfOr(Avg, "B", IsBelow(8)),
	)

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}

func TestPanelIDCanBeHooked(t *testing.T) {
	req := require.New(t)

	a := New("")

	a.HookPanelID("id")

	req.Equal("id", a.Builder.Rules[0].Annotations["__panelId__"])
}

func TestDashboardUIDCanBeHooked(t *testing.T) {
	req := require.New(t)

	a := New("")

	a.HookDashboardUID("uid")

	req.Equal("uid", a.Builder.Rules[0].Annotations["__dashboardUid__"])
}

func TestDatasourceUIDCanBeHooked(t *testing.T) {
	req := require.New(t)

	a := New(
		"",
		WithPrometheusQuery("A", "some prometheus query"),
		IfOr(Avg, "1", IsBelow(10)),
	)
	a.HookDatasourceUID("ds-uid")

	hooked := false
	for _, rule := range a.Builder.Rules {
		for i := range rule.GrafanaAlert.Data {
			query := &rule.GrafanaAlert.Data[i]

			if query.RefID == alertConditionRef {
				req.Equal("-100", query.DatasourceUID)
				req.Equal("-100", query.Model.Datasource.UID)
				continue
			}

			hooked = true

			req.Equal("ds-uid", query.DatasourceUID)
			req.Equal("ds-uid", query.Model.Datasource.UID)
		}
	}

	req.True(hooked)
}

func TestSummaryCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", Summary("summary content"))

	req.Equal("summary content", a.Builder.Rules[0].Annotations["summary"])
}

func TestDescriptionCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", Description("description content"))

	req.Equal("description content", a.Builder.Rules[0].Annotations["description"])
}

func TestRunbookCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", Runbook("runbook url"))

	req.Equal("runbook url", a.Builder.Rules[0].Annotations["runbook_url"])
}

func TestForIntervalCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", For("1m"))

	req.Equal("1m", a.Builder.Rules[0].For)
}

func TestFrequencyCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", EvaluateEvery("1m"))

	req.Equal("1m", a.Builder.Interval)
}

func TestErrorModeCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", OnExecutionError(ErrorKO))

	req.Equal(string(ErrorKO), a.Builder.Rules[0].GrafanaAlert.ExecutionErrorState)
}

func TestNoDataModeCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", OnNoData(NoDataAlerting))

	req.Equal(string(NoDataAlerting), a.Builder.Rules[0].GrafanaAlert.NoDataState)
}

func TestTagsCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", Tags(map[string]string{
		"severity": "warning",
	}))

	req.Len(a.Builder.Rules[0].Labels, 1)
	req.Equal("warning", a.Builder.Rules[0].Labels["severity"])
}

func TestConditionsCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", If(Avg, "1", IsBelow(10)))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}

func TestOrConditionsCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", IfOr(Avg, "1", IsBelow(10)))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}
