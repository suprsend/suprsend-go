package suprsend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type SubscriberListsService interface {
	GetAll(context.Context, *SubscriberListAllOptions) (*SubscriberListAll, error)
	Create(context.Context, *SubscriberListCreateInput) (*SubscriberList, error)
	Get(context.Context, string) (*SubscriberList, error)
	Add(context.Context, string, []string) (*Response, error)
	Remove(context.Context, string, []string) (*Response, error)
	Broadcast(context.Context, *SubscriberListBroadcast) (*Response, error)
}

type subscriberListsService struct {
	client             *Client
	_subscriberListUrl string
	_broadcastUrl      string
	//
	nonErrDefaultResponse *Response
}

var _ SubscriberListsService = &subscriberListsService{}

func newSubscriberListsService(client *Client) *subscriberListsService {
	bs := &subscriberListsService{
		client:             client,
		_subscriberListUrl: fmt.Sprintf("%sv1/subscriber_list/", client.baseUrl),
		_broadcastUrl:      fmt.Sprintf("%s%s/broadcast/", client.baseUrl, client.ApiKey),
		//
		nonErrDefaultResponse: &Response{Success: true, StatusCode: 201, Message: `{"success":true}`},
	}
	return bs
}

func (s *subscriberListsService) prepareQueryParams(opt *SubscriberListAllOptions) string {
	if opt == nil {
		opt = &SubscriberListAllOptions{}
	}
	opt.cleanParams()
	params := url.Values{}
	params.Add("limit", strconv.Itoa(opt.Limit))
	params.Add("offset", strconv.Itoa(opt.Offset))
	return params.Encode()
}

func (s *subscriberListsService) GetAll(ctx context.Context, opts *SubscriberListAllOptions) (*SubscriberListAll, error) {
	urlStr := fmt.Sprintf("%s?%s", s._subscriberListUrl, s.prepareQueryParams(opts))
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := s.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpResponse.StatusCode, string(responseBody))
	}
	var all SubscriberListAll
	err = json.Unmarshal(responseBody, &all)
	if err != nil {
		return nil, err
	}
	return &all, nil
}

func (s *subscriberListsService) validateListId(listId string) (string, error) {
	listId = strings.TrimSpace(listId)
	if listId == "" {
		return listId, errors.New("missing list_id")
	}
	return listId, nil
}

func (s *subscriberListsService) Create(ctx context.Context, createParams *SubscriberListCreateInput) (*SubscriberList, error) {
	var err error
	if createParams == nil {
		return nil, errors.New("missing payload")
	}
	createParams.ListId, err = s.validateListId(createParams.ListId)
	if err != nil {
		return nil, err
	}
	urlStr := s._subscriberListUrl
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("POST", urlStr, createParams)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := s.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpResponse.StatusCode, string(responseBody))
	}
	var sl SubscriberList
	err = json.Unmarshal(responseBody, &sl)
	if err != nil {
		return nil, err
	}
	return &sl, nil
}

func (b *subscriberListsService) listDetailAPIUrl(listId string) string {
	listId = url.QueryEscape(listId)
	return fmt.Sprintf("%s%s/", b._subscriberListUrl, listId)
}

func (s *subscriberListsService) Get(ctx context.Context, listId string) (*SubscriberList, error) {
	listId, err := s.validateListId(listId)
	if err != nil {
		return nil, err
	}
	urlStr := s.listDetailAPIUrl(listId)
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := s.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpResponse.StatusCode, string(responseBody))
	}
	var sl SubscriberList
	err = json.Unmarshal(responseBody, &sl)
	if err != nil {
		return nil, err
	}
	return &sl, nil
}

func (s *subscriberListsService) Add(ctx context.Context, listId string, distinctIds []string) (*Response, error) {
	listId, err := s.validateListId(listId)
	if err != nil {
		return nil, err
	}
	if len(distinctIds) == 0 {
		return s.nonErrDefaultResponse, nil
	}
	urlStr := fmt.Sprintf("%ssubscriber/add/", s.listDetailAPIUrl(listId))
	payload := map[string]interface{}{"distinct_ids": distinctIds}
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := s.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	suprResponse, err := s.formatAPIResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return suprResponse, nil
}

func (s *subscriberListsService) Remove(ctx context.Context, listId string, distinctIds []string) (*Response, error) {
	listId, err := s.validateListId(listId)
	if err != nil {
		return nil, err
	}
	if len(distinctIds) == 0 {
		return s.nonErrDefaultResponse, nil
	}
	urlStr := fmt.Sprintf("%ssubscriber/remove/", s.listDetailAPIUrl(listId))
	payload := map[string]interface{}{"distinct_ids": distinctIds}
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := s.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	suprResponse, err := s.formatAPIResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return suprResponse, nil
}

func (s *subscriberListsService) Broadcast(ctx context.Context, broadcastIns *SubscriberListBroadcast) (*Response, error) {
	if broadcastIns == nil {
		return nil, errors.New("missing payload")
	}
	broadcastBody, _, err := broadcastIns.getFinalJson(s.client)
	if err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("POST", s._broadcastUrl, broadcastBody)
	if err != nil {
		return nil, err
	}
	httpResponse, err := s.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	suprResponse, err := s.formatAPIResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return suprResponse, nil
}

func (s *subscriberListsService) formatAPIResponse(httpRes *http.Response) (*Response, error) {
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpRes.StatusCode, string(respBody))
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}
