package grabana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/stretchr/testify/require"
)

func TestConfigureAlertManager(t *testing.T) {
	req := require.New(t)

	endpointCalled := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpointCalled = true
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)
	manager := alertmanager.New()

	err := client.ConfigureAlertManager(context.TODO(), manager)

	req.NoError(err)
	req.True(endpointCalled)
}

func TestConfigureAlertManagerForwardsErrorOnFailure(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, `{
  "message": "something when wrong"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)
	manager := alertmanager.New()

	err := client.ConfigureAlertManager(context.TODO(), manager)

	req.Error(err)
	req.Contains(err.Error(), "something when wrong")
}
