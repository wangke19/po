package polarion

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetTestSteps(ctx context.Context, caseID string) ([]TestStep, error) {
	path := fmt.Sprintf("/projects/%s/workitems/%s/teststeps", c.project, caseID)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get test steps %s: %w", caseID, err)
	}

	var resp struct {
		Data []struct {
			Attributes struct {
				StepIndex      int    `json:"stepIndex"`
				Action         string `json:"action"`
				ExpectedResult string `json:"expectedResult"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	steps := make([]TestStep, len(resp.Data))
	for i, d := range resp.Data {
		steps[i] = TestStep{
			StepIndex:      d.Attributes.StepIndex,
			Action:         d.Attributes.Action,
			ExpectedResult: d.Attributes.ExpectedResult,
		}
	}
	return steps, nil
}

func (c *Client) AddTestStep(ctx context.Context, caseID string, in TestStepInput) ([]TestStep, error) {
	body := map[string]any{
		"data": []map[string]any{{
			"type": "teststeps",
			"attributes": map[string]any{
				"action":         in.Action,
				"expectedResult": in.ExpectedResult,
			},
		}},
	}
	path := fmt.Sprintf("/projects/%s/workitems/%s/teststeps", c.project, caseID)
	_, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("add test step %s: %w", caseID, err)
	}
	return c.GetTestSteps(ctx, caseID)
}
