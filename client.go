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

	"github.com/grafana-tools/sdk"
)

// ErrFolderNotFound is returned when the given folder can not be found.
var ErrFolderNotFound = errors.New("folder not found")

// ErrAlertChannelNotFound is returned when the given alert notification
// channel can not be found.
var ErrAlertChannelNotFound = errors.New("alert channel not found")

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

// Client represents a Grafana HTTP client.
type Client struct {
	http     *http.Client
	host     string
	apiToken string
}

// NewClient creates a new Grafana HTTP client.
func NewClient(http *http.Client, host string, apiToken string) *Client {
	return &Client{
		http:     http,
		host:     host,
		apiToken: apiToken,
	}
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

	resp, err := client.postJSON(ctx, "/api/folders", buf)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create folder: %s (HTTP status %d)", body, resp.StatusCode)
	}

	var folder Folder
	if err := decodeJSON(resp.Body, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// GetFolderByTitle finds a folder, given its title.
func (client *Client) GetFolderByTitle(ctx context.Context, title string) (*Folder, error) {
	resp, err := client.get(ctx, "/api/folders?limit=100")
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could list folders: %s (HTTP status %d)", body, resp.StatusCode)
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
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could lookup alert channels: %s (HTTP status %d)", body, resp.StatusCode)
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
func (client *Client) UpsertDashboard(ctx context.Context, folder *Folder, builder DashboardBuilder) (*Dashboard, error) {
	buf, err := json.Marshal(struct {
		Dashboard *sdk.Board `json:"dashboard"`
		FolderID  uint       `json:"folderId"`
		Overwrite bool       `json:"overwrite"`
	}{
		Dashboard: builder.board,
		FolderID:  folder.ID,
		Overwrite: true,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.postJSON(ctx, "/api/dashboards/db", buf)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create dashboard: %s", body)
	}

	var dashboard Dashboard
	if err := decodeJSON(resp.Body, &dashboard); err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (client Client) postJSON(ctx context.Context, path string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.url(path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.apiToken))

	return client.http.Do(request)
}

func (client Client) get(ctx context.Context, path string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.url(path), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.apiToken))

	return client.http.Do(request)
}

func (client Client) url(path string) string {
	return client.host + path
}

func decodeJSON(input io.Reader, data interface{}) error {
	return json.NewDecoder(input).Decode(data)
}
