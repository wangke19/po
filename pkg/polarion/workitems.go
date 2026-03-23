package polarion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) ListWorkItems(ctx context.Context, query string, limit int) ([]WorkItem, error) {
	path := fmt.Sprintf("/projects/%s/workitems?query=%s&pageSize=%d",
		c.project, url.QueryEscape(query), limit)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list work items: %w", err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Title       string `json:"title"`
				Type        string `json:"type"`
				Status      string `json:"status"`
				Author      string `json:"author"`
				Description string `json:"description"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	items := make([]WorkItem, len(resp.Data))
	for i, d := range resp.Data {
		items[i] = WorkItem{
			ID:          d.ID,
			Title:       d.Attributes.Title,
			Type:        d.Attributes.Type,
			Status:      d.Attributes.Status,
			Author:      d.Attributes.Author,
			Description: d.Attributes.Description,
			URL:         fmt.Sprintf("https://%s/polarion/#/project/%s/workitem?id=%s", extractHost(c.baseURL), c.project, d.ID),
		}
	}
	return items, nil
}

func (c *Client) GetWorkItem(ctx context.Context, id string) (*WorkItem, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s", c.project, id)
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
				Author      string `json:"author"`
				Description string `json:"description"`
			} `json:"attributes"`
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
		Author:      resp.Data.Attributes.Author,
		Description: resp.Data.Attributes.Description,
		URL:         fmt.Sprintf("https://%s/polarion/#/project/%s/workitem?id=%s", extractHost(c.baseURL), c.project, resp.Data.ID),
	}, nil
}

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

func (c *Client) UpdateWorkItem(ctx context.Context, id string, in WorkItemInput) (*WorkItem, error) {
	body := map[string]any{
		"data": map[string]any{
			"type": "workitems",
			"id":   id,
			"attributes": map[string]any{
				"title":  in.Title,
				"status": in.Status,
				"description": map[string]any{
					"type":  "text/html",
					"value": in.Description,
				},
			},
		},
	}
	path := fmt.Sprintf("/projects/%s/workitems/%s", c.project, id)
	_, err := c.makeRequest(ctx, "PATCH", path, body)
	if err != nil {
		return nil, fmt.Errorf("update work item %s: %w", id, err)
	}
	return c.GetWorkItem(ctx, id)
}

func (c *Client) DeleteWorkItem(ctx context.Context, id string) error {
	path := fmt.Sprintf("/projects/%s/workitems/%s", c.project, id)
	_, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("delete work item %s: %w", id, err)
	}
	return nil
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
