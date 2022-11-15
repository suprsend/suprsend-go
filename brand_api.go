package suprsend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

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
	//
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, errors.New(string(responseBody))
	}
	var brandList BrandList
	err = json.Unmarshal(responseBody, &brandList)
	if err != nil {
		return nil, err
	}
	return &brandList, nil
}

func (b *brandsService) brandAPIUrl(brandId string) string {
	brandId = url.QueryEscape(brandId)
	return fmt.Sprintf("%s%s/", b._url, brandId)
}

func (b *brandsService) Get(ctx context.Context, brandId string) (*Brand, error) {
	urlStr := b.brandAPIUrl(brandId)
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, errors.New(httpResponse.Status)
	}
	var brand Brand
	err = json.Unmarshal(responseBody, &brand)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (b *brandsService) Upsert(ctx context.Context, brandId string, payload *Brand) (*Brand, error) {
	urlStr := b.brandAPIUrl(brandId)
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, errors.New(httpResponse.Status)
	}
	var brand Brand
	err = json.Unmarshal(responseBody, &brand)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}
