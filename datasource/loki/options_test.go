package loki

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", Default())

	req.True(datasource.builder.IsDefault)
}

func TestBasicAuth(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", BasicAuth("joe", "lafrite"))

	req.True(*datasource.builder.BasicAuth)
	req.Equal("joe", *datasource.builder.BasicAuthUser)
	req.Equal("lafrite", *datasource.builder.BasicAuthPassword)
}

func TestTimeout(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", Timeout(30*time.Second))

	req.Equal(30, datasource.builder.JSONData.(map[string]interface{})["timeout"])
}

func TestSkipTlsVerify(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", SkipTLSVerify())

	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"])
}

func TestWithCertificate(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", WithCertificate("certificate-content"))

	req.Equal(false, datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"])
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"])
	req.Equal("certificate-content", datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"])
}

func TestWithCredentials(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", WithCredentials())

	req.True(datasource.builder.WithCredentials)
}

func TestForwardOauthIdentity(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", ForwardOauthIdentity())

	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"])
}

func TestForwardCookies(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", ForwardCookies("foo", "bar"))

	req.ElementsMatch([]string{"foo", "bar"}, datasource.builder.JSONData.(map[string]interface{})["keepCookies"])
}

func TestWithNodeGraph(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", MaximumLines(2000))

	req.Equal(2000, datasource.builder.JSONData.(map[string]interface{})["maxLines"])
}
