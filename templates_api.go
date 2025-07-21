package suprsend

import (
	"context"
	"fmt"
)

type TemplatesService interface {
	GetAllTemplates(context.Context, *TemplateOptions) (*TemplateList, error)
	GetDetails(context.Context, string) (*Template, error)
	GetChannelContent(context.Context, string, string) (*ChannelTemplateResponse, error)
}

var _ TemplatesService = &templatesService{}

type templatesService struct {
	client *Client
	_url   string
}

func newTemplatesService(client *Client) *templatesService {
	ts := &templatesService{
		client: client,
		_url:   fmt.Sprintf("%sv1/template/", client.baseUrl),
	}
	return ts
}

func (t *templatesService) GetAllTemplates(ctx context.Context, opts *TemplateOptions) (*TemplateList, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%s", t._url), opts.BuildQuery())

	request, err := t.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	resp := &TemplateList{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *templatesService) GetDetails(ctx context.Context, templateSlug string) (*Template, error) {
	urlStr := fmt.Sprintf("%s/%s/", t._url, templateSlug)

	request, err := t.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	resp := &Template{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *templatesService) GetChannelContent(ctx context.Context, templateSlug string, channelSlug string) (*ChannelTemplateResponse, error) {
	urlStr := fmt.Sprintf("%s/%s/channel/%s/", t._url, templateSlug, channelSlug)

	request, err := t.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := t.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	resp := &ChannelTemplateResponse{}
	err = t.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
