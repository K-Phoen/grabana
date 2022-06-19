package grabana

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/K-Phoen/sdk"

	"github.com/K-Phoen/grabana/alert"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/stretchr/testify/require"
)

func TestADashboardCanBeFoundByTitle(t *testing.T) {
	req := require.New(t)
	dashboardName := "Test dashboard"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "eErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "eErXDvCkyy",
    "title": "Test dashboard"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dash, err := client.GetDashboardByTitle(context.TODO(), strings.ToLower(dashboardName))

	req.NoError(err)
	req.Equal(dashboardName, dash.Title)
}

func TestAnExplicitErrorIsReturnedIfTheDashboardIsNotFoundByTitle(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "eErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "eErXDvCkyy",
    "title": "Test dashboard"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dash, err := client.GetDashboardByTitle(context.TODO(), "dashboard that do not exist")

	req.Error(err)
	req.Nil(dash)
	req.Equal(ErrDashboardNotFound, err)
}

func TestAnExplicitErrorIsReturnedIfTheDashboardIsNotFoundByTitleByGrafana(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `[]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dash, err := client.GetDashboardByTitle(context.TODO(), "dashboard that do not exist")

	req.Error(err)
	req.Nil(dash)
	req.Equal(ErrDashboardNotFound, err)
}

func TestGetDashboardByTitleCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dash, err := client.GetDashboardByTitle(context.TODO(), "does not matter")

	req.Error(err)
	req.Nil(dash)
}

func TestADashboardCanBeFoundByUID(t *testing.T) {
	req := require.New(t)
	dashboardUID := "lala-uid"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `{
  "id":1,
  "dashboard": {
    "id":1,
    "uid": "lala-uid",
    "title": "some title"
  }
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dash, err := client.rawDashboardByUID(context.TODO(), dashboardUID)

	req.NoError(err)
	req.Equal(dashboardUID, dash.UID)
}

func TestFetchingAnUnknownDashboardByUIDFailsCleanly(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	_, err := client.rawDashboardByUID(context.TODO(), "uid")

	req.Error(err)
	req.ErrorIs(err, ErrDashboardNotFound)
}

func TestDeleteDashboardWithNoAlerts(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Potential existing alerts retrieval
		if r.Method == http.MethodGet && r.URL.String() == "/api/ruler/grafana/api/v1/rules?dashboard_uid=some-uid" {
			_, _ = fmt.Fprintln(w, `{}`)
			return
		}

		// Dashboard deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/dashboards/uid/some-uid" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"title": "Production Overview"}`)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "oh noes, we should not get here", "method": "%s", "path": "%s"}\n`, r.Method, r.URL.String())
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some-uid")

	req.NoError(err)
}

func TestDeleteDashboardWithAlerts(t *testing.T) {
	req := require.New(t)
	firstAlertDeleted := false
	secondAlertDeleted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Potential existing alerts retrieval
		if r.Method == http.MethodGet && r.URL.String() == "/api/ruler/grafana/api/v1/rules?dashboard_uid=some-uid" {
			_, _ = fmt.Fprintln(w, `{
  "test ns 1": [
    {"name": "alert 1"}
  ],
  "test ns 2": [
    {"name": "alert 2"}
  ]
}`)
			return
		}

		// First alert deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/test%20ns%201/alert%201" {
			firstAlertDeleted = true
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Second alert deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/test%20ns%202/alert%202" {
			secondAlertDeleted = true
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Dashboard deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/dashboards/uid/some-uid" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"title": "Production Overview"}`)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "oh noes, we should not get here", "method": "%s", "path": "%s"}\n`, r.Method, r.URL.String())
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some-uid")

	req.NoError(err)
	req.True(firstAlertDeleted)
	req.True(secondAlertDeleted)
}

func TestDeleteDashboardCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintln(w, `{"message": "oh noes"}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some uid")

	req.Error(err)
}

func TestDeletingANonExistingDashboardReturnsSpecificError(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Potential existing alerts retrieval
		if r.Method == http.MethodGet && r.URL.String() == "/api/ruler/grafana/api/v1/rules?dashboard_uid=some-uid" {
			_, _ = fmt.Fprintln(w, `{}`)
			return
		}

		// Dashboard deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/dashboards/uid/some-uid" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintln(w, `{"message": "oh noes, does not exist"}`)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "oh noes, we should not get here", "method": "%s", "path": "%s"}\n`, r.Method, r.URL.String())
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some-uid")

	req.Equal(ErrDashboardNotFound, err)
}

