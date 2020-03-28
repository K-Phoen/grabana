package decoder

import (
	"io"

	builder "github.com/K-Phoen/grabana/dashboard"
	"gopkg.in/yaml.v2"
)

func UnmarshalYAML(input io.Reader) (builder.Builder, error) {
	decoder := yaml.NewDecoder(input)
	decoder.SetStrict(true)

	parsed := &DashboardModel{}
	if err := decoder.Decode(parsed); err != nil {
		return builder.Builder{}, err
	}

	return parsed.toDashboardBuilder()
}
