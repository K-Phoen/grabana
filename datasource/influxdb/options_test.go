package influxdb

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

func TestHTTPMethod(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", HTTPMethod(http.MethodPost))
	req.Equal(http.MethodPost, datasource.builder.JSONData.(map[string]interface{})["httpMethod"])
}

func TestAccessMode(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", AccessMode(Browser))
	req.Equal(string(Browser), datasource.builder.Access)
}

func TestKeepCookies(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", KeepCookies([]string{"foo", "bar"}))
	req.Equal([]string{"foo", "bar"}, datasource.builder.JSONData.(map[string]interface{})["keepCookies"])
}

func TestTimeout(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", Timeout(10*time.Second))
	req.Equal(10, datasource.builder.JSONData.(map[string]interface{})["timeout"])
}

func TestDatabase(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", Database("lolilol"))
	req.Equal(stringPtr("lolilol"), datasource.builder.Database)
}

func TestUser(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", User("lolilol"))
	req.Equal(stringPtr("lolilol"), datasource.builder.User)
}

func TestPassword(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", Password("whereisthepoulette"))
	req.Equal("whereisthepoulette", datasource.builder.SecureJSONData.(map[string]interface{})["password"])
}

func TestMinTimeInterval(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", MinTimeInterval(10*time.Second))
	req.Equal("10s", datasource.builder.JSONData.(map[string]interface{})["timeInterval"])
}

func TestMaxSeries(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", MaxSeries(43234))
	req.Equal(43234, datasource.builder.JSONData.(map[string]interface{})["maxSeries"])
}

func TestBasicAuth(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", BasicAuth("john", "doe"))
	req.Equal(boolPtr(true), datasource.builder.BasicAuth)
	req.Equal(stringPtr("john"), datasource.builder.BasicAuthUser)
	req.Equal(stringPtr("doe"), datasource.builder.BasicAuthPassword)
}

func TestWithCredentials(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", WithCredentials())
	req.Equal(true, datasource.builder.WithCredentials)
}

func TestSkipTLSVerify(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", SkipTLSVerify())
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"])
}

func TestForwardOauthIdentity(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", ForwardOauthIdentity())
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"])
}

func TestTLSClientAuth(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", TLSClientAuth("Foo", "bar"))
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsAuth"])
	req.Equal("Foo", datasource.builder.SecureJSONData.(map[string]interface{})["tlsClientCert"])
	req.Equal("bar", datasource.builder.SecureJSONData.(map[string]interface{})["tlsClientKey"])
}

func TestWithCACert(t *testing.T) {
	req := require.New(t)

	datasource := New("", "", WithCACert("bozo"))
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"])
	req.Equal("bozo", datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"])
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
