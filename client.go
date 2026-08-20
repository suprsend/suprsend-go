package suprsend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/suprsend/suprsend-go/signature"
)

const (
	AuthMethod_WsKeySecret string = "ws_key_secret"
	AuthMethod_ApiKey      string = "api_key"
)

type Client struct {
	// auth_methods: ws_key_secret / api_key
	AuthMethod string
	// -- For workspace key/secret clients
	ApiKey    string
	ApiSecret string
	// -- For HTTP API Key (Bearer) clients
	WorkspaceUid string
	HttpApiKey   string
	//
	Users           *usersService
	Tenants         *tenantsService
	Brands          *brandsService
	Objects         *objectsService
	SubscriberLists *subscriberListsService
	Workflows       *workflowsService
	Messages        *messagesService
	// todo: Deprecated: this
	BulkWorkflows *bulkWorkflowsService
	//
	BulkEvents *bulkEventsService
	BulkUsers  *bulkSubscribersService
	//
	baseUrl  string
	debug    bool
	timeout  int
	proxyUrl *url.URL
	//
	appInfo *AppInfo
	//
	userAgent       string
	clientUserAgent string
	//
	workflowTrigger *workflowTrigger
	eventCollector  *eventsCollector
	//
	httpClient    *http.Client
	commonHeaders map[string]string
}

func NewClient(apiKey string, apiSecret string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		AuthMethod: AuthMethod_WsKeySecret,
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
	}
	err := c.init(opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// NewWorkspaceClientWithAPIKey returns a client that authenticates with an HTTP API Key.
// It sends the key as a Bearer token, and the workspace uid in the X-SS-WSUID header.
//
// Get the workspace uid from SuprSend dashboard -> Settings -> General -> Workspace UID.
// Get the API Key from SuprSend dashboard -> Developers -> API Keys.
//
// The API Key alone identifies the workspace on the server. The workspace uid
// still fills the "env" body field, and the /{workspace}/trigger/ and
// /{workspace}/broadcast/ path segments.
func NewWorkspaceClientWithAPIKey(workspaceUid string, apiKey string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		AuthMethod:   AuthMethod_ApiKey,
		WorkspaceUid: workspaceUid,
		HttpApiKey:   apiKey,
	}
	err := c.init(opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) init(opts ...ClientOption) error {
	var err error
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return err
		}
	}
	if c.timeout <= 0 {
		c.timeout = 30
	}
	c.setDerivedBaseUrl()
	err = c.basicValidation()
	if err != nil {
		return err
	}
	if c.httpClient == nil {
		c.httpClient = defaultHTTPClient(c.debug, c.timeout, c.proxyUrl)
	}
	c.userAgent, c.clientUserAgent = buildUserAgent(c.appInfo)
	c.commonHeaders = map[string]string{
		"Content-Type":                 "application/json; charset=utf-8",
		"User-Agent":                   c.userAgent,
		"X-Suprsend-Client-User-Agent": c.clientUserAgent,
	}
	//
	c.Users = newUsersService(c)
	c.Tenants = newTenantsService(c)
	c.Brands = newBrandService(c)
	c.Objects = newObjectsService(c)
	c.Messages = newMessagesService(c)
	//
	c.Workflows = newWorkflowService(c)
	//
	c.SubscriberLists = newSubscriberListsService(c)
	c.BulkUsers = &bulkSubscribersService{client: c}
	c.BulkEvents = &bulkEventsService{client: c}
	c.BulkWorkflows = &bulkWorkflowsService{client: c}
	//
	c.workflowTrigger = newWorkflowTriggerInstance(c)
	c.eventCollector = newEventCollectorInstance(c)
	//
	return nil
}

func defaultHTTPClient(debug bool, timeout int, proxyUrl *url.URL) *http.Client {
	transport := http.DefaultTransport

	if proxyUrl != nil {
		transportClone := transport.(*http.Transport).Clone()
		transportClone.Proxy = http.ProxyURL(proxyUrl)
		transport = transportClone
	}

	if debug {
		return &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: LoggingRoundTripper{
				Proxied: transport,
			},
		}
	} else {
		return &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: transport,
		}
	}
}

func (c *Client) setDerivedBaseUrl() {
	baseUrl := c.baseUrl
	// if url not passed, set default url
	if baseUrl == "" {
		baseUrl = DEFAULT_URL
	}
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl + "/"
	}
	c.baseUrl = baseUrl
}

