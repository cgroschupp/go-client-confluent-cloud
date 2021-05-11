package confluentcloud

import (
	"fmt"
	"net/url"
)

type Environment struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
	Deactivated    bool   `json:"deactivated"`
}

type EnvironmentResponse struct {
	Account Environment `json:"account"`
}

type EnvironmentsResponse struct { 
	Accounts []Environment `json:"accounts"`
}

type EnvironmentCreateRequest struct {
	Account EnvironmentRequest `json:"account"`
}

type EnvironmentRequest struct {
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
}

func (c *Client) GetEnvironment(id string) (*Environment, error) {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s", id))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetResult(&EnvironmentResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("get environment: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*EnvironmentResponse).Account, nil
}

func (c *Client) ListEnvironments() ([]Environment, error) { 
	rel, err := url.Parse("accounts")
	if err != nil {
		return []Environment{}, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
	SetResult(&EnvironmentsResponse{}).
	SetError(&ErrorResponse{}).
	Get(u.String())

	if err != nil {
		return []Environment{}, err
	}

	if response.IsError() {
		return []Environment{}, fmt.Errorf("get environments: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return response.Result().(*EnvironmentsResponse).Accounts, nil
}

func (c *Client) CreateEnvironment(name string, organizationID int) (*Environment, error) {
	rel, err := url.Parse("accounts")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	request := EnvironmentRequest{
		Name:           name,
		OrganizationID: organizationID,
	}

	response, err := c.NewRequest().
		SetBody(&EnvironmentCreateRequest{Account: request}).
		SetResult(&EnvironmentResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("environment: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*EnvironmentResponse).Account, nil
}

func (c *Client) DeleteEnvironment(id string) error {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s", id))
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetError(&ErrorResponse{}).
		Delete(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("delete environment: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return nil
}

func (c *Client) UpdateEnvironment(id, newName string, organizationID int) (*Environment, error) {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s", id))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	request := EnvironmentRequest{
		Name:           newName,
		OrganizationID: organizationID,
	}

	response, err := c.NewRequest().
		SetBody(&EnvironmentCreateRequest{Account: request}).
		SetResult(&EnvironmentResponse{}).
		SetError(&ErrorResponse{}).
		Put(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("environment: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*EnvironmentResponse).Account, nil
}
