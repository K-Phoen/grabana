package prometheus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPrometheus(t *testing.T) {
	req := require.New(t)

	datasource, err := New("ds-name", "http://localhost:9090")

	req.NoError(err)
	req.Equal("ds-name", datasource.Name())
	req.Equal("http://localhost:9090", datasource.builder.URL)
	req.Equal("prometheus", datasource.builder.Type)
	req.NotNil(datasource.builder.JSONData)
	req.NotNil(datasource.builder.SecureJSONData)

	_, err = datasource.MarshalJSON()
	req.NoError(err)
}
