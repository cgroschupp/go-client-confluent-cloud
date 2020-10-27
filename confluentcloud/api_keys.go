package confluentcloud

import (
	"fmt"
	"log"
	"net/url"
)

type LogicalCluster struct {
	ID   string  `json:"id"`
	Type *string `json:"type,omitempty"`
}

type APIKey struct {
	Key             string           `json:"key"`
	Secret          string           `json:"secret"`
	HashedSecret    string           `json:"hashed_secret"`
	HashedFunction  string           `json:"hashed_function"`
	SASLMechanism   string           `json:"sasl_mechanism"`
	UserID          int              `json:"user_id"`
	Deactived       bool             `json:"deactived"`
	ID              int              `json:"id"`
	Description     string           `json:"description"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
	AccountID       string           `json:"account_id"`
	ServiceAccount  bool             `json:"service_account"`
}

type APIKeysResponse struct {
	APIKeys []APIKey `json:"api_keys"`
}

type APIKeyResponse struct {
	APIKey APIKey `json:"api_key"`
}
type ApiKeyCreateRequestW struct {
	APIKey *ApiKeyCreateRequest `json:"api_key"`
}
type ApiKeyCreateRequest struct {
	AccountID       string           `json:"accountId"`
	UserID          int              `json:"user_id,omitempty"`
	Description     string           `json:"description,omitempty"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
}

func (c *Client) CreateAPIKey(request *ApiKeyCreateRequest) (*APIKey, error) {
	rel, err := url.Parse("api_keys")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(&ApiKeyCreateRequestW{APIKey: request}).
		SetResult(&APIKeyResponse{}).
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

func (c *Client) DeleteAPIKey(id, account_id string, logical_clusters []LogicalCluster) error {
	rel, err := url.Parse(fmt.Sprintf("api_keys/%s", id))
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(
			map[string]interface{}{
				"api_key": map[string]interface{}{
					"id":               id,
					"accountId":        account_id,
					"logical_clusters": logical_clusters,
				},
			},
		).
		SetError(&ErrorResponse{}).
		Delete(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("delete api key: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	log.Printf("[DEBUG] DeleteAPIKey Success(%s, %s)", id, account_id)
	return nil
}

func (c *Client) ListAPIKeys(clusterID, accountID string) ([]APIKey, error) {
	rel, err := url.Parse("api_keys")
	if err != nil {
		return []APIKey{}, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetQueryParam("account_id", accountID).
		SetQueryParam("cluster_id", clusterID).
		SetResult(&APIKeysResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return []APIKey{}, err
	}

	if response.IsError() {
		return []APIKey{}, fmt.Errorf("clusters: %s", response.Error().(*ErrorResponse).Error.Message)
	}
	return response.Result().(*APIKeysResponse).APIKeys, nil
}
