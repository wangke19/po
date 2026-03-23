package config_test

import (
	"path/filepath"
	"testing"

	"github.com/wangke19/po/internal/config"
)

func TestReadWriteConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := config.New(filepath.Join(dir, "config.yml"))

	if err := cfg.SetHost("example.com", "PROJ1", true); err != nil {
		t.Fatal(err)
	}

	cfg2 := config.New(filepath.Join(dir, "config.yml"))
	host, err := cfg2.DefaultHost()
	if err != nil {
		t.Fatal(err)
	}
	if host != "example.com" {
		t.Errorf("got %q, want %q", host, "example.com")
	}

	proj, err := cfg2.DefaultProject("example.com")
	if err != nil {
		t.Fatal(err)
	}
	if proj != "PROJ1" {
		t.Errorf("got %q, want %q", proj, "PROJ1")
	}
}

func TestNormalizeHostname(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://polarion.example.com/", "polarion.example.com"},
		{"http://polarion.example.com", "polarion.example.com"},
		{"polarion.example.com", "polarion.example.com"},
	}
	for _, tc := range cases {
		got := config.NormalizeHostname(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeHostname(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestEnvVarOverride(t *testing.T) {
	t.Setenv("POLARION_URL", "https://env.example.com")
	t.Setenv("POLARION_TOKEN", "tok123")
	t.Setenv("POLARION_PROJECT", "ENVPROJ")

	cfg := config.New(filepath.Join(t.TempDir(), "config.yml"))

	host, _ := cfg.DefaultHost()
	if host != "env.example.com" {
		t.Errorf("env override: got %q, want %q", host, "env.example.com")
	}

	proj, _ := cfg.DefaultProject(host)
	if proj != "ENVPROJ" {
		t.Errorf("env override project: got %q, want %q", proj, "ENVPROJ")
	}
}
