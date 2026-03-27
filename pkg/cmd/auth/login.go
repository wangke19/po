package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

type loginOptions struct {
	hostname  string
	project   string
	withToken bool
	insecure  bool
}

// NewCmdLogin returns the 'auth login' command.
func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	opts := &loginOptions{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a Polarion instance",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogin(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.hostname, "hostname", "", "Polarion instance URL (required)")
	cmd.Flags().StringVar(&opts.project, "project", "", "Default Polarion project ID (required)")
	cmd.Flags().BoolVar(&opts.withToken, "with-token", false, "Read token from stdin")
	cmd.Flags().BoolVar(&opts.insecure, "insecure", false, "Skip TLS certificate verification")
	_ = cmd.MarkFlagRequired("hostname")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}

func runLogin(f *cmdutil.Factory, opts *loginOptions) error {
	hostname := config.NormalizeHostname(opts.hostname)

	var token string
	if opts.withToken {
		data, err := io.ReadAll(f.IOStreams.In)
		if err != nil {
			return fmt.Errorf("reading token from stdin: %w", err)
		}
		token = strings.TrimSpace(string(data))
	} else {
		_, _ = fmt.Fprint(f.IOStreams.Out, "Token: ")
		raw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("reading token: %w", err)
		}
		_, _ = fmt.Fprintln(f.IOStreams.Out)
		token = strings.TrimSpace(string(raw))
	}

	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	if err := validateToken(hostname, opts.project, token, opts.insecure); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := keyring.Set("po", hostname, token); err != nil {
		return fmt.Errorf("storing token in keyring: %w", err)
	}

	cfg, err := f.Config()
	if err != nil {
		return err
	}
	if err := cfg.SetHost(hostname, opts.project, true); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	_, _ = fmt.Fprintf(f.IOStreams.Out, "Logged in to %s as project %s\n", hostname, opts.project)
	return nil
}

func validateToken(hostname, project, token string, insecure bool) error {
	url := fmt.Sprintf("https://%s/polarion/rest/v1/projects/%s", hostname, project)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid token (HTTP 401)")
	}
	if resp.StatusCode == 404 {
		return fmt.Errorf("project %q not found (HTTP 404)", project)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned HTTP %d", resp.StatusCode)
	}
	return nil
}
