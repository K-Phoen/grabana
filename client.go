package grabana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
			request.Header.Add("Authorization", "Bearer "+token)
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

func (client Client) httpError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
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
