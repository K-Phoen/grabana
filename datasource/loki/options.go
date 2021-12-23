package loki

import "time"

// Default configures this datasource to be the default one.
func Default() Option {
	return func(datasource *Loki) {
		datasource.builder.IsDefault = true
	}
}

// Timeout sets the timeout for HTTP requests.
func Timeout(timeout time.Duration) Option {
	return func(datasource *Loki) {
		datasource.builder.JSONData.(map[string]interface{})["timeout"] = int(timeout.Seconds())
	}
}

// BasicAuth configures basic authentication for this datasource.
func BasicAuth(username string, password string) Option {
	return func(datasource *Loki) {
		yep := true
		datasource.builder.BasicAuth = &yep
		datasource.builder.BasicAuthUser = &username
		datasource.builder.BasicAuthPassword = &password
	}
}

// SkipTLSVerify disables verification of SSL certificates.
func SkipTLSVerify() Option {
	return func(datasource *Loki) {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = true
	}
}

// WithCertificate sets a self-signed certificate that can be verified against.
func WithCertificate(certificate string) Option {
	return func(datasource *Loki) {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = false
		datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"] = true
		datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"] = certificate
	}
}

// WithCredentials joins credentials such as cookies or auth headers to cross-site requests.
func WithCredentials() Option {
	return func(datasource *Loki) {
		datasource.builder.WithCredentials = true
	}
}

// ForwardOauthIdentity forward the user's upstream OAuth identity to the data
// source (Their access token gets passed along).
func ForwardOauthIdentity() Option {
	return func(datasource *Loki) {
		datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"] = true
	}
}

// ForwardCookies configures a list of cookies that should be forwarded to the
// datasource.
func ForwardCookies(cookies ...string) Option {
	return func(datasource *Loki) {
		datasource.builder.JSONData.(map[string]interface{})["keepCookies"] = cookies
	}
}

// MaximumLines sets the maximum number of lines returned by Loki (default: 1000).
// Increase this value to have a bigger result set for ad-hoc analysis.
// Decrease this limit if your browser becomes sluggish when displaying the
// log results.
func MaximumLines(max int) Option {
	return func(datasource *Loki) {
		datasource.builder.JSONData.(map[string]interface{})["maxLines"] = max
	}
}
