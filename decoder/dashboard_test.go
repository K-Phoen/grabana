package decoder

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type testCase struct {
	name                string
	yaml                string
	expectedGrafanaJSON string
}

func TestUnmarshalYAMLWithInvalidInput(t *testing.T) {
	_, err := UnmarshalYAML(bytes.NewBufferString(""))

	require.Error(t, err)
}

func TestUnmarshalYAML(t *testing.T) {
	testCases := []testCase{
		generalOptions(),
		tagAnnotations(),
		variables(),
		textPanel(),
		graphPanel(),
		singleStatPanel(),
		tablePanel(),
		graphPanelWithStackdriverTarget(),
		heatmapPanel(),
		graphPanelWithGraphiteTarget(),
		graphPanelWithInfluxdbTarget(),
		timeseriesPanel(),
		collapseRow(),
		logsPanel(),
		statPanel(),
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			builder, err := UnmarshalYAML(bytes.NewBufferString(tc.yaml))
			req.NoError(err)

			json, err := builder.MarshalJSON()
			req.NoError(err)

			req.JSONEq(dashboardFromFixtures(t, "testdata/"+tc.expectedGrafanaJSON), string(json), "test file "+tc.expectedGrafanaJSON)
		})
	}
}

func TestUnmarshalYAMLWithInvalidTimezone(t *testing.T) {
	payload := `timezone: invalid`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidTimezone, err)
}

func TestUnmarshalYAMLWithInvalidPanel(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - {}`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrPanelNotConfigured, err)
}

func TestUnmarshalYAMLWithInvalidVariable(t *testing.T) {
	payload := `
variables:
  - {}`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrVariableNotConfigured, err)
}

func TestUnmarshalYAMLWithNoTargetTable(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - table:
          title: Threads
          targets:
            - {}
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrTargetNotConfigured, err)
}

func TestUnmarshalYAMLWithNoTargetSingleStat(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - single_stat:
          title: Threads
          targets:
            - {}
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrTargetNotConfigured, err)
}

func TestUnmarshalYAMLWithInvalidSparklineModeSingleStat(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - single_stat:
          title: Threads
          sparkline: unknown-mode
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidSparkLineMode, err)
}

func TestUnmarshalYAMLWithSingleStatAndInvalidColoringTarget(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - single_stat:
          title: Heap Allocations
          datasource: prometheus-default
          targets:
            - prometheus:
                query: 'go_memstats_heap_alloc_bytes{job="prometheus"}'
          unit: bytes
          thresholds: ["26000000", "28000000"]
          color: ["value", "invalid target"]
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidColoringTarget, err)
}

func TestUnmarshalYAMLWithSingleStatAndInvalidValueType(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - single_stat:
          title: Heap Allocations
          datasource: prometheus-default
          targets:
            - prometheus:
                query: 'go_memstats_heap_alloc_bytes{job="prometheus"}'
          value_type: invalid
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidSingleStatValueType, err)
}

func TestUnmarshalYAMLWithNoTargetSingleGraph(t *testing.T) {
	payload := `
rows:
  - name: Prometheus
    panels:
      - graph:
          title: Threads
          targets:
            - {}
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrTargetNotConfigured, err)
}

func TestUnmarshalYAMLWithNoAlertThresholdGraph(t *testing.T) {
	payload := `
rows:
  - name: Test row
    panels:
      - graph:
          title: Heap allocations
          targets:
            - prometheus: { query: "go_memstats_heap_alloc_bytes" }
          alert:
            summary: Too many heap allocations
            evaluate_every: 1m
            for: 1m
            if:
              - { avg: A }
            targets:
              - prometheus: { ref: A, query: "go_memstats_heap_alloc_bytes" }
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrNoAlertThresholdDefined, err)
}

func TestUnmarshalYAMLWithInvalidLegendGraph(t *testing.T) {
	payload := `
rows:
  - name: Test row
    panels:
      - graph:
          title: Heap allocations
          legend: [invalid_attribute]
          targets:
            - prometheus: { query: "go_memstats_heap_alloc_bytes" }
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidLegendAttribute, err)
}

func TestUnmarshalYAMLWithInvalidAlertValueFunctionGraph(t *testing.T) {
	payload := `
rows:
  - name: Test row
    panels:
      - graph:
          title: Heap allocations
          targets:
            - prometheus: { query: "go_memstats_heap_alloc_bytes" }
          alert:
            summary: Too many heap allocations
            evaluate_every: 1m
            for: 1m
            if:
              - { BLOOPER: A, above: 23000000 }
            targets:
              - prometheus: { ref: A, query: "go_memstats_heap_alloc_bytes" }
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), (&yaml.TypeError{}).Error()))
}

func TestUnmarshalYAMLWithInvalidStackdriverAggregation(t *testing.T) {
	payload := `
rows:
  - name: Test row
    panels:
      - graph:
          title: Pubsub Ack msg count
          datasource: stackdriver-default
          targets:
            - stackdriver:
                type: delta
                metric: pubsub.googleapis.com/subscription/ack_message_count
                aggregation: invalid
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidStackdriverAggregation, err)
}

