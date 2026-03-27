package polarion

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetProject retrieves a project by ID.
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	path := fmt.Sprintf("/projects/%s", id)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get project %s: %w", id, err)
	}

	var resp struct {
		Data struct {
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

	return &Project{
		ID:          resp.Data.ID,
		Name:        resp.Data.Attributes.Name,
		Description: resp.Data.Attributes.Description,
	}, nil
}

// ListProjects returns all accessible projects.
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
