package loki

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLokiQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	lokiQuery := "rate({app=\"loki\"}[$__interval])"

	query := New("A", lokiQuery)

	builder := query.Builder

	req.Equal("A", builder.RefID)
	req.Equal("A", builder.Model.RefID)
	req.Equal("time_series", builder.Model.Format)
	req.Equal("loki", builder.Model.Datasource.Type)
	req.Equal(lokiQuery, builder.Model.Expr)
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
