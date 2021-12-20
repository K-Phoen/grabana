package prometheus

import (
	"net/http"
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

func TestAccessMode(t *testing.T) {
	testCases := []struct {
		mode     Access
		expected string
	}{
		{
			mode:     Proxy,
			expected: "proxy",
		},
		{
			mode:     Browser,
			expected: "direct",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.expected, func(t *testing.T) {
			req := require.New(t)

			datasource := New("", "", AccessMode(tc.mode))

			req.Equal(tc.expected, datasource.builder.Access)
		})
	}
}

func TestHTTPMethod(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", HTTPMethod(http.MethodGet))

	req.Equal(http.MethodGet, datasource.builder.JSONData.(map[string]interface{})["httpMethod"])
}

func TestScrapeInterval(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", ScrapeInterval(10*time.Second))

	req.Equal("10s", datasource.builder.JSONData.(map[string]interface{})["timeInterval"])
}

func TestQueryTimeout(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", QueryTimeout(30*time.Second))

	req.Equal("30s", datasource.builder.JSONData.(map[string]interface{})["queryTimeout"])
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
