package testrun_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wangke19/po/pkg/cmd/testrun"
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

// recordItem builds a test record matching the actual Polarion REST shape:
// - result and comment (nested value) are in attributes
// - testCase ID is in relationships
func recordItem(caseID, result, comment string) map[string]any {
	return map[string]any{
		"id": "TR-1/" + caseID,
		"attributes": map[string]any{
			"result":  result,
			"comment": map[string]any{"value": comment},
		},
		"relationships": map[string]any{
			"testCase": map[string]any{
				"data": map[string]any{"id": caseID},
			},
		},
	}
}

func makeRecordsHandler(records []map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"data": records})
	}
}

func TestRecords_allRecords(t *testing.T) {
	srv := httptest.NewServer(makeRecordsHandler([]map[string]any{
		recordItem("TC-1", "passed", "ok"),
		recordItem("TC-2", "failed", "bad"),
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := testrun.NewCmdRecords(f)
	cmd.SetArgs([]string{"TR-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if !strings.Contains(out, "TC-1") || !strings.Contains(out, "TC-2") {
		t.Errorf("expected both cases in output, got: %q", out)
	}
}

func TestRecords_filterByResult(t *testing.T) {
	srv := httptest.NewServer(makeRecordsHandler([]map[string]any{
		recordItem("TC-1", "passed", ""),
		recordItem("TC-2", "failed", ""),
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := testrun.NewCmdRecords(f)
	cmd.SetArgs([]string{"TR-1", "--result", "passed"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if !strings.Contains(out, "TC-1") {
		t.Errorf("expected TC-1 in output, got: %q", out)
	}
	if strings.Contains(out, "TC-2") {
		t.Errorf("TC-2 should be filtered out, got: %q", out)
	}
}

func TestRecords_filterByCase(t *testing.T) {
	srv := httptest.NewServer(makeRecordsHandler([]map[string]any{
		recordItem("TC-1", "passed", ""),
		recordItem("TC-2", "failed", ""),
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := testrun.NewCmdRecords(f)
	cmd.SetArgs([]string{"TR-1", "--case", "TC-2"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if strings.Contains(out, "TC-1") {
		t.Errorf("TC-1 should be filtered out, got: %q", out)
	}
	if !strings.Contains(out, "TC-2") {
		t.Errorf("expected TC-2 in output, got: %q", out)
	}
}

func TestRecords_notRun(t *testing.T) {
	srv := httptest.NewServer(makeRecordsHandler([]map[string]any{
		recordItem("TC-1", "passed", ""),
		recordItem("TC-2", "", ""),
		recordItem("TC-3", "", ""),
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := testrun.NewCmdRecords(f)
	cmd.SetArgs([]string{"TR-1", "--not-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	if strings.Contains(out, "TC-1") {
		t.Errorf("TC-1 (passed) should be filtered out, got: %q", out)
	}
	if !strings.Contains(out, "TC-2") || !strings.Contains(out, "TC-3") {
		t.Errorf("expected TC-2 and TC-3 in output, got: %q", out)
	}
}

func TestRecords_json(t *testing.T) {
	srv := httptest.NewServer(makeRecordsHandler([]map[string]any{
		recordItem("TC-1", "passed", "good"),
	}))
	defer srv.Close()

	f := newFactory(t, srv)
	cmd := testrun.NewCmdRecords(f)
	cmd.SetArgs([]string{"TR-1", "--json", ""})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := outputOf(f)
	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("output not valid JSON: %v\nout: %q", err, out)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result))
	}
}
