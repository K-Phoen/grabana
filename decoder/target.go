package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/target/prometheus"
)

var ErrTargetNotConfigured = fmt.Errorf("target not configured")

type Target struct {
	Prometheus *PrometheusTarget
}

type PrometheusTarget struct {
	Query  string
	Legend string
	Ref    string
}

func (t PrometheusTarget) toOptions() []prometheus.Option {
	var opts []prometheus.Option

	if t.Legend != "" {
		opts = append(opts, prometheus.Legend(t.Legend))
	}
	if t.Ref != "" {
		opts = append(opts, prometheus.Ref(t.Ref))
	}

	return opts
}
