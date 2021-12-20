package stackdriver

import "encoding/json"

// Default configures this datasource to be the default one.
func Default() Option {
	return func(datasource *Stackdriver) {
		datasource.builder.IsDefault = true
	}
}

// GCEAuthentication uses GCE default Service Account to authenticate to Stackdriver API.
func GCEAuthentication() Option {
	return func(datasource *Stackdriver) {
		datasource.builder.JSONData.(map[string]interface{})["authenticationType"] = "gce"
	}
}

// JWTAuthentication uses the given ServiceAccount key file to authenticate to Stackdriver API.
func JWTAuthentication(jwt string) Option {
	return func(datasource *Stackdriver) {
		parsedJwt := struct {
			ClientEmail    string `json:"client_email"`
			DefaultProject string `json:"project_id"`
			TokenURI       string `json:"token_uri"`
			PrivateKey     string `json:"private_key"`
		}{}

		// FIXME: we have no way of handling errors in options :(
		if err := json.Unmarshal([]byte(jwt), &parsedJwt); err != nil {
			return
		}

		datasource.builder.JSONData.(map[string]interface{})["authenticationType"] = "jwt"
		datasource.builder.JSONData.(map[string]interface{})["clientEmail"] = parsedJwt.ClientEmail
		datasource.builder.JSONData.(map[string]interface{})["defaultProject"] = parsedJwt.DefaultProject
		datasource.builder.JSONData.(map[string]interface{})["tokenUri"] = parsedJwt.TokenURI
		datasource.builder.SecureJSONData.(map[string]interface{})["privateKey"] = parsedJwt.PrivateKey
	}
}
