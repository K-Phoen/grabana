package grabana

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiTokenCanBeGiven(t *testing.T) {
	req := require.New(t)
	token := "foooo"

	serverCalled := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true

		req.Equal("Bearer "+token, r.Header.Get("Authorization"))
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, WithAPIToken(token))

	_, _ = client.CreateFolder(context.TODO(), "not relevant")

	req.True(serverCalled)
}

func TestBasicAuthCanBeSet(t *testing.T) {
	req := require.New(t)
	user := "joe"
	pass := "la frite"

	serverCalled := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true

		req.Equal("Basic am9lOmxhIGZyaXRl", r.Header.Get("Authorization"))
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, WithBasicAuth(user, pass))

	_, _ = client.CreateFolder(context.TODO(), "not relevant")

	req.True(serverCalled)
}
