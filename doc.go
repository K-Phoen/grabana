/*
Package grabana provides a developer-friendly way of creating Grafana dashboards.

Whether you prefer writing **code or YAML**, if you are looking for a way to
version your dashboards configuration or automate tedious and error-prone
creation of dashboards, this library is meant for you.

	builder := dashboard.New(
		"Awesome dashboard",
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
		dashboard.Row(
			"Prometheus",
			row.WithGraph(
				"HTTP Rate",
				graph.WithPrometheusTarget(
					"rate(promhttp_metric_handler_requests_total[$interval])",
					prometheus.Legend("{{handler}} - {{ code }}"),
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
				singlestat.WithPrometheusTarget("go_memstats_heap_alloc_bytes"),
			),
		),
		dashboard.Row(
			"Some text, because it might be useful",
			row.WithText(
				"Some awesome html?",
				text.HTML("<b>lalalala</b>"),
			),
		),
	)

For a more information visit https://github.com/K-Phoen/grabana
*/
package grabana
