package polarion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ListWorkItems searches for work items matching a query.
func (c *Client) ListWorkItems(ctx context.Context, query string, limit int) ([]WorkItem, error) {
	path := fmt.Sprintf("/projects/%s/workitems?query=%s&page%%5Bsize%%5D=%d&fields%%5Bworkitems%%5D=title,type,status",
		c.project, url.QueryEscape(query), limit)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list work items: %w", err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Title  string `json:"title"`
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"attributes"`
			Relationships struct {
				Author struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"author"`
			} `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	items := make([]WorkItem, len(resp.Data))
	for i, d := range resp.Data {
		items[i] = WorkItem{
			ID:     d.ID,
			Title:  d.Attributes.Title,
			Type:   d.Attributes.Type,
			Status: d.Attributes.Status,
			Author: d.Relationships.Author.Data.ID,
			URL:    fmt.Sprintf("https://%s/polarion/#/project/%s/workitem?id=%s", extractHost(c.baseURL), c.project, stripProject(d.ID)),
		}
	}
	return items, nil
}

// GetWorkItem retrieves a single work item by ID.
func (c *Client) GetWorkItem(ctx context.Context, id string) (*WorkItem, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s?fields%%5Bworkitems%%5D=title,type,status,description,author",
		c.project, stripProject(id))
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get work item %s: %w", id, err)
	}

	var resp struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Title       string `json:"title"`
				Type        string `json:"type"`
				Status      string `json:"status"`
				Description struct {
					Value string `json:"value"`
				} `json:"description"`
			} `json:"attributes"`
			Relationships struct {
				Author struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"author"`
			} `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &WorkItem{
		ID:          resp.Data.ID,
		Title:       resp.Data.Attributes.Title,
		Type:        resp.Data.Attributes.Type,
		Status:      resp.Data.Attributes.Status,
		Author:      resp.Data.Relationships.Author.Data.ID,
		Description: resp.Data.Attributes.Description.Value,
		URL:         fmt.Sprintf("https://%s/polarion/#/project/%s/workitem?id=%s", extractHost(c.baseURL), c.project, stripProject(resp.Data.ID)),
	}, nil
}

// CreateWorkItem creates a new work item.
func (c *Client) CreateWorkItem(ctx context.Context, in WorkItemInput) (*WorkItem, error) {
	body := map[string]any{
		"data": []map[string]any{{
			"type": "workitems",
			"attributes": map[string]any{
				"title":  in.Title,
				"type":   in.Type,
				"status": in.Status,
				"description": map[string]any{
					"type":  "text/html",
					"value": in.Description,
				},
			},
		}},
	}
	path := fmt.Sprintf("/projects/%s/workitems", c.project)
	data, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("create work item: %w", err)
	}
	var resp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse create response: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no data in create response")
	}
	return c.GetWorkItem(ctx, resp.Data[0].ID)
}

// UpdateWorkItem modifies an existing work item.
func (c *Client) UpdateWorkItem(ctx context.Context, id string, in WorkItemInput) (*WorkItem, error) {
	id = stripProject(id)
	attrs := map[string]any{}
	if in.Title != "" {
		attrs["title"] = in.Title
	}
	if in.Type != "" {
		attrs["type"] = in.Type
	}
	if in.Status != "" {
		attrs["status"] = in.Status
	}
	if in.Description != "" {
		attrs["description"] = map[string]any{
			"type":  "text/html",
			"value": in.Description,
		}
	}
	body := map[string]any{
		"data": map[string]any{
			"type":       "workitems",
			"id":         id,
			"attributes": attrs,
		},
	}
	path := fmt.Sprintf("/projects/%s/workitems/%s", c.project, id)
	_, err := c.makeRequest(ctx, "PATCH", path, body)
	if err != nil {
		return nil, fmt.Errorf("update work item %s: %w", id, err)
	}
	return c.GetWorkItem(ctx, id)
}

// DeleteWorkItem permanently removes a work item.
func (c *Client) DeleteWorkItem(ctx context.Context, id string) error {
	path := fmt.Sprintf("/projects/%s/workitems/%s", c.project, stripProject(id))
	_, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("delete work item %s: %w", id, err)
	}
	return nil
}

// stripProject removes the "PROJECT/" prefix from an ID like "OSE/OCP-123" → "OCP-123".
func stripProject(id string) string {
	if i := strings.IndexByte(id, '/'); i >= 0 {
		return id[i+1:]
	}
	return id
}

// extractHost returns the host portion from a URL like https://host/polarion/rest/v1
func extractHost(baseURL string) string {
	after := strings.TrimPrefix(baseURL, "https://")
	after = strings.TrimPrefix(after, "http://")
	if i := strings.IndexByte(after, '/'); i >= 0 {
		return after[:i]
	}
	return after
}
