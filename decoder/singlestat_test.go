package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/singlestat"

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
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardSingleStat{ValueType: tc.input}

			opt, err := panel.valueType()

			req.NoError(err)

			singleStat := singlestat.New("test")
			opt(singleStat)

			req.Equal(tc.input, singleStat.Builder.ValueName)
		})
	}
}
