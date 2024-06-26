package suprsend

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/suprsend/suprsend-go/signature"
	"golang.org/x/exp/maps"
)

type Client struct {
	ApiKey    string
	ApiSecret string
	//
	Users           *subscribersService
	Tenants         *tenantsService
	Brands          *brandsService
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
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
		sdkVersion: VERSION,
		userAgent:  fmt.Sprintf("suprsend/%s;go/%s", VERSION, runtime.Version()),
	}
	//
	var err error
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return nil, err
		}
	}
	if c.timeout <= 0 {
		c.timeout = 30
	}
	c.setDerivedBaseUrl()
	err = c.basicValidation()
	if err != nil {
		return nil, err
	}
	if c.httpClient == nil {
		c.httpClient = defaultHTTPClient(c.debug, c.timeout)
	}
	c.commonHeaders = map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"User-Agent":   c.userAgent,
	}
	//
	c.Users = &subscribersService{client: c}
	c.Tenants = newTenantsService(c)
	c.Brands = newBrandService(c)
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
	return c, nil
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
	if c.ApiKey == "" {
		return ErrMissingAPIKey
	}
	if c.ApiSecret == "" {
		return ErrMissingAPISecret
	}
	if c.baseUrl == "" {
		return ErrMissingBaseUrl
	}
	return nil
}

// todo: Deprecated: this
func (c *Client) TriggerWorkflow(wf *Workflow) (*Response, error) {
	return c.workflowTrigger.Trigger(wf)
}

func (c *Client) TrackEvent(event *Event) (*Response, error) {
	return c.eventCollector.Collect(event)
}

func (c *Client) prepareHttpRequest(httpMethod string, httpUrl string, httpBody interface{},
) (*http.Request, error) {
	// Headers
	headers := maps.Clone(c.commonHeaders)
	maps.Copy(headers, map[string]string{"Date": CurrentTimeFormatted()})
	//
	contentBody, sig, err := signature.GetRequestSignature(httpUrl, httpMethod, httpBody, headers, c.ApiSecret)
	if err != nil {
		return nil, err
	}
	headers["Authorization"] = fmt.Sprintf("%s:%s", c.ApiKey, sig)
	//
	request, err := http.NewRequest(httpMethod, httpUrl, bytes.NewBuffer(contentBody))
	if err != nil {
		return nil, err
	}
	// Add headers to request
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	return request, nil
}