func TestUnmarshalYAMLWithInvalidStackdriverAlignmentMethod(t *testing.T) {
	payload := `
rows:
  - name: Test row
    panels:
      - graph:
          title: Pubsub Ack msg count
          datasource: stackdriver-default
          targets:
            - stackdriver:
                type: delta
                metric: pubsub.googleapis.com/subscription/ack_message_count
                alignment: {method: invalid, period: stackdriver-auto}
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidStackdriverAlignment, err)
}

func generalOptions() testCase {
	yaml := `title: Awesome dashboard

editable: true
shared_crosshair: true
tags: [generated, yaml]
auto_refresh: 10s

time: ["now-6h", "now"]
timezone: utc
`

	return testCase{
		name:                "general options",
		yaml:                yaml,
		expectedGrafanaJSON: "general_options.json",
	}
}

func tagAnnotations() testCase {
	yaml := `title: Awesome dashboard

timezone: default

tags_annotations:
  - name: Deployments
    datasource: "-- Grafana --"
    color: "#5794F2"
    tags: ["deploy", "production"]
`

	return testCase{
		name:                "tag annotations",
		yaml:                yaml,
		expectedGrafanaJSON: "tag_annotations.json",
	}
}

func variables() testCase {
	yaml := `title: Awesome dashboard

timezone: browser

variables:
  - interval:
      name: interval
      label: Interval
      default: 30s
      values: ["30s", "1m"]
  - query:
      name: status
      label: HTTP status
      datasource: prometheus-default
      include_all: true
      default_all: true
      request: "label_values(prometheus_http_requests_total, code)"
  - const:
      name: percentile
      label: Percentile
      default: 50
      hide: label
      values_map:
        50th: "50"
  - custom:
      name: vX
      label: vX
      default: v1
      values_map:
        v1: v1
`

	return testCase{
		name:                "variables",
		yaml:                yaml,
		expectedGrafanaJSON: "variables.json",
	}
}

func textPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - text:
          description: Some description
          height: 400px
          span: 6
          transparent: true
          title: Some markdown?
          markdown: "*markdown*"
      - text:
          height: 400px
          span: 6
          title: Some html?
          html: "Some <b>awesome</b> html"
`

	return testCase{
		name:                "single row with text panels",
		yaml:                yaml,
		expectedGrafanaJSON: "text_panel.json",
	}
}

func graphPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - graph:
          title: Heap allocations
          description: Some description
          height: 400px
          span: 4
          transparent: true
          datasource: prometheus-default
          legend: [avg, current, min, max, as_table, no_null_series, no_zero_series]
          alert:
            summary: Too many heap allocations
            description: "Wow, a we're allocating a lot."
            evaluate_every: 1m
            for: 2m
            on_no_data: alerting
            on_execution_error: alerting
            tags:
              severity: super-critical-stop-the-world-now
            if:
              - { avg: A, above: 23000000 }
            targets:
              - prometheus: { ref: A, query: "go_memstats_heap_alloc_bytes" }
          axes:
            left: { unit: short, min: 0, max: 100, label: Requests }
            right: { hidden: true }
            bottom: { hidden: true }
          targets:
            - prometheus:
                query: "go_memstats_heap_alloc_bytes"
                legend: "{{job}}"
`

	return testCase{
		name:                "single row with single graph panel",
		yaml:                yaml,
		expectedGrafanaJSON: "graph_panel.json",
	}
}

func graphPanelWithStackdriverTarget() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - graph:
          title: Pubsub Ack msg count
          datasource: voi-stage-stackdriver
          targets:
            - stackdriver:
                legend: Ack-ed messages
                type: delta
                metric: pubsub.googleapis.com/subscription/ack_message_count
                aggregation: mean
                alignment: {method: delta, period: stackdriver-auto}
                filters:
                  eq:
                    resource.type: pubsub_subscription
`

	return testCase{
		name:                "single row with single graph panel and stackdriver target",
		yaml:                yaml,
		expectedGrafanaJSON: "graph_panel_stackdriver_target.json",
	}
}

func graphPanelWithGraphiteTarget() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - graph:
          title: Packets received
          datasource: graphite-test
          targets:
            - graphite:
                query: stats_counts.statsd.packets_received
