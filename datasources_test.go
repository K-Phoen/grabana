package grabana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/K-Phoen/grabana/datasource/prometheus"
	"github.com/stretchr/testify/require"
)

func TestDatasourceUpsertCanCreateANewDatasource(t *testing.T) {
	req := require.New(t)
	datasourcePosted := false

	datasource, err := prometheus.New("name", "address")
	req.NoError(err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		datasourcePosted = true

		req.Equal(http.MethodPost, r.Method)
		req.Equal("/api/datasources", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err = client.UpsertDatasource(context.TODO(), datasource)

	req.NoError(err)
	req.True(datasourcePosted)
}

func TestDatasourceUpsertCanUpdateADatasource(t *testing.T) {
	req := require.New(t)
	datasourceUpdated := false

	datasource, err := prometheus.New("name", "address")
	req.NoError(err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"id": 2}`)
			return
		}

		datasourceUpdated = true

		req.Equal(http.MethodPut, r.Method)
		req.Equal("/api/datasources/2", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err = client.UpsertDatasource(context.TODO(), datasource)

	req.NoError(err)
	req.True(datasourceUpdated)
}

func TestUpsertDatasourceForwardsErrorsOnFailure(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, `{
  "message": "something when wrong"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	datasource, err := prometheus.New("name", "address")
	req.NoError(err)

	err = client.UpsertDatasource(context.TODO(), datasource)

	req.Error(err)
	req.Contains(err.Error(), "something when wrong")
}

func TestDeleteDatasource(t *testing.T) {
	req := require.New(t)
	datasourceDeleted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"id": 2}`)
			return
		}

		datasourceDeleted = true

		req.Equal(http.MethodDelete, r.Method)
		req.Equal("/api/datasources/2", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDatasource(context.TODO(), "test-ds")

	req.NoError(err)
	req.True(datasourceDeleted)
}

func TestDeleteDatasourceReturnsKnownErrorIfDatasourceDoesNotExist(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDatasource(context.TODO(), "test-ds")

	req.Error(err)
	req.Equal(ErrDatasourceNotFound, err)
}

func TestDeleteDatasourceForwardsErrorOnFailure(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"id": 2}`)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, `{
  "message": "something when wrong"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDatasource(context.TODO(), "test-ds")

	req.Error(err)
	req.Contains(err.Error(), "something when wrong")
}

func TestGetDatasourceUIDByName(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{
  "uid": "some-uid"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	uid, err := client.GetDatasourceUIDByName(context.TODO(), "some-prometheus")

	req.NoError(err)
	req.Equal("some-uid", uid)
}

func TestGetDatasourceUIDByNameReturnsASpecificErrorIfDatasourceIsNotFound(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	uid, err := client.GetDatasourceUIDByName(context.TODO(), "some-prometheus")

	req.Error(err)
	req.Equal(ErrDatasourceNotFound, err)
	req.Empty(uid)
}
