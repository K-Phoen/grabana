package prometheus

import (
	"strings"
	"time"
)

type Access string

const (
	Proxy   Access = "proxy"
	Browser Access = "direct"
)

// Default configures this datasource to be the default one.
func Default() Option {
	return func(datasource *Prometheus) {
		datasource.builder.IsDefault = true
	}
}

// BasicAuth configures basic authentication for this datasource.
func BasicAuth(username string, password string) Option {
	return func(datasource *Prometheus) {
		yep := true
		datasource.builder.BasicAuth = &yep
		datasource.builder.BasicAuthUser = &username
		datasource.builder.BasicAuthPassword = &password
	}
}

// AccessMode controls how requests to the data source will be handled. Proxy
// should be the preferred way if nothing else is stated. Browser will let your
// browser send the requests (deprecated).
func AccessMode(mode Access) Option {
	return func(datasource *Prometheus) {
		datasource.builder.Access = string(mode)
	}
}

// HTTPMethod sets the method used to query Prometheus. POST is the recommended
// method as it allows bigger queries. Change this to GET if you have a
// Prometheus version older than 2.1 or if POST requests are restricted in your
// network.
func HTTPMethod(method string) Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["httpMethod"] = strings.ToUpper(method)
	}
}

// ScrapeInterval configures the scrape and evaluation interval. Should be set
// to the typical value in Prometheus (defaults to 15s).
func ScrapeInterval(interval time.Duration) Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["timeInterval"] = interval.String()
	}
}

// QueryTimeout sets the timeout for queries. Defaults to 60s
func QueryTimeout(timeout time.Duration) Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["queryTimeout"] = timeout.String()
	}
}

// SkipTLSVerify disables verification of SSL certificates.
func SkipTLSVerify() Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = true
	}
}

// WithCertificate sets a self-signed certificate that can be verified against.
func WithCertificate(certificate string) Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = false
		datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"] = true
		datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"] = certificate
	}
}

// WithCredentials joins credentials such as cookies or auth headers to cross-site requests.
func WithCredentials() Option {
	return func(datasource *Prometheus) {
		datasource.builder.WithCredentials = true
	}
}

// ForwardOauthIdentity forward the user's upstream OAuth identity to the data
// source (Their access token gets passed along).
func ForwardOauthIdentity() Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"] = true
	}
}

// ForwardCookies configures a list of cookies that should be forwarded to the
// datasource.
func ForwardCookies(cookies ...string) Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["keepCookies"] = cookies
	}
}

// Exemplars configures a list of exemplars on this datasource.
func Exemplars(exemplars ...Exemplar) Option {
	return func(datasource *Prometheus) {
		datasource.builder.JSONData.(map[string]interface{})["exemplarTraceIdDestinations"] = exemplars
	}
}