func TestDashboardsCanBeCreatedWithNoAlertAndNoPreviousAlerts(t *testing.T) {
	req := require.New(t)

	builder, err := dashboard.New(
		"Dashboard no alert",
		dashboard.Row(
			"Row",
			row.WithText("Some text", text.Markdown("Markdown")),
		),
	)
	req.NoError(err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dashboard creation
		if r.Method == http.MethodPost && r.URL.Path == "/api/dashboards/db" {
			_, _ = fmt.Fprintln(w, `{
  "id":      1,
  "uid":     "cIBgcSjkk",
  "url":     "/d/cIBgcSjkk/production-overview",
  "status":  "success",
  "version": 1,
  "slug":    "test-dashboard"
}`)
			return
		}

		// Raw dashboard retrieval after creation
		if r.Method == http.MethodGet && r.URL.Path == "/api/dashboards/uid/cIBgcSjkk" {
			_, _ = fmt.Fprintln(w, `{"dashboard": {
  "id": 1,
  "uid": "cIBgcSjkk",
  "slug": "test-dashboard",
  "title": "Dashboard no alert",
  "originalTitle": "",
  "tags": null,
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "sharedCrosshair": true,
  "panels": null,
  "rows": [
    {
      "title": "Row",
      "showTitle": true,
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "editable": false,
          "error": false,
          "gridPos": {},
          "id": 1,
          "isNew": false,
          "renderer": "flot",
          "span": 6,
          "title": "Some text",
          "transparent": false,
          "type": "text",
          "content": "Markdown",
          "mode": "markdown",
          "pageSize": 0,
          "scroll": false,
          "showHeader": false,
          "sort": {
            "col": 0,
            "desc": false
          },
          "styles": null,
          "fieldConfig": {
            "defaults": {
              "unit": "",
              "color": {
                "mode": ""
              },
              "thresholds": {
                "mode": "",
                "steps": null
              },
              "custom": {
                "axisPlacement": "",
                "barAlignment": 0,
                "drawStyle": "",
                "fillOpacity": 0,
                "gradientMode": "",
                "lineInterpolation": "",
                "lineWidth": 0,
                "pointSize": 0,
                "showPoints": "",
                "spanNulls": false,
                "hideFrom": {
                  "legend": false,
                  "tooltip": false,
                  "viz": false
                },
                "lineStyle": {
                  "fill": ""
                },
                "scaleDistribution": {
                  "type": ""
                },
                "stacking": {
                  "group": "",
                  "mode": ""
                },
                "thresholdsStyle": {
                  "mode": ""
                }
              }
            },
            "overrides": null
          },
          "options": {
            "content": "",
            "mode": ""
          }
        }
      ],
      "repeat": null
    }
  ],
  "templating": {
    "list": null
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 0,
  "links": null,
  "time": {
    "from": "now-3h",
    "to": "now"
  }
}}`)
			return
		}

		// Potential existing alerts retrieval
		if r.Method == http.MethodGet && r.URL.String() == "/api/ruler/grafana/api/v1/rules?dashboard_uid=cIBgcSjkk" {
			_, _ = fmt.Fprintln(w, `{}`)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "oh noes, we should not get here", "method": "%s", "path": "%s"}\n`, r.Method, r.URL.String())
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	board, err := client.UpsertDashboard(context.TODO(), &Folder{}, builder)

	req.NoError(err)
	req.NotNil(board)
}

func TestDashboardsCanBeCreatedWithNoAlertAndDeletesPreviousAlerts(t *testing.T) {
	req := require.New(t)

	builder, err := dashboard.New(
		"Dashboard no alert",
		dashboard.Row(
			"Row",
			row.WithTimeSeries(
				"HTTP Rate",
				timeseries.DataSource("Prometheus"),
				timeseries.WithPrometheusTarget(
					"sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
				),
			),
		),
	)
	req.NoError(err)

	dashboardPersisted := false
	firstAlertDeleted := false
	secondAlertDeleted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dashboard creation
		if r.Method == http.MethodPost && r.URL.Path == "/api/dashboards/db" {
			dashboardPersisted = true
			_, _ = fmt.Fprintln(w, `{
  "id":      1,
  "uid":     "cIBgcSjkk",
  "url":     "/d/cIBgcSjkk/production-overview",
  "status":  "success",
  "version": 1,
  "slug":    "test-dashboard"
}`)
			return
		}

		// Raw dashboard retrieval after creation
		if r.Method == http.MethodGet && r.URL.Path == "/api/dashboards/uid/cIBgcSjkk" {
			_, _ = fmt.Fprintln(w, `{"dashboard": {
  "id": 1,
  "uid": "cIBgcSjkk",
  "slug": "test-dashboard",
  "title": "Dashboard no alert",
  "originalTitle": "",
  "tags": null,
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "sharedCrosshair": true,
  "panels": null,
  "rows": [
    {
      "title": "Row",
      "showTitle": true,
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "datasource": "Prometheus",
          "editable": false,
          "error": false,
          "gridPos": {},
          "id": 1,
          "isNew": false,
          "span": 6,
          "title": "HTTP Rate",
          "transparent": false,
          "type": "timeseries",
          "targets": [
            {
              "refId": "",
              "expr": "sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
              "format": "time_series"
            }
          ],
          "options": {
            "legend": {
              "calcs": [],
              "displayMode": "list",
              "placement": "bottom"
            },
            "tooltip": {
              "mode": "single"
            }
          },
          "fieldConfig": {
            "defaults": {
              "unit": "",
              "color": {
                "mode": "palette-classic",
                "fixedColor": "green",
                "seriesBy": "last"
              },
              "thresholds": {
                "mode": "",
                "steps": null
              },
              "custom": {
                "axisPlacement": "auto",
                "barAlignment": 0,
                "drawStyle": "line",
                "fillOpacity": 25,
                "gradientMode": "opacity",
                "lineInterpolation": "linear",
                "lineWidth": 1,
                "pointSize": 5,
                "showPoints": "",
                "spanNulls": false,
                "hideFrom": {
                  "legend": false,
                  "tooltip": false,
                  "viz": false
                },
                "lineStyle": {
                  "fill": "solid"
                },
                "scaleDistribution": {
                  "type": "linear"
                },
                "stacking": {
                  "group": "",
                  "mode": ""
                },
                "thresholdsStyle": {
                  "mode": ""
                }
              }
            },
            "overrides": null
          }
        }
      ],
      "repeat": null
    }
  ],
  "templating": {
    "list": null
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 0,
  "links": null,
  "time": {
    "from": "now-3h",
    "to": "now"
  }
}}`)
			return
		}

		// Potential existing alerts retrieval
		if r.Method == http.MethodGet && r.URL.String() == "/api/ruler/grafana/api/v1/rules?dashboard_uid=cIBgcSjkk" {
			_, _ = fmt.Fprintln(w, `{
  "test ns 1": [
    {"name": "alert 1"}
  ],
  "test ns 2": [
    {"name": "alert 2"}
  ]
}`)
			return
		}

		// First alert deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/test%20ns%201/alert%201" {
			firstAlertDeleted = true
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Second alert deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/test%20ns%202/alert%202" {
			secondAlertDeleted = true
			w.WriteHeader(http.StatusAccepted)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "oh noes, we should not get here", "method": "%s", "path": "%s"}\n`, r.Method, r.URL.String())
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	board, err := client.UpsertDashboard(context.TODO(), &Folder{}, builder)

	req.NoError(err)
	req.NotNil(board)
	req.True(dashboardPersisted)
	req.True(firstAlertDeleted)
	req.True(secondAlertDeleted)
}

