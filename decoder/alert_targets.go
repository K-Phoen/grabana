package decoder

import (
	"fmt"
	"time"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/alert/queries/prometheus"
)

var ErrMissingRef = fmt.Errorf("target ref missing")

type AlertTarget struct {
	Prometheus *AlertPrometheus `yaml:",omitempty"`
}

func (t AlertTarget) toOption() (string, alert.Option, error) {
	if t.Prometheus != nil {
		return t.Prometheus.toOptions()
	}

	return "", nil, ErrTargetNotConfigured
}

type AlertPrometheus struct {
	Ref      string `yaml:",omitempty"`
	Query    string
	Legend   string `yaml:",omitempty"`
	Lookback string `yaml:",omitempty"`
}

func (t AlertPrometheus) toOptions() (string, alert.Option, error) {
	opts := []prometheus.Option{
		prometheus.Legend(t.Legend),
	}

	if t.Ref == "" {
		return "", nil, ErrMissingRef
	}

	if t.Lookback != "" {
		from, err := time.ParseDuration(t.Lookback)
		if err != nil {
			return "", nil, err
		}

		opts = append(opts, prometheus.TimeRange(from, 0))
	}

	return t.Ref, alert.WithPrometheusQuery(t.Ref, t.Query, opts...), nil
}
