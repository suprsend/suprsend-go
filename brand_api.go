package suprsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Brand has been renamed to Tenant. Brand is kept for backward-compatibilty.
// Use Tenant instead of Brand
type BrandsService interface {
	Get(context.Context, string) (*Brand, error)
	Upsert(context.Context, string, *Brand) (*Brand, error)
	List(context.Context, *BrandListOptions) (*BrandList, error)
}

type brandsService struct {
	client *Client
	_url   string
}

var _ BrandsService = &brandsService{}

func newBrandService(client *Client) *brandsService {
	bs := &brandsService{
		client: client,
		_url:   fmt.Sprintf("%sv1/brand/", client.baseUrl),
	}
	return bs
}

func (b *brandsService) prepareQueryParams(opt *BrandListOptions) string {
	if opt == nil {
		opt = &BrandListOptions{}
	}
	opt.cleanParams()
	params := url.Values{}
	params.Add("limit", strconv.Itoa(opt.Limit))
	params.Add("offset", strconv.Itoa(opt.Offset))
	return params.Encode()
}

func (b *brandsService) List(ctx context.Context, opts *BrandListOptions) (*BrandList, error) {
	urlStr := fmt.Sprintf("%s?%s", b._url, b.prepareQueryParams(opts))
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &BrandList{}
	err = b.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *brandsService) brandAPIUrl(brandId string) string {
	brandId = url.PathEscape(brandId)
	return fmt.Sprintf("%s%s/", b._url, brandId)
}

func (b *brandsService) Get(ctx context.Context, brandId string) (*Brand, error) {
	urlStr := b.brandAPIUrl(brandId)
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &Brand{}
	err = b.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *brandsService) Upsert(ctx context.Context, brandId string, payload *Brand) (*Brand, error) {
	urlStr := b.brandAPIUrl(brandId)
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &Brand{}
	err = b.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
