package graphite

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGraphiteQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	graphiteQuery := "aliasByMetric(stats_counts.*.*)"

	query := New("A", graphiteQuery)

	builder := query.Builder

	req.Equal("A", builder.RefID)
	req.Equal("A", builder.Model.RefID)
	req.Equal("time_series", builder.Model.Format)
	req.Equal("graphite", builder.Model.Datasource.Type)
	req.Equal(graphiteQuery, builder.Model.Expr)
	req.Equal(graphiteQuery, builder.Model.Target)
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
