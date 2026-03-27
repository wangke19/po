package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/zalando/go-keyring"
)

type apiOptions struct {
	endpoint  string
	method    string
	fields    []string
	headers   []string
	paginate  bool
	inputFile string
}

func NewCmdApi(f *cmdutil.Factory) *cobra.Command {
	opts := &apiOptions{}

	cmd := &cobra.Command{
		Use:   "api <endpoint>",
		Short: "Make an authenticated Polarion API request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.endpoint = args[0]
			return runApi(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.method, "method", "GET", "HTTP method: GET, POST, PATCH, DELETE")
	cmd.Flags().StringArrayVarP(&opts.fields, "raw-field", "f", nil, "key=value fields for request body")
	cmd.Flags().StringArrayVarP(&opts.headers, "header", "H", nil, "Additional request headers")
	cmd.Flags().BoolVar(&opts.paginate, "paginate", false, "Follow pagination (GET only)")
	cmd.Flags().StringVar(&opts.inputFile, "input", "", "File to use as request body")

	return cmd
}

func runApi(f *cmdutil.Factory, opts *apiOptions) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	host, err := cfg.DefaultHost()
	if err != nil {
		return err
	}

	project, err := cfg.DefaultProject(host)
	if err != nil && strings.Contains(opts.endpoint, "{project}") {
		return fmt.Errorf("no project configured: use POLARION_PROJECT or po auth login")
	}
	// Note: missing project is non-fatal when {project} not in endpoint

	token := os.Getenv("POLARION_TOKEN")
	if token == "" {
		token, err = keyring.Get("po", host)
		if err != nil {
			return fmt.Errorf("not authenticated: run po auth login")
		}
	}

	httpClient, err := f.HttpClient()
	if err != nil {
		return err
	}

	endpoint := strings.ReplaceAll(opts.endpoint, "{project}", project)
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	baseURL := buildBaseURL(host)

	if opts.paginate && opts.method != "GET" {
		return fmt.Errorf("--paginate is only valid with GET requests")
	}

	if opts.paginate {
		return runPaginated(context.Background(), f, httpClient, baseURL, endpoint, token, opts)
	}

	body, err := buildBody(opts)
	if err != nil {
		return err
	}

	respData, err := doRequest(context.Background(), httpClient, opts.method, baseURL+endpoint, token, opts.headers, body)
	if err != nil {
		return err
	}

	return printJSON(f, respData)
}

func buildBaseURL(host string) string {
	// If POLARION_URL is set, derive base from it preserving scheme
	if rawURL := os.Getenv("POLARION_URL"); rawURL != "" {
		// rawURL may be e.g. "http://127.0.0.1:8080" (no path) or "https://host"
		// Normalize: strip trailing slash, append REST path
		rawURL = strings.TrimRight(rawURL, "/")
		// If it already contains /polarion, use as-is base
		if strings.Contains(rawURL, "/polarion") {
			return rawURL
		}
		return rawURL + "/polarion/rest/v1"
	}
	return fmt.Sprintf("https://%s/polarion/rest/v1", host)
}

func buildBody(opts *apiOptions) (io.Reader, error) {
	if opts.inputFile != "" {
		f, err := os.Open(opts.inputFile)
		if err != nil {
			return nil, fmt.Errorf("open input file: %w", err)
		}
		defer func() { _ = f.Close() }()
		data, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("read input file: %w", err)
		}
		return bytes.NewReader(data), nil
	}
	if len(opts.fields) > 0 {
		m := make(map[string]string, len(opts.fields))
		for _, kv := range opts.fields {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid field %q: expected key=value", kv)
			}
			m[parts[0]] = parts[1]
		}
		data, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("marshal fields: %w", err)
		}
		return bytes.NewReader(data), nil
	}
	return nil, nil
}

func doRequest(ctx context.Context, client *http.Client, method, url, token string, extraHeaders []string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, h := range extraHeaders {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header %q: expected \"Key: Value\"", h)
		}
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func runPaginated(ctx context.Context, f *cmdutil.Factory, client *http.Client, baseURL, endpoint, token string, opts *apiOptions) error {
	var allData []json.RawMessage
	url := baseURL + endpoint

	for url != "" {
		data, err := doRequest(ctx, client, "GET", url, token, opts.headers, nil)
		if err != nil {
			return err
		}

		var envelope struct {
			Data  []json.RawMessage `json:"data"`
			Links *struct {
				Next string `json:"next"`
			} `json:"links"`
		}
		if err := json.Unmarshal(data, &envelope); err != nil {
			// Not a paginated envelope — print as-is and stop
			return printJSON(f, data)
		}

		allData = append(allData, envelope.Data...)
		if envelope.Links == nil || envelope.Links.Next == "" {
			break
		}
		url = envelope.Links.Next
	}

	combined, err := json.Marshal(allData)
	if err != nil {
		return fmt.Errorf("marshal paginated: %w", err)
	}
	return printJSON(f, combined)
}

func printJSON(f *cmdutil.Factory, data []byte) error {
	if f.IOStreams.IsTerminal() {
		var v any
		if err := json.Unmarshal(data, &v); err == nil {
			if pretty, err := json.MarshalIndent(v, "", "  "); err == nil {
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(pretty))
				return nil
			}
		}
	}
	_, _ = fmt.Fprintln(f.IOStreams.Out, string(data))
	return nil
}
