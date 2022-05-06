package stackdriver

import (
	"encoding/json"

	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/sdk"
)

var _ datasource.Datasource = Stackdriver{}

type Stackdriver struct {
	builder *sdk.Datasource
}

type Option func(datasource *Stackdriver) error

func New(name string, options ...Option) (Stackdriver, error) {
	stackdriver := &Stackdriver{
		builder: &sdk.Datasource{
			Name:           name,
			Type:           "stackdriver",
			Access:         "proxy",
			JSONData:       map[string]interface{}{},
			SecureJSONData: map[string]interface{}{},
		},
	}

	defaults := []Option{
		GCEAuthentication(),
	}

	for _, opt := range append(defaults, options...) {
		if err := opt(stackdriver); err != nil {
			return *stackdriver, err
		}
	}

	return *stackdriver, nil
}

func (datasource Stackdriver) Name() string {
	return datasource.builder.Name
}

func (datasource Stackdriver) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
