# Grabana

![CI](https://github.com/K-Phoen/grabana/workflows/CI/badge.svg) [![codecov](https://codecov.io/gh/K-Phoen/grabana/branch/master/graph/badge.svg)](https://codecov.io/gh/K-Phoen/grabana) [![GoDoc](https://godoc.org/github.com/K-Phoen/grabana?status.svg)](https://godoc.org/github.com/K-Phoen/grabana)

Grabana provides a developer-friendly way of creating Grafana dashboards.

If you are looking for a way to version your dashboards configuration or
automate tedious and error-prone creation of dashboards, this library is meant
for you.

## Design goals

* provide an understandable abstraction over dashboards configuration
* expose a developer-friendly API
* allow IDE assistance and auto-completion

## Example

Dashboard configuration:

```go
dashboard := grabana.NewDashboardBuilder(
    "Awesome dashboard",
    grabana.AutoRefresh("5s"),
    grabana.Tags([]string{"generated"}),
    grabana.VariableAsInterval(
        "interval",
        interval.Values([]string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"}),
    ),
    grabana.Row(
        "Prometheus",
        row.WithGraph(
            "HTTP Rate",
            graph.DataSource("prometheus-default"),
            graph.WithPrometheusTarget(
                "rate(prometheus_http_requests_total[30s])",
                prometheus.Legend("{{handler}} - {{ code }}"),
            ),
        ),
    ),
)
```

Dashboard creation:

```go
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

if _, err := client.UpsertDashboard(ctx, folder, dashboard); err != nil {
    fmt.Printf("Could not create dashboard: %s\n", err)
    os.Exit(1)
}
```

For a more complete example, see the [`example`](./cmd/example/) directory.

## License

This library is under the [MIT](LICENSE) license.
