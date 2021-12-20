package stackdriver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	req := require.New(t)

	datasource := New("", Default())

	req.True(datasource.builder.IsDefault)
}

func TestGCEAuthentication(t *testing.T) {
	req := require.New(t)

	datasource := New("", GCEAuthentication())

	req.Equal("gce", datasource.builder.JSONData.(map[string]interface{})["authenticationType"])
}

func TestJWTAuthentication(t *testing.T) {
	req := require.New(t)
	jwtKey := `{
  "type": "service_account",
  "project_id": "dark-sandbox",
  "private_key_id": "unused",
  "private_key": "-----BEGIN PRIVATE KEY-----\nSOMETHING_REALLY_SECRET_HERE\n-----END PRIVATE KEY-----\n",
  "client_email": "dark-operator@dark-sandbox.iam.gserviceaccount.com",
  "client_id": "unused",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/dark-operator%40dark-sandbox.iam.gserviceaccount.com"
}`

	datasource := New("", JWTAuthentication(jwtKey))

	req.Equal("jwt", datasource.builder.JSONData.(map[string]interface{})["authenticationType"])
	req.Equal("dark-operator@dark-sandbox.iam.gserviceaccount.com", datasource.builder.JSONData.(map[string]interface{})["clientEmail"])
	req.Equal("dark-sandbox", datasource.builder.JSONData.(map[string]interface{})["defaultProject"])
	req.Equal("https://oauth2.googleapis.com/token", datasource.builder.JSONData.(map[string]interface{})["tokenUri"])
	req.Equal("-----BEGIN PRIVATE KEY-----\nSOMETHING_REALLY_SECRET_HERE\n-----END PRIVATE KEY-----\n", datasource.builder.SecureJSONData.(map[string]interface{})["privateKey"])
}

func TestJWTAuthenticationUnfortunatelyFailsSilentlyIfJWTIsInvalid(t *testing.T) {
	req := require.New(t)
	jwtKey := `{invalid json`

	datasource := New("", JWTAuthentication(jwtKey))

	req.Equal("gce", datasource.builder.JSONData.(map[string]interface{})["authenticationType"])
}
