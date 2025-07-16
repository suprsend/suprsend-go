package suprsend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/suprsend/suprsend-go/signature"
	"golang.org/x/exp/maps"
)

const (
	AuthMethod_WsKeySecret string = "ws_key_secret"
)

type Client struct {
	// auth_methods: ws_key_secret
	AuthMethod string
	// -- For workspace key/secret clients
	ApiKey    string
	ApiSecret string
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
	baseUrl  string
	debug    bool
	timeout  int
	proxyUrl string
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
		c.httpClient = defaultHTTPClient(c.debug, c.timeout, c.proxyUrl)
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

func defaultHTTPClient(debug bool, timeout int, proxyUrl string) *http.Client {
	transport := &http.Transport{}

	if proxyUrl == "" {
		proxyUrl = os.Getenv("HTTP_PROXY")
	}

	if proxyUrl != "" {
		log.Printf("Proxy url found: %s\n", proxyUrl)
		parsed, err := url.Parse(proxyUrl)
		if err != nil {
			log.Printf("Invalid HTTP_PROXY: %v\n", err)
		} else {
			transport.Proxy = http.ProxyURL(parsed)
		}
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
	if !slices.Contains([]string{AuthMethod_WsKeySecret}, c.AuthMethod) {
		return ErrInvalidAuthMethod
	}
	if c.AuthMethod == AuthMethod_WsKeySecret {
		if c.ApiKey == "" {
			return ErrMissingAPIKey
		}
		if c.ApiSecret == "" {
			return ErrMissingAPISecret
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
