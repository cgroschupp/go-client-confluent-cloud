package confluentcloud

import (
	"fmt"
	"net/url"
	"time"
)

type SchemaRegistry struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	KafkaClusterID    string    `json:"kafka_cluster_id"`
	Endpoint          string    `json:"endpoint"`
	Created           time.Time `json:"created"`
	Modified          time.Time `json:"modified"`
	Status            string    `json:"status"`
	PhysicalClusterID string    `json:"physical_cluster_id"`
	AccountID         string    `json:"account_id"`
	OrganizationID    int       `json:"organization_id"`
	MaxSchemas        int       `json:"max_schemas"`
}

type SchemaRegistriesResponse struct {
	Error    interface{}      `json:"error"`
	Clusters []SchemaRegistry `json:"clusters"`
}

func (c *Client) ListSchemaRegistries(accountID string) ([]SchemaRegistry, error) {
	rel, err := url.Parse("schema_registries")
	if err != nil {
		return []SchemaRegistry{}, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetQueryParam("account_id", accountID).
		SetResult(&SchemaRegistriesResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return []SchemaRegistry{}, err
	}

	if response.IsError() {
		return []SchemaRegistry{}, fmt.Errorf("clusters: %s", response.Error().(*ErrorResponse).Error.Message)
	}
	return response.Result().(*SchemaRegistriesResponse).Clusters, nil
}

type CreateSchemaRegistryRequest struct {
	Config CreateSchemaRegistryConfig `json:"config"`
}

type CreateSchemaRegistryConfig struct {
	Name            string `json:"name"`
	AccountID       string `json:"account_id"`
	KafkaClusterID  string `json:"kafka_cluster_id"`
	Location        string `json:"location"`
	ServiceProvider string `json:"service_provider"`
}

type CreateSchemaRegistryResponse struct {
	Error            interface{}      `json:"error"`
	ValidationErrors ValidationErrors `json:"validation_errors"`
	Cluster          SchemaRegistry   `json:"cluster"`
	Credentials      interface{}      `json:"credentials"`
}
type ValidationErrors struct {
}

func (c *Client) CreateSchemaRegistry(request CreateSchemaRegistryConfig) (*SchemaRegistry, error) {
	rel, err := url.Parse("schema_registries")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(&CreateSchemaRegistryRequest{Config: request}).
		SetResult(&CreateSchemaRegistryResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("clusters: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*CreateSchemaRegistryResponse).Cluster, nil
}
