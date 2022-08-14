package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/gauge"
	"github.com/K-Phoen/grabana/row"
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

	builder, err := dashboard.New(
		"Grabana - Gauge example",
		dashboard.AutoRefresh("30s"),
		dashboard.Time("now-30m", "now"),
		dashboard.Row(
			"Kubernetes",
			row.WithGauge(
				"Cluster Pod Usage",
				gauge.Span(6),
				gauge.Height("400px"),
				gauge.Unit("percentunit"),
				gauge.Decimals(2),
				gauge.AbsoluteThresholds([]gauge.ThresholdStep{
					{Color: "#299c46"},
					{Color: "rgba(237, 129, 40, 0.89)", Value: float64Ptr(0.8)},
					{Color: "#d44a3a", Value: float64Ptr(0.9)},
				}),
				gauge.WithPrometheusTarget(
					"sum(kube_pod_info{}) / sum(kube_node_status_allocatable{resource=\"pods\"})",
				),
			),
		),
	)
	if err != nil {
		fmt.Printf("Could not build dashboard: %s\n", err)
		os.Exit(1)
	}

	dash, err := client.UpsertDashboard(ctx, folder, builder)
	if err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("The deed is done:\n%s\n", os.Args[1]+dash.URL)
}

func float64Ptr(input float64) *float64 {
	return &input
}
