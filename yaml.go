package grabana

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

func UnmarshalYAML(input io.Reader) (DashboardBuilder, error) {
	decoder := yaml.NewDecoder(input)
	decoder.SetStrict(true)

	parsed := &dashboardYaml{}
	if err := decoder.Decode(parsed); err != nil {
		return DashboardBuilder{}, err
	}

	fmt.Printf("parsed: %#v\n", parsed)

	return parsed.toDashboardBuilder()
}

type dashboardYaml struct {
	Title           string   `yaml:"title"`
	Editable        bool     `yaml:"editable"`
	SharedCrosshair bool     `yaml:"shared_crosshair"`
	Tags            []string `yaml:"tags"`
	AutoRefresh     string   `yaml:"autorefresh"`
}

func (dashboard *dashboardYaml) toDashboardBuilder() (DashboardBuilder, error) {
	opts := []DashboardBuilderOption{
		dashboard.sharedCrossHair(),
	}

	if len(dashboard.Tags) != 0 {
		opts = append(opts, Tags(dashboard.Tags))
	}

	if dashboard.AutoRefresh != "" {
		opts = append(opts, AutoRefresh(dashboard.AutoRefresh))
	}

	return NewDashboardBuilder(dashboard.Title, opts...), nil
}

func (dashboard *dashboardYaml) sharedCrossHair() DashboardBuilderOption {
	if dashboard.SharedCrosshair {
		return SharedCrossHair()
	}

	return DefaultTooltip()
}
