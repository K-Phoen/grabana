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

func TestUnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name                string
		yaml                string
		expectedGrafanaJSON string
	}{
		generalOptions(),
		tagAnnotations(),
		textPanel(),
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
          title: Some markdown?
          markdown: "*markdown*"
      - text:
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
