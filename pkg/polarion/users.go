package polarion

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	data, err := c.makeRequest(ctx, "GET", "/users/current", nil)
	if err != nil {
		return nil, fmt.Errorf("get current user: %w", err)
	}

	var resp struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &User{
		ID:    resp.Data.ID,
		Name:  resp.Data.Attributes.Name,
		Email: resp.Data.Attributes.Email,
	}, nil
}
