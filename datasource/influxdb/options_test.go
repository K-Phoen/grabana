package influxdb

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

func TestHTTPMethod(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", HTTPMethod(http.MethodPost))

	req.NoError(err)
	req.Equal(http.MethodPost, datasource.builder.JSONData.(map[string]interface{})["httpMethod"])
}

func TestHTTPMethodRejectsInvalidValues(t *testing.T) {
	req := require.New(t)

	_, err := New("", "", HTTPMethod(http.MethodPut))

	req.Error(err)
	req.ErrorIs(err, errors.ErrInvalidArgument)
}

func TestAccessMode(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", AccessMode(Browser))

	req.NoError(err)
	req.Equal(string(Browser), datasource.builder.Access)
}

func TestKeepCookies(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", KeepCookies([]string{"foo", "bar"}))

	req.NoError(err)
	req.Equal([]string{"foo", "bar"}, datasource.builder.JSONData.(map[string]interface{})["keepCookies"])
}

func TestTimeout(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", Timeout(10*time.Second))

	req.NoError(err)
	req.Equal(10, datasource.builder.JSONData.(map[string]interface{})["timeout"])
}

func TestDatabase(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", Database("lolilol"))

	req.NoError(err)
	req.Equal(stringPtr("lolilol"), datasource.builder.Database)
}

func TestUser(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", User("lolilol"))

	req.NoError(err)
	req.Equal(stringPtr("lolilol"), datasource.builder.User)
}

func TestPassword(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", Password("whereisthepoulette"))

	req.NoError(err)
	req.Equal("whereisthepoulette", datasource.builder.SecureJSONData.(map[string]interface{})["password"])
}

func TestMinTimeInterval(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", MinTimeInterval(10*time.Second))

	req.NoError(err)
	req.Equal("10s", datasource.builder.JSONData.(map[string]interface{})["timeInterval"])
}

func TestMaxSeries(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", MaxSeries(43234))

	req.NoError(err)
	req.Equal(43234, datasource.builder.JSONData.(map[string]interface{})["maxSeries"])
}

func TestBasicAuth(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", BasicAuth("john", "doe"))

	req.NoError(err)
	req.Equal(boolPtr(true), datasource.builder.BasicAuth)
	req.Equal(stringPtr("john"), datasource.builder.BasicAuthUser)
	req.Equal(stringPtr("doe"), datasource.builder.BasicAuthPassword)
}

func TestWithCredentials(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", WithCredentials())

	req.NoError(err)
	req.Equal(true, datasource.builder.WithCredentials)
}

func TestSkipTLSVerify(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", SkipTLSVerify())

	req.NoError(err)
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"])
}

func TestForwardOauthIdentity(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", ForwardOauthIdentity())

	req.NoError(err)
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"])
}

func TestTLSClientAuth(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", TLSClientAuth("Foo", "bar"))

	req.NoError(err)
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsAuth"])
	req.Equal("Foo", datasource.builder.SecureJSONData.(map[string]interface{})["tlsClientCert"])
	req.Equal("bar", datasource.builder.SecureJSONData.(map[string]interface{})["tlsClientKey"])
}

func TestWithCACert(t *testing.T) {
	req := require.New(t)

	datasource, err := New("", "", WithCACert("bozo"))

	req.NoError(err)
	req.Equal(true, datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"])
	req.Equal("bozo", datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"])
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
