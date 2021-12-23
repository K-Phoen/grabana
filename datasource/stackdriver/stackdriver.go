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

type Option func(datasource *Stackdriver)

func New(name string, options ...Option) Stackdriver {
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
		opt(stackdriver)
	}

	return *stackdriver
}

func (datasource Stackdriver) Name() string {
	return datasource.builder.Name
}

func (datasource Stackdriver) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
