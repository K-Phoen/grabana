package grabana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	builder "github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/datasource/prometheus"
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

func TestFoldersCanBeCreated(t *testing.T) {
	req := require.New(t)
	folderName := "Test folder"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "uid": "nErXDvCkzz",
  "id": 1,
  "title": "Test folder"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.CreateFolder(context.TODO(), folderName)

	req.NoError(err)
	req.Equal(folderName, folder.Title)
}

func TestFoldersCreationCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{
  "message": "The folder has been changed by someone else",
  "status": "version-mismatch"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.CreateFolder(context.TODO(), "")

	req.Error(err)
	req.Nil(folder)
}

func TestAFolderCanBeFoundByTitle(t *testing.T) {
	req := require.New(t)
	folderName := "Test folder"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "nErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "nErXDvCkyy",
    "title": "Test folder"
  },
  {
    "id":3,
    "uid": "nErXDvCkxx",
    "title": "Department XYZ"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.GetFolderByTitle(context.TODO(), strings.ToLower(folderName))

	req.NoError(err)
	req.Equal(folderName, folder.Title)
}

func TestAnExplicitErrorIsReturnedIfTheFolderIsNotFound(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "nErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "nErXDvCkyy",
    "title": "Test folder"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.GetFolderByTitle(context.TODO(), "folder that do not exist")

	req.Error(err)
	req.Nil(folder)
	req.Equal(ErrFolderNotFound, err)
}

func TestGetFolderByTitleCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.GetFolderByTitle(context.TODO(), "folder that do not exist")

	req.Error(err)
	req.Nil(folder)
}

func TestAnAlertChannelCanBeFoundByName(t *testing.T) {
	req := require.New(t)
	name := "Team B"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
  {
    "id": 1,
    "uid": "team-a-email-notifier",
    "name": "Team A",
    "type": "email"
  },
  {
    "id": 2,
    "uid": "team-b-email-notifier",
    "name": "Team B",
    "type": "email"
  },
  {
    "id": 1,
    "uid": "team-c-email-notifier",
    "name": "Team C",
    "type": "email"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	channel, err := client.GetAlertChannelByName(context.TODO(), strings.ToLower(name))

	req.NoError(err)
	req.Equal(name, channel.Name)
	req.Equal("email", channel.Type)
}

func TestAnExplicitErrorIsReturnedIfTheChannelIsNotFound(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
   {
    "id": 1,
    "uid": "team-a-email-notifier",
    "name": "Team A",
    "type": "email"
  },
  {
    "id": 2,
    "uid": "team-b-email-notifier",
    "name": "Team B",
    "type": "email"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	channel, err := client.GetAlertChannelByName(context.TODO(), "channel that do not exist")

	req.Error(err)
	req.Nil(channel)
	req.Equal(ErrAlertChannelNotFound, err)
}

func TestGetAlertChannelByNameCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.GetAlertChannelByName(context.TODO(), "channel that do not exist")

	req.Error(err)
	req.Nil(folder)
}

func TestDashboardsCanBeCreated(t *testing.T) {
	req := require.New(t)
	dashboard := builder.New("Dashboard name")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "id":      1,
  "uid":     "cIBgcSjkk",
  "url":     "/d/cIBgcSjkk/production-overview",
  "status":  "success",
  "version": 1,
  "slug":    "production-overview"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	board, err := client.UpsertDashboard(context.TODO(), &Folder{}, dashboard)

	req.NoError(err)
	req.NotNil(board)
}

func TestDashboardsCreationCanFail(t *testing.T) {
	req := require.New(t)
	dashboard := builder.New("Dashboard name")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{
  "message": "The folder has been changed by someone else",
  "status": "version-mismatch"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	board, err := client.UpsertDashboard(context.TODO(), &Folder{}, dashboard)

	req.Error(err)
	req.Nil(board)
}

func TestDeleteDashboard(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"title": "Production Overview"}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some uid")

	req.NoError(err)
}

func TestDeleteDashboardCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{"message": "oh noes"}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some uid")

	req.Error(err)
}

func TestDeletingANonExistingDashboardReturnsSpecificError(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"message": "oh noes, does not exist"}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some uid")

	req.Equal(ErrDashboardNotFound, err)
}

func TestDatasourceUpsertCanCreateANewDatasource(t *testing.T) {
	req := require.New(t)
	datasourcePosted := false

	datasource := prometheus.New("name", "address")
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
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.UpsertDatasource(context.TODO(), datasource)

	req.NoError(err)
	req.True(datasourcePosted)
}

func TestDatasourceUpsertCanUpdateADatasource(t *testing.T) {
	req := require.New(t)
	datasourceUpdated := false

	datasource := prometheus.New("name", "address")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Datasource ID by name" call
		if strings.HasPrefix(r.URL.Path, "/api/datasources/id/") {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"id": 2}`)
			return
		}

		datasourceUpdated = true

		req.Equal(http.MethodPut, r.Method)
		req.Equal("/api/datasources/2", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.UpsertDatasource(context.TODO(), datasource)

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
		fmt.Fprintln(w, `{
  "message": "something when wrong"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.UpsertDatasource(context.TODO(), prometheus.New("name", "address"))

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
			fmt.Fprintln(w, `{"id": 2}`)
			return
		}

		datasourceDeleted = true

		req.Equal(http.MethodDelete, r.Method)
		req.Equal("/api/datasources/2", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
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
			fmt.Fprintln(w, `{"id": 2}`)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{
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
		fmt.Fprintln(w, `{
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

func TestCreateAPIKey(t *testing.T) {
	req := require.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"name":"mykey","key":"eyJrIjoiWHZiSWd3NzdCYUZnNUtibE9obUpESmE3bzJYNDRIc0UiLCJuIjoibXlrZXkiLCJpZCI6MX1=","id":1}`)
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
			fmt.Fprintln(w, `{}`)
			return
		}

		// API keys list call
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": 2, "name": "foo"}]`)
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
			fmt.Fprintln(w, `{}`)
			return
		}

		// API keys list call
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": 2, "name": "foo"}]`)
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
			fmt.Fprintln(w, `{
  "message": "something when wrong"
}`)
			return
		}

		// API keys list call
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": 2, "name": "foo"}]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteAPIKeyByName(context.TODO(), "foo")

	req.Error(err)
	req.Contains(err.Error(), "something when wrong")
}
