package polarion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) ListTestRuns(ctx context.Context, query string, limit int) ([]TestRun, error) {
	path := fmt.Sprintf("/projects/%s/testruns?pageSize=%d", c.project, limit)
	if query != "" {
		path += "&query=" + url.QueryEscape(query)
	}
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list test runs: %w", err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Title    string `json:"title"`
				Status   string `json:"status"`
				Template string `json:"templateId"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	runs := make([]TestRun, len(resp.Data))
	for i, d := range resp.Data {
		runs[i] = TestRun{
			ID:       d.ID,
			Title:    d.Attributes.Title,
			Status:   d.Attributes.Status,
			Template: d.Attributes.Template,
			URL:      fmt.Sprintf("https://%s/polarion/#/project/%s/testrun?id=%s", extractHost(c.baseURL), c.project, d.ID),
		}
	}
	return runs, nil
}

func (c *Client) GetTestRun(ctx context.Context, id string) (*TestRun, error) {
	path := fmt.Sprintf("/projects/%s/testruns/%s", c.project, id)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get test run %s: %w", id, err)
	}

	var resp struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Title    string `json:"title"`
				Status   string `json:"status"`
				Template string `json:"templateId"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &TestRun{
		ID:       resp.Data.ID,
		Title:    resp.Data.Attributes.Title,
		Status:   resp.Data.Attributes.Status,
		Template: resp.Data.Attributes.Template,
		URL:      fmt.Sprintf("https://%s/polarion/#/project/%s/testrun?id=%s", extractHost(c.baseURL), c.project, resp.Data.ID),
	}, nil
}

func (c *Client) CreateTestRun(ctx context.Context, in TestRunInput) (*TestRun, error) {
	body := map[string]any{
		"data": []map[string]any{{
			"type": "testruns",
			"attributes": map[string]any{
				"title":      in.Title,
				"templateId": in.Template,
			},
		}},
	}
	path := fmt.Sprintf("/projects/%s/testruns", c.project)
	data, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("create test run: %w", err)
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
	return c.GetTestRun(ctx, resp.Data[0].ID)
}

func (c *Client) UpdateTestRun(ctx context.Context, id string, in TestRunInput) (*TestRun, error) {
	attrs := map[string]any{}
	if in.Title != "" {
		attrs["title"] = in.Title
	}
	if in.Template != "" {
		attrs["templateId"] = in.Template
	}
	body := map[string]any{
		"data": map[string]any{
			"type":       "testruns",
			"id":         id,
			"attributes": attrs,
		},
	}
	path := fmt.Sprintf("/projects/%s/testruns/%s", c.project, id)
	_, err := c.makeRequest(ctx, "PATCH", path, body)
	if err != nil {
		return nil, fmt.Errorf("update test run %s: %w", id, err)
	}
	return c.GetTestRun(ctx, id)
}

func (c *Client) DeleteTestRun(ctx context.Context, id string) error {
	path := fmt.Sprintf("/projects/%s/testruns/%s", c.project, id)
	_, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("delete test run %s: %w", id, err)
	}
	return nil
}

func (c *Client) GetTestRunRecords(ctx context.Context, runID string) ([]TestRecord, error) {
	path := fmt.Sprintf("/projects/%s/testruns/%s/testrecords", c.project, runID)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get test run records %s: %w", runID, err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Result  string `json:"result"`
				Comment struct {
					Value string `json:"value"`
				} `json:"comment"`
			} `json:"attributes"`
			Relationships struct {
				TestCase struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"testCase"`
			} `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	records := make([]TestRecord, len(resp.Data))
	for i, d := range resp.Data {
		records[i] = TestRecord{
			CaseID:  d.Relationships.TestCase.Data.ID,
			Result:  d.Attributes.Result,
			Comment: d.Attributes.Comment.Value,
		}
	}
	return records, nil
}

func (c *Client) UpdateTestRunStatus(ctx context.Context, runID, status string) (*TestRun, error) {
	body := map[string]any{
		"data": map[string]any{
			"type": "testruns",
			"id":   runID,
			"attributes": map[string]any{
				"status": status,
			},
		},
	}
	path := fmt.Sprintf("/projects/%s/testruns/%s", c.project, runID)
	_, err := c.makeRequest(ctx, "PATCH", path, body)
	if err != nil {
		return nil, fmt.Errorf("update test run status %s: %w", runID, err)
	}
	return c.GetTestRun(ctx, runID)
}

func (c *Client) AddTestRecord(ctx context.Context, runID, caseID string, result TestResult) error {
	body := map[string]any{
		"data": []map[string]any{{
			"type": "testrecords",
			"attributes": map[string]any{
				"result":  result.Result,
				"comment": map[string]any{"type": "text/plain", "value": result.Comment},
			},
			"relationships": map[string]any{
				"testCase": map[string]any{
					"data": map[string]any{
						"type": "testcases",
						"id":   caseID,
					},
				},
			},
		}},
	}
	path := fmt.Sprintf("/projects/%s/testruns/%s/testrecords", c.project, runID)
	_, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return fmt.Errorf("add test record %s/%s: %w", runID, caseID, err)
	}
	return nil
}

func (c *Client) UpdateTestRunResult(ctx context.Context, runID, caseID string, result TestResult) error {
	body := map[string]any{
		"data": map[string]any{
			"type": "testrecords",
			"attributes": map[string]any{
				"result":  result.Result,
				"comment": map[string]any{"type": "text/plain", "value": result.Comment},
			},
		},
	}
	path := fmt.Sprintf("/projects/%s/testruns/%s/testrecords/%s/%s",
		c.project, runID, c.project, caseID)
	_, err := c.makeRequest(ctx, "PATCH", path, body)
	if err != nil {
		return fmt.Errorf("update test run result %s/%s: %w", runID, caseID, err)
	}
	return nil
}
