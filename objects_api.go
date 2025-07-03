package suprsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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
	GetPreference(context.Context, ObjectIdentifier, *ObjectPreferenceOptions) (*ObjectPreferenceResponse, error)
	GetAllCategoriesPreference(context.Context, ObjectIdentifier, *ObjectPreferenceOptions) (*CursorListApiResponse, error)
	GetGlobalChannelsPreferences(context.Context, ObjectIdentifier, *ObjectGlobalPreferenceOptions) (*ObjectGlobalChannelPreferencesResponse, error)
	GetCategoryPreference(context.Context, ObjectIdentifier, string, *ObjectCategoryPreferenceOptions) (*ObjectCategoryPreferenceResponse, error)
	UpdateCategoryPreference(context.Context, ObjectIdentifier, string, ObjectUpdateCategoryPreferenceBody, *ObjectCategoryUpdatePreferenceOptions) (*ObjectCategoryPreferenceResponse, error)
	UpdateGlobalChannelsPreferences(context.Context, ObjectIdentifier, ObjectGlobalChannelPreferenceUpdateBody, *ObjectGlobalPreferenceOptions) (*ObjectGlobalChannelPreferencesResponse, error)
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

type ObjectPreferenceOptions struct {
	TenantId           string `json:"tenant_id"`
	ShowOptOutChannels bool   `json:"show_opt_out_channels"`
	Tags               string `json:"tags"`
}

type ObjectGlobalPreferenceOptions struct {
	TenantId string `json:"tenant_id"`
}

type ObjectCategoryPreferenceOptions struct {
	TenantId           string `json:"tenant_id"`
	ShowOptOutChannels bool   `json:"show_opt_out_channels"`
}

type ObjectCategoryUpdatePreferenceOptions struct {
	TenantId string `json:"tenant_id"`
}

type ObjectUpdateCategoryPreferenceBody struct {
	Preference     string   `json:"preference"`
	OptOutChannels []string `json:"opt_out_channels"`
}

type ObjectGlobalChannelPreferenceUpdateBody struct {
	ChannelPreferences []UserChannelPreferenceIn `json:"channel_preferences"`
}

type ObjectPreferenceResponse struct {
	Sections           []any `json:"sections"`
	ChannelPreferences []any `json:"channel_preferences"`
}

type ObjectGlobalChannelPreferencesResponse struct {
	ChannelPreferences []any `json:"channel_preferences"`
}

