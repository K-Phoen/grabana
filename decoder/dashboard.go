package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
)

var ErrPanelNotConfigured = fmt.Errorf("panel not configured")
var ErrInvalidTimezone = fmt.Errorf("invalid timezone")

type DashboardModel struct {
	Title           string
	Editable        bool
	SharedCrosshair bool `yaml:"shared_crosshair"`
	Tags            []string
	AutoRefresh     string `yaml:"auto_refresh"`

	Time     [2]string
	Timezone string `yaml:",omitempty"`

	TagsAnnotation []dashboard.TagAnnotation `yaml:"tags_annotations,omitempty"`
	Variables      []DashboardVariable       `yaml:",omitempty"`
	ExternalLinks  []DashboardExternalLink   `yaml:"external_links,omitempty"`

	Rows []DashboardRow
}

func (d *DashboardModel) ToBuilder() (dashboard.Builder, error) {
	emptyDashboard := dashboard.Builder{}
	opts := []dashboard.Option{
		d.editable(),
		d.sharedCrossHair(),
	}

	if len(d.Tags) != 0 {
		opts = append(opts, dashboard.Tags(d.Tags))
	}

	for _, externalLink := range d.ExternalLinks {
		opts = append(opts, dashboard.ExternalLinks(externalLink.toModel()))
	}

	if d.AutoRefresh != "" {
		opts = append(opts, dashboard.AutoRefresh(d.AutoRefresh))
	}

	for _, tagAnnotation := range d.TagsAnnotation {
		opts = append(opts, dashboard.TagsAnnotation(tagAnnotation))
	}

	if d.Time[0] != "" && d.Time[1] != "" {
		opts = append(opts, dashboard.Time(d.Time[0], d.Time[1]))
	}

	switch d.Timezone {
	case "":
	case "default":
		opts = append(opts, dashboard.Timezone(dashboard.DefaultTimezone))
	case "utc":
		opts = append(opts, dashboard.Timezone(dashboard.UTC))
	case "browser":
		opts = append(opts, dashboard.Timezone(dashboard.Browser))
	default:
		return emptyDashboard, ErrInvalidTimezone
	}

	for _, variable := range d.Variables {
		opt, err := variable.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	for _, r := range d.Rows {
		opt, err := r.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	return dashboard.New(d.Title, opts...)
}

func (d *DashboardModel) sharedCrossHair() dashboard.Option {
	if d.SharedCrosshair {
		return dashboard.SharedCrossHair()
	}

	return dashboard.DefaultTooltip()
}

func (d *DashboardModel) editable() dashboard.Option {
	if d.Editable {
		return dashboard.Editable()
	}

	return dashboard.ReadOnly()
}

type DashboardPanel struct {
	Graph      *DashboardGraph      `yaml:",omitempty"`
	Table      *DashboardTable      `yaml:",omitempty"`
	SingleStat *DashboardSingleStat `yaml:"single_stat,omitempty"`
	Stat       *DashboardStat       `yaml:"stat,omitempty"`
	Text       *DashboardText       `yaml:",omitempty"`
	Heatmap    *DashboardHeatmap    `yaml:",omitempty"`
	TimeSeries *DashboardTimeSeries `yaml:"timeseries,omitempty"`
	Logs       *DashboardLogs       `yaml:"logs,omitempty"`
	Gauge      *DashboardGauge      `yaml:"gauge,omitempty"`
}

func (panel DashboardPanel) toOption() (row.Option, error) {
	if panel.Graph != nil {
		return panel.Graph.toOption()
	}
	if panel.TimeSeries != nil {
		return panel.TimeSeries.toOption()
	}
	if panel.Table != nil {
		return panel.Table.toOption()
	}
	if panel.SingleStat != nil {
		return panel.SingleStat.toOption()
	}
	if panel.Stat != nil {
		return panel.Stat.toOption()
	}
	if panel.Text != nil {
		return panel.Text.toOption(), nil
	}
	if panel.Heatmap != nil {
		return panel.Heatmap.toOption()
	}
	if panel.Logs != nil {
		return panel.Logs.toOption()
	}
	if panel.Gauge != nil {
		return panel.Gauge.toOption()
	}

	return nil, ErrPanelNotConfigured
}
