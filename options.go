package suprsend

import (
	"net/http"
	"strings"
)

type ClientOption func(c *Client) error

func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *Client) error {
		c.baseUrl = strings.TrimSpace(baseUrl)
		return nil
	}
}

func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

func WithTimeout(timeoutInSeconds int) ClientOption {
	return func(c *Client) error {
		c.timeout = timeoutInSeconds
		return nil
	}
}
