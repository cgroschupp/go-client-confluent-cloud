package confluentcloud

import (
	"fmt"
	"net/url"
)

type ACLCreateRequestW = []ACLCreateRequest
type ACLCreateRequest struct {
	Pattern *Pattern `json:"pattern"`
	Entry   *Entry   `json:"entry"`
}
type Pattern struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}
type Entry struct {
	Principal      string `json:"principal"`
	Operation      string `json:"operation"`
	Host           string `json:"host"`
	PermissionType string `json:"permissionType"`
}

type CreateACLResponse = []ACLCreateRequest

type PatternFilter struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}
type EntryFilter struct {
	Principal      string `json:"principal"`
	Operation      string `json:"operation"`
	Host           string `json:"host"`
	PermissionType string `json:"permissionType"`
}
type ACLDeleteRequestW = []ACLDeleteRequest
type ACLDeleteRequest struct {
	PatternFilter *PatternFilter `json:"patternFilter"`
	EntryFilter   *EntryFilter   `json:"entryFilter"`
}
type DeleteACLResponse = []ACLDeleteRequest

func (c *Client) CreateACLs(apiEndpoint *url.URL, clusterID string, aclCreateRequestW *ACLCreateRequestW) (interface{}, error) {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return nil, err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	rel, err := url.Parse("acls")
	if err != nil {
		return nil, err
	}

	u := cc.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetAuthToken(*token).
		SetBody(aclCreateRequestW).
		SetResult(&CreateACLResponse{}).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("create_acls: %s", response.Body())
	}

	result := response.Result().(*CreateACLResponse)
	return *result, nil
}

func (c *Client) DeleteACLs(apiEndpoint *url.URL, clusterID string, aclDeleteRequestW *ACLDeleteRequestW) (interface{}, error) {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return nil, err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	rel, err := url.Parse("acls/delete")
	if err != nil {
		return nil, err
	}

	u := cc.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetAuthToken(*token).
		SetBody(aclDeleteRequestW).
		SetResult(&DeleteACLResponse{}).
		Delete(u.String())
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("delete_acls: %s", response.Body())
	}

	result := response.Result().(*DeleteACLResponse)
	return *result, nil
}