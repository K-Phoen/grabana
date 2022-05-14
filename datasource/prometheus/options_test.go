package prometheus

import (
	"net/http"
	"testing"
	"time"

	"github.com/K-Phoen/grabana/errors"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", Default())

	req.NoError(err)
	req.True(datasource.builder.IsDefault)
}

func TestBasicAuth(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", BasicAuth("joe", "lafrite"))

	req.NoError(err)
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

			datasource, err := New("", "", AccessMode(tc.mode))

			req.NoError(err)
			req.Equal(tc.expected, datasource.builder.Access)
		})
	}
}

func TestHTTPMethod(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", HTTPMethod(http.MethodGet))

	req.NoError(err)
	req.Equal(http.MethodGet, datasource.builder.JSONData.(map[string]interface{})["httpMethod"])
}

func TestHTTPMethodRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	_, err := New("", "", HTTPMethod(http.MethodPut))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestScrapeInterval(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", ScrapeInterval(10*time.Second))

	req.NoError(err)
	req.Equal("10s", datasource.builder.JSONData.(map[string]interface{})["timeInterval"])
}

func TestQueryTimeout(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", QueryTimeout(30*time.Second))

	req.NoError(err)
	req.Equal("30s", datasource.builder.JSONData.(map[string]interface{})["queryTimeout"])
}

func TestSkipTlsVerify(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", SkipTLSVerify())

	req.NoError(err)
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"])
}

func TestWithCertificate(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", WithCertificate("certificate-content"))

	req.NoError(err)
	req.Equal(false, datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"])
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"])
	req.Equal("certificate-content", datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"])
}

func TestWithCredentials(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", WithCredentials())

	req.NoError(err)
	req.True(datasource.builder.WithCredentials)
}

func TestForwardOauthIdentity(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", ForwardOauthIdentity())

	req.NoError(err)
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"])
}

func TestForwardCookies(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", ForwardCookies("foo", "bar"))

	req.NoError(err)
	req.ElementsMatch([]string{"foo", "bar"}, datasource.builder.JSONData.(map[string]interface{})["keepCookies"])
}

func TestExemplars(t *testing.T) {
	req := require.New(t)

	traceIDExemplar := Exemplar{
		LabelName:     "traceID",
		DatasourceUID: "tempo",
	}

	datasource, err := New("", "", Exemplars(traceIDExemplar))

	req.NoError(err)
	req.ElementsMatch([]Exemplar{traceIDExemplar}, datasource.builder.JSONData.(map[string]interface{})["exemplarTraceIdDestinations"])
}
