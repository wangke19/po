package comment_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wangke19/po/pkg/cmd/comment"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/iostreams"
	"github.com/wangke19/po/pkg/polarion"
)

func newFactory(t *testing.T, srv *httptest.Server, stdin string) *cmdutil.Factory {
	t.Helper()
	client := polarion.NewClient(srv.URL, "test-token", "TEST", http.DefaultClient)
	var out bytes.Buffer
	ios := &iostreams.IOStreams{
		Out:    &out,
		ErrOut: &out,
		In:     io.NopCloser(strings.NewReader(stdin)),
	}
	return &cmdutil.Factory{
		IOStreams:       ios,
		PolarionClient: func() (*polarion.Client, error) { return client, nil },
	}
}

func outputOf(f *cmdutil.Factory) string {
	return f.IOStreams.Out.(*bytes.Buffer).String()
}

// listResponse matches ListComments response: author is nested under attributes.author.id
func listResponse(items []map[string]any) map[string]any {
	return map[string]any{"data": items}
}

func commentItem(id, authorID, created, text string) map[string]any {
	return map[string]any{
		"id": id,
		"attributes": map[string]any{
			"text":    text,
			"created": created,
			"author":  map[string]any{"id": authorID},
		},
	}
}

// addResponse matches AddComment response: same structure as list but wrapped in array
func addResponse(id, authorID, created, text string) map[string]any {
	return map[string]any{
		"data": []map[string]any{commentItem(id, authorID, created, text)},
	}
}

func TestListComments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(listResponse([]map[string]any{
			commentItem("CMT-1", "alice", "2026-01-01", "First"),
			commentItem("CMT-2", "bob", "2026-01-02", "Second"),
		}))
	}))
	defer srv.Close()

	f := newFactory(t, srv, "")
	cmd := comment.NewCmdList(f)
	cmd.SetArgs([]string{"WI-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if !strings.Contains(out, "alice") || !strings.Contains(out, "bob") {
		t.Errorf("expected both authors, got: %q", out)
	}
}

func TestAddComment_body(t *testing.T) {
	var gotText string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)
		data := req["data"].([]any)[0].(map[string]any)
		attrs := data["attributes"].(map[string]any)
		gotText = attrs["text"].(string)
		json.NewEncoder(w).Encode(addResponse("CMT-1", "jdoe", "2026-01-01", gotText))
	}))
	defer srv.Close()

	f := newFactory(t, srv, "")
	cmd := comment.NewCmdAdd(f)
	cmd.SetArgs([]string{"WI-1", "--body", "Hello world"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if gotText != "Hello world" {
		t.Errorf("expected %q, got %q", "Hello world", gotText)
	}
}

func TestAddComment_stdin(t *testing.T) {
	var gotText string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)
		data := req["data"].([]any)[0].(map[string]any)
		attrs := data["attributes"].(map[string]any)
		gotText = attrs["text"].(string)
		json.NewEncoder(w).Encode(addResponse("CMT-1", "jdoe", "2026-01-01", gotText))
	}))
	defer srv.Close()

	f := newFactory(t, srv, "from stdin\n")
	cmd := comment.NewCmdAdd(f)
	cmd.SetArgs([]string{"WI-1", "--body", "-"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if gotText != "from stdin" {
		t.Errorf("expected %q, got %q", "from stdin", gotText)
	}
}
