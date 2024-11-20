package suprsend

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type ObjectRequest struct {
	ObjectType string `json:"object_type"`
	ObjectId   string `json:"object_id"`
}

type ObjectsService interface {
	List(context.Context, string, *CursorPaginationListOptions) (*CursorPaginationList, error)
	Get(context.Context, *ObjectRequest) (map[string]any, error)
	Upsert(context.Context, *ObjectRequest, map[string]any) (map[string]any, error)
	Edit(context.Context, *ObjectRequest, map[string]any) (map[string]any, error)
	Delete(context.Context, *ObjectRequest) (map[string]any, error)
	BulkDelete(context.Context, string, map[string]any) (map[string]any, error)
	GetSubscriptions(context.Context, *ObjectRequest, *CursorPaginationListOptions) (*CursorPaginationList, error)
	CreateSubscriptions(context.Context, *ObjectRequest, map[string]any) (map[string]any, error)
	DeleteSubscriptions(context.Context, *ObjectRequest, map[string]any) (map[string]any, error)
	GetInstance(*ObjectRequest) (IObject, error)
}

type objectsService struct {
	client *Client
	_url   string
}

var _ ObjectsService = &objectsService{}

func newObjectsService(client *Client) *objectsService {
	os := &objectsService{
		client: client,
		_url:   fmt.Sprintf("%sv1/object/", client.baseUrl),
	}
	return os
}

func (o *objectsService) prepareQueryParams(opt *CursorPaginationListOptions) string {
	if opt == nil {
		opt = &CursorPaginationListOptions{}
	}
	opt.cleanParams()
	params := url.Values{}
	params.Add("limit", strconv.Itoa(opt.Limit))
	if opt.Before != "" {
		params.Add("before", opt.Before)
	}
	if opt.After != "" {
		params.Add("after", opt.After)
	}
	return params.Encode()
}

func (o *objectsService) validateObjectEntityId(entityId string) (string, error) {
	entityId = strings.TrimSpace(entityId)
	if entityId == "" {
		return "", fmt.Errorf("missing entity_id")
	}
	return entityId, nil
}

func (O *ObjectRequest) validateObjectRequest() error {
	if O.ObjectType == "" {
		return fmt.Errorf("missing object_type")
	}
	O.ObjectType = strings.TrimSpace(O.ObjectType)
	if O.ObjectType == "" {
		return fmt.Errorf("missing object_type")
	}
	//
	if O.ObjectId == "" {
		return fmt.Errorf("missing object_id")
	}
	O.ObjectId = strings.TrimSpace(O.ObjectId)
	if O.ObjectId == "" {
		return fmt.Errorf("missing object_id")
	}
	return nil
}

func (o *objectsService) List(ctx context.Context, objectType string, opts *CursorPaginationListOptions) (*CursorPaginationList, error) {
	objectType, err := o.validateObjectEntityId(objectType)
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s/?%s", o._url, objectType, o.prepareQueryParams(opts))
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericListApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) objectAPIUrl(obj *ObjectRequest) (string, error) {
	err := obj.validateObjectRequest()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s/%s/", o._url, obj.ObjectType, obj.ObjectId), nil
}

