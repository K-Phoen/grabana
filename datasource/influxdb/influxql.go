package influxdb

import (
	"encoding/json"
	"net/http"

	"github.com/K-Phoen/sdk"
)

type InfluxQL struct {
	builder *sdk.Datasource
}

type Option func(datasource *InfluxQL) error

func New(name, url string, options ...Option) (InfluxQL, error) {
	datasource := InfluxQL{
		builder: &sdk.Datasource{
			Name:   name,
			Type:   "influxdb",
			Access: "proxy",
			URL:    url,
			JSONData: map[string]interface{}{
				"version": "InfluxQL",
			},
			SecureJSONData: map[string]interface{}{},
		},
	}

	defaultOptions := []Option{
		HTTPMethod(http.MethodGet),
		AccessMode(Proxy),
		MaxSeries(1000),
	}

	for _, opt := range append(defaultOptions, options...) {
		if err := opt(&datasource); err != nil {
			return datasource, err
		}
	}

	return datasource, nil
}

func (datasource InfluxQL) Name() string {
	return datasource.builder.Name
}

func (datasource InfluxQL) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
