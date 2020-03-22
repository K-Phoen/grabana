package decoder

import (
	"github.com/K-Phoen/grabana/target/prometheus"
)

type target struct {
	Prometheus *prometheusTarget
}

type prometheusTarget struct {
	Query  string
	Legend string
	Ref    string
}

func (t prometheusTarget) toOptions() []prometheus.Option {
	var opts []prometheus.Option

	if t.Legend != "" {
		opts = append(opts, prometheus.Legend(t.Legend))
	}
	if t.Ref != "" {
		opts = append(opts, prometheus.Ref(t.Ref))
	}

	return opts
}