func (o *objectsService) Get(ctx context.Context, obj *ObjectRequest) (map[string]any, error) {
	urlStr, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) Upsert(ctx context.Context, obj *ObjectRequest, payload map[string]any) (map[string]any, error) {
	urlStr, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) Edit(ctx context.Context, obj *ObjectRequest, payload map[string]any) (map[string]any, error) {
	urlStr, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("PATCH", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) Delete(ctx context.Context, obj *ObjectRequest) (map[string]any, error) {
	urlStr, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("DELETE", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) bulkOpsUrl(objectType string) (string, error) {
	objectType, err := o.validateObjectEntityId(objectType)
	if err != nil {
		return "", err
	}
	objectType = url.QueryEscape(objectType)
	return fmt.Sprintf("%sv1/bulk/object/%s/", o.client.baseUrl, objectType), nil
}

func (o *objectsService) BulkDelete(ctx context.Context, objectType string, payload map[string]any) (map[string]any, error) {
	urlStr, err := o.bulkOpsUrl(objectType)
	if err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("DELETE", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) GetSubscriptions(ctx context.Context, obj *ObjectRequest, opts *CursorPaginationListOptions) (*CursorPaginationList, error) {
	_url, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%ssubscription/?%s", _url, o.prepareQueryParams(opts))
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericListApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) CreateSubscriptions(ctx context.Context, obj *ObjectRequest, payload map[string]any) (map[string]any, error) {
	_url, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%ssubscription/", _url)
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) DeleteSubscriptions(ctx context.Context, obj *ObjectRequest, payload map[string]any) (map[string]any, error) {
	_url, err := o.objectAPIUrl(obj)
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%ssubscription/", _url)
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("DELETE", urlStr, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) GetInstance(obj *ObjectRequest) (IObject, error) {
	obj.ObjectType = strings.TrimSpace(obj.ObjectType)
	obj.ObjectId = strings.TrimSpace(obj.ObjectId)
	return newObject(o.client, obj), nil
}

type IObject interface {
	Save() (map[string]any, error)
	//
	AppendKV(string, interface{})
	Append(map[string]interface{})
	SetKV(string, interface{})
	Set(map[string]interface{})
	SetOnceKV(string, interface{})
	SetOnce(map[string]interface{})
	IncrementKV(string, interface{})
	Increment(map[string]interface{})
	RemoveKV(string, interface{})
	Remove(map[string]interface{})
	Unset([]string)
	//
	SetPreferredLanguage(string)
	SetTimezone(string)
	//
	AddEmail(value string)
	RemoveEmail(value string)
	//
	AddSms(value string)
	RemoveSms(value string)
	//
	AddWhatsapp(value string)
	RemoveWhatsapp(value string)
	//
	AddAndroidpush(value, provider string)
	RemoveAndroidpush(value, provider string)
	//
	AddIospush(value, provider string)
	RemoveIospush(value, provider string)
	//
	AddWebpush(value map[string]interface{}, provider string)
	RemoveWebpush(value map[string]interface{}, provider string)
	//
	AddSlack(value map[string]interface{})
	RemoveSlack(value map[string]interface{})
	//
	AddMSTeams(value map[string]interface{})
	RemoveMSTeams(value map[string]interface{})
}

var _ IObject = &object{}

type object struct {
	client      *Client
	objectType  string
	objectId    string
	_url        string
	_superProps map[string]interface{}
	//
	_errors    []string
	_warnings  []string
	operations []map[string]interface{}
	//
	_helper *objectHelper
}

func newObject(client *Client, obj *ObjectRequest) IObject {
	o := &object{
		client:     client,
		objectType: obj.ObjectType,
		objectId:   obj.ObjectId,
	}
	// object url
	o._url = fmt.Sprintf("%sv1/object/%s/%s/", client.baseUrl, o.objectType, o.objectId)
	o._superProps = map[string]interface{}{"$ss_sdk_version": client.userAgent}
	o._helper = newObjectHelper()
	return o
}

func (o *object) validateBody() error {
	if len(o._warnings) > 0 {
		msg := fmt.Sprintf("[object: %s/%s] %s", o.objectType, o.objectId, strings.Join(o._warnings, "\n"))
		log.Println("WARNING:", msg)
	}
	if len(o._errors) > 0 {
		msg := fmt.Sprintf("[object: %s/%s] %s", o.objectType, o.objectId, strings.Join(o._errors, "\n"))
		log.Println("ERROR:", msg)
	}
	return nil
}

func (o *object) Save() (map[string]any, error) {
	if err := o.validateBody(); err != nil {
		return nil, err
	}
	// -----
	payload := map[string]any{
		"operations": o.operations,
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("PATCH", o._url, payload)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp, err := ParseGenericApiResponse(err, httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *object) _collectPayload() {
	resp := o._helper.getIdentityEvent()
	if len(resp.errors) > 0 {
		o._errors = append(o._errors, resp.errors...)
	}
	if len(resp.info) > 0 {
		o._warnings = append(o._warnings, resp.info...)
	}
	if len(resp.payload) > 0 {
		o.operations = append(o.operations, resp.payload)
	}
}

/*
Usage:
 1. append(k, v)
 2. append({k1: v1, k2, v2})
*/
func (o *object) AppendKV(k string, v interface{}) {
	caller := "appendKV"
	o._helper.appendKV(k, v, map[string]interface{}{}, caller)
	o._collectPayload()
}

func (o *object) Append(kvMap map[string]interface{}) {
	caller := "append"
	for k, v := range kvMap {
		o._helper.appendKV(k, v, kvMap, caller)

	}
	o._collectPayload()
}

/*
Usage:
 1. SetKV(k, v)
 2. Set({k1: v1, k2, v2})
*/
func (o *object) SetKV(k string, v interface{}) {
	caller := "setKV"
	o._helper.setKV(k, v, map[string]interface{}{}, caller)
	o._collectPayload()
}

func (o *object) Set(kvMap map[string]interface{}) {
	caller := "set"
	for k, v := range kvMap {
		o._helper.setKV(k, v, kvMap, caller)

	}
	o._collectPayload()
}

/*
Usage:
 1. SetOnceKV(k, v)
 2. SetOnce({k1: v1, k2, v2})
*/
func (o *object) SetOnceKV(k string, v interface{}) {
	caller := "set_onceKV"
	o._helper.setOnceKV(k, v, map[string]interface{}{}, caller)
	o._collectPayload()
}

func (o *object) SetOnce(kvMap map[string]interface{}) {
	caller := "set_once"
	for k, v := range kvMap {
		o._helper.setOnceKV(k, v, kvMap, caller)

	}
	o._collectPayload()
}

/*
Usage:
 1. IncrementKV(k, v)
 2. Increment({k1: v1, k2, v2})
*/
func (o *object) IncrementKV(k string, v interface{}) {
	caller := "incrementKV"
	o._helper.incrementKV(k, v, map[string]interface{}{}, caller)
	o._collectPayload()
}

func (o *object) Increment(kvMap map[string]interface{}) {
	caller := "increment"
	for k, v := range kvMap {
		o._helper.incrementKV(k, v, kvMap, caller)

	}
	o._collectPayload()
}

/*
Usage:
 1. RemoveKV(k, v)
 2. Remove({k1: v1, k2, v2})
*/
func (o *object) RemoveKV(k string, v interface{}) {
	caller := "removeKV"
	o._helper.removeKV(k, v, map[string]interface{}{}, caller)
	o._collectPayload()
}

func (o *object) Remove(kvMap map[string]interface{}) {
	caller := "remove"
	for k, v := range kvMap {
		o._helper.removeKV(k, v, kvMap, caller)

	}
	o._collectPayload()
}

/*
Usage:
 1. unset(k)
 2. unset([k1, k2])
*/
func (o *object) Unset(keys []string) {
	caller := "unset"
	for _, key := range keys {
		o._helper.unsetKey(key, caller)
	}
	o._collectPayload()
}

// ----------------------- Preferred language

func (o *object) SetPreferredLanguage(langCode string) {
	caller := "set_preferred_language"
	o._helper.setPreferredLanguage(langCode, caller)
	o._collectPayload()
}

// SetTimezone : set IANA supported timezone as subscriber property
func (o *object) SetTimezone(timezone string) {
	caller := "set_timezone"
	o._helper.setTimezone(timezone, caller)
	o._collectPayload()
}

// ------------------------ Email

func (o *object) AddEmail(value string) {
	caller := "add_email"
	o._helper.addEmail(value, caller)
	o._collectPayload()
}

func (o *object) RemoveEmail(value string) {
	caller := "remove_email"
	o._helper.removeEmail(value, caller)
	o._collectPayload()
}

// ------------------------ SMS

func (o *object) AddSms(value string) {
	caller := "add_sms"
	o._helper.addSms(value, caller)
	o._collectPayload()
}

func (o *object) RemoveSms(value string) {
	caller := "remove_sms"
	o._helper.removeSms(value, caller)
	o._collectPayload()
}

// ------------------------ Whatsapp

func (o *object) AddWhatsapp(value string) {
	caller := "add_whatsapp"
	o._helper.addWhatsapp(value, caller)
	o._collectPayload()
}

func (o *object) RemoveWhatsapp(value string) {
	caller := "remove_whatsapp"
	o._helper.removeWhatsapp(value, caller)
	o._collectPayload()
}

// ------------------------ Androidpush [providers: fcm / xiaomi / oppo]

func (o *object) AddAndroidpush(value, provider string) {
	caller := "add_androidpush"
	o._helper.addAndroidpush(value, provider, caller)
	o._collectPayload()
}

func (o *object) RemoveAndroidpush(value, provider string) {
	caller := "remove_androidpush"
	o._helper.removeAndroidpush(value, provider, caller)
	o._collectPayload()
}

// ------------------------ Iospush [providers: apns]

func (o *object) AddIospush(value, provider string) {
	caller := "add_iospush"
	o._helper.addIospush(value, provider, caller)
	o._collectPayload()
}

func (o *object) RemoveIospush(value, provider string) {
	caller := "remove_iospush"
	o._helper.removeIospush(value, provider, caller)
	o._collectPayload()
}

// ------------------------ Webpush [providers: vapid]

func (o *object) AddWebpush(value map[string]interface{}, provider string) {
	caller := "add_webpush"
	o._helper.addWebpush(value, provider, caller)
	o._collectPayload()
}

func (o *object) RemoveWebpush(value map[string]interface{}, provider string) {
	caller := "remove_webpush"
	o._helper.removeWebpush(value, provider, caller)
	o._collectPayload()
}

// ------------------------ Slack

func (o *object) AddSlack(value map[string]interface{}) {
	caller := "add_slack"
	o._helper.addSlack(value, caller)
	o._collectPayload()
}

func (o *object) RemoveSlack(value map[string]interface{}) {
	caller := "remove_slack"
	o._helper.removeSlack(value, caller)
	o._collectPayload()
}

// ------------------------ MS Teams

func (o *object) AddMSTeams(value map[string]interface{}) {
	caller := "add_ms_teams"
	o._helper.addMSTeams(value, caller)
	o._collectPayload()
}

func (o *object) RemoveMSTeams(value map[string]interface{}) {
	caller := "remove_ms_teams"
	o._helper.removeMSTeams(value, caller)
	o._collectPayload()
}
