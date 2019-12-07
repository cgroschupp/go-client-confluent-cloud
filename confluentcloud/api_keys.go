package confluentcloud

import (
	"fmt"
	"net/url"
)

type LogicalCluster struct {
	ID   string  `json:"id"`
	Type *string `json:"type"`
}

type APIKey struct {
	Key             string           `json:"key"`
	Secret          string           `json:"secret"`
	HashedSecret    string           `json:"hashed_secret"`
	HashedFunction  string           `json:"hashed_function"`
	SASLMechanism   string           `json:"sasl_mechanism"`
	UserID          string           `json:"user_id"`
	Deactived       bool             `json:"deactived"`
	ID              int              `json:"id"`
	Description     string           `json:"description"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
	AccountID       string           `json:"account_id"`
	ServiceAccount  bool             `json:"service_account"`
}

type APIKeyResponse struct {
	APIKey APIKey `json:"api_key"`
}

type ApiKeyCreateRequest struct {
	AccountID       string           `json:"accountId"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
}

func (c *Client) CreateAPIKey(request *ApiKeyCreateRequest) (*APIKey, error) {
	rel, err := url.Parse("api_keys")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(request).
		SetResult(&ClusterResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("api_keys: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*APIKeyResponse).APIKey, nil
}
