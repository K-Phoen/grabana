package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/variable/constant"
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
		grabana.WithTags([]string{"generated"}),
		grabana.WithTagsAnnotation(grabana.TagAnnotation{
			Name:       "Deployments",
			Datasource: "-- Grafana --",
			IconColor:  "#5794F2",
			Tags:       []string{"deploy", "production"},
		}),
		grabana.WithVariableAsConst(
			"percentile",
			constant.WithLabel("Percentile"),
			constant.WithValues([]constant.Value{
				{Text: "50", Value: "50"},
				{Text: "75", Value: "75"},
				{Text: "80", Value: "80"},
				{Text: "85", Value: "85"},
				{Text: "90", Value: "90"},
				{Text: "95", Value: "95"},
				{Text: "99", Value: "99"},
			}),
			constant.WithDefault("80"),
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
				text.Markdown("# Title\n\nFor markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)\n"),
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
