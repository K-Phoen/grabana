package decoder

import (
	"io"

	"github.com/K-Phoen/grabana/dashboard"
	"gopkg.in/yaml.v3"
)

func UnmarshalYAML(input io.Reader) (dashboard.Builder, error) {
	decoder := yaml.NewDecoder(input)

	parsed := &DashboardModel{}
	if err := decoder.Decode(parsed); err != nil {
		return dashboard.Builder{}, err
	}

	return parsed.toDashboardBuilder()
}
