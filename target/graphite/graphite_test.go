package graphite_test

import (
	"testing"

	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/stretchr/testify/require"
)

func TestQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	target := graphite.New("stats_counts.statsd.packets_received")

	req.Equal("stats_counts.statsd.packets_received", target.Builder.Target)
}

func TestRefCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := graphite.New("", graphite.Ref("A"))

	req.Equal("A", target.Builder.RefID)
}

func TestTargetCanBeHidden(t *testing.T) {
	req := require.New(t)

	target := graphite.New("", graphite.Hide())

	req.True(target.Builder.Hide)
}
