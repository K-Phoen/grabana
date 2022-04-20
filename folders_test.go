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

func TestFoldersCanBeCreated(t *testing.T) {
	req := require.New(t)
	folderName := "Test folder"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `{
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
		_, _ = fmt.Fprintln(w, `{
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
		_, _ = fmt.Fprintln(w, `[
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
		_, _ = fmt.Fprintln(w, `[
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
		_, _ = fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	folder, err := client.GetFolderByTitle(context.TODO(), "folder that do not exist")

	req.Error(err)
	req.Nil(folder)
}