func (c *Client) basicValidation() error {
	if !slices.Contains([]string{AuthMethod_WsKeySecret, AuthMethod_ApiKey}, c.AuthMethod) {
		return ErrInvalidAuthMethod
	}
	if c.AuthMethod == AuthMethod_WsKeySecret {
		if c.ApiKey == "" {
			return ErrMissingAPIKey
		}
		if c.ApiSecret == "" {
			return ErrMissingAPISecret
		}
	} else if c.AuthMethod == AuthMethod_ApiKey {
		if c.WorkspaceUid == "" {
			return ErrMissingWorkspaceUid
		}
		if c.HttpApiKey == "" {
			return ErrMissingAPIKey
		}
	}
	if c.baseUrl == "" {
		return ErrMissingBaseUrl
	}
	return nil
}

// getWsIdentifierValue returns the value that identifies the workspace in a
// request path (e.g /{workspace_key}/broadcast/) and in the "env" body field.
func (c *Client) getWsIdentifierValue() string {
	if c.AuthMethod == AuthMethod_WsKeySecret {
		return c.ApiKey
	} else if c.AuthMethod == AuthMethod_ApiKey {
		return c.WorkspaceUid
	}
	return ""
}

// todo: Deprecated: this
func (c *Client) TriggerWorkflow(wf *Workflow) (*Response, error) {
	return c.workflowTrigger.TriggerWithContext(context.Background(), wf)
}

func (c *Client) TrackEvent(event *Event) (*Response, error) {
	return c.eventCollector.CollectWithContext(context.Background(), event)
}

func (c *Client) TrackEventWithContext(ctx context.Context, event *Event) (*Response, error) {
	return c.eventCollector.CollectWithContext(ctx, event)
}

func (c *Client) prepareHttpRequest(ctx context.Context, httpMethod string, httpUrl string, httpBody any,
) (*http.Request, error) {
	// Headers
	headers := maps.Clone(c.commonHeaders)
	//
	var request *http.Request
	if c.AuthMethod == AuthMethod_WsKeySecret {
		headers["Date"] = CurrentTimeFormatted()
		contentBody, sig, err := signature.GetRequestSignature(httpUrl, httpMethod, httpBody, headers, c.ApiSecret)
		if err != nil {
			return nil, &Error{Err: err}
		}
		headers["Authorization"] = fmt.Sprintf("%s:%s", c.ApiKey, sig)
		//
		request, err = http.NewRequestWithContext(ctx, httpMethod, httpUrl, bytes.NewBuffer(contentBody))
		if err != nil {
			return nil, &Error{Err: err}
		}
	} else if c.AuthMethod == AuthMethod_ApiKey {
		var contentBody []byte
		if httpMethod == "GET" || signature.SafeCheckNil(httpBody) {
			contentBody = []byte("")
		} else {
			cBytes, err := json.Marshal(httpBody)
			if err != nil {
				return nil, &Error{Err: fmt.Errorf("failed to marshal content: %w", err)}
			}
			contentBody = cBytes
		}
		headers["Authorization"] = fmt.Sprintf("Bearer %s", c.HttpApiKey)
		headers["X-SS-WSUID"] = c.WorkspaceUid
		//
		var err error
		request, err = http.NewRequestWithContext(ctx, httpMethod, httpUrl, bytes.NewBuffer(contentBody))
		if err != nil {
			return nil, &Error{Err: err}
		}
	} else {
		return nil, ErrInvalidAuthMethod
	}
	// Add headers to request
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	return request, nil
}

func (c *Client) parseApiResponse(httpResponse *http.Response, respPtr any) error {
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return &Error{Err: err}
	}
	if httpResponse.StatusCode >= 400 {
		var serr Error
		err = json.Unmarshal(responseBody, &serr)
		if err != nil {
			return &Error{Code: httpResponse.StatusCode, Message: string(responseBody)}
		}
		return &serr
	}
	// In some APIs (e.g http DELETE), we don't need to parse the response body
	// To skip response body parsing, Caller can just pass nil as the response pointer
	if respPtr == nil {
		return nil
	} else {
		err = json.Unmarshal(responseBody, respPtr)
		if err != nil {
			return &Error{Err: err}
		}
	}
	return nil
}
