package confluentcloud

import (
	"fmt"
	"net/url"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthSuccessResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login() error {
	rel, err := url.Parse("sessions")
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetBody(AuthRequest{Email: c.email, Password: c.password}).
		SetResult(&AuthSuccessResponse{}).
		Post(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("login: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	if response.IsSuccess() {
		c.token = response.Result().(*AuthSuccessResponse).Token
	}
	return nil
}
