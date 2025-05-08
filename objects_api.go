package suprsend

import (
	"context"
	"fmt"
	"net/url"
)

type ObjectIdentifier struct {
	ObjectType string `json:"object_type"`
	Id         string `json:"id"`
}

type ObjectsService interface {
	List(context.Context, string, *CursorListApiOptions) (*CursorListApiResponse, error)
	Get(context.Context, ObjectIdentifier) (map[string]any, error)
	Upsert(context.Context, ObjectIdentifier, map[string]any) (map[string]any, error)
	Edit(context.Context, ObjectEditRequest) (map[string]any, error)
	Delete(context.Context, ObjectIdentifier) error
	BulkDelete(context.Context, string, ObjectBulkDeletePayload) error
	GetSubscriptions(context.Context, ObjectIdentifier, *CursorListApiOptions) (*CursorListApiResponse, error)
	CreateSubscriptions(context.Context, ObjectIdentifier, map[string]any) (map[string]any, error)
	DeleteSubscriptions(context.Context, ObjectIdentifier, map[string]any) error
	GetEditInstance(ObjectIdentifier) ObjectEdit
}

type objectsService struct {
	client   *Client
	_url     string
	_bulkUrl string
}

var _ ObjectsService = &objectsService{}

func newObjectsService(client *Client) *objectsService {
	os := &objectsService{
		client:   client,
		_url:     fmt.Sprintf("%sv1/object/", client.baseUrl),
		_bulkUrl: fmt.Sprintf("%sv1/bulk/object/", client.baseUrl),
	}
	return os
}

func (o *objectsService) List(ctx context.Context, objectType string, opts *CursorListApiOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%s%s/", o._url, url.PathEscape(objectType)), opts.BuildQuery())
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
	resp := &CursorListApiResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) objectDetailAPIUrl(objectType, id string) string {
	return fmt.Sprintf(
		"%s%s/%s/",
		o._url,
		url.PathEscape(objectType),
		url.PathEscape(id),
	)
}

func (o *objectsService) Get(ctx context.Context, obj ObjectIdentifier) (map[string]any, error) {
	urlStr := o.objectDetailAPIUrl(obj.ObjectType, obj.Id)
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
	resp := map[string]any{}
	err = o.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) Upsert(ctx context.Context, obj ObjectIdentifier, payload map[string]any) (map[string]any, error) {
	urlStr := o.objectDetailAPIUrl(obj.ObjectType, obj.Id)
	if payload == nil {
		payload = map[string]any{}
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
	//
	resp := map[string]any{}
	err = o.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// either (identifier + payload) OR editInstance must be provided
type ObjectEditRequest struct {
	Identifier *ObjectIdentifier
	// {"operations": [{"$set": {"prop1": "val1"}, {"$append": {"$email": "abc@test.com"}}]}
	Payload map[string]any
	//
	EditInstance ObjectEdit
}

func (o *objectsService) Edit(ctx context.Context, req ObjectEditRequest) (map[string]any, error) {
	var urlStr string
	var payload map[string]any
	if req.EditInstance != nil {
		oe := req.EditInstance.(*objectEdit)
		oe.validateBody()
		payload = oe.GetPayload()
		urlStr = o.objectDetailAPIUrl(oe.objectType, oe.objectId)
	} else {
		payload = req.Payload
		if payload == nil {
			payload = map[string]any{}
		}
		urlStr = o.objectDetailAPIUrl(req.Identifier.ObjectType, req.Identifier.Id)
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("PATCH", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := map[string]any{}
	err = o.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) Delete(ctx context.Context, obj ObjectIdentifier) error {
	urlStr := o.objectDetailAPIUrl(obj.ObjectType, obj.Id)
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("DELETE", urlStr, nil)
	if err != nil {
		return err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	//
	err = o.client.parseApiResponse(httpResponse, nil)
	if err != nil {
		return err
	}
	return nil
}

// payload:  {"object_ids": ["id1", "id2"]}
type ObjectBulkDeletePayload struct {
	ObjectIds []string `json:"object_ids"`
}

func (o *objectsService) BulkDelete(ctx context.Context, objectType string, payload ObjectBulkDeletePayload) error {
	urlStr := fmt.Sprintf("%s%s/", o._bulkUrl, url.PathEscape(objectType))
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("DELETE", urlStr, payload)
	if err != nil {
		return err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	//
	err = o.client.parseApiResponse(httpResponse, nil)
	if err != nil {
		return err
	}
	return nil
}

func (o *objectsService) GetSubscriptions(ctx context.Context, obj ObjectIdentifier, opts *CursorListApiOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%ssubscription/", o.objectDetailAPIUrl(obj.ObjectType, obj.Id)), opts.BuildQuery())
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
	resp := &CursorListApiResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Create Subscription Payload
//
//	payload: {
//		"recipients": ["distinct_id1", {"object_type": "type1", "id": "id1"},],
//		"properties": {"type": "admin"},
//		"parent_object_properties: {}, // if value non-null, does upsert on parent-object too.
//	}
func (o *objectsService) CreateSubscriptions(ctx context.Context, obj ObjectIdentifier, payload map[string]any) (map[string]any, error) {
	urlStr := fmt.Sprintf("%ssubscription/", o.objectDetailAPIUrl(obj.ObjectType, obj.Id))
	if payload == nil {
		payload = map[string]any{}
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
	//
	resp := map[string]any{}
	err = o.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete Subscription Payload
//
//	payload: {
//		"recipients": ["distinct_id1", {"object_type": "type1", "id": "id1"},]
//	}
func (o *objectsService) DeleteSubscriptions(ctx context.Context, obj ObjectIdentifier, payload map[string]any) error {
	urlStr := fmt.Sprintf("%ssubscription/", o.objectDetailAPIUrl(obj.ObjectType, obj.Id))
	if payload == nil {
		payload = map[string]any{}
	}
	// prepare http.Request object
	request, err := o.client.prepareHttpRequest("DELETE", urlStr, payload)
	if err != nil {
		return err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	//
	err = o.client.parseApiResponse(httpResponse, nil)
	if err != nil {
		return err
	}
	return nil
}

func (o *objectsService) GetObjectsSubscribedTo(ctx context.Context, obj ObjectIdentifier, opts *CursorListApiOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%ssubscribed_to/object/", o.objectDetailAPIUrl(obj.ObjectType, obj.Id)), opts.BuildQuery())
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
	resp := &CursorListApiResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) GetEditInstance(obj ObjectIdentifier) ObjectEdit {
	return newObjectEdit(o.client, obj)
}
