package grabana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/K-Phoen/grabana/datasource"
)

// ErrDatasourceNotFound is returned when the given datasource can not be found.
var ErrDatasourceNotFound = errors.New("datasource not found")

const defaultDatasourceKey = "$grabana_default_datasource_key$"

// UpsertDatasource creates or replaces a datasource.
func (client *Client) UpsertDatasource(ctx context.Context, datasource datasource.Datasource) error {
	buf, err := json.Marshal(datasource)
	if err != nil {
		return err
	}

	id, err := client.getDatasourceIDByName(ctx, datasource.Name())
	if err != nil && !errors.Is(err, ErrDatasourceNotFound) {
		return err
	}

	method := http.MethodPost
	url := "/api/datasources"
	if id != 0 {
		method = http.MethodPut
		url = fmt.Sprintf("/api/datasources/%d", id)
	}

	resp, err := client.sendJSON(ctx, method, url, buf)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return client.httpError(resp)
	}

	return nil
}

// DeleteDatasource deletes a datasource given its name.
func (client *Client) DeleteDatasource(ctx context.Context, name string) error {
	id, err := client.getDatasourceIDByName(ctx, name)
	if err != nil {
		return err
	}

	resp, err := client.delete(ctx, fmt.Sprintf("/api/datasources/%d", id))
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return ErrDatasourceNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return client.httpError(resp)
	}

	return nil
}

// GetDatasourceUIDByName finds a datasource UID given its name.
func (client *Client) GetDatasourceUIDByName(ctx context.Context, name string) (string, error) {
	resp, err := client.get(ctx, "/api/datasources/name/"+name)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrDatasourceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return "", client.httpError(resp)
	}

	response := struct {
		UID string `json:"uid"`
	}{}
	if err := decodeJSON(resp.Body, &response); err != nil {
		return "", err
	}

	return response.UID, nil
}

// datasourcesUIDMap builds a map of datasources UIDs indexed by their name.
func (client *Client) datasourcesUIDMap(ctx context.Context) (map[string]string, error) {
	resp, err := client.get(ctx, "/api/datasources")
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var datasources []struct {
		UID       string `json:"uid"`
		Name      string `json:"name"`
		IsDefault bool   `json:"isDefault"`
	}
	if err := decodeJSON(resp.Body, &datasources); err != nil {
		return nil, err
	}

	datasourcesMap := make(map[string]string, len(datasources))
	for _, ds := range datasources {
		datasourcesMap[ds.Name] = ds.UID

		if ds.IsDefault {
			datasourcesMap[defaultDatasourceKey] = ds.UID
		}
	}

	return datasourcesMap, nil
}

// getDatasourceIDByName finds a datasource, given its name.
func (client *Client) getDatasourceIDByName(ctx context.Context, name string) (int, error) {
	resp, err := client.get(ctx, "/api/datasources/id/"+name)
	if err != nil {
		return 0, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return 0, ErrDatasourceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return 0, client.httpError(resp)
	}

	response := struct {
		ID int `json:"id"`
	}{}
	if err := decodeJSON(resp.Body, &response); err != nil {
		return 0, err
	}

	return response.ID, nil
}
