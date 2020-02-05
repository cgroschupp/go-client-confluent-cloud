package confluentcloud

import (
	"fmt"
	"net/url"
)

type ServiceAccount struct {
	ID          int    `json:"id"`
	Name        string `json:"service_name"`
	Description string `json:"service_description"`
}

type ServiceAccountsResponse struct {
	ServiceAccounts []ServiceAccount `json:"users"`
}

type ServiceAccountResponse struct {
	ServiceAccount ServiceAccount `json:"user"`
}
type ServiceAccountCreateRequestW struct {
	ServiceAccount *ServiceAccountCreateRequest `json:"user"`
}
type ServiceAccountCreateRequest struct {
	Name        string `json:"service_name"`
	Description string `json:"service_description"`
}
type ServiceAccountDeleteRequestW struct {
	ServiceAccount ServiceAccountDeleteRequest `json:"user"`
}
type ServiceAccountDeleteRequest struct {
	ID int `json:"id"`
}

func (c *Client) CreateServiceAccount(request *ServiceAccountCreateRequest) (*ServiceAccount, error) {
	rel, err := url.Parse("service_accounts")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(&ServiceAccountCreateRequestW{ServiceAccount: request}).
		SetResult(&ServiceAccountResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("service_accounts: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*ServiceAccountResponse).ServiceAccount, nil
}

func (c *Client) ListServiceAccounts(clusterID, accountID string) ([]ServiceAccount, error) {
	rel, err := url.Parse("service_accounts")
	if err != nil {
		return []ServiceAccount{}, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetQueryParam("account_id", accountID).
		SetQueryParam("cluster_id", clusterID).
		SetResult(&ServiceAccountsResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return []ServiceAccount{}, err
	}

	if response.IsError() {
		return []ServiceAccount{}, fmt.Errorf("service_accounts: %s", response.Error().(*ErrorResponse).Error.Message)
	}
	return response.Result().(*ServiceAccountsResponse).ServiceAccounts, nil
}

func (c *Client) DeleteServiceAccount(id int) error {
	rel, err := url.Parse("service_accounts")
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)

	request := ServiceAccountDeleteRequest{
		ID: id,
	}

	response, err := c.NewRequest().
		SetError(&ErrorResponse{}).
		SetBody(&ServiceAccountDeleteRequestW{ServiceAccount: request}).
		Delete(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("delete service account: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return nil
}
