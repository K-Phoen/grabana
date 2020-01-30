package grabana

import (
	"github.com/grafana-tools/sdk"
)

type Dashboard struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	URL string `json:"url"`
}

type DashboardBuilderOption func(dashboard *DashboardBuilder)

type DashboardBuilder struct {
	board *sdk.Board
}

func NewDashboardBuilder(title string, options ...DashboardBuilderOption) DashboardBuilder {
	board := sdk.NewBoard(title)
	board.ID = 0
	board.Timezone = ""

	builder := &DashboardBuilder{board: board}

	for _, opt := range dashboardDefaults() {
		opt(builder)
	}

	for _, opt := range options {
		opt(builder)
	}

	return *builder
}

func dashboardDefaults() []DashboardBuilderOption {
	return []DashboardBuilderOption{
		WithDefaultTimePicker(),
		WithDefaultTime(),
		WithSharedCrossHair(),
	}
}

func Editable() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Editable = true
	}
}

func ReadOnly() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Editable = false
	}
}

func WithRow(title string, options ...RowOption) DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		row := &Row{builder: builder.board.AddRow(title)}

		for _, opt := range rowDefaults() {
			opt(row)
		}

		for _, opt := range options {
			opt(row)
		}
	}
}

func WithDefaultTime() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Time = sdk.Time{
			From: "now-3h",
			To:   "now",
		}
	}
}

func WithDefaultTimePicker() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.Timepicker = sdk.Timepicker{
			RefreshIntervals: []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"},
			TimeOptions:      []string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"},
		}
	}
}

func WithSharedCrossHair() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.SharedCrosshair = true
	}
}

func WithoutSharedCrossHair() DashboardBuilderOption {
	return func(builder *DashboardBuilder) {
		builder.board.SharedCrosshair = false
	}
}
