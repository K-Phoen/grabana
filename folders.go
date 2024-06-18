package grabana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ErrFolderNotFound is returned when the given folder can not be found.
var ErrFolderNotFound = errors.New("folder not found")

// Folder represents a dashboard folder.
// See https://grafana.com/docs/grafana/latest/reference/dashboard_folders/
type Folder struct {
	ID    uint   `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
}

// FindOrCreateFolder returns the folder by its name or creates it if it doesn't exist.
func (client *Client) FindOrCreateFolder(ctx context.Context, name string) (*Folder, error) {
	folder, err := client.GetFolderByTitle(ctx, name)
	if err != nil && !errors.Is(err, ErrFolderNotFound) {
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
	resp, err := client.get(ctx, fmt.Sprintf("/api/search?type=dash-folder&query=%s", url.QueryEscape(title)))
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
