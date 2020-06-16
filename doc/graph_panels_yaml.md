# Graph panels

> The main panel in Grafana is named Graph. It provides a very rich set of
> graphing options.
>
> — https://grafana.com/docs/grafana/latest/features/panels/graph/

```yaml
rows:
  - name: "Graph panels row"
    panels:
      - graph:
        title: HTTP Rate
        height: 400px
        span: 16
        datasource: prometheus-default
        targets:
          - prometheus:
              query: "rate(promhttp_metric_handler_requests_total[$interval])"
              legend: "{{handler}} - {{ code }}"
        axes:
          left: { unit: short, min: 0, label: Requests }
          right: { hidden: true }

      - graph:
        title: Heap allocations
        height: 400px
        span: 16
        datasource: prometheus-default
        targets:
          - prometheus:
              query: "go_memstats_heap_alloc_bytes"
              legend: "{{job}}"
              ref: A
              #hidden: true # useful for queries only referenced in alerts
        # Valid values are: hide, as_table, to_the_right, min, max, avg, current, total, no_null_series, no_zero_series
        legend: [avg, current, no_null_series, no_zero_series]
        alert:
          title: Too many heap allocations
          evaluate_every: 1m
          for: 1m
          # UID of the notification channel
          notify: "P-N3fxuZz"
          # UIDs of the notification channels
          #notifications: ["P-N3fxuZz"]
          message: "Wow, a we're allocating a lot."
          # Valid values are: no_data, alerting, keep_state, ok
          on_no_data: alerting
          # Valid values are: alerting, keep_state
          on_execution_error: alerting
          if:
            - operand: and
              # valid `func` values are: avg, sum, count, last, min, max, median, diff, percent_diff
              value: {func: avg, ref: A, from: 1m, to: now}
              threshold: {above: 23000000}
              # threshold: {has_no_value: true}
              # threshold: {below: 23000000}
              # threshold: {outside_range: [23000000, 26000000]}
              # threshold: {within_range: [23000000, 26000000]}
```

## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
