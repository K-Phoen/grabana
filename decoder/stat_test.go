package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/stat"
	"github.com/stretchr/testify/require"
)

func TestStatValidValueTypes(t *testing.T) {
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

			panel := DashboardStat{ValueType: tc.input}

			opt, err := panel.valueType()

			req.NoError(err)

			statPanel, err := stat.New("")
			req.NoError(err)

			req.NoError(opt(statPanel))

			req.Equal(tc.expected, statPanel.Builder.StatPanel.Options.ReduceOptions.Calcs[0])
		})
	}
}

func TestStatValidColorMode(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "background", expected: "background"},
		{input: "value", expected: "value"},
		{input: "none", expected: "none"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardStat{ColorMode: tc.input}

			opt, err := panel.colorMode()

			req.NoError(err)

			statPanel, err := stat.New("")
			req.NoError(err)

			req.NoError(opt(statPanel))

			req.Equal(tc.expected, statPanel.Builder.StatPanel.Options.ColorMode)
		})
	}
}
