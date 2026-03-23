package polarion_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wangke19/po/pkg/polarion"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *polarion.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return polarion.NewClient(srv.URL, "test-token", "TEST", http.DefaultClient)
}

func TestListWorkItems(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing auth header")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "WI-1", "attributes": map[string]any{"title": "Test case 1", "type": "testcase", "status": "draft"}},
			},
		})
	})

	items, err := client.ListWorkItems(context.Background(), "type:testcase", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ID != "WI-1" {
		t.Errorf("got id %q", items[0].ID)
	}
}

func TestGetWorkItem_notFound(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	_, err := client.GetWorkItem(context.Background(), "WI-999")
	if err == nil {
		t.Error("expected error for 404")
	}
}

func TestCreateWorkItem(t *testing.T) {
	callCount := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// POST - return created ID
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{"id": "WI-2"}},
			})
		} else {
			// GET - return full item
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"id": "WI-2",
					"attributes": map[string]any{"title": "New case", "type": "testcase", "status": "draft"},
				},
			})
		}
	})

	item, err := client.CreateWorkItem(context.Background(), polarion.WorkItemInput{
		Title: "New case",
		Type:  "testcase",
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "WI-2" {
		t.Errorf("got id %q", item.ID)
	}
}

func TestCreateWorkItem_emptyResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	_, err := client.CreateWorkItem(context.Background(), polarion.WorkItemInput{Title: "X", Type: "testcase"})
	if err == nil {
		t.Error("expected error for empty response")
	}
}

func TestListTestRuns(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "TR-1", "attributes": map[string]any{"title": "Sprint 1 Run", "status": "inprogress"}},
			},
		})
	})
	runs, err := client.ListTestRuns(context.Background(), "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 || runs[0].ID != "TR-1" {
		t.Errorf("unexpected runs: %+v", runs)
	}
}
