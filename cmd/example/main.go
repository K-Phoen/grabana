package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana/variable/interval"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run -mod=vendor main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], os.Args[2])

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

	dashboard := grabana.NewDashboardBuilder(
		"Awesome dashboard",
		grabana.AutoRefresh("5s"),
		grabana.WithTags([]string{"generated"}),
		grabana.WithTagsAnnotation(grabana.TagAnnotation{
			Name:       "Deployments",
			Datasource: "-- Grafana --",
			IconColor:  "#5794F2",
			Tags:       []string{"deploy", "production"},
		}),
		grabana.WithVariableAsInterval(
			"interval",
			interval.Values([]string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"}),
		),
		grabana.WithVariableAsConst(
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
			constant.Default("80th"),
		),
		grabana.WithVariableAsCustom(
			"vX",
			custom.Multi(),
			custom.IncludeAll(),
			custom.Values(map[string]string{
				"v1": "v1",
				"v2": "v2",
			}),
			custom.Default("v2"),
		),
		grabana.WithRow(
			"Prometheus",
			row.WithGraph(
				"HTTP Rate",
				graph.Height("400px"),
				graph.Span(12),
				graph.DataSource("prometheus-default"),
				graph.WithPrometheusTarget(
					"rate(prometheus_http_requests_total[30s])",
					prometheus.WithLegend("{{handler}} - {{ code }}"),
				),
			),
		),
		grabana.WithRow(
			"Some text, because it might be useful",
			row.WithText(
				"Some awesome text?",
				text.Markdown("Markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)\n${percentile}"),
			),
			row.WithText(
				"Some awesome html?",
				text.HTML("<b>lalalala</b>"),
			),
		),
	)
	if _, err := client.UpsertDashboard(ctx, folder, dashboard); err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("The deed is done.")
}