func TestDashboardsCanBeCreatedWithNewAlertsAndDeletesPreviousAlerts(t *testing.T) {
	req := require.New(t)

	builder, err := dashboard.New(
		"Dashboard no alert",
		dashboard.Row(
			"Row",
			row.WithTimeSeries(
				"Heap allocations",
				timeseries.DataSource("Prometheus"),
				timeseries.WithPrometheusTarget(
					"sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
				),
				timeseries.Alert(
					"Too many heap allocations",
					alert.WithPrometheusQuery(
						"A",
						"sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
					),
					alert.If(alert.Avg, "A", alert.IsAbove(3)),
				),
			),
		),
	)
	req.NoError(err)

	dashboardPersisted := false
	firstAlertDeleted := false
	secondAlertDeleted := false
	newAlertDeletionAttempt := false
	newAlertCreated := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dashboard creation
		if r.Method == http.MethodPost && r.URL.Path == "/api/dashboards/db" {
			dashboardPersisted = true
			_, _ = fmt.Fprintln(w, `{
  "id":      1,
  "uid":     "cIBgcSjkk",
  "url":     "/d/cIBgcSjkk/production-overview",
  "status":  "success",
  "version": 1,
  "slug":    "test-dashboard"
}`)
			return
		}

		// Raw dashboard retrieval after creation
		if r.Method == http.MethodGet && r.URL.Path == "/api/dashboards/uid/cIBgcSjkk" {
			_, _ = fmt.Fprintln(w, `{"dashboard": {
  "id": 1,
  "uid": "cIBgcSjkk",
  "slug": "test-dashboard",
  "title": "Dashboard no alert",
  "originalTitle": "",
  "tags": null,
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "sharedCrosshair": true,
  "panels": null,
  "rows": [
    {
      "title": "Row",
      "showTitle": true,
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "datasource": "Prometheus",
          "editable": false,
          "error": false,
          "gridPos": {},
          "id": 1,
          "isNew": false,
          "span": 6,
          "title": "Heap allocations",
          "transparent": false,
          "type": "timeseries",
          "targets": [
            {
              "refId": "",
              "expr": "sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
              "format": "time_series"
            }
          ],
          "options": {
            "legend": {
              "calcs": [],
              "displayMode": "list",
              "placement": "bottom"
            },
            "tooltip": {
              "mode": "single"
            }
          },
          "fieldConfig": {
            "defaults": {
              "unit": "",
              "color": {
                "mode": "palette-classic",
                "fixedColor": "green",
                "seriesBy": "last"
              },
              "thresholds": {
                "mode": "",
                "steps": null
              },
              "custom": {
                "axisPlacement": "auto",
                "barAlignment": 0,
                "drawStyle": "line",
                "fillOpacity": 25,
                "gradientMode": "opacity",
                "lineInterpolation": "linear",
                "lineWidth": 1,
                "pointSize": 5,
                "showPoints": "",
                "spanNulls": false,
                "hideFrom": {
                  "legend": false,
                  "tooltip": false,
                  "viz": false
                },
                "lineStyle": {
                  "fill": "solid"
                },
                "scaleDistribution": {
                  "type": "linear"
                },
                "stacking": {
                  "group": "",
                  "mode": ""
                },
                "thresholdsStyle": {
                  "mode": ""
                }
              }
            },
            "overrides": null
          }
        }
      ],
      "repeat": null
    }
  ],
  "templating": {
    "list": null
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 0,
  "links": null,
  "time": {
    "from": "now-3h",
    "to": "now"
  }
}}`)
			return
		}

		// Potential existing alerts retrieval
		if r.Method == http.MethodGet && r.URL.String() == "/api/ruler/grafana/api/v1/rules?dashboard_uid=cIBgcSjkk" {
			_, _ = fmt.Fprintln(w, `{
  "test ns 1": [
    {"name": "alert 1"}
  ],
  "test ns 2": [
    {"name": "alert 2"}
  ]
}`)
			return
		}

		// First alert deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/test%20ns%201/alert%201" {
			firstAlertDeleted = true
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Second alert deletion
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/test%20ns%202/alert%202" {
			secondAlertDeleted = true
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Retrieve known datasources
		if r.Method == http.MethodGet && r.URL.String() == "/api/datasources" {
			_, _ = fmt.Fprintln(w, `[
  {
    "uid": "loki-uid",
    "name": "Loki",
    "isDefault": false
  },
  {
    "uid": "prom-uid",
    "name": "Prometheus",
    "isDefault": true
  }
]`)
			return
		}

		// new alert deletion attempt
		if r.Method == http.MethodDelete && r.URL.String() == "/api/ruler/grafana/api/v1/rules/Folder%20name/Heap%20allocations" {
			newAlertDeletionAttempt = true
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// new alert creation
		if r.Method == http.MethodPost && r.URL.String() == "/api/ruler/grafana/api/v1/rules/Folder%20name" {
			newAlertCreated = true
			body, err := ioutil.ReadAll(r.Body)
			req.NoError(err)

			req.JSONEq(`{
  "name": "Heap allocations",
  "interval": "1m",
  "rules": [
    {
      "for": "5m",
      "grafana_alert": {
        "title": "Heap allocations",
        "condition": "_alert_condition_",
        "no_data_state": "NoData",
        "exec_err_state": "Alerting",
        "data": [
          {
            "refId": "_alert_condition_",
            "queryType": "",
            "datasourceUid": "-100",
            "model": {
              "refId": "_alert_condition_",
              "type": "classic_conditions",
              "datasource": {
                "uid": "-100",
                "type": "__expr__"
              },
              "hide": false,
              "conditions": [
                {
                  "type": "query",
                  "evaluator": {
                    "params": [
                      3
                    ],
                    "type": "gt"
                  },
                  "operator": {
                    "type": "and"
                  },
                  "query": {
                    "params": [
                      "A"
                    ]
                  },
                  "reducer": {
                    "type": "avg"
                  }
                }
              ]
            }
          },
          {
            "refId": "A",
            "queryType": "",
            "relativeTimeRange": {
              "from": 600,
              "to": 0
            },
            "datasourceUid": "prom-uid",
            "model": {
              "refId": "A",
              "expr": "sum(go_memstats_heap_alloc_bytes{app!=\"\"}) by (app)",
              "format": "time_series",
              "datasource": {
                "uid": "prom-uid",
                "type": "prometheus"
              },
              "intervalMs": 15000,
              "hide": false
            }
          }
        ]
      },
      "annotations": {
        "__dashboardUid__": "cIBgcSjkk",
        "__panelId__": "1",
        "summary": "Too many heap allocations"
      }
    }
  ]
}`, string(body))

			w.WriteHeader(http.StatusAccepted)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"message": "oh noes, we should not get here", "method": "%s", "path": "%s"}\n`, r.Method, r.URL.String())
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	board, err := client.UpsertDashboard(context.TODO(), &Folder{
		Title: "Folder name",
	}, builder)

	req.NoError(err)
	req.NotNil(board)
	req.True(dashboardPersisted)
	req.True(firstAlertDeleted)
	req.True(secondAlertDeleted)
	req.True(newAlertDeletionAttempt)
	req.True(newAlertCreated)
}

func TestClient_panelIDByTitle_panelInBoard(t *testing.T) {
	req := require.New(t)

	board := sdk.NewBoard("board title")
	panel := sdk.NewLogs("Logs panel")
	panel.ID = 42

	board.Panels = append(board.Panels, panel)

	req.Equal("42", panelIDByTitle(board, "Logs panel"))
	req.Equal("", panelIDByTitle(board, "not found"))
}

func TestClient_panelIDByTitle_panelInRow(t *testing.T) {
	req := require.New(t)

	board := sdk.NewBoard("board title")
	panel := sdk.NewHeatmap("Heamtap panel")
	panel.ID = 24
	rowPanel := board.AddRow("Row title")

	rowPanel.Panels = append(rowPanel.Panels, *panel)

	req.Equal("24", panelIDByTitle(board, "Heamtap panel"))
	req.Equal("", panelIDByTitle(board, "not found"))
}
