package grabana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrAPIKeyNotFound is returned when the given API key can not be found.
var ErrAPIKeyNotFound = errors.New("API key not found")

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
