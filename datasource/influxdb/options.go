package influxdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/K-Phoen/grabana/errors"
)

type Access string

const (
	Proxy   Access = "proxy"
	Browser Access = "direct"
)

// Default configures this datasource to be the default one.
func Default() Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.IsDefault = true

		return nil
	}
}

// HTTPMethod defines the Method used to query the database (GET or POST HTTP verb).
// The POST verb allows heavy queries that would return an error using the GET verb.
// Default is GET.
func HTTPMethod(method string) Option {
	normalizedMethod := strings.ToUpper(method)

	if normalizedMethod != "GET" && normalizedMethod != "POST" {
		return invalidArgument(fmt.Errorf("HTTP method must be GET or POST: %w", errors.ErrInvalidArgument))
	}

	return setJSONData("httpMethod", normalizedMethod)
}

// AccessMode controls how requests to the data source will be handled. Proxy
// should be the preferred way if nothing else is stated. Browser will let your
// browser send the requests (deprecated).
func AccessMode(mode Access) Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.Access = string(mode)

		return nil
	}
}

// KeepCookies controls the cookies that will be forwarded to the data source.
// All other cookies will be deleted.
func KeepCookies(cookies []string) Option {
	return setJSONData("keepCookies", cookies)
}

// Timeout sets the timeout for HTTP requests.
func Timeout(timeout time.Duration) Option {
	return setJSONData("timeout", int(timeout.Seconds()))
}

// Database sets the ID of the bucket you want to query from,
// copied from the Buckets page of the InfluxDB UI.
func Database(database string) Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.Database = &database

		return nil
	}
}

// User sets username to use to sign into InfluxDB.
func User(user string) Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.User = &user

		return nil
	}
}

// Password sets token you use to query the selected bucked,
// copied from the Tokens page of the InfluxDB UI.
func Password(password string) Option {
	return setSecureJSONData("password", password)
}

// MinTimeInterval defines a lower limit for the auto group by time interval.
// Recommended to be set to write frequency, for example 1m if your data is written every minute.
func MinTimeInterval(interval time.Duration) Option {
	return setJSONData("timeInterval", interval.String())
}

// MaxSeries limits the number of series/tables that Grafana processes.
// Lower this number to prevent abuse, and increase it if you have lots of small time series
// and not all are shown. Defaults to 1000.
func MaxSeries(max int) Option {
	return setJSONData("maxSeries", max)
}

// BasicAuth configures basic authentication for this datasource.
func BasicAuth(username string, password string) Option {
	return func(datasource *InfluxQL) error {
		yep := true
		datasource.builder.BasicAuth = &yep
		datasource.builder.BasicAuthUser = &username
		datasource.builder.BasicAuthPassword = &password

		return nil
	}
}

// WithCredentials joins credentials such as cookies or auth headers to cross-site requests.
func WithCredentials() Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.WithCredentials = true

		return nil
	}
}

// SkipTLSVerify disables verification of SSL certificates.
func SkipTLSVerify() Option {
	return setJSONData("tlsSkipVerify", true)
}

// ForwardOauthIdentity forward the user's upstream OAuth identity to the datasource.
func ForwardOauthIdentity() Option {
	return setJSONData("oauthPassThru", true)
}

// TLSClientAuth enables TLS client side authentication. Expects PEM encoded content.
func TLSClientAuth(cert string, key string) Option {
	return multiOption(
		setJSONData("tlsAuth", true),
		setSecureJSONData("tlsClientCert", cert),
		setSecureJSONData("tlsClientKey", key),
	)
}

// WithCACert allows to provide a PEM encoded CA certificate to trust for this data source.
func WithCACert(cert string) Option {
	return multiOption(
		setJSONData("tlsAuthWithCACert", true),
		setSecureJSONData("tlsCACert", cert),
	)
}

func multiOption(opts ...Option) Option {
	return func(datasource *InfluxQL) error {
		for _, opt := range opts {
			if err := opt(datasource); err != nil {
				return err
			}
		}

		return nil
	}
}

func setJSONData(key string, value interface{}) Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.JSONData.(map[string]interface{})[key] = value

		return nil
	}
}

func setSecureJSONData(key string, value interface{}) Option {
	return func(datasource *InfluxQL) error {
		datasource.builder.SecureJSONData.(map[string]interface{})[key] = value

		return nil
	}
}

func invalidArgument(err error) Option {
	return func(datasource *InfluxQL) error {
		return err
	}
}
