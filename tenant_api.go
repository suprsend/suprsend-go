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
	UpdateCategoryPreference(context.Context, string, string, TenantPreferenceCategoryUpdateBody, *TenantPreferenceCategoryOptions) (*TenantCategoryPreferencesResponse, error)
	GetAllCategoriesPreference(context.Context, string) (*TenantCategoryPreferencesResponse, error)
}

type tenantsService struct {
	client *Client
	_url   string
}

var _ TenantsService = &tenantsService{}

type TenantCategoryMeta struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type TenantCategoryPreferencesResponse struct {
	Meta    TenantCategoryMeta `json:"meta"`
	Results []TenantCategory   `json:"results"`
}

type TenantCategory struct {
	Name                     string   `json:"name"`
	Category                 string   `json:"category"`
	Description              string   `json:"description"`
	RootCategory             string   `json:"root_category"`
	DefaultPreference        string   `json:"default_preference"`
	DefaultMandatoryChannels []string `json:"default_mandatory_channels"`
	VisibleToSubscriber      bool     `json:"visible_to_subscriber"`
	Preference               string   `json:"preference"`
	MandatoryChannels        []string `json:"mandatory_channels"`
	BlockedChannels          []string `json:"blocked_channels"`
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

func (t *tenantsService) UpdateCategoryPreference(ctx context.Context, tenantId, category string, body TenantPreferenceCategoryUpdateBody, opts *TenantPreferenceCategoryOptions) (*TenantCategoryPreferencesResponse, error) {
	if strings.TrimSpace(tenantId) == "" {
		return nil, &Error{Message: "tenant_id is required"}
	}

	if strings.TrimSpace(category) == "" {
		return nil, &Error{Message: "category is required"}
	}

	urlStr := fmt.Sprintf("%s%s/preference/category/%s/", t._url, url.PathEscape(strings.TrimSpace(tenantId)), url.PathEscape(strings.TrimSpace(category)))

	query := url.Values{}
	if opts != nil {
		if len(opts.Tags) != 0 {
			query.Set("tags", strings.Join(opts.Tags, ","))
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := t.client.prepareHttpRequest("PATCH", urlStr, body)
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

func (t *tenantsService) GetAllCategoriesPreference(ctx context.Context, tenantId string) (*TenantCategoryPreferencesResponse, error) {
	if strings.TrimSpace(tenantId) == "" {
		return nil, &Error{Message: "tenant_id is required"}
	}

	urlStr := fmt.Sprintf("%s%s/category", t._url, url.PathEscape(strings.TrimSpace(tenantId)))

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

type TenantPreferenceCategoryOptions struct {
	Tags []string `json:"tags"`
}

type TenantPreferenceCategoryUpdateBody struct {
	Preference          string   `json:"preference"`
	VisibleToSubscriber bool     `json:"visible_to_subscriber"`
	MandatoryChannels   []string `json:"mandatory_channels"`
	BlockedChannels     []string `json:"blocked_channels"`
}
