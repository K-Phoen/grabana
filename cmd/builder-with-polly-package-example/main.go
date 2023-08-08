package main

import (
	"fmt"
	"os"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/packages"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "Usage: go run main.go path-to-package\n")
		os.Exit(1)
	}

	pkg, err := packages.LoadFile(os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, "could not load package: %w\n", err)
		os.Exit(1)
	}

	builder, err := dashboard.New(
		"Awesome MySQL dashboard",
		dashboard.UID("test-packages-package"),
		dashboard.AutoRefresh("30s"),
		dashboard.Time("now-30m", "now"),
		dashboard.Tags([]string{"generated", "packages", "package"}),

		dashboard.Row(
			"Some stuff",

			// Complete panel taken from a package
			row.WithPanelFromPackage(pkg, packages.PanelRef("mysql-open-tables")),

			// Manually-defined panel, with a target taken from a package
			row.WithTimeSeries(
				"Tables definitions",
				timeseries.Span(6),
				timeseries.Height("400px"),
				timeseries.DataSource("Prometheus"),
				timeseries.WithTargetFromPackage(pkg, packages.TargetRef("mysql-table-definition-cache-2")),
			),

			// Completely manually-defined panel
			row.WithTimeSeries(
				"HTTP Rate",
				timeseries.Span(6),
				timeseries.Height("400px"),
				timeseries.DataSource("Prometheus"),
				timeseries.WithPrometheusTarget(
					"sum(rate(promhttp_metric_handler_requests_total[$interval])) by (code)",
					prometheus.Legend("{{ code }}"),
				),
			),
		),
	)
	if err != nil {
		fmt.Printf("Could not build dashboard: %s\n", err)
		os.Exit(1)
	}

	dashboardJSON, err := builder.MarshalIndentJSON()
	if err != nil {
		fmt.Printf("Could not marshal dashboard: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(dashboardJSON))
}
