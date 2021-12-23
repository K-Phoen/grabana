package tempo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTempo(t *testing.T) {
	req := require.New(t)

	datasource := New("ds-tempo", "http://localhost:3100")

	req.Equal("ds-tempo", datasource.Name())
	req.Equal("http://localhost:3100", datasource.builder.URL)
	req.Equal("tempo", datasource.builder.Type)
	req.NotNil(datasource.builder.JSONData)
	req.NotNil(datasource.builder.SecureJSONData)

	_, err := datasource.MarshalJSON()
	req.NoError(err)
}
