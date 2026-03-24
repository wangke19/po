package polarion

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (c *Client) GetTestSteps(ctx context.Context, caseID string) ([]TestStep, error) {
	caseID = stripProject(caseID)
	path := fmt.Sprintf("/projects/%s/workitems/%s/teststeps?fields%%5Bteststeps%%5D=keys,values,index", c.project, caseID)
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get test steps %s: %w", caseID, err)
	}

	var resp struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Keys   []string `json:"keys"`
				Values []struct {
					Value string `json:"value"`
				} `json:"values"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	steps := make([]TestStep, len(resp.Data))
	for i, d := range resp.Data {
		// extract step index from ID like "OSE/OCP-123/1"
		idx := 0
		if parts := strings.Split(d.ID, "/"); len(parts) > 0 {
			idx, _ = strconv.Atoi(parts[len(parts)-1])
		}
		// map keys to values: keys=["step","expectedResult"], values=[{value:...},{value:...}]
		var action, expectedResult string
		for j, k := range d.Attributes.Keys {
			if j >= len(d.Attributes.Values) {
				break
			}
			switch k {
			case "step":
				action = d.Attributes.Values[j].Value
			case "expectedResult":
				expectedResult = d.Attributes.Values[j].Value
			}
		}
		steps[i] = TestStep{
			StepIndex:      idx,
			Action:         action,
			ExpectedResult: expectedResult,
		}
	}
	return steps, nil
}

func (c *Client) DeleteTestStep(ctx context.Context, caseID string, stepIndex int) ([]TestStep, error) {
	caseID = stripProject(caseID)
	path := fmt.Sprintf("/projects/%s/workitems/%s/teststeps/%d", c.project, caseID, stepIndex)
	_, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, fmt.Errorf("delete test step %s/%d: %w", caseID, stepIndex, err)
	}
	return c.GetTestSteps(ctx, caseID)
}

func (c *Client) UpdateTestStep(ctx context.Context, caseID string, stepIndex int, in TestStepInput) ([]TestStep, error) {
	caseID = stripProject(caseID)
	attrs := map[string]any{}
	if in.Action != "" {
		attrs["action"] = in.Action
	}
	if in.ExpectedResult != "" {
		attrs["expectedResult"] = in.ExpectedResult
	}
	body := map[string]any{
		"data": map[string]any{
			"type":       "teststeps",
			"attributes": attrs,
		},
	}
	path := fmt.Sprintf("/projects/%s/workitems/%s/teststeps/%d", c.project, caseID, stepIndex)
	_, err := c.makeRequest(ctx, "PATCH", path, body)
	if err != nil {
		return nil, fmt.Errorf("update test step %s/%d: %w", caseID, stepIndex, err)
	}
	return c.GetTestSteps(ctx, caseID)
}

func (c *Client) AddTestStep(ctx context.Context, caseID string, in TestStepInput) ([]TestStep, error) {
	caseID = stripProject(caseID)
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
