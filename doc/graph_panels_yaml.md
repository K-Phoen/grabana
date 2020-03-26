# Graph panels

> The main panel in Grafana is simply named Graph. It provides a very rich set
> of graphing options.
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
        alert:
          title: Too many heap allocations
          evaluate_every: 1m
          for: 1m
          # ID of the notification channel
          notify: 1
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
```

## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
