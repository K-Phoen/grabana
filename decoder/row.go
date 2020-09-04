package decoder

import (
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
)

type DashboardRow struct {
	Name   string
	Repeat string `yaml:"repeat_for"`
	Panels []DashboardPanel
}

func (r DashboardRow) toOption() (dashboard.Option, error) {
	opts := []row.Option{}

	if r.Repeat != "" {
		opts = append(opts, row.RepeatFor(r.Repeat))
	}

	for _, panel := range r.Panels {
		opt, err := panel.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return dashboard.Row(r.Name, opts...), nil
}
