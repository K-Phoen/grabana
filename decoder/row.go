package decoder

import (
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
)

// DashboardRow represents a dashboard row.
type DashboardRow struct {
	Name      string
	Repeat    string `yaml:"repeat_for,omitempty"`
	Collapse  bool   `yaml:",omitempty"`
	HideTitle bool   `yaml:"hide_title,omitempty"`
	Panels    []DashboardPanel
}

func (r DashboardRow) toOption() (dashboard.Option, error) {
	opts := []row.Option{}

	if r.Repeat != "" {
		opts = append(opts, row.RepeatFor(r.Repeat))
	}
	if r.Collapse {
		opts = append(opts, row.Collapse())
	}
	if r.HideTitle {
		opts = append(opts, row.HideTitle())
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
