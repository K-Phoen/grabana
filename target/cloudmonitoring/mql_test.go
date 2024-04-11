package cloudmonitoring

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/assert"
)

func TestMQL(t *testing.T) {
	for _, testCase := range []struct {
		desc       string
		query      string
		options    []MQLOption
		wantTarget *sdk.Target
	}{
		{
			desc:  "default",
			query: "project_id=\"blublu\"",
			wantTarget: &sdk.Target{
				QueryType: "timeSeriesQuery",
				TimeSeriesQuery: &sdk.StackdriverTimeSeriesQuery{
					ProjectName: testProjectName,
					Query:       "project_id=\"blublu\"",
				},
			},
		},
		{
			desc:    "all options",
			query:   "project_id=\"blublu\"",
			options: []MQLOption{GraphPeriod("120s"), MQLAliasBy("Bozo")},
			wantTarget: &sdk.Target{
				QueryType: "timeSeriesQuery",
				AliasBy:   "Bozo",
				TimeSeriesQuery: &sdk.StackdriverTimeSeriesQuery{
					ProjectName: testProjectName,
					Query:       "project_id=\"blublu\"",
					GraphPeriod: "120s",
				},
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			assert.Equal(
				t,
				testCase.wantTarget,
				NewMQL(
					testProjectName,
					testCase.query,
					testCase.options...,
				).Target(),
			)
		})
	}
}
