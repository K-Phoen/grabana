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

func TestDeleteAlertGroup(t *testing.T) {
	req := require.New(t)
	groupDeleted := false
	groupNamespace := "group namespace"
	groupName := "group name"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupDeleted = true

		req.Equal(http.MethodDelete, r.Method)
		req.Equal("/api/ruler/grafana/api/v1/rules/group namespace/group name", r.URL.Path)

		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteAlertGroup(context.TODO(), groupNamespace, groupName)

	req.NoError(err)
	req.True(groupDeleted)
}

func TestDeleteAlertGroupThatDoesNotExist(t *testing.T) {
	req := require.New(t)
	groupDeleted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupDeleted = true

		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteAlertGroup(context.TODO(), "ns", "name")

	req.Error(err)
	req.ErrorIs(err, ErrAlertNotFound)
	req.True(groupDeleted)
}
