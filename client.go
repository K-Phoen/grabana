package grabana

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/sdk"
)

// ErrFolderNotFound is returned when the given folder can not be found.
var ErrFolderNotFound = errors.New("folder not found")

// ErrDashboardNotFound is returned when the given dashboard can not be found.
var ErrDashboardNotFound = errors.New("dashboard not found")

// ErrDatasourceNotFound is returned when the given datasource can not be found.
var ErrDatasourceNotFound = errors.New("datasource not found")

// ErrAPIKeyNotFound is returned when the given API key can not be found.
var ErrAPIKeyNotFound = errors.New("API key not found")

// ErrAlertChannelNotFound is returned when the given alert notification
// channel can not be found.
var ErrAlertChannelNotFound = errors.New("alert channel not found")

// APIKeyRole represents a role given to an API key.
type APIKeyRole uint8

const (
	AdminRole APIKeyRole = iota
	EditorRole
	ViewerRole
)

func (role APIKeyRole) MarshalJSON() ([]byte, error) {
	var s string
	switch role {
	case ViewerRole:
		s = "Viewer"
	case EditorRole:
		s = "Editor"
	case AdminRole:
		s = "Admin"
	default:
		s = "None"
	}

	return json.Marshal(s)
}

// CreateAPIKeyRequest represents a request made to the API key creation endpoint.
type CreateAPIKeyRequest struct {
	Name          string     `json:"name"`
	Role          APIKeyRole `json:"role"`
	SecondsToLive int        `json:"secondsToLive"`
}

// APIKey represents an API key.
type APIKey struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Dashboard represents a Grafana dashboard.
type Dashboard struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	URL string `json:"url"`
}

// Folder represents a dashboard folder.
// See https://grafana.com/docs/grafana/latest/reference/dashboard_folders/
type Folder struct {
	ID    uint   `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
}

// Option represents an option that can be used to configure a client.
type Option func(client *Client)

type requestModifier func(request *http.Request)

// Client represents a Grafana HTTP client.
type Client struct {
	http             *http.Client
	host             string
	requestModifiers []requestModifier
}

// NewClient creates a new Grafana HTTP client, using an API token.
func NewClient(http *http.Client, host string, options ...Option) *Client {
	client := &Client{
		http: http,
		host: host,
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

// WithAPIToken sets up the client to use the given token to authenticate.
func WithAPIToken(token string) Option {
	return func(client *Client) {
		client.requestModifiers = append(client.requestModifiers, func(request *http.Request) {
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		})
	}
}

// WithBasicAuth sets up the client to use the given credentials to authenticate.
func WithBasicAuth(username string, password string) Option {
	return func(client *Client) {
		client.requestModifiers = append(client.requestModifiers, func(request *http.Request) {
			request.SetBasicAuth(username, password)
		})
	}
}

func (client *Client) modifyRequest(request *http.Request) {
	for _, modifier := range client.requestModifiers {
		modifier(request)
	}
}

// CreateAPIKey creates a new API key.
func (client *Client) CreateAPIKey(ctx context.Context, request CreateAPIKeyRequest) (string, error) {
	buf, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/auth/keys", buf)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", client.httpError(resp)
	}

	var response struct {
		Key string `json:"key"`
	}
	if err := decodeJSON(resp.Body, &response); err != nil {
		return "", err
	}

	return response.Key, nil
}

// DeleteAPIKeyByName deletes an API key given its name.
func (client *Client) DeleteAPIKeyByName(ctx context.Context, name string) error {
	apiKeys, err := client.APIKeys(ctx)
	if err != nil {
		return err
	}

	keyToDelete, ok := apiKeys[name]
	if !ok {
		return ErrAPIKeyNotFound
	}

	resp, err := client.delete(ctx, fmt.Sprintf("/api/auth/keys/%d", keyToDelete.ID))
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return ErrAPIKeyNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return client.httpError(resp)
	}

	return nil
}

// APIKeys lists active API keys.
func (client *Client) APIKeys(ctx context.Context) (map[string]APIKey, error) {
	resp, err := client.get(ctx, "/api/auth/keys")
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var keys []APIKey
	if err := decodeJSON(resp.Body, &keys); err != nil {
		return nil, err
	}

	keysMap := make(map[string]APIKey, len(keys))
	for _, key := range keys {
		keysMap[key.Name] = key
	}

	return keysMap, nil
}

// FindOrCreateFolder returns the folder by its name or creates it if it doesn't exist.
func (client *Client) FindOrCreateFolder(ctx context.Context, name string) (*Folder, error) {
	folder, err := client.GetFolderByTitle(ctx, name)
	if err != nil && err != ErrFolderNotFound {
		return nil, fmt.Errorf("could not find or create folder: %w", err)
	}
	if folder == nil {
		folder, err = client.CreateFolder(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("could not find create folder: %w", err)
		}
	}

	return folder, nil
}

// CreateFolder creates a dashboard folder.
// See https://grafana.com/docs/grafana/latest/reference/dashboard_folders/
func (client *Client) CreateFolder(ctx context.Context, name string) (*Folder, error) {
	buf, err := json.Marshal(struct {
		Title string `json:"title"`
	}{
		Title: name,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/folders", buf)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var folder Folder
	if err := decodeJSON(resp.Body, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// GetFolderByTitle finds a folder, given its title.
func (client *Client) GetFolderByTitle(ctx context.Context, title string) (*Folder, error) {
	resp, err := client.get(ctx, "/api/folders?limit=1000")
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var folders []Folder
	if err := decodeJSON(resp.Body, &folders); err != nil {
		return nil, err
	}

	for i := range folders {
		if strings.EqualFold(folders[i].Title, title) {
			return &folders[i], nil
		}
	}

	return nil, ErrFolderNotFound
}

// GetAlertChannelByName finds an alert notification channel, given its name.
func (client *Client) GetAlertChannelByName(ctx context.Context, name string) (*alert.Channel, error) {
	resp, err := client.get(ctx, "/api/alert-notifications")
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var channels []alert.Channel
	if err := decodeJSON(resp.Body, &channels); err != nil {
		return nil, err
	}

	for i := range channels {
		if strings.EqualFold(channels[i].Name, name) {
			return &channels[i], nil
		}
	}

	return nil, ErrAlertChannelNotFound
}

// UpsertDashboard creates or replaces a dashboard, in the given folder.
func (client *Client) UpsertDashboard(ctx context.Context, folder *Folder, builder dashboard.Builder) (*Dashboard, error) {
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

// UpsertDatasource creates or replaces a datasource.
func (client *Client) UpsertDatasource(ctx context.Context, datasource datasource.Datasource) error {
	buf, err := json.Marshal(datasource)
	if err != nil {
		return err
	}

	id, err := client.getDatasourceIDByName(ctx, datasource.Name())
	if err != nil && err != ErrDatasourceNotFound {
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

func (client Client) httpError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("could not query grafana: %s (HTTP status %d)", body, resp.StatusCode)
}

func (client Client) delete(ctx context.Context, path string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, client.url(path), nil)
	if err != nil {
		return nil, err
	}

	client.modifyRequest(request)

	return client.http.Do(request)
}

func (client Client) sendJSON(ctx context.Context, method string, path string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, client.url(path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	client.modifyRequest(request)

	return client.http.Do(request)
}

func (client Client) get(ctx context.Context, path string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.url(path), nil)
	if err != nil {
		return nil, err
	}

	client.modifyRequest(request)

	return client.http.Do(request)
}

func (client Client) url(path string) string {
	return client.host + path
}

func decodeJSON(input io.Reader, data interface{}) error {
	return json.NewDecoder(input).Decode(data)
}
