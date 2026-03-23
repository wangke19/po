package polarion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (c *Client) ListAttachments(ctx context.Context, workItemID string) ([]Attachment, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s/attachments", c.project, workItemID)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list attachments %s: %w", workItemID, err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				FileName    string `json:"fileName"`
				Title       string `json:"title"`
				ContentType string `json:"contentType"`
				Length      int64  `json:"length"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	attachments := make([]Attachment, len(resp.Data))
	for i, d := range resp.Data {
		attachments[i] = Attachment{
			ID:          d.ID,
			FileName:    d.Attributes.FileName,
			Title:       d.Attributes.Title,
			ContentType: d.Attributes.ContentType,
			Size:        d.Attributes.Length,
		}
	}
	return attachments, nil
}

func (c *Client) UploadAttachment(ctx context.Context, workItemID, fileName string, content io.Reader) (*Attachment, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s/attachments", c.project, workItemID)
	data, err := c.makeMultipartRequest(ctx, path, "file", fileName, content)
	if err != nil {
		return nil, fmt.Errorf("upload attachment %s: %w", workItemID, err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				FileName    string `json:"fileName"`
				Title       string `json:"title"`
				ContentType string `json:"contentType"`
				Length      int64  `json:"length"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no data in upload response")
	}
	d := resp.Data[0]
	return &Attachment{
		ID:          d.ID,
		FileName:    d.Attributes.FileName,
		Title:       d.Attributes.Title,
		ContentType: d.Attributes.ContentType,
		Size:        d.Attributes.Length,
	}, nil
}

func (c *Client) DownloadAttachment(ctx context.Context, workItemID, attachmentID string) (io.ReadCloser, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s/attachments/%s/content", c.project, workItemID, attachmentID)
	url := strings.TrimRight(c.baseURL, "/") + path

	req, err := newGetRequest(ctx, url, c.token)
	if err != nil {
		return nil, fmt.Errorf("download attachment %s: %w", attachmentID, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download attachment %s: %w", attachmentID, err)
	}
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("download attachment HTTP %d", resp.StatusCode)
	}
	return resp.Body, nil
}
