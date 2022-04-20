package grabana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAPIKey(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{"name":"mykey","key":"eyJrIjoiWHZiSWd3NzdCYUZnNUtibE9obUpESmE3bzJYNDRIc0UiLCJuIjoibXlrZXkiLCJpZCI6MX1=","id":1}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	token, err := client.CreateAPIKey(context.TODO(), CreateAPIKeyRequest{
		Name: "mykey",
		Role: AdminRole,
	})

	req.NoError(err)
	req.Equal("eyJrIjoiWHZiSWd3NzdCYUZnNUtibE9obUpESmE3bzJYNDRIc0UiLCJuIjoibXlrZXkiLCJpZCI6MX1=", token)
}

func TestDeleteDeleteAPIKeyByName(t *testing.T) {
	req := require.New(t)
	deleted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deletion call
		if strings.HasPrefix(r.URL.Path, "/api/auth/keys/") {
			deleted = true
			req.Equal(http.MethodDelete, r.Method)
			req.Equal("/api/auth/keys/2", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{}`)
			return
		}

		// API keys list call
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `[{"id": 2, "name": "foo"}]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteAPIKeyByName(context.TODO(), "foo")

	req.NoError(err)
	req.True(deleted)
}

func TestDeleteAPIKeyByNameReturnsKnownErrorIfDatasourceDoesNotExist(t *testing.T) {
	req := require.New(t)

	deleted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deletion call
		if strings.HasPrefix(r.URL.Path, "/api/auth/keys/") {
			deleted = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{}`)
			return
		}

		// API keys list call
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `[{"id": 2, "name": "foo"}]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteAPIKeyByName(context.TODO(), "unknown")

	req.Error(err)
	req.Equal(ErrAPIKeyNotFound, err)
	req.False(deleted)
}

func TestDeleteAPIKeyByNameForwardsErrorOnFailure(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deletion call
		if strings.HasPrefix(r.URL.Path, "/api/auth/keys/") {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintln(w, `{
  "message": "something when wrong"
}`)
			return
		}

		// API keys list call
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `[{"id": 2, "name": "foo"}]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteAPIKeyByName(context.TODO(), "foo")

	req.Error(err)
	req.Contains(err.Error(), "something when wrong")
}
