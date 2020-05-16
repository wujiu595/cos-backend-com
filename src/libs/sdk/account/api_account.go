package account

import (
	"context"
)

func (c *Client) Me(ctx context.Context) (*UsersResult, error) {
	var res UsersResult
	err := c.client.Call(ctx, &res, "GET", c.config.Host+"/users/me")
	if err != nil {
		return nil, err
	}
	return &res, nil
}
