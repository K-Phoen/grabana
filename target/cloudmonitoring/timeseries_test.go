package cloudmonitoring

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/assert"
)

const (
	testProjectName = "joelafrite"
	testMetricType  = "pubsub.googleapis.com/subscription/num_undelivered_messages"
)

func TestNewTimeSeries(t *testing.T) {
	for _, testCase := range []struct {
		desc       string
		options    []TimeSeriesOption
		wantTarget *sdk.Target
	}{
		{
			desc: "default values",
			wantTarget: &sdk.Target{
				QueryType: "timeSeriesList",
				TimeSeriesList: &sdk.StackdriverTimeSeriesList{
					ProjectName: testProjectName,
					Filters: []string{
						"metric.type",
						"=",
						testMetricType,
					},
				},
			},
		},
		{
			desc: "multiple filters",
			options: []TimeSeriesOption{
				Filter("potato", FilterOperatorEqual, "patata"),
			},
			wantTarget: &sdk.Target{
				QueryType: "timeSeriesList",
				TimeSeriesList: &sdk.StackdriverTimeSeriesList{
					ProjectName: testProjectName,
					Filters: []string{
						"metric.type",
						"=",
						testMetricType,
						"AND",
						"potato",
						"=",
						"patata",
					},
				},
			},
		},
		{
			desc: "all options",
			options: []TimeSeriesOption{
				CrossSeriesReducer(ReduceCountTrue),
				PerSeriesAligner(AlignMax, "12d"),
				GroupBy("foo"),
				GroupBy("bar"),
				GroupBy("biz"),
				View("lafritte"),
				Title("banana"),
				SecondaryCrossSeriesReducer(ReduceCountFalse),
				SecondaryPerSeriesAligner(AlignCount, "12s"),
				SecondaryGroupBy("far"),
				SecondaryGroupBy("fuz"),
				SecondaryGroupBy("fiz"),
				Preprocessor(PreprocessDelta),
				TimeSeriesAliasBy("pwet"),
			},
			wantTarget: &sdk.Target{
				QueryType: "timeSeriesList",
				AliasBy:   "pwet",
				TimeSeriesList: &sdk.StackdriverTimeSeriesList{
					ProjectName:        testProjectName,
					CrossSeriesReducer: string(ReduceCountTrue),
					PerSeriesAligner:   string(AlignMax),
					AlignmentPeriod:    "12d",
					GroupBys: []string{
						"foo",
						"bar",
						"biz",
					},
					View:                        "lafritte",
					Title:                       "banana",
					SecondaryCrossSeriesReducer: string(ReduceCountFalse),
					SecondaryPerSeriesAligner:   string(AlignCount),
					SecondaryAlignmentPeriod:    "12s",
					SecondaryGroupBys: []string{
						"far",
						"fuz",
						"fiz",
					},
					Preprocessor: string(PreprocessDelta),
					Filters: []string{
						"metric.type",
						"=",
						testMetricType,
					},
				},
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			assert.Equal(t,
				testCase.wantTarget,
				NewTimeSeries(testProjectName, testMetricType, testCase.options...).target,
			)
		})
	}
}
