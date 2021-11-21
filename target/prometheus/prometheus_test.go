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

	target := prometheus.New("", prometheus.Legend(legend))

	req.Equal(legend, target.LegendFormat)
}

func TestRefCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := prometheus.New("", prometheus.Ref("A"))

	req.Equal("A", target.Ref)
	req.False(target.Hidden)
}

func TestTargetCanBeHidden(t *testing.T) {
	req := require.New(t)

	target := prometheus.New("", prometheus.Hide())

	req.True(target.Hidden)
}

func TestTargetCanBeSetAsInstant(t *testing.T) {
	req := require.New(t)

	target := prometheus.New("", prometheus.Instant())

	req.True(target.Instant)
}

func TestFormatCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := prometheus.New("", prometheus.Format(prometheus.FormatTimeSeries))

	req.Equal("time_series", target.Format)
}

func TestIntervalFactorCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := prometheus.New("", prometheus.IntervalFactor(1))

	req.Equal(1, target.IntervalFactor)
}
