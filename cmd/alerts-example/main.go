package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/sdk"

	"github.com/K-Phoen/grabana"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], grabana.WithAPIToken(os.Args[2]))

	nope := false

	alertDefinition := sdk.Alert{
		Name:     "[grabana] GrafanaDashboards reconciliations status",
		Interval: "1m",
		Rules: []sdk.AlertRule{
			{
				For: "5m",
				Annotations: map[string]string{
					"__dashboardUid__": "dark-reconciliations",
					"__panelId__":      "1",
				},
				Labels: map[string]string{
					"owner":         "platform",
					"__managedBy__": "grabana",
				},
				GrafanaAlert: &sdk.GrafanaAlert{
					Title:               "[grabana] GrafanaDashboards reconciliations status",
					Condition:           "B",
					NoDataState:         "NoData",
					ExecutionErrorState: "Alerting",
					Data: []sdk.AlertQuery{
						{
							RefID:         "A",
							DatasourceUID: "PBFA97CFB590B2093",
							QueryType:     "",
							RelativeTimeRange: &sdk.AlertRelativeTimeRange{
								From: 3600,
								To:   0,
							},
							Model: sdk.AlertModel{
								RefID:        "A",
								Expr:         "sum(increase(controller_runtime_reconcile_total{controller=\"grafanadashboard\"}[1m])) by (result)",
								Format:       "time_series",
								LegendFormat: "{{ result }}",
								Datasource: sdk.AlertDatasourceRef{
									UID:  "PBFA97CFB590B2093",
									Type: "prometheus",
								},
								Interval:   "",
								IntervalMs: 15000,
							},
						},

						{
							RefID:         "B",
							DatasourceUID: "-100",
							QueryType:     "",
							Model: sdk.AlertModel{
								RefID: "B",
								Type:  "classic_conditions",
								Hide:  &nope,
								Datasource: sdk.AlertDatasourceRef{
									UID:  "-100",
									Type: "__expr__",
								},
								Conditions: []sdk.AlertCondition{
									{
										Type:      "query",
										Evaluator: sdk.AlertEvaluator{Type: "gt", Params: []float64{3}},
										Operator:  sdk.AlertOperator{Type: "and"},
										Query:     sdk.AlertConditionQueryRef{Params: []string{"A"}},
										Reducer:   sdk.AlertReducer{Type: "last", Params: []string{}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := client.AddAlert(ctx, "DARK", alertDefinition); err != nil {
		fmt.Printf("Could not add alert: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("The deed is done.")
}
