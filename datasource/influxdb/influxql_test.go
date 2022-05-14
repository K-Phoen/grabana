package influxdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInfluxQL(t *testing.T) {
	req := require.New(t)

	datasource, err := New("ds-name", "http://localhost:9090")

	req.NoError(err)
	req.Equal("ds-name", datasource.Name())
	req.Equal("http://localhost:9090", datasource.builder.URL)
	req.Equal("influxdb", datasource.builder.Type)
	req.Equal("proxy", datasource.builder.Access)
	req.Equal(
		map[string]interface{}{
			"version":    "InfluxQL",
			"httpMethod": "GET",
			"maxSeries":  1000,
		},
		datasource.builder.JSONData,
	)
	req.NotNil(datasource.builder.SecureJSONData)

	_, err = datasource.MarshalJSON()
	req.NoError(err)
}
