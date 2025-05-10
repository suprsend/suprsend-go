package suprsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type TenantsService interface {
	Get(context.Context, string) (*Tenant, error)
	Upsert(context.Context, string, *Tenant) (*Tenant, error)
	List(context.Context, *TenantListOptions) (*TenantList, error)
	Delete(context.Context, string) error
}

type tenantsService struct {
	client *Client
	_url   string
}

var _ TenantsService = &tenantsService{}

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