`

	return testCase{
		name:                "single row with single graph panel and graphite target",
		yaml:                yaml,
		expectedGrafanaJSON: "graph_panel_graphite_target.json",
	}
}

func graphPanelWithInfluxdbTarget() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - graph:
          title: Dummy
          datasource: influxdb-test
          targets:
            - influxdb:
                query: buckets()
`

	return testCase{
		name:                "single row with single graph panel and influxdb target",
		yaml:                yaml,
		expectedGrafanaJSON: "graph_panel_influxdb_target.json",
	}
}

func singleStatPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - single_stat:
          title: Heap Allocations
          description: Some description
          height: 400px
          span: 4
          transparent: true
          datasource: prometheus-default
          targets:
            - prometheus:
                query: 'go_memstats_heap_alloc_bytes{job="prometheus"}'
          unit: bytes
          value_type: current
          value_font_size: '120%'
          prefix_font_size: '80%'
          postfix_font_size: '80%'
          sparkline: bottom
          thresholds: ["26000000", "28000000"]
          color: ["value", "background"]
          colors: ["green", "yellow", "red"]
`

	return testCase{
		name:                "single row with one singlestat panel",
		yaml:                yaml,
		expectedGrafanaJSON: "singlestat_panel.json",
	}
}

func statPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Kubelet
    panels:
      - stat:
          title: HTTP requests
          description: Some description
          height: 400px
          span: 4
          transparent: true
          datasource: prometheus-default
          targets:
            - prometheus:
                query: "count(kubelet_http_requests_total) by (method, path)"
                legend: "{{ method }} - {{ path }}"
          orientation: horizontal
          text: value_and_name
          sparkline: true
          unit: short
          decimals: 2
          title_font_size: 100
          value_font_size: 150
          color_mode: background
          thresholds:
            - {color: green}
            - {value: 1, color: orange}
            - {value: 4, color: red}
`

	return testCase{
		name:                "single row with one stat panel",
		yaml:                yaml,
		expectedGrafanaJSON: "stat_panel.json",
	}
}

func tablePanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - table:
          title: Threads
          description: Threads here
          height: 400px
          span: 4
          transparent: true
          datasource: prometheus-default
          targets:
            - prometheus:
                query: "go_threads"
          hidden_columns: ["Time"]
          time_series_aggregations:
            - label: AVG
              type: avg
            - label: Current
              type: current
`

	return testCase{
		name:                "single row with single table panel",
		yaml:                yaml,
		expectedGrafanaJSON: "table_panel.json",
	}
}

func heatmapPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
    - heatmap:
        title: Reconciliation Performance
        description: Does it perform?
        span: 12
        datasource: $datasource
        data_format: time_series_buckets
        hide_zero_buckets: false
        highlight_cards: true
        targets:
        - prometheus:
            query: sum(increase(argocd_app_reconcile_bucket{namespace=~"$namespace"}[$interval])) by (le)
            legend: '{{le}}'
            format: heatmap
            interval_factor: 10
        tooltip:
          show: true
          showhistogram: true
          decimals: 0
`

	return testCase{
		name:                "single row with heatmap panel",
		yaml:                yaml,
		expectedGrafanaJSON: "heatmap_panel.json",
	}
}

func timeseriesPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
    - timeseries:
        title: Total Request per Second
        description: Does it perform?
        span: 12
        links:
        - { title: linky, url: http://linky }
        targets:
        - prometheus:
            query: "go_memstats_heap_alloc_bytes"
`

	return testCase{
		name:                "single row with timeseries panel",
		yaml:                yaml,
		expectedGrafanaJSON: "timeseries_panel.json",
	}
}

func logsPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
    - logs:
        title: Kubernetes logs
        description: Everything okay?
        span: 12
        visualization:
          deduplication: exact
        targets:
        - loki:
            query: "{namespace=\"default\"}"
`

	return testCase{
		name:                "single row with logs panel",
		yaml:                yaml,
		expectedGrafanaJSON: "logs_panel.json",
	}
}

func collapseRow() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test collapsed row
    collapse: true
    panels:
      - text:
          description: Some description
          height: 400px
          span: 6
          transparent: true
          title: Some markdown?
          markdown: "*markdown*"
`

	return testCase{
		name:                "single collapsed row",
		yaml:                yaml,
		expectedGrafanaJSON: "single_collapsed_row.json",
	}
}

func dashboardFromFixtures(t *testing.T, path string) string {
	req := require.New(t)

	payload, err := os.ReadFile(path)
	req.NoError(err)

	return string(payload)
}
