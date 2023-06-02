package cloudwatch_test

import (
	"testing"

	"github.com/K-Phoen/grabana/target/cloudwatch"
	"github.com/stretchr/testify/require"
)

func TestQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	dimensions := map[string]string{
		"QueueName": "test-queue",
	}
	statistics := []string{"Sum"}
	namespace := "AWS/SQS"
	metricName := "NumberOfMessagesReceived"
	period := "30"
	region := "us-east-1"

	target := cloudwatch.New(dimensions, statistics, namespace, metricName, period, region)

	req.Equal(dimensions, target.Builder.Dimensions)
	req.Equal(statistics, target.Builder.Statistics)
	req.Equal(namespace, target.Builder.Namespace)
	req.Equal(metricName, target.Builder.MetricName)
	req.Equal(period, target.Builder.Period)
	req.Equal(region, target.Builder.Region)
}

func TestRefCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := cloudwatch.New(nil, nil, "", "", "", "", cloudwatch.Ref("A"))

	req.Equal("A", target.Builder.RefID)
}

func TestRefCanBeHidden(t *testing.T) {
	req := require.New(t)

	target := cloudwatch.New(nil, nil, "", "", "", "", cloudwatch.Hide())

	req.True(target.Builder.Hide)
}
