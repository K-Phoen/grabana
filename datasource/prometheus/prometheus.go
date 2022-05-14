package prometheus

import (
	"encoding/json"
	"net/http"

	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/sdk"
)

var _ datasource.Datasource = Prometheus{}

type Exemplar struct {
	// The name of the field in the labels object that should be used to get the traceID.
	LabelName string `json:"name"`

	// The data source the exemplar is going to navigate to.
	// Set this value for internal exemplar links.
	DatasourceUID string `json:"datasourceUid"`

	// The URL of the trace backend the user would go to see its trace.
	// Set this value for external exemplar links.
	URL string `json:"url,omitempty"`
}

type Prometheus struct {
	builder *sdk.Datasource
}

type Option func(datasource *Prometheus) error

func New(name string, url string, options ...Option) (Prometheus, error) {
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
		if err := opt(prometheus); err != nil {
			return *prometheus, err
		}
	}

	return *prometheus, nil
}

func (datasource Prometheus) Name() string {
	return datasource.builder.Name
}

func (datasource Prometheus) MarshalJSON() ([]byte, error) {
	return json.Marshal(datasource.builder)
}
