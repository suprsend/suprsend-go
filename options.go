package suprsend

import (
	"net/http"
	"net/url"
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

func WithProxyUrl(proxyUrl string) ClientOption {
	return func(c *Client) error {
		parsed, err := url.Parse(proxyUrl)
		if err != nil {
			return &Error{Code: 404, Message: "invalid proxyUrl", Err: err}
		}
		c.proxyUrl = parsed
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
