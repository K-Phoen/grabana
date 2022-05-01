package decoder

import (
	"testing"
	"time"

	alertBuilder "github.com/K-Phoen/grabana/alert"
	"github.com/stretchr/testify/require"
)

func TestDecodingAlertTargetFailsIfNoTargetIsProvided(t *testing.T) {
	target := AlertTarget{}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrTargetNotConfigured, err)
}

func TestDecodingAPrometheusTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Prometheus: &AlertPrometheus{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingAPrometheusTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Prometheus: &AlertPrometheus{
			Ref:      "A",
			Query:    "prom-query",
			Legend:   "{{ code }}",
			Lookback: "15m",
		},
	}

	opt, err := target.toOption()

	req.NoError(err)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 1)

	promQuery := alert.Builder.Rules[0].GrafanaAlert.Data[0]

	req.Equal("A", promQuery.RefID)
	req.Equal("A", promQuery.Model.RefID)
	req.Equal("prom-query", promQuery.Model.Expr)
	req.Equal("{{ code }}", promQuery.Model.LegendFormat)
	req.Equal(int((15 * time.Minute).Seconds()), promQuery.RelativeTimeRange.From)
}

func TestDecodingAPrometheusTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Prometheus: &AlertPrometheus{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingALokiTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Loki: &AlertLoki{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingALokiTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Loki: &AlertLoki{
			Ref:      "A",
			Query:    "loki-query",
			Legend:   "{{ status }}",
			Lookback: "15m",
		},
	}

	opt, err := target.toOption()

	req.NoError(err)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 1)

	promQuery := alert.Builder.Rules[0].GrafanaAlert.Data[0]

	req.Equal("A", promQuery.RefID)
	req.Equal("A", promQuery.Model.RefID)
	req.Equal("loki-query", promQuery.Model.Expr)
	req.Equal("{{ status }}", promQuery.Model.LegendFormat)
	req.Equal(int((15 * time.Minute).Seconds()), promQuery.RelativeTimeRange.From)
}

func TestDecodingALokiTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Loki: &AlertLoki{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingAGraphiteTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Graphite: &AlertGraphite{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingAGraphiteTarget(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Graphite: &AlertGraphite{
			Ref:      "A",
			Query:    "graphite-query",
			Lookback: "15m",
		},
	}

	opt, err := target.toOption()

	req.NoError(err)

	alert := alertBuilder.New("", opt)

	req.Len(alert.Builder.Rules, 1)
	req.Len(alert.Builder.Rules[0].GrafanaAlert.Data, 1)

	promQuery := alert.Builder.Rules[0].GrafanaAlert.Data[0]

	req.Equal("A", promQuery.RefID)
	req.Equal("A", promQuery.Model.RefID)
	req.Equal("graphite-query", promQuery.Model.Expr)
	req.Equal(int((15 * time.Minute).Seconds()), promQuery.RelativeTimeRange.From)
}

func TestDecodingAGraphiteTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Graphite: &AlertGraphite{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}

func TestDecodingAStackdriverTargetWithNoRefFails(t *testing.T) {
	target := AlertTarget{
		Stackdriver: &AlertStackdriver{},
	}

	_, err := target.toOption()

	require.Error(t, err)
	require.Equal(t, ErrMissingRef, err)
}

func TestDecodingAStackdriverTargetWithInvalidLookback(t *testing.T) {
	req := require.New(t)

	target := AlertTarget{
		Stackdriver: &AlertStackdriver{
			Ref:      "A",
			Lookback: "not a valid duration",
		},
	}

	_, err := target.toOption()

	req.ErrorIs(err, ErrInvalidLookback)
}
