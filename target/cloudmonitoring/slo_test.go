package cloudmonitoring

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/assert"
)

func TestSLO(t *testing.T) {
	for _, testCase := range []struct {
		desc       string
		options    []SLOOption
		wantTarget *sdk.Target
	}{
		{
			desc: "default",
			wantTarget: &sdk.Target{
				QueryType: "slo",
				SLOQuery: &sdk.StackdriverSLOQuery{
					ProjectName: testProjectName,
				},
			},
		},
		{
			desc: "all options",
			options: []SLOOption{
				SLOPerSeriesAligner(AlignPercentChange, "10s"),
				SLOAliasBy("banana"),
				SelectorName("bah"),
				ServiceRef("abc", "service"),
				SLORef("cde", "slo"),
				LookbackPeriod("120m"),
			},
			wantTarget: &sdk.Target{
				QueryType: "slo",
				SLOQuery: &sdk.StackdriverSLOQuery{
					ProjectName:      testProjectName,
					AlignmentPeriod:  "10s",
					PerSeriesAligner: string(AlignPercentChange),
					AliasBy:          "banana",
					SelectorName:     "bah",
					ServiceID:        "abc",
					ServiceName:      "service",
					SLOID:            "cde",
					SLOName:          "slo",
					LookbackPeriod:   "120m",
				},
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			assert.Equal(
				t,
				testCase.wantTarget,
				NewSLO(
					testProjectName,
					testCase.options...,
				).Target(),
			)
		})
	}
}
