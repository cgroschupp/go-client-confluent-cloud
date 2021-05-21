package confluentcloud

import (
	"fmt"
	"net/url"
)

type ACLDeleteRequestW = []ACLDeleteRequest
type ACLDeleteRequest struct {
	PatternFilter *PatternFilter `json:"patternFilter"`
	EntryFilter   *EntryFilter   `json:"entryFilter"`
}
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

type DeleteACLResponse = []ACLDeleteRequest

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