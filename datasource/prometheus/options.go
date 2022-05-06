package prometheus

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
	return func(datasource *Prometheus) error {
		datasource.builder.IsDefault = true

		return nil
	}
}

// BasicAuth configures basic authentication for this datasource.
func BasicAuth(username string, password string) Option {
	return func(datasource *Prometheus) error {
		yep := true
		datasource.builder.BasicAuth = &yep
		datasource.builder.BasicAuthUser = &username
		datasource.builder.BasicAuthPassword = &password

		return nil
	}
}

// AccessMode controls how requests to the data source will be handled. Proxy
// should be the preferred way if nothing else is stated. Browser will let your
// browser send the requests (deprecated).
func AccessMode(mode Access) Option {
	return func(datasource *Prometheus) error {
		datasource.builder.Access = string(mode)

		return nil
	}
}

// HTTPMethod sets the method used to query Prometheus. POST is the recommended
// method as it allows bigger queries. Change this to GET if you have a
// Prometheus version older than 2.1 or if POST requests are restricted in your
// network.
func HTTPMethod(method string) Option {
	return func(datasource *Prometheus) error {
		normalizedMethod := strings.ToUpper(method)

		if normalizedMethod != "GET" && normalizedMethod != "POST" {
			return fmt.Errorf("HTTP method must be GET or POST: %w", errors.ErrInvalidArgument)
		}

		datasource.builder.JSONData.(map[string]interface{})["httpMethod"] = normalizedMethod

		return nil
	}
}

// ScrapeInterval configures the scrape and evaluation interval. Should be set
// to the typical value in Prometheus (defaults to 15s).
func ScrapeInterval(interval time.Duration) Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["timeInterval"] = interval.String()

		return nil
	}
}

// QueryTimeout sets the timeout for queries. Defaults to 60s
func QueryTimeout(timeout time.Duration) Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["queryTimeout"] = timeout.String()

		return nil
	}
}

// SkipTLSVerify disables verification of SSL certificates.
func SkipTLSVerify() Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = true

		return nil
	}
}

// WithCertificate sets a self-signed certificate that can be verified against.
func WithCertificate(certificate string) Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = false
		datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"] = true
		datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"] = certificate

		return nil
	}
}

// WithCredentials joins credentials such as cookies or auth headers to cross-site requests.
func WithCredentials() Option {
	return func(datasource *Prometheus) error {
		datasource.builder.WithCredentials = true

		return nil
	}
}

// ForwardOauthIdentity forward the user's upstream OAuth identity to the data
// source (Their access token gets passed along).
func ForwardOauthIdentity() Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"] = true

		return nil
	}
}

// ForwardCookies configures a list of cookies that should be forwarded to the
// datasource.
func ForwardCookies(cookies ...string) Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["keepCookies"] = cookies

		return nil
	}
}

// Exemplars configures a list of exemplars on this datasource.
func Exemplars(exemplars ...Exemplar) Option {
	return func(datasource *Prometheus) error {
		datasource.builder.JSONData.(map[string]interface{})["exemplarTraceIdDestinations"] = exemplars

		return nil
	}
}
