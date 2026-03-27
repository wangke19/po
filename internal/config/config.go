// Package config manages Polarion CLI configuration and authentication.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type hostEntry struct {
	DefaultProject string `yaml:"default_project"`
	VerifySSL      bool   `yaml:"verify_ssl"`
}

type configFile struct {
	Hosts map[string]hostEntry `yaml:"hosts"`
}

// Config holds Polarion CLI configuration.
type Config struct {
	path string
	data configFile
}

// New creates a new Config instance.
func New(path string) *Config {
	c := &Config{path: path}
	_ = c.load()
	return c
}

// DefaultConfigPath returns the default configuration file path.
func DefaultConfigPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "po", "config.yml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "po", "config.yml")
}

// NormalizeHostname normalizes a Polarion hostname.
func NormalizeHostname(host string) string {
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	return strings.TrimRight(host, "/")
}

func (c *Config) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &c.data)
}

func (c *Config) save() error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(&c.data)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}

// SetHost configures a Polarion host.
func (c *Config) SetHost(hostname, project string, verifySSL bool) error {
	hostname = NormalizeHostname(hostname)
	if c.data.Hosts == nil {
		c.data.Hosts = make(map[string]hostEntry)
	}
	c.data.Hosts[hostname] = hostEntry{DefaultProject: project, VerifySSL: verifySSL}
	return c.save()
}

// RemoveHost removes a configured host.
func (c *Config) RemoveHost(hostname string) error {
	hostname = NormalizeHostname(hostname)
	delete(c.data.Hosts, hostname)
	return c.save()
}

// DefaultHost returns the default Polarion host.
func (c *Config) DefaultHost() (string, error) {
	if v := os.Getenv("POLARION_URL"); v != "" {
		return NormalizeHostname(v), nil
	}
	for h := range c.data.Hosts {
		return h, nil
	}
	return "", fmt.Errorf("not logged in to any Polarion instance; run: po auth login")
}

// DefaultProject returns the default project for a host.
func (c *Config) DefaultProject(hostname string) (string, error) {
	if v := os.Getenv("POLARION_PROJECT"); v != "" {
		return v, nil
	}
	hostname = NormalizeHostname(hostname)
	if e, ok := c.data.Hosts[hostname]; ok {
		return e.DefaultProject, nil
	}
	return "", fmt.Errorf("no project configured for %s", hostname)
}

// VerifySSL returns whether SSL verification is enabled for a host.
func (c *Config) VerifySSL(hostname string) bool {
	if v := os.Getenv("POLARION_VERIFY_SSL"); v == "false" {
		return false
	}
	hostname = NormalizeHostname(hostname)
	if e, ok := c.data.Hosts[hostname]; ok {
		return e.VerifySSL
	}
	return true
}

// Hosts returns all configured hosts.
func (c *Config) Hosts() []string {
	hosts := make([]string, 0, len(c.data.Hosts))
	for h := range c.data.Hosts {
		hosts = append(hosts, h)
	}
	return hosts
}
