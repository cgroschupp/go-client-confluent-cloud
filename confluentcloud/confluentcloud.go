package confluentcloud

import (
	"net/url"

	resty "github.com/go-resty/resty/v2"
)

const (
	defaultBaseURL = "https://confluent.cloud/api/"
	libraryVersion = "0.1"
	userAgent      = "go-confluent-cloud " + libraryVersion
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string
	email     string
	password  string
	token     string
	client    *resty.Client
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}

func NewClient(email, password string) *Client {
	rc := resty.New()
	rc.SetDebug(true)
	return NewClientWithRestyClient(email, password, rc)
}

func NewClientWithRestyClient(email, password string, restyClient *resty.Client) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{BaseURL: baseURL, email: email, password: password, UserAgent: userAgent}
	c.client = restyClient
	return c
}

func (c *Client) NewRequest() *resty.Request {
	return c.client.R().
		SetHeader("User-Agent", c.UserAgent).
		SetError(&ErrorResponse{})
}
