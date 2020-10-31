package confluentcloud

import (
	"fmt"
	"net/url"
	"time"
)

const schemaRegistryApiResource = "schema_registries?account_id=%s"
const schemaRegistryName = "account schema-registry"

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

type SchemaRegistryResponse struct {
	Error interface{}         `json:"error"`
	Clusters []SchemaRegistry `json:"clusters"`
}

type SchemaRegistryCreateResponse struct {
	Error interface{}      `json:"error"`
	Cluster SchemaRegistry `json:"cluster"`
}

type SchemaRegistryCreateRequest struct {
	Config SchemaRegistryRequest `json:"config"`
}

type SchemaRegistryRequest struct {
	AccountID           string `json:"account_id"`
	KafkaClusterID      string `json:"kafka_cluster_id"`
	Location            string `json:"location"`
	Name                string `json:"name"`
	ServiceProvider     string `json:"service_provider"`
}

func (c *Client) GetSchemaRegistry(id string) (*SchemaRegistry, error) {
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
		return nil, fmt.Errorf("Get schema registry error: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	schema_clusters := response.Result().(*SchemaRegistryResponse).Clusters

	for i:= 0; i < len(schema_clusters); i++ {
		if schema_clusters[i].Name == schemaRegistryName {
			return &schema_clusters[i], nil
		}
	}

	return nil, nil
}

func (c *Client) CreateSchemaRegistry(accountID string, location string, serviceProvider string) (*SchemaRegistry, error) {
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

	schema_cluster, err := c.GetSchemaRegistry(accountID)

	if err != nil {
		return nil, err
	}

	if schema_cluster != nil {
		return schema_cluster, nil
	}

	u := c.BaseURL.ResolveReference(rel)

	request := SchemaRegistryRequest{
		KafkaClusterID: "",
		Name: schemaRegistryName,
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