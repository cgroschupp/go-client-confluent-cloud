package confluentcloud

import (
	"fmt"
	"net/url"
)

type ConnectorConfig = map[string]string

type ConnectorID struct {
	ID     string `json:"id"`
	IDType string `json:"id_type"`
}

type ConnectorTask struct {
	ConnectorName string `json:"connector"`
	TaskNo        int    `json:"task"`
}

type ConnectorInfo struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
	Tasks  []ConnectorTask   `json:"tasks"`
}

type Connector struct {
	ID          ConnectorID `json:"id"`
	Info        ConnectorInfo
	Name        string `json:"info.name"`
	Description string `json:"service_description"`
}

type CreateConnectorRequest struct {
	Name   string          `json:"name"`
	Config ConnectorConfig `json:"config"`
}

type ListConnectorsResponse = map[string]Connector

func (c *Client) ListConnectors(account_id, cluster_id string) ([]Connector, error) {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s/clusters/%s/connectors", account_id, cluster_id))
	if err != nil {
		return []Connector{}, err
	}

	u := c.BaseURL.ResolveReference(rel)
	response, err := c.NewRequest().
		SetQueryParam("expand", "id,info").
		SetResult(&ListConnectorsResponse{}).
		SetError(&ErrorResponse{}).
		Get(u.String())

	if err != nil {
		return []Connector{}, err
	}

	if response.IsError() {
		return []Connector{}, fmt.Errorf("connectors: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	result := response.Result().(*ListConnectorsResponse)
	list := make([]Connector, 0, len(*result))

	for _, v := range *result {
		list = append(list, v)
	}

	return list, nil
}

func (c *Client) CreateConnector(account_id, cluster_id, name string, config ConnectorConfig) (*ConnectorInfo, error) {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s/clusters/%s/connectors", account_id, cluster_id))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(&CreateConnectorRequest{Name: name, Config: config}).
		SetResult(&ConnectorInfo{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("connectors: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return response.Result().(*ConnectorInfo), nil
}

func (c *Client) UpdateConnectorConfig(account_id, cluster_id, name string, config ConnectorConfig) (*ConnectorInfo, error) {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s/clusters/%s/connectors/%s/config", account_id, cluster_id, name))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(&config).
		SetResult(&ConnectorInfo{}).
		SetError(&ErrorResponse{}).
		Put(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("connectors: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return response.Result().(*ConnectorInfo), nil

}

func (c *Client) DeleteConnector(account_id, cluster_id, name string) error {
	rel, err := url.Parse(fmt.Sprintf("accounts/%s/clusters/%s/connectors/%s", account_id, cluster_id, name))
	if err != nil {
		return err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetResult(&ConnectorInfo{}).
		SetError(&ErrorResponse{}).
		Delete(u.String())

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("connectors: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return nil
}
