package cloudwatch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", Default())

	req.NoError(err)
	req.True(datasource.builder.IsDefault)
}

func TestDefaultAuth(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", DefaultAuth())

	req.NoError(err)
	req.Equal("default", datasource.builder.JSONData.(map[string]interface{})["authType"])
}

func TestAccessSecretAuth(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", AccessSecretAuth("access", "secret"))

	req.NoError(err)
	req.Equal("keys", datasource.builder.JSONData.(map[string]interface{})["authType"])
	req.Equal("access", datasource.builder.SecureJSONData.(map[string]interface{})["accessKey"])
	req.Equal("secret", datasource.builder.SecureJSONData.(map[string]interface{})["secretKey"])
}

func TestDefaultRegion(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", DefaultRegion("eu-north-1"))

	req.NoError(err)
	req.Equal("eu-north-1", datasource.builder.JSONData.(map[string]interface{})["defaultRegion"])
}

func TestAssumeRoleARN(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", AssumeRoleARN("arn:aws:iam:role"))

	req.NoError(err)
	req.Equal("arn:aws:iam:role", datasource.builder.JSONData.(map[string]interface{})["assumeRoleArn"])
}

func TestExternalID(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", ExternalID("external-id"))

	req.NoError(err)
	req.Equal("external-id", datasource.builder.JSONData.(map[string]interface{})["externalId"])
}

func TestEndpoint(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", Endpoint("endpoint"))

	req.NoError(err)
	req.Equal("endpoint", datasource.builder.JSONData.(map[string]interface{})["endpoint"])
}

func TestCustomMetricsNamespaces(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", CustomMetricsNamespaces("ns1", "ns2"))

	req.NoError(err)
	req.Equal("ns1,ns2", datasource.builder.JSONData.(map[string]interface{})["customMetricsNamespaces"])
}
