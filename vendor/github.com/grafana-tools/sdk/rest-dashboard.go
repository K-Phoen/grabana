package sdk

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DefaultFolderId is the id of the general folder
// that is pre-created and cannot be removed.
const DefaultFolderId = 0

// BoardProperties keeps metadata of a dashboard.
type BoardProperties struct {
	IsStarred  bool      `json:"isStarred,omitempty"`
	IsHome     bool      `json:"isHome,omitempty"`
	IsSnapshot bool      `json:"isSnapshot,omitempty"`
	Type       string    `json:"type,omitempty"`
	CanSave    bool      `json:"canSave"`
	CanEdit    bool      `json:"canEdit"`
	CanStar    bool      `json:"canStar"`
	Slug       string    `json:"slug"`
	Expires    time.Time `json:"expires"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	UpdatedBy  string    `json:"updatedBy"`
	CreatedBy  string    `json:"createdBy"`
	Version    int       `json:"version"`
}

// GetDashboardByUID loads a dashboard and its metadata from Grafana by dashboard uid.
//
// Reflects GET /api/dashboards/uid/:uid API call.
func (r *Client) GetDashboardByUID(ctx context.Context, uid string) (Board, BoardProperties, error) {
	return r.getDashboard(ctx, "uid/"+uid)
}

// GetDashboardBySlug loads a dashboard and its metadata from Grafana by dashboard slug.
//
// For dashboards from a filesystem set "file/" prefix for slug. By default dashboards from
// a database assumed. Database dashboards may have "db/" prefix or may have not, it will
// be appended automatically.
//
// Reflects GET /api/dashboards/db/:slug API call.
// Deprecated: since Grafana v5 you should use uids. Use GetDashboardByUID() for that.
func (r *Client) GetDashboardBySlug(ctx context.Context, slug string) (Board, BoardProperties, error) {
	path := setPrefix(slug)
	return r.getDashboard(ctx, path)
}

// getDashboard loads a dashboard from Grafana instance along with metadata for a dashboard.
// For dashboards from a filesystem set "file/" prefix for slug. By default dashboards from
// a database assumed. Database dashboards may have "db/" prefix or may have not, it will
// be appended automatically.
//
// Reflects GET /api/dashboards/db/:slug API call.
func (r *Client) getDashboard(ctx context.Context, path string) (Board, BoardProperties, error) {
	raw, bp, err := r.getRawDashboard(ctx, path)
	if err != nil {
		return Board{}, BoardProperties{}, errors.Wrap(err, "get raw dashboard")
	}
	var (
		result struct {
			Meta  BoardProperties `json:"meta"`
			Board Board           `json:"dashboard"`
		}
	)
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&result.Board); err != nil {
		return Board{}, BoardProperties{}, errors.Wrap(err, "unmarshal board")
	}
	return result.Board, bp, err
}

// GetRawDashboard loads a dashboard JSON from Grafana instance along with metadata for a dashboard.
// Contrary to GetDashboard() it not unpack loaded JSON to Board structure. Instead it
// returns it as byte slice. It guarantee that data of dashboard returned untouched by conversion
// with Board so no matter how properly fields from a current version of Grafana mapped to
// our Board fields. It useful for backuping purposes when you want a dashboard exactly with
// same data as it exported by Grafana.
//
// For dashboards from a filesystem set "file/" prefix for slug. By default dashboards from
// a database assumed. Database dashboards may have "db/" prefix or may have not, it will
// be appended automatically.
//
// Reflects GET /api/dashboards/db/:slug API call.
// Deprecated: since Grafana v5 you should use uids. Use GetRawDashboardByUID() for that.
func (r *Client) getRawDashboard(ctx context.Context, path string) ([]byte, BoardProperties, error) {
	var (
		raw    []byte
		result struct {
			Meta  BoardProperties `json:"meta"`
			Board json.RawMessage `json:"dashboard"`
		}
		code int
		err  error
	)
	if raw, code, err = r.get(ctx, fmt.Sprintf("api/dashboards/%s", path), nil); err != nil {
		return nil, BoardProperties{}, err
	}
	if code != 200 {
		return nil, BoardProperties{}, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&result); err != nil {
		return nil, BoardProperties{}, errors.Wrap(err, "unmarshal board")
	}
	return []byte(result.Board), result.Meta, err
}

// GetRawDashboardByUID loads a dashboard and its metadata from Grafana by dashboard uid.
//
// Reflects GET /api/dashboards/uid/:uid API call.
func (r *Client) GetRawDashboardByUID(ctx context.Context, uid string) ([]byte, BoardProperties, error) {
	return r.getRawDashboard(ctx, "uid/"+uid)
}

// GetRawDashboardBySlug loads a dashboard and its metadata from Grafana by dashboard slug.
//
// For dashboards from a filesystem set "file/" prefix for slug. By default dashboards from
// a database assumed. Database dashboards may have "db/" prefix or may have not, it will
// be appended automatically.
//
// Reflects GET /api/dashboards/db/:slug API call.
// Deprecated: since Grafana v5 you should use uids. Use GetRawDashboardByUID() for that.
func (r *Client) GetRawDashboardBySlug(ctx context.Context, slug string) ([]byte, BoardProperties, error) {
	path := setPrefix(slug)
	return r.getRawDashboard(ctx, path)
}

// FoundBoard keeps result of search with metadata of a dashboard.
type FoundBoard struct {
	ID        uint     `json:"id"`
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

// SearchDashboards search dashboards by substring of their title. It allows restrict the result set with
// only starred dashboards and only for tags (logical OR applied to multiple tags).
//
// Reflects GET /api/search API call.
func (r *Client) SearchDashboards(ctx context.Context, query string, starred bool, tags ...string) ([]FoundBoard, error) {
	var (
		raw    []byte
		boards []FoundBoard
		code   int
		err    error
	)
	u := url.URL{}
	q := u.Query()
	if query != "" {
		q.Set("query", query)
	}
	if starred {
		q.Set("starred", "true")
	}
	for _, tag := range tags {
		q.Add("tag", tag)
	}
	if raw, code, err = r.get(ctx, "api/search", q); err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	err = json.Unmarshal(raw, &boards)
	return boards, err
}

// SetDashboardParams contains the extra parameteres
// that affects where and how the dashboard will be stored
type SetDashboardParams struct {
	FolderID  int
	Overwrite bool
}

// SetDashboard updates existing dashboard or creates a new one.
// Set dasboard ID to nil to create a new dashboard.
// Set overwrite to true if you want to overwrite existing dashboard with
// newer version or with same dashboard title.
// Grafana only can create or update a dashboard in a database. File dashboards
// may be only loaded with HTTP API but not created or updated.
//
// Reflects POST /api/dashboards/db API call.
func (r *Client) SetDashboard(ctx context.Context, board Board, params SetDashboardParams) (StatusMessage, error) {
	var (
		isBoardFromDB bool
		newBoard      struct {
			Dashboard Board `json:"dashboard"`
			FolderID  int   `json:"folderId"`
			Overwrite bool  `json:"overwrite"`
		}
		raw  []byte
		resp StatusMessage
		code int
		err  error
	)
	if board.Slug, isBoardFromDB = cleanPrefix(board.Slug); !isBoardFromDB {
		return StatusMessage{}, errors.New("only database dashboard (with 'db/' prefix in a slug) can be set")
	}
	newBoard.Dashboard = board
	newBoard.FolderID = params.FolderID
	newBoard.Overwrite = params.Overwrite
	if !params.Overwrite {
		newBoard.Dashboard.ID = 0
	}
	if raw, err = json.Marshal(newBoard); err != nil {
		return StatusMessage{}, err
	}
	if raw, code, err = r.post(ctx, "api/dashboards/db", nil, raw); err != nil {
		return StatusMessage{}, err
	}
	if err = json.Unmarshal(raw, &resp); err != nil {
		return StatusMessage{}, err
	}
	if code != 200 {
		return resp, fmt.Errorf("HTTP error %d: returns %s", code, *resp.Message)
	}
	return resp, nil
}

// SetRawDashboard updates existing dashboard or creates a new one.
// Contrary to SetDashboard() it accepts raw JSON instead of Board structure.
// Grafana only can create or update a dashboard in a database. File dashboards
// may be only loaded with HTTP API but not created or updated.
//
// Reflects POST /api/dashboards/db API call.
func (r *Client) SetRawDashboard(ctx context.Context, raw []byte) (StatusMessage, error) {
	var (
		rawResp []byte
		resp    StatusMessage
		code    int
		err     error
		buf     bytes.Buffer
		plain   = make(map[string]interface{})
	)
	if err = json.Unmarshal(raw, &plain); err != nil {
		return StatusMessage{}, err
	}
	// TODO(axel) fragile place, refactor it
	plain["id"] = 0
	raw, _ = json.Marshal(plain)
	buf.WriteString(`{"dashboard":`)
	buf.Write(raw)
	buf.WriteString(`, "overwrite": true}`)
	if rawResp, code, err = r.post(ctx, "api/dashboards/db", nil, buf.Bytes()); err != nil {
		return StatusMessage{}, err
	}
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return StatusMessage{}, err
	}
	if code != 200 {
		return StatusMessage{}, fmt.Errorf("HTTP error %d: returns %s", code, *resp.Message)
	}
	return resp, nil
}

// DeleteDashboard deletes dashboard that selected by slug string.
// Grafana only can delete a dashboard in a database. File dashboards
// may be only loaded with HTTP API but not deteled.
//
// Reflects DELETE /api/dashboards/db/:slug API call.
func (r *Client) DeleteDashboard(ctx context.Context, slug string) (StatusMessage, error) {
	var (
		isBoardFromDB bool
		raw           []byte
		reply         StatusMessage
		err           error
	)
	if slug, isBoardFromDB = cleanPrefix(slug); !isBoardFromDB {
		return StatusMessage{}, errors.New("only database dashboards (with 'db/' prefix in a slug) can be removed")
	}
	if raw, _, err = r.delete(ctx, fmt.Sprintf("api/dashboards/db/%s", slug)); err != nil {
		return StatusMessage{}, err
	}
	err = json.Unmarshal(raw, &reply)
	return reply, err
}

// implicitly use dashboards from Grafana DB not from a file system
func setPrefix(slug string) string {
	if strings.HasPrefix(slug, "db") {
		return slug
	}
	if strings.HasPrefix(slug, "file") {
		return slug
	}
	return fmt.Sprintf("db/%s", slug)
}

// assume we use database dashboard by default
func cleanPrefix(slug string) (string, bool) {
	if strings.HasPrefix(slug, "db") {
		return slug[3:], true
	}
	if strings.HasPrefix(slug, "file") {
		return slug[3:], false
	}
	return fmt.Sprintf("%s", slug), true
}
