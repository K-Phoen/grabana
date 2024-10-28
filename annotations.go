package grabana

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AnnotationResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

func (client *Client) AddAnnotation(ctx context.Context, text string, tags []string) (*AnnotationResponse, error) {
	buf, err := json.Marshal(struct {
		Time    int64    `json:"title"`
		Updated int64    `json:"updated"`
		Text    string   `json:"text"`
		Tags    []string `json:"tags"`
	}{
		Time:    time.Now().UnixMilli(),
		Updated: time.Now().UnixMilli(),
		Text:    text,
		Tags:    tags,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/annotations", buf)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var respann AnnotationResponse
	if err := decodeJSON(resp.Body, &respann); err != nil {
		return nil, err
	}

	return &respann, nil
}
