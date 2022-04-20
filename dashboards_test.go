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

func TestADashboardCanBeFoundByTitle(t *testing.T) {
	req := require.New(t)
	dashboardName := "Test dashboard"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "eErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "eErXDvCkyy",
    "title": "Test dashboard"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dashboard, err := client.GetDashboardByTitle(context.TODO(), strings.ToLower(dashboardName))

	req.NoError(err)
	req.Equal(dashboardName, dashboard.Title)
}

func TestAnExplicitErrorIsReturnedIfTheDashboardIsNotFound(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "eErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "eErXDvCkyy",
    "title": "Test dashboard"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dashboard, err := client.GetDashboardByTitle(context.TODO(), "dashboard that do not exist")

	req.Error(err)
	req.Nil(dashboard)
	req.Equal(ErrDashboardNotFound, err)
}

func TestGetDashboardByTitleCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	dashboard, err := client.GetDashboardByTitle(context.TODO(), "does not matter")

	req.Error(err)
	req.Nil(dashboard)
}

func TestDashboardsCanBeCreated(t *testing.T) {
	req := require.New(t)
	dashboard := builder.New("Dashboard name")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `{
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
		_, _ = fmt.Fprintln(w, `{
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
		_, _ = fmt.Fprintln(w, `{"title": "Production Overview"}`)
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
		_, _ = fmt.Fprintln(w, `{"message": "oh noes"}`)
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
		_, _ = fmt.Fprintln(w, `{"message": "oh noes, does not exist"}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL)

	err := client.DeleteDashboard(context.TODO(), "some uid")

	req.Equal(ErrDashboardNotFound, err)
}
