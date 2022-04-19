package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/alert"
	prometheusAlert "github.com/K-Phoen/grabana/alert/queries/prometheus"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], grabana.WithAPIToken(os.Args[2]))

	// create the folder holding the dashboard for the service
	folder, err := client.FindOrCreateFolder(ctx, "Test Folder")
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}

	builder := dashboard.New(
		"Awesome dashboard test",
		dashboard.UID("test-dashboard-alerts"),
		dashboard.AutoRefresh("30s"),
		dashboard.Tags([]string{"generated"}),
		dashboard.TagsAnnotation(dashboard.TagAnnotation{
			Name:       "Deployments",
			Datasource: "-- Grafana --",
			IconColor:  "#5794F2",
			Tags:       []string{"deploy", "production"},
		}),
		dashboard.VariableAsInterval(
			"interval",
			interval.Values([]string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"}),
			interval.Default("30s"),
		),
		dashboard.VariableAsQuery(
			"status",
			query.DataSource("Prometheus"),
			query.Request("label_values(prometheus_http_requests_total, code)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsConst(
			"percentile",
			constant.Label("Percentile"),
			constant.Values(map[string]string{
				"50th": "50",
				"75th": "75",
				"80th": "80",
				"85th": "85",
				"90th": "90",
				"95th": "95",
				"99th": "99",
			}),
			constant.Default("80"),
		),
		dashboard.VariableAsCustom(
			"vX",
			custom.Multi(),
			custom.IncludeAll(),
			custom.Values(map[string]string{
				"v1": "v1",
				"v2": "v2",
			}),
			custom.Default("v2"),
		),
		dashboard.Row(
			"Prometheus",
			row.WithGraph(
				"HTTP Rate",
				graph.Span(6),
				graph.Height("400px"),
				graph.DataSource("Prometheus"),
				graph.WithPrometheusTarget(
					"sum(rate(promhttp_metric_handler_requests_total[$interval])) by (app, code)",
					prometheus.Legend("{{ app }} - {{ code }}"),
				),
			),
			row.WithTimeSeries(
				"Heap allocations",
				timeseries.Span(6),
				timeseries.Height("400px"),
				timeseries.DataSource("Prometheus"),
				timeseries.WithPrometheusTarget("sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)", prometheus.Legend("{{ app }}")),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
				timeseries.Legend(timeseries.Last, timeseries.AsTable),
				timeseries.Alert(
					fmt.Sprintf("Too many heap allocations"),
					alert.Summary("Yup, too much of {{ app }}"),
					alert.Runbook("https://google.com"),
					alert.Tags(map[string]string{
						"service": "amazing-service",
						"owner":   "team-b",
					}),
					alert.If(
						alert.Avg("A"),
						alert.IsAbove(3),
					),
					alert.Queries(
						alert.WithPrometheusQuery(
							"A",
							"sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
							prometheusAlert.Legend("{{ app }}"),
						),
					),
				),
			),
			row.WithTable(
				"Threads",
				table.WithPrometheusTarget("sum(go_threads{app!=\"\"}) by (app)", prometheus.Legend("{{ app }}")),
				table.HideColumn("Time"),
				table.AsTimeSeriesAggregations([]table.Aggregation{
					{Label: "AVG", Type: table.AVG},
					{Label: "Current", Type: table.Current},
				}),
			),
			row.WithSingleStat(
				"Heap Allocations",
				singlestat.Unit("bytes"),
				singlestat.ColorValue(),
				singlestat.WithPrometheusTarget("sum(go_memstats_heap_alloc_bytes)"),
				singlestat.Thresholds([2]string{"26000000", "28000000"}),
			),
		),
		dashboard.Row(
			"Some text, because it might be useful",
			row.WithText(
				"Some awesome text?",
				text.Markdown("Markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)\n${percentile}"),
			),
			row.WithText(
				"Some awesome html?",
				text.HTML("Some awesome html?"),
			),
		),
	)
	if _, err := client.UpsertDashboard(ctx, folder, builder); err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("The deed is done.")
}
