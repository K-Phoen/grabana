package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/grabana/target/cloudwatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidValueTypes(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{input: "min"},
		{input: "max"},
		{input: "avg"},
		{input: "current"},
		{input: "total"},
		{input: "first"},
		{input: "delta"},
		{input: "diff"},
		{input: "range"},
		{input: "name"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardSingleStat{ValueType: tc.input}

			opt, err := panel.valueType()

			req.NoError(err)

			singleStat, err := singlestat.New("test")
			req.NoError(err)

			req.NoError(opt(singleStat))

			req.Equal(tc.input, singleStat.Builder.SinglestatPanel.ValueName)
		})
	}
}

func TestSparkLineModes(t *testing.T) {
	testCases := []struct {
		input string
		err   error
	}{
		{input: "", err: nil},
		{input: "bottom", err: nil},
		{input: "full", err: nil},
		{input: "invalid", err: ErrInvalidSparkLineMode},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardSingleStat{SparkLine: tc.input}

			_, err := panel.toOption()

			if tc.err == nil {
				req.NoError(err)
			} else {
				req.Equal(tc.err, err)
			}
		})
	}
}

func TestSinglestatCloudwatchTarget(t *testing.T) {
	panel := DashboardSingleStat{
		Title: "Singlestat target test",
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
