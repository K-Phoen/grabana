package prometheus

import (
	"encoding/json"
	"net/http"

	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/sdk"
)

var _ datasource.Datasource = Prometheus{}

type Prometheus struct {
	builder *sdk.Datasource
}

type Option func(datasource *Prometheus)

func New(name string, url string, options ...Option) Prometheus {
	prometheus := &Prometheus{
		builder: &sdk.Datasource{
			Name:           name,
			Type:           "prometheus",
			Access:         "proxy",
			URL:            url,
			JSONData:       map[string]interface{}{},
			SecureJSONData: map[string]interface{}{},
		},
	}

	defaults := []Option{
		HTTPMethod(http.MethodPost),
		AccessMode(Proxy),
	}

	for _, opt := range append(defaults, options...) {
		opt(prometheus)
	}

	return *prometheus
}

func (datasource Prometheus) Name() string {
	return datasource.builder.Name
}

func (datasource Prometheus) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
