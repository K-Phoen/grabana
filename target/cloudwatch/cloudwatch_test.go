package cloudwatch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCloudwatchTargetCanBeCreated(t *testing.T) {
	req := require.New(t)
	metricName := "Cloudwatch"
	namespace := "Test"

	target := New(metricName, namespace)

	req.Equal(metricName, target.MetricName)
	req.Equal(namespace, target.Namespace)
}

func TestLegendCanBeConfigured(t *testing.T) {
	req := require.New(t)
	legend := "lala"

	target := New("", "", Legend(legend))

	req.Equal(legend, target.LegendFormat)
}

func TestRefCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := New("", "", Ref("A"))

	req.Equal("A", target.Ref)
}

func TestRegionCanBeSet(t *testing.T) {
	req := require.New(t)

	target := New("", "", Region("C"))

	req.Equal("C", target.Region)
}

func TestStatisticCanBeSet(t *testing.T) {
	req := require.New(t)
	statistics := []string{"A", "B", "C"}

	target := New("", "", Statistic(statistics))

	req.Equal(statistics, target.Statistics)
}

func TestDimensionsCanBeSet(t *testing.T) {
	req := require.New(t)
	dimensions := map[string]string{
		"one": "1",
	}

	target := New("", "", Dimensions(dimensions))

	req.Equal(dimensions, target.Dimensions)
}
