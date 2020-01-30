package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/text"
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
		grabana.WithRow(
			"Prometheus",
			grabana.WithGraph(
				"HTTP Rate",
				graph.DataSource("prometheus-default"),
				graph.WithPrometheusTarget(graph.PrometheusTarget{
					RefID:        "A",
					Expr:         "rate(prometheus_http_requests_total[30s])",
					Format:       "time_series",
					LegendFormat: "{{handler}} - {{ code }}",
				}),
			),
		),
		grabana.WithRow(
			"Some text, because it might be useful",
			grabana.WithText(
				"Some awesome text?",
				text.Markdown("# Title\n\nFor markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)\n"),
			),
			grabana.WithText(
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
