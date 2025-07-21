package suprsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type TenantsService interface {
	Get(context.Context, string) (*Tenant, error)
	Upsert(context.Context, string, *Tenant) (*Tenant, error)
	List(context.Context, *TenantListOptions) (*TenantList, error)
	Delete(context.Context, string) error
	GetAllCategoriesPreference(context.Context, string) (*TenantCategoryPreferencesResponse, error)
	UpdateCategoryPreference(context.Context, string, string, TenantPreferenceCategoryUpdateBody, *TenantPreferenceCategoryOptions) (*TenantCategory, error)
}

type tenantsService struct {
	client *Client
	_url   string
}

var _ TenantsService = &tenantsService{}

type TenantCategoryPreferencesResponse struct {
	Meta    *ListApiMetaInfo `json:"meta"`
	Results []TenantCategory `json:"results"`
}

type TenantCategory struct {
	Name                     string   `json:"name"`
	Category                 string   `json:"category"`
	Description              string   `json:"description"`
	RootCategory             string   `json:"root_category"`
	DefaultPreference        string   `json:"default_preference"`
	DefaultMandatoryChannels []string `json:"default_mandatory_channels"`
	VisibleToSubscriber      bool     `json:"visible_to_subscriber"`
	Preference               *string  `json:"preference"`
	MandatoryChannels        []string `json:"mandatory_channels"`
	BlockedChannels          []string `json:"blocked_channels"`
	Tags                     []string `json:"tags"`
	EffectiveTags            []string `json:"effective_tags"`
}

type TenantPreferenceCategoryOptions struct {
	Tags []string `json:"tags"`
}

func (opts *TenantPreferenceCategoryOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if len(opts.Tags) != 0 {
			query.Set("tags", strings.Join(opts.Tags, ","))
		}
	}

	return query.Encode()
}

type TenantPreferenceCategoryUpdateBody struct {
	Preference          string   `json:"preference, omitempty"`
	VisibleToSubscriber *bool    `json:"visible_to_subscriber, omitempty"`
	MandatoryChannels   []string `json:"mandatory_channels, omitempty"`
	BlockedChannels     []string `json:"blocked_channels, omitempty"`
}

func newTenantsService(client *Client) *tenantsService {
	ts := &tenantsService{
		client: client,
		_url:   fmt.Sprintf("%sv1/tenant/", client.baseUrl),
	}
	return ts
}

func (t *tenantsService) prepareQueryParams(opt *TenantListOptions) string {
	if opt == nil {
		opt = &TenantListOptions{}
	}
	opt.cleanParams()
	params := url.Values{}
	params.Add("limit", strconv.Itoa(opt.Limit))
	params.Add("offset", strconv.Itoa(opt.Offset))
	return params.Encode()
}

func (t *tenantsService) List(ctx context.Context, opts *TenantListOptions) (*TenantList, error) {
	urlStr := fmt.Sprintf("%s?%s", t._url, t.prepareQueryParams(opts))
	// prepare http.Request object
	request, err := t.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &TenantList{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *tenantsService) tenantAPIUrl(tenantId string) string {
	tenantId = url.PathEscape(tenantId)
	return fmt.Sprintf("%s%s/", t._url, tenantId)
}

func (t *tenantsService) Get(ctx context.Context, tenantId string) (*Tenant, error) {
	urlStr := t.tenantAPIUrl(tenantId)
	// prepare http.Request object
	request, err := t.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &Tenant{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *tenantsService) Upsert(ctx context.Context, tenantId string, payload *Tenant) (*Tenant, error) {
	urlStr := t.tenantAPIUrl(tenantId)
	// prepare http.Request object
	request, err := t.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &Tenant{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *tenantsService) Delete(ctx context.Context, tenantId string) error {
	urlStr := t.tenantAPIUrl(tenantId)
	// prepare http.Request object
	request, err := t.client.prepareHttpRequest("DELETE", urlStr, nil)
	if err != nil {
		return err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	//
	err = t.client.parseApiResponse(httpResponse, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t *tenantsService) GetAllCategoriesPreference(ctx context.Context, tenantId string) (*TenantCategoryPreferencesResponse, error) {
	urlStr := fmt.Sprintf("%s%s/category/", t._url, url.PathEscape(strings.TrimSpace(tenantId)))

	request, err := t.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &TenantCategoryPreferencesResponse{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *tenantsService) UpdateCategoryPreference(ctx context.Context, tenantId, category string, body TenantPreferenceCategoryUpdateBody, opts *TenantPreferenceCategoryOptions) (*TenantCategory, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%s%s/category/%s/", t._url, url.PathEscape(strings.TrimSpace(tenantId)), url.PathEscape(strings.TrimSpace(category))), opts.BuildQuery())

	request, err := t.client.prepareHttpRequest("PATCH", urlStr, body)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	resp := &TenantCategory{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
