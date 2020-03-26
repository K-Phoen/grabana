# Variables

> Variables allows for more interactive and dynamic dashboards. Instead of
> hard-coding things like server, application and sensor name in your metric
> queries you can use variables in their place. Variables are shown as
> dropdown select boxes at the top of the dashboard. These dropdowns
> allow changing the data being displayed in your dashboard.
>
> — https://grafana.com/docs/grafana/latest/reference/templating/

```yaml
variables:
  - interval:
      name: interval
      label: Interval
      values: ["30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"]
  - query:
      name: status
      label: HTTP status
      datasource: prometheus-default
      request: "label_values(prometheus_http_requests_total, code)"
  - const:
      name: percentile
      label: Percentile
      default: 80
      values_map:
        50th: "50"
        75th: "75"
        80th: "80"
        85th: "85"
        90th: "90"
        95th: "95"
        99th: "99"
  - custom:
      name: vX
      default: v2
      values_map:
        v1: v1
        v2: v2
```

## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
