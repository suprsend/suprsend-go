package suprsend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/suprsend/suprsend-go/signature"
	"golang.org/x/exp/maps"
)

const (
	AuthMethod_WsKeySecret  string = "ws_key_secret"
	AuthMethod_ServiceToken string = "service_token"
)

type Client struct {
	// auth_methods: ws_key_secret / service_token
	AuthMethod string
	// -- For workspace key/secret clients
	ApiKey    string
	ApiSecret string
	// -- For service token clients
	ServiceToken string
	WorkspaceUid string
	//
	Users           *usersService
	Tenants         *tenantsService
	Brands          *brandsService
	Objects         *objectsService
	SubscriberLists *subscriberListsService
	Workflows       *workflowsService
	// todo: Deprecated: this
	BulkWorkflows *bulkWorkflowsService
	//
	BulkEvents *bulkEventsService
	BulkUsers  *bulkSubscribersService
	//
	baseUrl string
	debug   bool
	timeout int
	//
	sdkVersion string
	userAgent  string
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

func NewServiceTokenClient(token string, workspaceUid string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		AuthMethod:   AuthMethod_ServiceToken,
		ServiceToken: token,
		WorkspaceUid: workspaceUid,
	}
	err := c.init(opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) init(opts ...ClientOption) error {
	c.sdkVersion = VERSION
	c.userAgent = fmt.Sprintf("suprsend/%s;go/%s", VERSION, runtime.Version())
	//
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
		c.httpClient = defaultHTTPClient(c.debug, c.timeout)
	}
	c.commonHeaders = map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"User-Agent":   c.userAgent,
	}
	//
	c.Users = newUsersService(c)
	c.Tenants = newTenantsService(c)
	c.Brands = newBrandService(c)
	c.Objects = newObjectsService(c)
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

func defaultHTTPClient(debug bool, timeout int) *http.Client {
	if debug {
		return &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: LoggingRoundTripper{
				Proxied: http.DefaultTransport,
			},
		}
	} else {
		return &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
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
	if !slices.Contains([]string{AuthMethod_WsKeySecret, AuthMethod_ServiceToken}, c.AuthMethod) {
		return ErrInvalidAuthMethod
	}
	if c.AuthMethod == AuthMethod_WsKeySecret {
		if c.ApiKey == "" {
			return ErrMissingAPIKey
		}
		if c.ApiSecret == "" {
			return ErrMissingAPISecret
		}
	} else if c.AuthMethod == AuthMethod_ServiceToken {
		if c.ServiceToken == "" {
			return ErrMissingServiceToken
		}
		if c.WorkspaceUid == "" {
			return ErrMissingWorkspaceUid
		}
	}
	if c.baseUrl == "" {
		return ErrMissingBaseUrl
	}
	return nil
}

func (c *Client) getWsIdentifierValue() string {
	if c.AuthMethod == AuthMethod_WsKeySecret {
		return c.ApiKey
	} else if c.AuthMethod == AuthMethod_ServiceToken {
		return c.WorkspaceUid
	}
	return ""
}

// todo: Deprecated: this
func (c *Client) TriggerWorkflow(wf *Workflow) (*Response, error) {
	return c.workflowTrigger.Trigger(wf)
}

func (c *Client) TrackEvent(event *Event) (*Response, error) {
	return c.eventCollector.Collect(event)
}

func (c *Client) prepareHttpRequest(httpMethod string, httpUrl string, httpBody any,
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
		request, err = http.NewRequest(httpMethod, httpUrl, bytes.NewBuffer(contentBody))
		if err != nil {
			return nil, &Error{Err: err}
		}
	} else if c.AuthMethod == AuthMethod_ServiceToken {
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
		headers["Authorization"] = fmt.Sprintf("ServiceToken %s", c.ServiceToken)
		headers["X-SS-WSUID"] = c.WorkspaceUid
		//
		var err error
		request, err = http.NewRequest(httpMethod, httpUrl, bytes.NewBuffer(contentBody))
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
