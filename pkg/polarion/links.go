package polarion

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) ListLinks(ctx context.Context, workItemID string) ([]WorkItemLink, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s/linkedworkitems", c.project, workItemID)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list links %s: %w", workItemID, err)
	}

	var resp struct {
		Data []struct {
			Attributes struct {
				Role string `json:"role"`
			} `json:"attributes"`
			Relationships struct {
				WorkItem struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"workItem"`
			} `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	links := make([]WorkItemLink, len(resp.Data))
	for i, d := range resp.Data {
		links[i] = WorkItemLink{
			TargetID: d.Relationships.WorkItem.Data.ID,
			Role:     d.Attributes.Role,
		}
	}
	return links, nil
}

func (c *Client) AddLink(ctx context.Context, workItemID, targetID, role string) error {
	body := map[string]any{
		"data": []map[string]any{{
			"type": "linkedworkitems",
			"attributes": map[string]any{
				"role": role,
			},
			"relationships": map[string]any{
				"workItem": map[string]any{
					"data": map[string]any{
						"type": "workitems",
						"id":   targetID,
					},
				},
			},
		}},
	}
	path := fmt.Sprintf("/projects/%s/workitems/%s/linkedworkitems", c.project, workItemID)
	_, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return fmt.Errorf("add link %s->%s: %w", workItemID, targetID, err)
	}
	return nil
}

func (c *Client) RemoveLink(ctx context.Context, workItemID, targetID, role string) error {
	path := fmt.Sprintf("/projects/%s/workitems/%s/linkedworkitems/%s/%s/%s",
		c.project, workItemID, c.project, targetID, role)
	_, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("remove link %s->%s: %w", workItemID, targetID, err)
	}
	return nil
}
