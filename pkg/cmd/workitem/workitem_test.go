package workitem_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wangke19/po/pkg/cmd/workitem"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/iostreams"
	"github.com/wangke19/po/pkg/polarion"
)

func newFactory(t *testing.T, srv *httptest.Server) *cmdutil.Factory {
	t.Helper()
	client := polarion.NewClient(srv.URL, "test-token", "TEST", http.DefaultClient)
	var out bytes.Buffer
	return &cmdutil.Factory{
		IOStreams:      &iostreams.IOStreams{Out: &out, ErrOut: &out},
		PolarionClient: func() (*polarion.Client, error) { return client, nil },
	}
}

func outputOf(f *cmdutil.Factory) string {
	return f.IOStreams.Out.(*bytes.Buffer).String()
}

func TestListWorkItems_text(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "WI-1", "attributes": map[string]any{"title": "First", "type": "testcase", "status": "draft"}},
				{"id": "WI-2", "attributes": map[string]any{"title": "Second", "type": "testcase", "status": "approved"}},
			},
		})
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := workitem.NewCmdList(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if !strings.Contains(out, "WI-1") || !strings.Contains(out, "WI-2") {
		t.Errorf("expected both IDs in output, got: %q", out)
	}
}

func TestListWorkItems_json(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "WI-3", "attributes": map[string]any{"title": "JSON item", "type": "testcase", "status": "draft"}},
			},
		})
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := workitem.NewCmdList(f)
	cmd.SetArgs([]string{"--json", "id,title"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("output not valid JSON: %v\nout: %q", err, out)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0]["id"] != "WI-3" {
		t.Errorf("wrong id: %v", result[0]["id"])
	}
	if _, ok := result[0]["status"]; ok {
		t.Error("status field should be filtered out")
	}
}

func TestCreateWorkItem_withStatus(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// Verify status is sent in request body (data is an array)
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			data := body["data"].([]any)[0].(map[string]any)
			attrs := data["attributes"].(map[string]any)
			if attrs["status"] != "draft" {
				t.Errorf("expected status=draft in request, got: %v", attrs["status"])
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{"id": "WI-10"}},
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"id": "WI-10",
					"attributes": map[string]any{
						"title":       "My item",
						"type":        "testcase",
						"status":      "draft",
						"description": map[string]any{"value": ""},
					},
					"relationships": map[string]any{
						"author": map[string]any{"data": map[string]any{"id": ""}},
					},
					"links": map[string]any{"self": "http://example.com/WI-10"},
				},
			})
		}
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := workitem.NewCmdCreate(f)
	cmd.SetArgs([]string{"--type", "testcase", "--title", "My item", "--status", "draft"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if !strings.Contains(out, "WI-10") {
		t.Errorf("expected WI-10 in output, got: %q", out)
	}
}

func TestViewWorkItem_text(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "WI-5",
				"attributes": map[string]any{
					"title":       "View me",
					"type":        "testcase",
					"status":      "approved",
					"description": map[string]any{"value": "A description"},
				},
				"relationships": map[string]any{
					"author": map[string]any{
						"data": map[string]any{"id": "jdoe"},
					},
				},
				"links": map[string]any{"self": "http://example.com/WI-5"},
			},
		})
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := workitem.NewCmdView(f)
	cmd.SetArgs([]string{"WI-5"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	for _, want := range []string{"WI-5", "View me", "approved", "jdoe", "A description"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestDeleteWorkItem_requiresConfirm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := workitem.NewCmdDelete(f)
	cmd.SetArgs([]string{"WI-1"}) // no --confirm
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "confirm") {
		t.Errorf("expected confirm error, got: %v", err)
	}
}
