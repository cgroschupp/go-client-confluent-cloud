package confluentcloud

import (
	"fmt"
	"net/url"
)

const schemaRegistryApiResource = "schema_registries?account_id=%s"

type SchemaRegistryResponse struct {
	Clusters []Cluster `json:"clusters"`
}

type SchemaRegistryCreateResponse struct {
	Cluster Cluster `json:"cluster"`
}

type SchemaRegistryCreateRequest struct {
	Config SchemaRegistryRequest `json:"config"`
}

type SchemaRegistryRequest struct {
	AccountID           string `json:"account_id"`
	KafkaClusterID		string `json:"kafka_cluster_id"`
	Location			string `json:"location"`
	Name				string `json:"name"`
	ServiceProvider		string `json:"service_provider"`
}

func (c *Client) GetSchemaRegistries(id string) (*[]Cluster, error) {
	rel, err := url.Parse(fmt.Sprintf(schemaRegistryApiResource, id))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetResult(&SchemaRegistryResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return nil, err
	}
 
	if response.IsError() {
		return nil, fmt.Errorf("get environment: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*SchemaRegistryResponse).Clusters, nil
}

func (c *Client) CreateSchemaRegistry(accountID string, location string, serviceProvider string) (*Cluster, error) {
	rel, err := url.Parse(fmt.Sprintf(schemaRegistryApiResource, accountID))
	if err != nil {
		return nil, err
	}

	kafka_clusters, err := c.ListClusters(accountID)

	if err != nil {
		return nil, err
	}

	if len(kafka_clusters) == 0 {
		return nil, fmt.Errorf("No cluster found. Cannot enable schema registry. Environment: %s", accountID)
	}

	u := c.BaseURL.ResolveReference(rel)

	request := SchemaRegistryRequest{
		KafkaClusterID: "",
		Name: "account schema-registry",
		AccountID: accountID,
		Location: location,
		ServiceProvider: serviceProvider,
	}

	response, err := c.NewRequest().
		SetBody(&SchemaRegistryCreateRequest{Config: request}).
		SetResult(&SchemaRegistryCreateResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("Schema registry: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	cluster := &response.Result().(*SchemaRegistryCreateResponse).Cluster

	if cluster == nil {
		return nil, fmt.Errorf("Schema registry not created. Environment: %s", accountID)
	}

	return cluster, nil
}
