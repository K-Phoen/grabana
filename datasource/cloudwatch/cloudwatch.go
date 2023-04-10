package cloudwatch

import (
	"encoding/json"

	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/sdk"
)

var _ datasource.Datasource = CloudWatch{}

type CloudWatch struct {
	builder *sdk.Datasource
}

type Option func(datasource *CloudWatch) error

func New(name string, options ...Option) (CloudWatch, error) {
	cloudwatch := &CloudWatch{
		builder: &sdk.Datasource{
			Name:           name,
			Type:           "cloudwatch",
			Access:         "proxy",
			JSONData:       map[string]interface{}{},
			SecureJSONData: map[string]interface{}{},
		},
	}

	defaults := []Option{
		DefaultAuth(),
	}

	for _, opt := range append(defaults, options...) {
		if err := opt(cloudwatch); err != nil {
			return *cloudwatch, err
		}
	}

	return *cloudwatch, nil
}

func (datasource CloudWatch) Name() string {
	return datasource.builder.Name
}

func (datasource CloudWatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
