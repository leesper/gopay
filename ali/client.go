package ali

import (
	"crypto/tls"
	"net/http"
)

// Config contains all configuration info.
type Config struct {
	AppID     string
	AppKey    string
	NotifyURL string
	SandBox   bool
}

// Client handles all transactions.
type Client struct {
	config    Config
	tlsClient http.Client
}

// NewClient returns a *Client for Alipay.
func NewClient(cfg Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := http.Client{Transport: tr}
	return &Client{
		config:    cfg,
		tlsClient: client,
	}
}

// CreateOrder creates order from Alipay.
func (c *Client) CreateOrder() {}

// QueryOrder queries order from Alipay.
func (c *Client) QueryOrder() {}
