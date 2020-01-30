package grabana

import (
	"github.com/grafana-tools/sdk"
)

type Dashboard struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	URL string `json:"url"`
}

type DashboardBuilder struct {
	board *sdk.Board
}
