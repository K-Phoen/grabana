package decoder

import (
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
)

type dashboardRow struct {
	Name   string
	Panels []dashboardPanel
}

func (r dashboardRow) toOption() (dashboard.Option, error) {
	opts := []row.Option{}

	for _, panel := range r.Panels {
		opt, err := panel.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return dashboard.Row(r.Name, opts...), nil
}
