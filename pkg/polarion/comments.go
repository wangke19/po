package polarion

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListComments returns all comments for a work item.
func (c *Client) ListComments(ctx context.Context, workItemID string) ([]Comment, error) {
	workItemID = stripProject(workItemID)
	path := fmt.Sprintf("/projects/%s/workitems/%s/comments?fields%%5Bworkitem_comments%%5D=title,author,created", c.project, workItemID)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list comments %s: %w", workItemID, err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Title   string `json:"title"`
				Created string `json:"created"`
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

	comments := make([]Comment, len(resp.Data))
	for i, d := range resp.Data {
		comments[i] = Comment{
			ID:      d.ID,
			Author:  d.Relationships.Author.Data.ID,
			Created: d.Attributes.Created,
			Body:    d.Attributes.Title,
		}
	}
	return comments, nil
}

// AddComment adds a new comment to a work item.
func (c *Client) AddComment(ctx context.Context, workItemID, body string) (*Comment, error) {
	workItemID = stripProject(workItemID)
	reqBody := map[string]any{
		"data": []map[string]any{{
			"type": "comments",
			"attributes": map[string]any{
				"text": body,
			},
		}},
	}
	path := fmt.Sprintf("/projects/%s/workitems/%s/comments", c.project, workItemID)
	data, err := c.makeRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("add comment %s: %w", workItemID, err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Text    string `json:"text"`
				Created string `json:"created"`
				Author  struct {
					ID string `json:"id"`
				} `json:"author"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no data in add comment response")
	}
	d := resp.Data[0]
	return &Comment{
		ID:      d.ID,
		Author:  d.Attributes.Author.ID,
		Created: d.Attributes.Created,
		Body:    d.Attributes.Text,
	}, nil
}
