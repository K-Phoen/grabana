package loki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoki(t *testing.T) {
	req := require.New(t)

	datasource := New("ds-loki", "http://localhost:3100")

	req.Equal("ds-loki", datasource.Name())
	req.Equal("http://localhost:3100", datasource.builder.URL)
	req.Equal("loki", datasource.builder.Type)
	req.NotNil(datasource.builder.JSONData)
	req.NotNil(datasource.builder.SecureJSONData)

	_, err := datasource.MarshalJSON()
	req.NoError(err)
}
