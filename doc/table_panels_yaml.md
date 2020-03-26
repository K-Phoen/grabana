# Table panels

> The table panel is very flexible, supporting both multiple modes for time
> series as well as for table, annotation and raw JSON data. It also provides
> date formatting and value formatting and coloring options.
>
> — https://grafana.com/docs/grafana/latest/features/panels/table_panel/

```yaml
rows:
  - name: "Table panels row"
    panels:
      - table:
        title: Threads
        span: 12
        height: 400px
        datasource: prometheus-default
        targets:
          - prometheus:
              query: "go_threads"
        # hides the column having a label matching the given pattern.
        hidden_columns: ["Time"]
        time_series_aggregations:
          - label: AVG
            # valid types are: avg, count, current, min, max
            type: avg
          - label: Current
            type: current
```

## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
