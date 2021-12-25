package jaeger

import "time"

// Default configures this datasource to be the default one.
func Default() Option {
	return func(datasource *Jaeger) {
		datasource.builder.IsDefault = true
	}
}

// Timeout sets the timeout for HTTP requests.
func Timeout(timeout time.Duration) Option {
	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["timeout"] = int(timeout.Seconds())
	}
}

// BasicAuth configures basic authentication for this datasource.
func BasicAuth(username string, password string) Option {
	return func(datasource *Jaeger) {
		yep := true
		datasource.builder.BasicAuth = &yep
		datasource.builder.BasicAuthUser = &username
		datasource.builder.BasicAuthPassword = &password
	}
}

// SkipTLSVerify disables verification of SSL certificates.
func SkipTLSVerify() Option {
	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = true
	}
}

// WithCertificate sets a self-signed certificate that can be verified against.
func WithCertificate(certificate string) Option {
	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["tlsSkipVerify"] = false
		datasource.builder.JSONData.(map[string]interface{})["tlsAuthWithCACert"] = true
		datasource.builder.SecureJSONData.(map[string]interface{})["tlsCACert"] = certificate
	}
}

// WithCredentials joins credentials such as cookies or auth headers to cross-site requests.
func WithCredentials() Option {
	return func(datasource *Jaeger) {
		datasource.builder.WithCredentials = true
	}
}

// ForwardOauthIdentity forward the user's upstream OAuth identity to the data
// source (Their access token gets passed along).
func ForwardOauthIdentity() Option {
	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["oauthPassThru"] = true
	}
}

// ForwardCookies configures a list of cookies that should be forwarded to the
// datasource.
func ForwardCookies(cookies ...string) Option {
	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["keepCookies"] = cookies
	}
}

// WithNodeGraph enables the Node Graph visualization in the trace viewer.
func WithNodeGraph() Option {
	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["nodeGraph"] = map[string]interface{}{
			"enabled": true,
		}
	}
}

// TraceToLogs defines how to navigate from a trace span to the selected datasource logs.
func TraceToLogs(logsDatasourceUID string, options ...TraceToLogsOption) Option {
	settings := map[string]interface{}{
		"datasourceUid": logsDatasourceUID,
	}

	for _, opt := range options {
		opt(settings)
	}

	return func(datasource *Jaeger) {
		datasource.builder.JSONData.(map[string]interface{})["tracesToLogs"] = settings
	}
}

// Tags defines tags that will be used in the Loki query.
// Default tags: 'cluster', 'hostname', 'namespace', 'pod'.
func Tags(tags ...string) TraceToLogsOption {
	return func(settings map[string]interface{}) {
		settings["tags"] = tags
	}
}

// SpanStartShift shifts the start time of the span.
// Default 0 (Time units can be used here, for example: 5s, 1m, 3h)
func SpanStartShift(shift time.Duration) TraceToLogsOption {
	return func(settings map[string]interface{}) {
		settings["spanStartTimeShift"] = shift.String()
	}
}

// SpanEndShift shifts the start time of the span.
// Default 0 (Time units can be used here, for example: 5s, 1m, 3h)
func SpanEndShift(shift time.Duration) TraceToLogsOption {
	return func(settings map[string]interface{}) {
		settings["spanEndTimeShift"] = shift.String()
	}
}

// FilterByTrace filters logs by Trace ID. Appends '|=<trace id>' to the query.
func FilterByTrace() TraceToLogsOption {
	return func(settings map[string]interface{}) {
		settings["filterByTraceID"] = true
	}
}

// FilterBySpan filters logs by Trace ID. Appends '|=<trace id>' to the query.
func FilterBySpan() TraceToLogsOption {
	return func(settings map[string]interface{}) {
		settings["filterBySpanID"] = true
	}
}
