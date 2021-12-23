package jaeger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJaeger(t *testing.T) {
	req := require.New(t)

	datasource := New("ds-jaeger", "http://localhost:16686")

	req.Equal("ds-jaeger", datasource.Name())
	req.Equal("http://localhost:16686", datasource.builder.URL)
	req.Equal("jaeger", datasource.builder.Type)
	req.NotNil(datasource.builder.JSONData)
	req.NotNil(datasource.builder.SecureJSONData)

	_, err := datasource.MarshalJSON()
	req.NoError(err)
}
