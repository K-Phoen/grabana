package grabana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	builder "github.com/K-Phoen/grabana/dashboard"
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
