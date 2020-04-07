package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/axis"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run -mod=vendor main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], grabana.WithAPIToken(os.Args[2]))

	// create the folder holding the dashboard for the service
	folder, err := client.GetFolderByTitle(ctx, "Test Folder")
	if err != nil && err != grabana.ErrFolderNotFound {
		fmt.Printf("Could not create folder: %s\n", err)
		os.Exit(1)
	}
	if folder == nil {
		folder, err = client.CreateFolder(ctx, "Test Folder")
		if err != nil {
			fmt.Printf("Could not create folder: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Folder created (id: %d, uid: %s)\n", folder.ID, folder.UID)
	}

	channel, err := client.GetAlertChannelByName(ctx, "Mail")
	if err != nil {
		fmt.Printf("Could not find notification channel 'Mail': %s\n", err)
		os.Exit(1)
	}

	builder := dashboard.New(
		"Awesome dashboard",
		dashboard.AutoRefresh("5s"),
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
		),
		dashboard.VariableAsQuery(
			"status",
			query.DataSource("prometheus-default"),
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
				graph.DataSource("prometheus-default"),
				graph.WithPrometheusTarget(
					"rate(promhttp_metric_handler_requests_total[$interval])",
					prometheus.Legend("{{handler}} - {{ code }}"),
				),
			),
			row.WithGraph(
				"Heap allocations",
				graph.Span(6),
				graph.Height("400px"),
				graph.DataSource("prometheus-default"),
				graph.WithPrometheusTarget("go_memstats_heap_alloc_bytes", prometheus.Ref("A")),
				graph.LeftYAxis(axis.Unit("bytes"), axis.Label("memory"), axis.Min(0)),
				graph.Legend(graph.Current, graph.NoNullSeries, graph.NoZeroSeries, graph.AsTable),
				graph.Alert(
					"Too many heap allocations",
					alert.If(
						alert.And,
						alert.Avg("A", "1m", "now"),
						alert.IsAbove(35000000),
					),
					alert.EvaluateEvery("1m"),
					alert.For("1m"),
					alert.Notify(channel),
				),
			),
			row.WithTable(
				"Threads",
				table.WithPrometheusTarget("go_threads"),
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
				singlestat.WithPrometheusTarget("go_memstats_heap_alloc_bytes"),
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
