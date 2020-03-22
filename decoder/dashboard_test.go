package decoder

import (
	"bytes"
	"fmt"
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
	testCases := []struct {
		name                string
		yaml                string
		expectedGrafanaJSON string
	}{
		generalOptions(),
		tagAnnotations(),
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

			fmt.Printf("json:\n%s\n", json)

			req.JSONEq(tc.expectedGrafanaJSON, string(json))
		})
	}
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
						"format": "time",
						"logBase": 1,
						"show": true
					},
					"yaxes": [
						{
							"format": "short",
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
					"valueName": "avg"
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
