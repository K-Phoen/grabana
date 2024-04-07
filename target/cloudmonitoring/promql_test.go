package cloudmonitoring

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/assert"
)

func TestPromQL(t *testing.T) {
	for _, testCase := range []struct {
		desc       string
		expr       string
		options    []PromQLOption
		wantTarget *sdk.Target
	}{
		{
			desc: "default",
			expr: "uptime{foo=\"bar\"}",
			wantTarget: &sdk.Target{
				QueryType: "promQL",
				TimeSeriesQuery: &sdk.StackdriverTimeSeriesQuery{
					ProjectName: testProjectName,
				},
				PromQLQuery: &sdk.StackdriverPromQLQuery{
					ProjectName: testProjectName,
					Expr:        "uptime{foo=\"bar\"}",
					Step:        "10s",
				},
			},
		},
		{
			desc:    "with min step",
			expr:    "uptime{foo=\"bar\"}",
			options: []PromQLOption{MinStep("120s")},
			wantTarget: &sdk.Target{
				QueryType: "promQL",
				TimeSeriesQuery: &sdk.StackdriverTimeSeriesQuery{
					ProjectName: testProjectName,
				},
				PromQLQuery: &sdk.StackdriverPromQLQuery{
					ProjectName: testProjectName,
					Expr:        "uptime{foo=\"bar\"}",
					Step:        "120s",
				},
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			assert.Equal(
				t,
				testCase.wantTarget,
				NewPromQL(
					testProjectName,
					testCase.expr,
					testCase.options...,
				).Target(),
			)
		})
	}
}
