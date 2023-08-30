package cloudwatch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCloudWatchQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	cloudwatchQuery := `AVG(SEARCH('{"AWS/ElasticBeanstalk","EnvironmentName","InstanceId"} MetricName="CPUUser" "EnvironmentName"="prod-activity-service-blue"', 'Sum', 60))`

	query := New("A", cloudwatchQuery)

	builder := query.Builder

	req.Equal("A", builder.RefID)
	req.Equal("A", builder.Model.RefID)
	req.Equal("time_series", builder.Model.Format)
	req.Equal("cloudwatch", builder.Model.Datasource.Type)
	req.Equal(cloudwatchQuery, builder.Model.Expr)
	req.Equal(cloudwatchQuery, builder.Model.Target)
	req.NotEqual(0, builder.RelativeTimeRange.From)
	req.Equal(0, builder.RelativeTimeRange.To)
}

func TestTimeRangeCanBeSet(t *testing.T) {
	req := require.New(t)

	query := New("A", "", TimeRange(5*time.Minute, 0))

	builder := query.Builder

	req.NotEqual((5 * time.Minute).Seconds(), builder.RelativeTimeRange.From)
	req.Equal(0, builder.RelativeTimeRange.To)
}

func TestLegendCanBeSet(t *testing.T) {
	req := require.New(t)

	query := New("A", "", Legend("legend"))

	builder := query.Builder

	req.Equal("legend", builder.Model.LegendFormat)
}
