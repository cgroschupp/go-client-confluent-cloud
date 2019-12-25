package confluentcloud

import (
	"fmt"
	"net/url"
)

func (c *Client) GetCloudMetadata(accountID string) (*EnvMetadataResponse, error) {
	rel, err := url.Parse("env_metadata")
	if err != nil {
		return &EnvMetadataResponse{}, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetResult(&EnvMetadataResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return &EnvMetadataResponse{}, err
	}

	if response.IsError() {
		return &EnvMetadataResponse{}, fmt.Errorf("clusters: %s", response.Error().(*ErrorResponse).Error.Message)
	}
	return response.Result().(*EnvMetadataResponse), nil
}

type EnvMetadataResponse struct {
	//Error                   interface{}               `json:"error"`
	Clouds []Clouds `json:"clouds"`
	//Status                  interface{}               `json:"status"`
	SchemaRegistryLocations []SchemaRegistryLocations `json:"schema_registry_locations"`
}

type Regions struct {
	ID                 string        `json:"id"`
	Cloud              string        `json:"cloud"`
	Zones              []interface{} `json:"zones"`
	Name               string        `json:"name"`
	IsSchedulable      bool          `json:"is_schedulable"`
	IsMultizoneEnabled bool          `json:"is_multizone_enabled"`
}

type Clouds struct {
	ID      string    `json:"id"`
	Regions []Regions `json:"regions"`
	Name    string    `json:"name"`
}

type SchemaRegistryLocations struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ClusterID       string `json:"cluster_id"`
	ServiceProvider string `json:"service_provider"`
}
