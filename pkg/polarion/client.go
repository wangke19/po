package polarion

import "net/http"

type Client struct {
	baseURL    string
	token      string
	project    string
	httpClient *http.Client
}

func NewClient(baseURL, token, project string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		project:    project,
		httpClient: httpClient,
	}
}
