package decoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
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
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			builder, err := UnmarshalYAML(bytes.NewBufferString(tc.yaml))
			req.NoError(err)

			json, err := builder.MarshalJSON()
			req.NoError(err)

			req.JSONEq(tc.expectedGrafanaJSON, string(json))
		})
	}
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
          alert:
            title: Too many heap allocations
            evaluate_every: 1m
            for: 1m
            if:
              - operand: and
                value: {func: avg, ref: A, from: 1m, to: now}

          targets:
            - prometheus: { query: "go_memstats_heap_alloc_bytes" }
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrNoAlertThresholdDefined, err)
}

func TestUnmarshalYAMLWithInvalidAlertValueFunctionGraph(t *testing.T) {
	payload := `
rows:
  - name: Test row
    panels:
      - graph:
          title: Heap allocations
          alert:
            title: Too many heap allocations
            evaluate_every: 1m
            for: 1m
            if:
              - operand: and
                value: {func: BLOOPER, ref: A, from: 1m, to: now}
                threshold: {above: 23000000}
          targets:
            - prometheus: { query: "go_memstats_heap_alloc_bytes" }
`

	_, err := UnmarshalYAML(bytes.NewBufferString(payload))

	require.Error(t, err)
	require.Equal(t, ErrInvalidAlertValueFunc, err)
}

func generalOptions() testCase {
	yaml := `title: Awesome dashboard

editable: true
shared_crosshair: true
tags: [generated, yaml]
auto_refresh: 10s
`
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": ["generated", "yaml"],
	"style": "dark",
	"timezone": "",
	"editable": true,
	"hideControls": false,
	"sharedCrosshair": true,
	"templating": {"list": null},
	"annotations": {"list": null},
	"links": null,
	"panels": null,
	"rows": [],
	"refresh": "10s",
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "general options",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}

func tagAnnotations() testCase {
	yaml := `title: Awesome dashboard

tags_annotations:
  - name: Deployments
    datasource: "-- Grafana --"
    color: "#5794F2"
    tags: ["deploy", "production"]
`
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": null,
	"style": "dark",
	"timezone": "",
	"editable": false,
	"hideControls": false,
	"sharedCrosshair": false,
	"templating": {
		"list": null
	},
	"annotations": {
		"list": [{
			"datasource": "-- Grafana --",
			"enable": true,
			"iconColor": "#5794F2",
			"iconSize": 0,
			"name": "Deployments",
			"query": "",
			"showLine": false,
			"lineColor": "",
			"tags": ["deploy", "production"],
			"tagsField": "",
			"textField": "",
			"type": "tags"
		}]
	},
	"links": null,
	"panels": null,
	"rows": [],
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "tag annotations",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}

func variables() testCase {
	yaml := `title: Awesome dashboard

variables:
  - interval:
      name: interval
      label: Interval
      default: 30s
      values: ["30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"]
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
      values_map:
        50th: "50"
  - custom:
      name: vX
      label: vX
      default: v1
      values_map:
        v1: v1
`
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": null,
	"style": "dark",
	"timezone": "",
	"editable": false,
	"hideControls": false,
	"sharedCrosshair": false,
	"templating": {
		"list": [
			{
				"name": "interval",
				"type": "interval",
				"datasource": null,
				"refresh": false,
				"options": null,
				"includeAll": false,
				"allFormat": "",
				"allValue": "",
				"multi": false,
				"multiFormat": "",
				"query": "10m,12h,1h,1m,30m,30s,5m,6h",
				"regex": "",
				"current": {
					"text": "30s",
					"value": "30s"
				},
				"label": "Interval",
				"hide": 0,
				"sort": 0
			},
			{
				"name": "status",
				"type": "query",
				"datasource": "prometheus-default",
				"refresh": 1,
				"options": [
					{
						"text": "All",
						"value": "$__all",
						"selected": false
					}
				],
				"includeAll": true,
				"allFormat": "",
				"allValue": "",
				"multi": false,
				"multiFormat": "",
				"query": "label_values(prometheus_http_requests_total, code)",
				"regex": "",
				"current": {
					"text": "All",
					"value": "$__all"
				},
				"label": "HTTP status",
				"hide": 0,
				"sort": 0
			},
			{
				"name": "percentile",
				"type": "constant",
				"datasource": null,
				"refresh": false,
				"options": [
					{
						"selected": false,
						"text": "50th",
						"value": "50"
					}
				],
				"includeAll": false,
				"allFormat": "",
				"allValue": "",
				"multi": false,
				"multiFormat": "",
				"query": "50",
				"regex": "",
				"current": {
					"text": "50th",
					"value": "50"
				},
				"label": "Percentile",
				"hide": 0,
				"sort": 0
			},
			{
				"name": "vX",
				"type": "custom",
				"datasource": null,
				"refresh": false,
				"options": [
					{
						"text": "v1",
						"value": "v1",
						"selected": false
					}
				],
				"includeAll": false,
				"allFormat": "",
				"allValue": "",
				"multi": false,
				"multiFormat": "",
				"query": "v1",
				"regex": "",
				"current": {
					"text": "v1",
					"value": "v1"
				},
				"label": "vX",
				"hide": 0,
				"sort": 0
			}
		]
	},
	"annotations": {"list": null},
	"links": null,
	"panels": null,
	"rows": [],
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "variables",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}

func textPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - text:
          height: 400px
          span: 6
          title: Some markdown?
          markdown: "*markdown*"
      - text:
          height: 400px
          span: 6
          title: Some html?
          html: "Some <b>awesome</b> html"
`
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": null,
	"style": "dark",
	"timezone": "",
	"editable": false,
	"hideControls": false,
	"sharedCrosshair": false,
	"templating": {"list": null},
	"annotations": {"list": null},
	"links": null,
	"panels": null,
	"rows": [
		{
			"title": "Test row",
			"collapse": false,
			"editable": true,
			"height": "250px",
			"repeat": null,
			"showTitle": true,
			"panels": [
				{
					"type": "text",
					"mode": "markdown",
					"content": "*markdown*",
					"editable": false,
					"error": false,
					"gridPos": {},
					"id": 1,
					"isNew": false,
					"pageSize": 0,
					"scroll": false,
					"renderer": "flot",
					"showHeader": false,
					"sort": {"col": 0, "desc": false},
					"span": 6,
					"height": "400px",
					"styles": null,
					"title": "Some markdown?",
					"transparent": false
				},
				{
					"type": "text",
					"mode": "html",
					"content": "Some <b>awesome</b> html",
					"editable": false,
					"error": false,
					"gridPos": {},
					"id": 2,
					"isNew": false,
					"scroll": false,
					"pageSize": 0,
					"renderer": "flot",
					"showHeader": false,
					"sort": {"col": 0, "desc": false},
					"span": 6,
					"height": "400px",
					"styles": null,
					"title": "Some html?",
					"transparent": false
				}
			]
		}
	],
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "single row with text panels",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}

func graphPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - graph:
          title: Heap allocations
          height: 400px
          span: 4
          datasource: prometheus-default
          alert:
            title: Too many heap allocations
            evaluate_every: 1m
            for: 1m
            notify: 1
            message: "Wow, a we're allocating a lot."
            on_no_data: alerting
            on_execution_error: alerting
            if:
              - operand: and
                value: {func: avg, ref: A, from: 1m, to: now}
                threshold: {above: 23000000}
          axes:
            left: { unit: short, min: 0, max: 100, label: Requests }
            right: { hidden: true }
            bottom: { hidden: true }
          targets:
            - prometheus:
                query: "go_memstats_heap_alloc_bytes"
                legend: "{{job}}"
                ref: A
`
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": null,
	"style": "dark",
	"timezone": "",
	"editable": false,
	"hideControls": false,
	"sharedCrosshair": false,
	"templating": {"list": null},
	"annotations": {"list": null},
	"links": null,
	"panels": null,
	"rows": [
		{
			"title": "Test row",
			"collapse": false,
			"editable": true,
			"height": "250px",
			"repeat": null,
			"showTitle": true,
			"panels": [
				{
					"type": "graph",
					"datasource": "prometheus-default",
					"editable": true,
					"error": false,
					"height": "400px",
					"gridPos": {},
					"id": 3,
					"isNew": false,
					"renderer": "flot",
					"span": 4,
					"fill": 1,
					"title": "Heap allocations",
					"aliasColors": {},
					"alert": {
						"conditions": [
							{
								"evaluator": {
									"params": [23000000],
									"type": "gt"
								},
								"operator": {"type": "and"},
								"query": {"params": ["A", "1m", "now"]},
								"reducer": {"type": "avg"},
								"type": "query"
							}
						],
						"executionErrorState": "alerting",
						"for": "1m",
						"frequency": "1m",
						"handler": 1,
						"message": "Wow, a we're allocating a lot.",
						"name": "Too many heap allocations",
						"noDataState": "alerting",
						"notifications": [
							{
								"disableResolveMessage": false,
								"frequency": "",
								"id": 1,
								"isDefault": false,
								"name": "",
								"sendReminder": false,
								"settings": null,
								"type": ""
							}
						]
					},
					"bars": false,
					"points": false,
					"stack": false,
					"steppedLine": false,
					"lines": true,
					"linewidth": 1,
					"pointradius": 5,
					"percentage": false,
					"nullPointMode": "null as zero",
					"legend": {
						"alignAsTable": false,
						"avg": false,
						"current": false,
						"hideEmpty": true,
						"hideZero": true,
						"max": false,
						"min": false,
						"rightSide": false,
						"show": true,
						"total": false,
						"values": false
					},
					"targets": [
						{
							"refId": "A",
							"expr": "go_memstats_heap_alloc_bytes",
							"legendFormat": "{{job}}",
							"format": "time_series"
						}
					],
					"tooltip": {
						"shared": true,
						"value_type": "",
						"sort": 2
					},
					"x-axis": true,
					"y-axis": true,
					"xaxis": {
						"format": "short",
						"logBase": 1,
						"show": false
					},
					"yaxes": [
						{
							"format": "short",
							"label": "Requests",
							"min": 0,
							"max": 100,
							"logBase": 1,
							"show": true
						},
						{
							"format": "short",
							"logBase": 1,
							"show": false
						}
					],
					"transparent": false
				}
			]
		}
	],
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "single row with single graph panel",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}

func singleStatPanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - single_stat:
          title: Heap Allocations
          height: 400px
          span: 4
          datasource: prometheus-default
          targets:
            - prometheus:
                query: 'go_memstats_heap_alloc_bytes{job="prometheus"}'
          unit: bytes
          value_type: current
          sparkline: bottom
          thresholds: ["26000000", "28000000"]
          color: ["value", "background"]
          colors: ["green", "yellow", "red"]
`
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": null,
	"style": "dark",
	"timezone": "",
	"editable": false,
	"hideControls": false,
	"sharedCrosshair": false,
	"templating": {"list": null},
	"annotations": {"list": null},
	"links": null,
	"panels": null,
	"rows": [
		{
			"title": "Test row",
			"collapse": false,
			"editable": true,
			"height": "250px",
			"repeat": null,
			"showTitle": true,
			"panels": [
				{
					"datasource": "prometheus-default",
					"editable": true,
					"error": false,
					"gridPos": {},
					"id": 4,
					"isNew": false,
					"renderer": "flot",
					"span": 4,
					"height": "400px",
					"title": "Heap Allocations",
					"transparent": false,
					"type": "singlestat",
					"colors": [
						"green",
						"yellow",
						"red"
					],
					"colorValue": true,
					"colorBackground": true,
					"decimals": 0,
					"format": "bytes",
					"gauge": {
						"maxValue": 0,
						"minValue": 0,
						"show": false,
						"thresholdLabels": false,
						"thresholdMarkers": false
					},
					"mappingType": 1,
					"mappingTypes": [
						{
							"name": "value to text",
							"value": 1
						},
						{
							"name": "range to text",
							"value": 2
						}
					],
					"nullPointMode": "",
					"sparkline": {
						"show": true,
						"fillColor": "rgba(31, 118, 189, 0.18)",
						"lineColor": "rgb(31, 120, 193)"
					},
					"targets": [
						{
							"refId": "",
							"expr": "go_memstats_heap_alloc_bytes{job=\"prometheus\"}",
							"format": "time_series"
						}
					],
					"thresholds": "26000000,28000000",
					"valueFontSize": "100%",
					"valueMaps": [
						{
							"op": "=",
							"text": "N/A",
							"value": "null"
						}
					],
					"valueName": "current"
				}
			]
		}
	],
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "single row with single graph panel",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}

func tablePanel() testCase {
	yaml := `title: Awesome dashboard

rows:
  - name: Test row
    panels:
      - table:
          title: Threads
          height: 400px
          span: 4
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
	json := `{
	"slug": "",
	"title": "Awesome dashboard",
	"originalTitle": "",
	"tags": null,
	"style": "dark",
	"timezone": "",
	"editable": false,
	"hideControls": false,
	"sharedCrosshair": false,
	"templating": {"list": null},
	"annotations": {"list": null},
	"links": null,
	"panels": null,
	"rows": [
		{
			"title": "Test row",
			"collapse": false,
			"editable": true,
			"height": "250px",
			"repeat": null,
			"showTitle": true,
			"panels": [
				{
					"datasource": "prometheus-default",
					"editable": true,
					"error": false,
					"gridPos": {},
					"height": "400px",
					"id": 5,
					"isNew": false,
					"renderer": "flot",
					"span": 4,
					"title": "Threads",
					"transparent": false,
					"type": "table",
					"columns": [
						{
							"text": "AVG",
							"value": "avg"
						},
						{
							"text": "Current",
							"value": "current"
						}
					],
					"styles": [
						{
							"alias": "",
							"pattern": "/.*/",
							"type": "string"
						},
						{
							"alias": null,
							"pattern": "Time",
							"type": "hidden"
						}
					],
					"transform": "timeseries_aggregations",
					"targets": [
						{
							"refId": "",
							"expr": "go_threads",
							"format": "time_series"
						}
					],
					"scroll": false
				}
			]
		}
	],
	"time": {"from": "now-3h", "to": "now"},
	"timepicker": {
		"refresh_intervals": ["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],
		"time_options": ["5m","15m","1h","6h","12h","24h","2d","7d","30d"]
	},
	"schemaVersion": 0,
	"version": 0
}`

	return testCase{
		name:                "single row with single graph panel",
		yaml:                yaml,
		expectedGrafanaJSON: json,
	}
}
