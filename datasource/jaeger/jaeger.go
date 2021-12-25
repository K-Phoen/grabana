package jaeger

import (
	"encoding/json"

	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/sdk"
)

var _ datasource.Datasource = Jaeger{}

type Jaeger struct {
	builder *sdk.Datasource
}

type Option func(datasource *Jaeger)
type TraceToLogsOption func(settings map[string]interface{})

func New(name string, url string, options ...Option) Jaeger {
	jaeger := &Jaeger{
		builder: &sdk.Datasource{
			Name:           name,
			Type:           "jaeger",
			Access:         "proxy",
			URL:            url,
			JSONData:       map[string]interface{}{},
			SecureJSONData: map[string]interface{}{},
		},
	}

	for _, opt := range options {
		opt(jaeger)
	}

	return *jaeger
}

func (datasource Jaeger) Name() string {
	return datasource.builder.Name
}

func (datasource Jaeger) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
