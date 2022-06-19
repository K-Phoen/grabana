package grabana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/sdk"
)

// ErrDashboardNotFound is returned when the given dashboard can not be found.
var ErrDashboardNotFound = errors.New("dashboard not found")

// Dashboard represents a Grafana dashboard.
type Dashboard struct {
	ID          int      `json:"id"`
	UID         string   `json:"uid"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	IsStarred   bool     `json:"isStarred"`
	FolderID    int      `json:"folderId"`
	FolderUID   string   `json:"folderUid"`
	FolderTitle string   `json:"folderTitle"`
	FolderURL   string   `json:"folderUrl"`
}

// GetDashboardByTitle finds a dashboard, given its title.
func (client *Client) GetDashboardByTitle(ctx context.Context, title string) (*Dashboard, error) {
	resp, err := client.get(ctx, "/api/search?type=dash-db&query="+url.QueryEscape(title))
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var dashboards []Dashboard
	if err := decodeJSON(resp.Body, &dashboards); err != nil {
		return nil, err
	}

	if len(dashboards) == 0 {
		return nil, ErrDashboardNotFound
	}

	for i := range dashboards {
		if strings.EqualFold(dashboards[i].Title, title) {
			return &dashboards[i], nil
		}
	}

	return nil, ErrDashboardNotFound
}

// GetDashboardByUID finds a dashboard, given its UID.
func (client *Client) rawDashboardByUID(ctx context.Context, uid string) (*sdk.Board, error) {
	resp, err := client.get(ctx, "/api/dashboards/uid/"+url.PathEscape(uid))
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrDashboardNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	response := struct {
		Board sdk.Board `json:"dashboard"`
	}{}

	if err := decodeJSON(resp.Body, &response); err != nil {
		return nil, err
	}

	return &response.Board, nil
}

// UpsertDashboard creates or replaces a dashboard, in the given folder.
func (client *Client) UpsertDashboard(ctx context.Context, folder *Folder, builder dashboard.Builder) (*Dashboard, error) {
	// first pass: save the new dashboard
	dashboardModel, err := client.persistDashboard(ctx, folder, builder)
	if err != nil {
		return nil, err
	}

	dashboardFromGrafana, err := client.rawDashboardByUID(ctx, dashboardModel.UID)
	if err != nil {
		return nil, err
	}

	// second pass: delete existing alerts associated to that dashboard
	alertRefs, err := client.listAlertsForDashboard(ctx, dashboardModel.UID)
	if err != nil {
		return nil, fmt.Errorf("could not prepare deletion of previous alerts for dashboard: %w", err)
	}
	for _, ref := range alertRefs {
		if err := client.DeleteAlertGroup(ctx, ref.Namespace, ref.RuleGroup); err != nil {
			return nil, fmt.Errorf("could not delete previous alerts for dashboard: %w", err)
		}
	}

	// third pass: create new alerts
	alerts := builder.Alerts()

	// If there are no alerts to create, we can return early
	if len(alerts) == 0 {
		return dashboardModel, nil
	}

	datasourcesMap, err := client.datasourcesUIDMap(ctx)
	if err != nil {
		return nil, err
	}

	for i := range alerts {
		alert := *alerts[i]

		alert.HookDashboardUID(dashboardFromGrafana.UID)
		alert.HookPanelID(panelIDByTitle(dashboardFromGrafana, alert.Builder.Name))

		if err := client.AddAlert(ctx, folder.Title, alert, datasourcesMap); err != nil {
			return nil, fmt.Errorf("could not add new alerts for dashboard: %w", err)
		}
	}

	return dashboardModel, nil
}

func (client *Client) persistDashboard(ctx context.Context, folder *Folder, builder dashboard.Builder) (*Dashboard, error) {
	buf, err := json.Marshal(struct {
		Dashboard *sdk.Board `json:"dashboard"`
		FolderID  uint       `json:"folderId"`
		Overwrite bool       `json:"overwrite"`
	}{
		Dashboard: builder.Internal(),
		FolderID:  folder.ID,
		Overwrite: true,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/dashboards/db", buf)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var model Dashboard
	if err := decodeJSON(resp.Body, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// DeleteDashboard deletes a dashboard given its UID.
func (client *Client) DeleteDashboard(ctx context.Context, uid string) error {
	// first: delete existing alerts associated to that dashboard
	alertRefs, err := client.listAlertsForDashboard(ctx, uid)
	if err != nil {
		return fmt.Errorf("could not prepare deletion of alerts for dashboard: %w", err)
	}
	for _, ref := range alertRefs {
		if err := client.DeleteAlertGroup(ctx, ref.Namespace, ref.RuleGroup); err != nil {
			return fmt.Errorf("could not delete alerts for dashboard: %w", err)
		}
	}

	// then: delete the dashboard itself
	resp, err := client.delete(ctx, "/api/dashboards/uid/"+uid)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return ErrDashboardNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return client.httpError(resp)
	}

	return nil
}

func panelIDByTitle(board *sdk.Board, title string) string {
	for _, row := range board.Rows {
		for _, panel := range row.Panels {
			if panel.Title == title {
				return fmt.Sprintf("%d", panel.ID)
			}
		}
	}

	for _, panel := range board.Panels {
		if panel.Title == title {
			return fmt.Sprintf("%d", panel.ID)
		}
	}

	return ""
}
