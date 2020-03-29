# Singlestat panels

> The Singlestat Panel allows you to show the one main summary stat of a SINGLE
> series. It reduces the series into a single number (by looking at the max,
> min, average, or sum of values in the series). Singlestat also provides
> thresholds to color the stat or the Panel background. It can also translate
> the single number into a text value, and show a sparkline summary of the series.
>
> — https://grafana.com/docs/grafana/latest/features/panels/singlestat/

```yaml
rows:
  - name: "Graph panels row"
    panels:
      - single_stat:
          title: Heap Allocations
          span: 12
          height: 400px
          datasource: prometheus-default
          targets:
            - prometheus:
                query: 'go_memstats_heap_alloc_bytes{job="prometheus"}'
          unit: bytes
          # valid values are: min, max, avg, current, total, first, delta, diff, range
          value_type: avg
          thresholds: ["26000000", "28000000"]
          # valid values are: value, background
          color: ["value"]
          # valid values are: bottom, full
          sparkline: bottom
```

## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
