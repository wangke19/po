package polarion

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	data, err := c.makeRequest(ctx, "GET", "/projects", nil)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	projects := make([]Project, len(resp.Data))
	for i, d := range resp.Data {
		projects[i] = Project{
			ID:          d.ID,
			Name:        d.Attributes.Name,
			Description: d.Attributes.Description,
		}
	}
	return projects, nil
}
