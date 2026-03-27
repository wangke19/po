package api_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmd/api"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/iostreams"
)

func TestApiCmd_injectsAuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"data":"ok"}`))
	}))
	defer srv.Close()

	t.Setenv("POLARION_TOKEN", "test-bearer-token")
	t.Setenv("POLARION_PROJECT", "PROJ")
	t.Setenv("POLARION_URL", srv.URL)

	var out bytes.Buffer
	ios := &iostreams.IOStreams{Out: &out, ErrOut: &out}

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return config.New(t.TempDir() + "/config.yml"), nil
		},
		HTTPClient: func() (*http.Client, error) { return http.DefaultClient, nil },
	}

	cmd := api.NewCmdApi(f)
	cmd.SetArgs([]string{"/projects/{project}/workitems"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !strings.HasPrefix(gotAuth, "Bearer ") {
		t.Errorf("expected Bearer auth header, got: %q", gotAuth)
	}
	if gotAuth != "Bearer test-bearer-token" {
		t.Errorf("wrong token in auth header: %q", gotAuth)
	}
}

func TestApiCmd_projectSubstitution(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	t.Setenv("POLARION_TOKEN", "tok")
	t.Setenv("POLARION_PROJECT", "MY_PROJ")
	t.Setenv("POLARION_URL", srv.URL)

	var out bytes.Buffer
	ios := &iostreams.IOStreams{Out: &out, ErrOut: &out}

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return config.New(t.TempDir() + "/config.yml"), nil
		},
		HTTPClient: func() (*http.Client, error) { return http.DefaultClient, nil },
	}

	cmd := api.NewCmdApi(f)
	cmd.SetArgs([]string{"/projects/{project}/workitems"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !strings.Contains(gotPath, "MY_PROJ") {
		t.Errorf("expected MY_PROJ in path, got: %q", gotPath)
	}
}
