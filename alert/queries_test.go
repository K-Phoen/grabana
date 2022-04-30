package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrometheusQueriesCanBeAdded(t *testing.T) {
	req := require.New(t)

	a := New("", WithPrometheusQuery("A", "some prometheus query"))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}

func TestGraphiteQueriesCanBeAdded(t *testing.T) {
	req := require.New(t)

	a := New("", WithGraphiteQuery("A", "some graphite query"))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}

func TestLokiQueriesCanBeAdded(t *testing.T) {
	req := require.New(t)

	a := New("", WithLokiQuery("A", "some loki query"))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}

func TestStackdriverQueriesCanBeAdded(t *testing.T) {
	req := require.New(t)

	a := New("", WithStackdriverQuery("A", "some stackdriver query"))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}

func TestInfluxDBQueriesCanBeAdded(t *testing.T) {
	req := require.New(t)

	a := New("", WithInfluxDBQuery("A", "some influxdb query"))

	req.Len(a.Builder.Rules[0].GrafanaAlert.Data, 1)
}
