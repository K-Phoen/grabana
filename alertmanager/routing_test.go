package alertmanager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPolicy(t *testing.T) {
	req := require.New(t)

	policy := Policy("team-a")

	req.Equal("team-a", policy.builder.Receiver)
	req.Empty(policy.builder.ObjectMatchers)
}

func TestTagEq(t *testing.T) {
	req := require.New(t)

	policy := Policy("team-a", TagEq("owner", "team-a"))

	req.Len(policy.builder.ObjectMatchers, 1)

	matcher := policy.builder.ObjectMatchers[0]
	req.Equal("owner", matcher[0])
	req.Equal("=", matcher[1])
	req.Equal("team-a", matcher[2])
}

func TestTagNeq(t *testing.T) {
	req := require.New(t)

	policy := Policy("team-a", TagNeq("severity", "P4"))

	req.Len(policy.builder.ObjectMatchers, 1)

	matcher := policy.builder.ObjectMatchers[0]
	req.Equal("severity", matcher[0])
	req.Equal("!=", matcher[1])
	req.Equal("P4", matcher[2])
}

func TestTagMatches(t *testing.T) {
	req := require.New(t)

	policy := Policy("team-a", TagMatches("owner", "(infra\\-)?platform"))

	req.Len(policy.builder.ObjectMatchers, 1)

	matcher := policy.builder.ObjectMatchers[0]
	req.Equal("owner", matcher[0])
	req.Equal("=~", matcher[1])
	req.Equal("(infra\\-)?platform", matcher[2])
}

func TestTagNotMatches(t *testing.T) {
	req := require.New(t)

	policy := Policy("team-a", TagNotMatches("severity", "P[345]"))

	req.Len(policy.builder.ObjectMatchers, 1)

	matcher := policy.builder.ObjectMatchers[0]
	req.Equal("severity", matcher[0])
	req.Equal("!~", matcher[1])
	req.Equal("P[345]", matcher[2])
}
