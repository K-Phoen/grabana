package stackdriver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFiltersCanBeSet(t *testing.T) {
	testCases := []struct {
		desc     string
		opts     []FilterOption
		expected []string
	}{
		{
			desc: "simple eq",
			opts: []FilterOption{
				Eq("property", "value"),
			},
			expected: []string{"property", "=", "value"},
		},
		{
			desc: "simple neq",
			opts: []FilterOption{
				Neq("property", "value"),
			},
			expected: []string{"property", "!=", "value"},
		},
		{
			desc: "simple regex",
			opts: []FilterOption{
				Matches("property", "regex"),
			},
			expected: []string{"property", "=~", "regex"},
		},
		{
			desc: "simple NOT regex",
			opts: []FilterOption{
				NotMatches("property", "regex"),
			},
			expected: []string{"property", "!=~", "regex"},
		},

		{
			desc: "simple AND",
			opts: []FilterOption{
				Eq("property", "value"),
				Neq("other-property", "other-value"),
			},
			expected: []string{"property", "=", "value", "AND", "other-property", "!=", "other-value"},
		},
		{
			desc: "multiple AND",
			opts: []FilterOption{
				Eq("property", "value"),
				Neq("other-property", "other-value"),
				Matches("last-property", "last-value"),
			},
			expected: []string{"property", "=", "value", "AND", "other-property", "!=", "other-value", "AND", "last-property", "=~", "last-value"},
		},
	}

	//nolint: scopelint
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			req := require.New(t)

			query := Delta("", "", Filter(test.opts...))

			req.Equal(test.expected, query.Builder.Model.MetricQuery.Filters)
		})
	}
}
