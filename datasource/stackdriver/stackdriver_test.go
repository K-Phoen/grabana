package stackdriver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStackdriver(t *testing.T) {
	req := require.New(t)

	datasource := New("ds-name")

	req.Equal("ds-name", datasource.Name())
	req.Equal("stackdriver", datasource.builder.Type)
	req.NotNil(datasource.builder.JSONData)
	req.NotNil(datasource.builder.SecureJSONData)

	_, err := datasource.MarshalJSON()
	req.NoError(err)
}
