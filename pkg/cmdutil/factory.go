package cmdutil

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/iostreams"
	"github.com/wangke19/po/pkg/polarion"
	"github.com/zalando/go-keyring"
)

type Factory struct {
	AppVersion     string
	IOStreams      *iostreams.IOStreams
	Config         func() (*config.Config, error)
	HttpClient     func() (*http.Client, error)
	PolarionClient func() (*polarion.Client, error)
}

func New(version string) *Factory {
	f := &Factory{
		AppVersion: version,
		IOStreams:  iostreams.System(),
	}

	f.Config = func() (*config.Config, error) {
		return config.New(config.DefaultConfigPath()), nil
	}

	f.HttpClient = func() (*http.Client, error) {
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}
		// host error is non-fatal here: VerifySSL defaults to true when host is empty
		host, _ := cfg.DefaultHost()
		verifySSL := cfg.VerifySSL(host)
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !verifySSL}
		return &http.Client{Transport: transport}, nil
	}

	f.PolarionClient = func() (*polarion.Client, error) {
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}
		host, err := cfg.DefaultHost()
		if err != nil {
			return nil, err
		}
		project, err := cfg.DefaultProject(host)
		if err != nil {
			return nil, err
		}

		token := os.Getenv("POLARION_TOKEN")
		if token == "" {
			var kerr error
			token, kerr = keyring.Get("po", host)
			if kerr != nil {
				return nil, fmt.Errorf("not authenticated for %s: run po auth login", host)
			}
		}

		httpClient, err := f.HttpClient()
		if err != nil {
			return nil, err
		}

		baseURL := fmt.Sprintf("https://%s/polarion/rest/v1", host)
		return polarion.NewClient(baseURL, token, project, httpClient), nil
	}

	return f
}
