package prometheus_test

import (
	"testing"

	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/stretchr/testify/require"
)

func TestNewPrometheusTargetCanBeCreated(t *testing.T) {
	req := require.New(t)
	query := "rate(prometheus_http_requests_total[30s])"

	target := prometheus.New(query)

	req.Equal(query, target.Expr)
	req.Equal("time_series", target.Format)
}

func TestLegendCanBeConfigured(t *testing.T) {
	req := require.New(t)
	legend := "{{ code }} - {{ path }}"

	target := prometheus.New("", prometheus.WithLegend(legend))

	req.Equal(legend, target.LegendFormat)
}
