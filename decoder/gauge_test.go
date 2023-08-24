package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/gauge"
	"github.com/K-Phoen/grabana/target/cloudwatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGaugeValidValueTypes(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "min", expected: "min"},
		{input: "max", expected: "max"},
		{input: "avg", expected: "mean"},
		{input: "count", expected: "count"},
		{input: "total", expected: "sum"},
		{input: "range", expected: "range"},
		{input: "first", expected: "first"},
		{input: "first_non_null", expected: "firstNotNull"},
		{input: "last", expected: "last"},
		{input: "last_non_null", expected: "lastNotNull"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardGauge{ValueType: tc.input}

			opt, err := panel.valueType()

			req.NoError(err)

			gaugePanel, err := gauge.New("")
			req.NoError(err)

			req.NoError(opt(gaugePanel))

			req.Equal(tc.expected, gaugePanel.Builder.GaugePanel.Options.ReduceOptions.Calcs[0])
		})
	}
}

func TestGaugeCloudwatchTarget(t *testing.T) {

	panel := DashboardGauge{
		Title: "cloudwatch target test",
		Targets: []Target{
			{
				Cloudwatch: &CloudwatchTarget{
					QueryParams: cloudwatch.QueryParams{
						Dimensions: map[string]string{
							"Name": "test",
						},
					},
				},
			},
		},
	}

	option, err := panel.target(panel.Targets[0])
	assert.NoError(t, err)
	assert.NotNil(t, option)
}
