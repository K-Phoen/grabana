package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
)

var ErrInvalidSortOrder = fmt.Errorf("invalid sort order")
var ErrInvalidDeduplicationStrategy = fmt.Errorf("invalid deduplication strategy")

type DashboardLogs struct {
	Title         string
	Description   string              `yaml:",omitempty"`
	Span          float32             `yaml:",omitempty"`
	Height        string              `yaml:",omitempty"`
	Transparent   bool                `yaml:",omitempty"`
	Datasource    string              `yaml:",omitempty"`
	Repeat        string              `yaml:",omitempty"`
	Links         DashboardPanelLinks `yaml:",omitempty"`
	Targets       []LogsTarget        `yaml:",omitempty"`
	Visualization *LogsVisualization  `yaml:",omitempty"`
}

type LogsTarget struct {
	Loki *LokiTarget `yaml:",omitempty"`
}

func (panel DashboardLogs) toOption() (row.Option, error) {
	opts := []logs.Option{}

	if panel.Description != "" {
		opts = append(opts, logs.Description(panel.Description))
	}
	if panel.Span != 0 {
		opts = append(opts, logs.Span(panel.Span))
	}
	if panel.Height != "" {
		opts = append(opts, logs.Height(panel.Height))
	}
	if panel.Transparent {
		opts = append(opts, logs.Transparent())
	}
	if panel.Datasource != "" {
		opts = append(opts, logs.DataSource(panel.Datasource))
	}
	if panel.Repeat != "" {
		opts = append(opts, logs.Repeat(panel.Repeat))
	}
	if len(panel.Links) != 0 {
		opts = append(opts, logs.Links(panel.Links.toModel()...))
	}
	for _, t := range panel.Targets {
		opt, err := panel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	vizOpts, err := panel.Visualization.toOptions()
	if err != nil {
		return nil, err
	}

	opts = append(opts, vizOpts...)

	return row.WithLogs(panel.Title, opts...), nil
}

func (panel DashboardLogs) target(t LogsTarget) (logs.Option, error) {
	if t.Loki != nil {
		return logs.WithLokiTarget(t.Loki.Query, t.Loki.toOptions()...), nil
	}

	return nil, ErrTargetNotConfigured
}

type LogsVisualization struct {
	Time           bool   `yaml:",omitempty"`
	UniqueLabels   bool   `yaml:"unique_labels,omitempty"`
	CommonLabels   bool   `yaml:"common_labels,omitempty"`
	WrapLines      bool   `yaml:"wrap_lines,omitempty"`
	PrettifyJSON   bool   `yaml:"prettify_json,omitempty"`
	HideLogDetails bool   `yaml:"hide_log_details,omitempty"`
	Order          string `yaml:",omitempty"`
	Deduplication  string `yaml:",omitempty"`
}

func (viz *LogsVisualization) toOptions() ([]logs.Option, error) {
	if viz == nil {
		return nil, nil
	}

	opts := []logs.Option{}

	if viz.Time {
		opts = append(opts, logs.Time())
	}
	if viz.UniqueLabels {
		opts = append(opts, logs.UniqueLabels())
	}
	if viz.CommonLabels {
		opts = append(opts, logs.CommonLabels())
	}
	if viz.WrapLines {
		opts = append(opts, logs.WrapLines())
	}
	if viz.PrettifyJSON {
		opts = append(opts, logs.PrettifyJSON())
	}
	if viz.HideLogDetails {
		opts = append(opts, logs.HideLogDetails())
	}

	if viz.Order != "" {
		switch viz.Order {
		case "asc":
			opts = append(opts, logs.Order(logs.Asc))
		case "desc":
			opts = append(opts, logs.Order(logs.Desc))
		default:
			return nil, ErrInvalidSortOrder
		}
	}

	if viz.Deduplication != "" {
		switch viz.Deduplication {
		case "none":
			opts = append(opts, logs.Deduplication(logs.None))
		case "exact":
			opts = append(opts, logs.Deduplication(logs.Exact))
		case "numbers":
			opts = append(opts, logs.Deduplication(logs.Numbers))
		case "signature":
			opts = append(opts, logs.Deduplication(logs.Signature))
		default:
			return nil, ErrInvalidDeduplicationStrategy
		}
	}

	return opts, nil
}
