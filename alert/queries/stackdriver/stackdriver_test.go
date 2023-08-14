package stackdriver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProjectCanBeConfigured(t *testing.T) {
	req := require.New(t)
	
	project := "some-gcp-project-id"

	query := Delta("A", "some metric", Project(project))

	builder := query.Builder

	req.Equal(project, builder.Model.MetricQuery.ProjectName)
}

func TestDeltaQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	metric := "some metric type for delta"

	query := Delta("A", metric)

	builder := query.Builder

	req.Equal("A", builder.RefID)
	req.Equal("A", builder.Model.RefID)
	req.Equal("metrics", builder.Model.QueryType)
	req.Equal("time_series", builder.Model.Format)
	req.Equal("DELTA", builder.Model.MetricQuery.MetricKind)
	req.Equal(metric, builder.Model.MetricQuery.MetricType)
	req.Equal("stackdriver", builder.Model.Datasource.Type)
}

func TestGaugeQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	metric := "some metric type for gauge"

	query := Gauge("A", metric)

	builder := query.Builder

	req.Equal("A", builder.RefID)
	req.Equal("A", builder.Model.RefID)
	req.Equal("metrics", builder.Model.QueryType)
	req.Equal("time_series", builder.Model.Format)
	req.Equal("GAUGE", builder.Model.MetricQuery.MetricKind)
	req.Equal(metric, builder.Model.MetricQuery.MetricType)
	req.Equal("stackdriver", builder.Model.Datasource.Type)
}

func TestCumulativeQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	metric := "some metric type"

	query := Cumulative("A", metric)

	builder := query.Builder

	req.Equal("A", builder.RefID)
	req.Equal("A", builder.Model.RefID)
	req.Equal("metrics", builder.Model.QueryType)
	req.Equal("time_series", builder.Model.Format)
	req.Equal("CUMULATIVE", builder.Model.MetricQuery.MetricKind)
	req.Equal(metric, builder.Model.MetricQuery.MetricType)
	req.Equal("stackdriver", builder.Model.Datasource.Type)
}

func TestTimeRangeCanBeSet(t *testing.T) {
	req := require.New(t)

	query := Delta("A", "some metricTpe", TimeRange(5*time.Minute, 0))

	builder := query.Builder

	req.NotEqual((5 * time.Minute).Seconds(), builder.RelativeTimeRange.From)
	req.Equal(0, builder.RelativeTimeRange.To)
}

func TestLegendCanBeSet(t *testing.T) {
	req := require.New(t)

	query := Delta("A", "some metricTpe", Legend("legend"))

	builder := query.Builder

	req.Equal("legend", builder.Model.MetricQuery.AliasBy)
}

func TestAggregationCanBeSet(t *testing.T) {
	req := require.New(t)
	reducers := []Reducer{
		ReduceNone,
		ReduceMean,
		ReduceMin,
		ReduceMax,
		ReduceSum,
		ReduceStdDev,
		ReduceCount,
		ReduceCountTrue,
		ReduceCountFalse,
		ReduceCountFractionTrue,
		ReducePercentile99,
		ReducePercentile95,
		ReducePercentile50,
		ReducePercentile05,
	}

	for _, reducer := range reducers {
		query := Delta("", "", Aggregation(reducer))

		req.Equal(string(reducer), query.Builder.Model.MetricQuery.CrossSeriesReducer)
	}
}

func TestAlignmentCanBeSet(t *testing.T) {
	req := require.New(t)
	aligners := []Aligner{
		AlignNone,
		AlignDelta,
		AlignRate,
		AlignNextOlder,
		AlignMin,
		AlignMax,
		AlignMean,
		AlignCount,
		AlignSum,
		AlignStdDev,
		AlignCountTrue,
		AlignCountFalse,
		AlignFractionTrue,
		AlignPercentile99,
		AlignPercentile95,
		AlignPercentile50,
		AlignPercentile05,
		AlignPercentChange,
	}

	for _, aligner := range aligners {
		query := Delta("", "", Alignment(aligner, AlignmentStackdriverAuto))

		req.Equal(string(aligner), query.Builder.Model.MetricQuery.PerSeriesAligner)
		req.Equal("stackdriver-auto", query.Builder.Model.MetricQuery.AlignmentPeriod)
	}
}

func TestGroupBysCanBeSet(t *testing.T) {
	req := require.New(t)

	query := Delta("", "", GroupBys("field", "other"))

	req.ElementsMatch(query.Builder.Model.MetricQuery.GroupBys, []string{"field", "other"})
}

func TestPreprocessorCanBeSet(t *testing.T) {
	req := require.New(t)
	preprocessors := []PreprocessorMethod{
		PreprocessDelta,
		PreprocessRate,
	}

	for _, preprocessor := range preprocessors {
		query := Delta("", "", Preprocessor(preprocessor))

		req.Equal(string(preprocessor), query.Builder.Model.MetricQuery.Preprocessor)
	}
}
