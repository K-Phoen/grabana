package influxdb_test

import (
	"testing"

	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/stretchr/testify/require"
)

func TestQueriesCanBeCreated(t *testing.T) {
	req := require.New(t)

	target := influxdb.New("buckets()")

	req.Equal("buckets()", target.Builder.Query)
}

func TestRefCanBeConfigured(t *testing.T) {
	req := require.New(t)

	target := influxdb.New("", influxdb.Ref("A"))

	req.Equal("A", target.Builder.RefID)
}

func TestTargetCanBeHidden(t *testing.T) {
	req := require.New(t)

	target := influxdb.New("", influxdb.Hide())

	req.True(target.Builder.Hide)
}
