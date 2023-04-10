package cloudwatch

import (
	"strings"
)

// DefaultAuth relies on the AWS SDK default authentication to authenticate to the CloudWatch service.
func DefaultAuth() Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["authType"] = "default"

		return nil
	}
}

// AccessSecretAuth relies on an access and secret key to authenticate to the CloudWatch service.
func AccessSecretAuth(accessKey string, secretKey string) Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["authType"] = "keys"
		datasource.builder.SecureJSONData.(map[string]interface{})["accessKey"] = accessKey
		datasource.builder.SecureJSONData.(map[string]interface{})["secretKey"] = secretKey

		return nil
	}
}

// Default configures this datasource to be the default one.
func Default() Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.IsDefault = true

		return nil
	}
}

// DefaultRegion sets the default region to use.
// Example: eu-north-1.
func DefaultRegion(region string) Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["defaultRegion"] = region

		return nil
	}
}

// AssumeRoleARN specifies the ARN of a role to assume.
// Format: arn:aws:iam:*
func AssumeRoleARN(roleARN string) Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["assumeRoleArn"] = roleARN

		return nil
	}
}

// ExternalID specifies the external identifier of a role to assume in another account.
func ExternalID(externalID string) Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["externalId"] = externalID

		return nil
	}
}

// Endpoint specifies a custom endpoint for the CloudWatch service.
func Endpoint(endpoint string) Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["endpoint"] = endpoint

		return nil
	}
}

// CustomMetricsNamespaces specifies a list of namespaces for custom metrics.
func CustomMetricsNamespaces(namespaces ...string) Option {
	return func(datasource *CloudWatch) error {
		datasource.builder.JSONData.(map[string]interface{})["customMetricsNamespaces"] = strings.Join(namespaces, ",")

		return nil
	}
}
