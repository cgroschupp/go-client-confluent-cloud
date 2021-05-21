package confluentcloud

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

const (
	baseURLSuffix = "2.0/kafka/"
)

type KafkaClusterClient struct {
	KafkaApiEndpoint *url.URL
	BaseURLSuffix    string
	BaseURL          *url.URL
	client           *resty.Client
	token            string
}

type AccessTokenRequest struct{}
type AccessTokenResponse struct {
	Token string `json:"token"`
}

// NewKafkaClusterClient constructs a new client to connect to the relevant Kafka Confluent Cluster in order to query ACLs
//
// kafkaApiEndpoint and clusterID can be retrieved from the ID and APIEndpoint fields within Cluster
func NewKafkaClusterClient(kafkaApiEndpoint *url.URL, clusterID string, token string) *KafkaClusterClient {
	_baseURL := fmt.Sprintf("%s/%s%s/", kafkaApiEndpoint, baseURLSuffix, clusterID)
	baseURL, _ := url.Parse(_baseURL)

	client := resty.New()
	client.SetDebug(true)
	c := &KafkaClusterClient{KafkaApiEndpoint: kafkaApiEndpoint, BaseURL: baseURL, BaseURLSuffix: baseURLSuffix}
	c.client = client
	c.token = token

	return c
}

func (c *KafkaClusterClient) NewKafkaClusterRequest() *resty.Request {
	return c.client.R()
}

// GetKafkaClusterAccessToken retrieves the token required to authenticate the NewKafkaClusterClient
//
// This hits the standard confluent.cloud endpoint
func (c *Client) GetKafkaClusterAccessToken() (*string, error) {
	rel, err := url.Parse("access_tokens")
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetBody(AccessTokenRequest{}).
		SetResult(&AccessTokenResponse{}).
		SetError(&ErrorResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("access_tokens: %s", response.Error().(*ErrorResponse).Error.Message)
	}

	return &response.Result().(*AccessTokenResponse).Token, nil
}
