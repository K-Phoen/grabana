package grabana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// ErrAlertNotFound is returned when the requested alert can not be found.
var ErrAlertNotFound = errors.New("alert not found")

type alertRef struct {
	Namespace string
	RuleGroup string
}

// ConfigureAlertManager updates the alert manager configuration.
func (client *Client) ConfigureAlertManager(ctx context.Context, manager *alertmanager.Manager) error {
	buf, err := manager.MarshalIndentJSON()
	if err != nil {
		return err
	}

	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/alertmanager/grafana/config/api/v1/alerts", buf)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		return client.httpError(resp)
	}

	return nil
}

// AddAlert creates an alert group within a given namespace.
func (client *Client) AddAlert(ctx context.Context, namespace string, alertDefinition alert.Alert, datasourcesMap map[string]string) error {
	// Find out which datasource the alert depends on, and inject its UID into the sdk definition
	datasource := defaultDatasourceKey
	if alertDefinition.Datasource != "" {
		datasource = alertDefinition.Datasource
	}

	datasourceUID := datasourcesMap[datasource]
	if datasourceUID == "" {
		return fmt.Errorf("could not infer datasource UID from its name: %s", datasource)
	}

	alertDefinition.HookDatasourceUID(datasourceUID)

	// Before we can add this alert, we need to delete any other alert that might exist for this dashboard and panel
	if err := client.DeleteAlertGroup(ctx, namespace, alertDefinition.Builder.Name); err != nil && !errors.Is(err, ErrAlertNotFound) {
		return fmt.Errorf("could not delete existing alerts: %w", err)
	}

	buf, err := json.Marshal(alertDefinition.Builder)
	if err != nil {
		return err
	}

	// Save the alert!
	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/ruler/grafana/api/v1/rules/"+url.PathEscape(namespace), buf)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		return client.httpError(resp)
	}

	return nil
}

// DeleteAlertGroup deletes an alert group.
func (client *Client) DeleteAlertGroup(ctx context.Context, namespace string, groupName string) error {
	deleteURL := fmt.Sprintf("/api/ruler/grafana/api/v1/rules/%s/%s", url.PathEscape(namespace), url.PathEscape(groupName))
	resp, err := client.delete(ctx, deleteURL)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return ErrAlertNotFound
	}
	if resp.StatusCode != http.StatusAccepted {
		return client.httpError(resp)
	}

	return nil
}

// listAlertsForDashboard fetches a list of alerts linked to the given dashboard.
func (client *Client) listAlertsForDashboard(ctx context.Context, dashboardUID string) ([]alertRef, error) {
	resp, err := client.get(ctx, "/api/ruler/grafana/api/v1/rules?dashboard_uid="+url.QueryEscape(dashboardUID))
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var alerts map[string][]sdk.Alert
	if err := decodeJSON(resp.Body, &alerts); err != nil {
		return nil, err
	}

	var refs []alertRef

	for namespace := range alerts {
		for _, a := range alerts[namespace] {
			refs = append(refs, alertRef{
				Namespace: namespace,
				RuleGroup: a.Name,
			})
		}
	}

	return refs, nil
}
