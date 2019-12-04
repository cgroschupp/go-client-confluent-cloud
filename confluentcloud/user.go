package confluentcloud

import (
	"fmt"
	"net/url"
)

type AccountMessage struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
}

type UserInfoRequest struct {
	Account      AccountMessage `json:"account"`
	Organization interface{}    `json:"organization"`
	User         interface{}    `json:"user"`
}

func (c *Client) Me() (*UserInfoRequest, error) {
	rel, err := url.Parse("me")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetResult(&UserInfoRequest{}).
		Get(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("me: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return response.Result().(*UserInfoRequest), nil
}
