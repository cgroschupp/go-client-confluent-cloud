package confluentcloud

import (
	"fmt"
	"net/url"
)

type ClustersResponse struct {
	Clusters []Cluster `json:"clusters"`
}

type ClusterCreateConfig struct {
	Name            string `json:"name"`
	AccountID       string `json:"accountId"`
	Storage         int    `json:"storage"`
	Region          string `json:"region"`
	ServiceProvider string `json:"serviceProvider"`
}
type ClusterCreateRequest struct {
	Config ClusterCreateConfig `json:"config"`
}

type Cluster struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	AccountID       string `json:"account_id"`
	NetworkIngress  int    `json:"network_ingress"`
	NetworkEgress   int    `json:"network_egress"`
	Storage         int    `json:"storage"`
	Durability      string `json:"durability"`
	Status          string `json:"status"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	ServiceProvider string `json:"service_provider"`
	OrganizationID  int    `json:"organization_id"`
	Enterprise      bool   `json:"enterprise"`
	Type            string `json:"type"`
	APIEndpoint     string `json:"api_endpoint"`
	InternalProxy   bool   `json:"internal_proxy"`
	Dedicated       bool   `json:"dedicated"`
}
type ClusterResponse struct {
	Cluster Cluster `json:"cluster"`
}

func (c *Client) ListClusters(account_id string) ([]Cluster, error) {
	rel, err := url.Parse("clusters")
	if err != nil {
		return []Cluster{}, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetQueryParam("account_id", account_id).
		SetResult(&ClustersResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return []Cluster{}, err
	}

	if response.IsError() {
		return []Cluster{}, fmt.Errorf("clusters: %s", response.Error().(*ErrorResponse).Error.Message)
	}
	return response.Result().(*ClustersResponse).Clusters, nil
}

func (c *Client) CreateCluster(name, region, cloud, account_id string) (*Cluster, error) {
	rel, err := url.Parse("clusters")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	request := ClusterCreateConfig{
		Name:            name,
		Region:          region,
		ServiceProvider: cloud,
		Storage:         5000,
		AccountID:       account_id,
	}

	response, err := c.NewRequest().
		SetBody(&ClusterCreateRequest{Config: request}).
		SetResult(&ClusterResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("clusters: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*ClusterResponse).Cluster, nil
}

func (c *Client) DeleteCluster(id, account_id string) error {
	rel, err := url.Parse(fmt.Sprintf("clusters/%s", id))
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)

	data, _ := c.GetCluster(id, account_id)

	response, err := c.NewRequest().
		SetBody(&ClusterResponse{Cluster: *data}).
		SetError(&ErrorResponse{}).
		Delete(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("delete cluster: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return nil
}

func (c *Client) GetCluster(id, account_id string) (*Cluster, error) {
	rel, err := url.Parse(fmt.Sprintf("clusters/%s", id))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	fmt.Println(rel.String())

	response, err := c.NewRequest().
		SetResult(&ClusterResponse{}).
		SetQueryParam("account_id", account_id).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("get cluster: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*ClusterResponse).Cluster, nil
}

func (c *Client) UpdateCluster(id, account_id, name string) error {
	rel, err := url.Parse(fmt.Sprintf("clusters/%s", id))
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)

	data, err := c.GetCluster(id, account_id)

	if err != nil {
		return err
	}

	data.Name = name

	response, err := c.NewRequest().
		SetBody(&ClusterResponse{Cluster: *data}).
		SetError(&ErrorResponse{}).
		Put(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("update cluster: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return nil
}