type ObjectCategoryPreferenceResponse struct {
	Name               string `json:"name"`
	Category           string `json:"category"`
	Description        string `json:"description"`
	OriginalPreference string `json:"original_preference"`
	Preference         string `json:"preference"`
	IsEditable         bool   `json:"is_editable"`
	Channels           []any  `json:"channels"`
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

func (o *objectsService) GetPreference(ctx context.Context, obj ObjectIdentifier, opts *ObjectPreferenceOptions) (*ObjectPreferenceResponse, error) {
	if strings.TrimSpace(obj.Id) == "" {
		return nil, &Error{Message: "object_id is required"}
	}
	if strings.TrimSpace(obj.ObjectType) == "" {
		return nil, &Error{Message: "object_type is required"}
	}

	urlStr := fmt.Sprintf("%s%s/%s/preference/", o._url, url.PathEscape(strings.TrimSpace(obj.ObjectType)), url.PathEscape(strings.TrimSpace(obj.Id)))

	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.Tags != "" {
			query.Set("tags", opts.Tags)
		}
		if strconv.FormatBool(opts.ShowOptOutChannels) != "" {
			query.Set("show_opt_out_channels", strconv.FormatBool(opts.ShowOptOutChannels))
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &ObjectPreferenceResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) GetAllCategoriesPreference(ctx context.Context, obj ObjectIdentifier, opts *ObjectPreferenceOptions) (*CursorListApiResponse, error) {
	if strings.TrimSpace(obj.Id) == "" {
		return nil, &Error{Message: "object_id is required"}
	}
	if strings.TrimSpace(obj.ObjectType) == "" {
		return nil, &Error{Message: "object_type is required"}
	}

	urlStr := fmt.Sprintf("%s%s/%s/preference/category", o._url, url.PathEscape(strings.TrimSpace(obj.ObjectType)), url.PathEscape(strings.TrimSpace(obj.Id)))

	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.Tags != "" {
			query.Set("tags", opts.Tags)
		}
		if strconv.FormatBool(opts.ShowOptOutChannels) != "" {
			query.Set("show_opt_out_channels", strconv.FormatBool(opts.ShowOptOutChannels))
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &CursorListApiResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) GetGlobalChannelsPreferences(ctx context.Context, obj ObjectIdentifier, opts *ObjectGlobalPreferenceOptions) (*ObjectGlobalChannelPreferencesResponse, error) {
	if strings.TrimSpace(obj.Id) == "" {
		return nil, &Error{Message: "object_id is required"}
	}
	if strings.TrimSpace(obj.ObjectType) == "" {
		return nil, &Error{Message: "object_type is required"}
	}

	urlStr := fmt.Sprintf("%s%s/%s/preference/channel_preference", o._url, url.PathEscape(strings.TrimSpace(obj.ObjectType)), url.PathEscape(strings.TrimSpace(obj.Id)))

	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &ObjectGlobalChannelPreferencesResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) GetCategoryPreference(ctx context.Context, obj ObjectIdentifier, category string, opts *ObjectCategoryPreferenceOptions) (*ObjectCategoryPreferenceResponse, error) {
	if strings.TrimSpace(obj.Id) == "" {
		return nil, &Error{Message: "object_id is required"}
	}
	if strings.TrimSpace(obj.ObjectType) == "" {
		return nil, &Error{Message: "object_type is required"}
	}
	if strings.TrimSpace(category) == "" {
		return nil, &Error{Message: "category is required"}
	}

	urlStr := fmt.Sprintf("%s%s/%s/preference/category/%s", o._url, url.PathEscape(strings.TrimSpace(obj.ObjectType)), url.PathEscape(strings.TrimSpace(obj.Id)), url.PathEscape(strings.TrimSpace(category)))

	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if strconv.FormatBool(opts.ShowOptOutChannels) != "" {
			query.Set("show_opt_out_channels", strconv.FormatBool(opts.ShowOptOutChannels))
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := o.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &ObjectCategoryPreferenceResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) UpdateCategoryPreference(ctx context.Context, obj ObjectIdentifier, category string, body ObjectUpdateCategoryPreferenceBody, opts *ObjectCategoryUpdatePreferenceOptions) (*ObjectCategoryPreferenceResponse, error) {
	if strings.TrimSpace(obj.Id) == "" {
		return nil, &Error{Message: "object_id is required"}
	}
	if strings.TrimSpace(obj.ObjectType) == "" {
		return nil, &Error{Message: "object_type is required"}
	}
	if strings.TrimSpace(category) == "" {
		return nil, &Error{Message: "category is required"}
	}

	urlStr := fmt.Sprintf("%s%s/%s/preference/category/%s", o._url, url.PathEscape(strings.TrimSpace(obj.ObjectType)), url.PathEscape(strings.TrimSpace(obj.Id)), url.PathEscape(strings.TrimSpace(category)))

	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := o.client.prepareHttpRequest("PATCH", urlStr, body)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &ObjectCategoryPreferenceResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *objectsService) UpdateGlobalChannelsPreferences(ctx context.Context, obj ObjectIdentifier, body ObjectGlobalChannelPreferenceUpdateBody, opts *ObjectGlobalPreferenceOptions) (*ObjectGlobalChannelPreferencesResponse, error) {
	if strings.TrimSpace(obj.Id) == "" {
		return nil, &Error{Message: "object_id is required"}
	}
	if strings.TrimSpace(obj.ObjectType) == "" {
		return nil, &Error{Message: "object_type is required"}
	}

	urlStr := fmt.Sprintf("%s%s/%s/preference/channel_preference", o._url, url.PathEscape(strings.TrimSpace(obj.ObjectType)), url.PathEscape(strings.TrimSpace(obj.Id)))

	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	request, err := o.client.prepareHttpRequest("PATCH", urlStr, body)
	if err != nil {
		return nil, err
	}
	httpResponse, err := o.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	resp := &ObjectGlobalChannelPreferencesResponse{}
	err = o.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
