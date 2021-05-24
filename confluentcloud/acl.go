package confluentcloud

import (
	"fmt"
	"net/url"
)

type ACLRequest struct {
	PatternFilter *ListPatternFilter `json:"patternFilter"`
	EntryFilter   *ListEntryFilter   `json:"entryFilter"`
}
type ListPatternFilter struct {
	ResourceType string `json:"resourceType"`
	PatternType  string `json:"patternType"`
}
type ListEntryFilter struct {
	Operation      string `json:"operation"`
	Host           string `json:"host"`
	PermissionType string `json:"permissionType"`
}
type ListACLResponse = []ACL
type ACL struct {
	Pattern Pattern `json:"pattern"`
	Entry   Entry   `json:"entry"`
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

type ACLCreateRequestW = []ACLCreateRequest
type ACLCreateRequest struct {
	Pattern *Pattern `json:"pattern"`
	Entry   *Entry   `json:"entry"`
}
type CreateACLResponse = []ACLCreateRequest

type DeletePatternFilter struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}
type DeleteEntryFilter struct {
	Principal      string `json:"principal"`
	Operation      string `json:"operation"`
	Host           string `json:"host"`
	PermissionType string `json:"permissionType"`
}
type ACLDeleteRequestW = []ACLDeleteRequest
type ACLDeleteRequest struct {
	PatternFilter *DeletePatternFilter `json:"patternFilter"`
	EntryFilter   *DeleteEntryFilter   `json:"entryFilter"`
}
type DeleteACLResponse = []ACLDeleteRequest

func (c *Client) ListACLs(apiEndpoint *url.URL, clusterID string, aclRequest *ACLRequest) ([]ACL, error) {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return nil, err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	// cannot use url.Parse due to colon being interpreted as scheme
	suffix := url.URL{
		Path: "acls:search",
	}
	u := cc.BaseURL.ResolveReference(&suffix)
	response, err := cc.NewKafkaClusterRequest().
		SetAuthToken(*token).
		SetBody(aclRequest).
		SetResult(&ListACLResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	// response is raw, cannot parse from JSON
	if response.IsError() {
		return nil, fmt.Errorf("list_acls: %s", response.Body())
	}

	result := response.Result().(*ListACLResponse)
	return *result, nil
}

func (c *Client) CreateACLs(apiEndpoint *url.URL, clusterID string, aclCreateRequestW *ACLCreateRequestW) error {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	rel, err := url.Parse("acls")
	if err != nil {
		return err
	}

	u := cc.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetAuthToken(*token).
		SetBody(aclCreateRequestW).
		SetResult(&CreateACLResponse{}).
		Post(u.String())
	if err != nil {
		return err
	}
	if response.IsError() {
		return fmt.Errorf("create_acls: %s", response.Body())
	}

	return nil
}

func (c *Client) DeleteACLs(apiEndpoint *url.URL, clusterID string, aclDeleteRequestW *ACLDeleteRequestW) error {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	rel, err := url.Parse("acls/delete")
	if err != nil {
		return err
	}

	u := cc.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetAuthToken(*token).
		SetBody(aclDeleteRequestW).
		SetResult(&DeleteACLResponse{}).
		Delete(u.String())
	if err != nil {
		return err
	}
	if response.IsError() {
		return fmt.Errorf("delete_acls: %s", response.Body())
	}

	return nil
}
