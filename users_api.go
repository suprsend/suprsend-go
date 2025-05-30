package suprsend

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type UsersService interface {
	List(context.Context, *CursorListApiOptions) (*CursorListApiResponse, error)
	Get(context.Context, string) (map[string]any, error)
	Upsert(context.Context, string, map[string]any) (map[string]any, error)
	AsyncEdit(context.Context, UserEdit) (*Response, error)
	Edit(context.Context, UserEditRequest) (map[string]any, error)
	Merge(context.Context, string, UserMergeRequest) (map[string]any, error)
	Delete(context.Context, string) error
	BulkDelete(context.Context, UserBulkDeletePayload) error
	GetObjectsSubscribedTo(context.Context, string, *CursorListApiOptions) (*CursorListApiResponse, error)
	GetListsSubscribedTo(context.Context, string, *CursorListApiOptions) (*CursorListApiResponse, error)
	//
	GetEditInstance(string) UserEdit
	GetBulkEditInstance() BulkUsersEdit
	// Old accessor method (to be deprecated)
	GetInstance(string) Subscriber
	GetUserPreferences(context.Context, string, *UserPreferencesOptions) (*UserPreferencesResponse, error)
}

type usersService struct {
	client   *Client
	_url     string
	_bulkUrl string
}

var _ UsersService = &usersService{}

func newUsersService(client *Client) *usersService {
	us := &usersService{
		client:   client,
		_url:     fmt.Sprintf("%sv1/user/", client.baseUrl),
		_bulkUrl: fmt.Sprintf("%sv1/bulk/user/", client.baseUrl),
	}
	return us
}

func (u *usersService) List(ctx context.Context, opts *CursorListApiOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(u._url, opts.BuildQuery())
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &CursorListApiResponse{}
	err = u.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) userDetailAPIUrl(distinctId string) string {
	return fmt.Sprintf(
		"%s%s/",
		u._url,
		url.PathEscape(strings.TrimSpace(distinctId)),
	)
}

func (u *usersService) Get(ctx context.Context, distinctId string) (map[string]any, error) {
	urlStr := u.userDetailAPIUrl(distinctId)
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := map[string]any{}
	err = u.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) Upsert(ctx context.Context, distinctId string, payload map[string]any) (map[string]any, error) {
	urlStr := u.userDetailAPIUrl(distinctId)
	if payload == nil {
		payload = map[string]any{}
	}
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := map[string]any{}
	err = u.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) AsyncEdit(ctx context.Context, editInstance UserEdit) (*Response, error) {
	ue := editInstance.(*userEdit)
	ue.validateBody()
	payload := ue.GetAsyncPayload()
	_, _, err := ue.validatePayloadSize(payload)
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%sevent/", u.client.baseUrl)
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp, err := u.asyncAPIResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) asyncAPIResponse(httpRes *http.Response) (*Response, error) {
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, &Error{Err: err}
	}
	if httpRes.StatusCode >= 400 {
		return nil, &Error{Code: httpRes.StatusCode, Message: string(respBody)}
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}

// Either (distinct_id & payload) OR edit_instance should be provided
type UserEditRequest struct {
	DistinctId string
	// {"operations": [{"$set": {"prop1": "val1"}, {"$append": {"$email": "abc@test.com"}}]}
	Payload map[string]any
	//
	EditInstance UserEdit
}

func (u *usersService) Edit(ctx context.Context, req UserEditRequest) (map[string]any, error) {
	var urlStr string
	var payload map[string]any
	if req.EditInstance != nil {
		ue := req.EditInstance.(*userEdit)
		ue.validateBody()
		payload = ue.GetPayload()
		urlStr = u.userDetailAPIUrl(ue.distinctId)
	} else {
		payload = req.Payload
		if payload == nil {
			payload = map[string]any{}
		}
		urlStr = u.userDetailAPIUrl(req.DistinctId)
	}
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("PATCH", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := map[string]any{}
	err = u.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type UserMergeRequest struct {
	FromUserId string `json:"from_user_id"`
}

func (u *usersService) Merge(ctx context.Context, distinctId string, payload UserMergeRequest) (map[string]any, error) {
	urlStr := fmt.Sprintf("%smerge/", u.userDetailAPIUrl(distinctId))
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("POST", urlStr, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := map[string]any{}
	err = u.client.parseApiResponse(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) Delete(ctx context.Context, distinctId string) error {
	urlStr := u.userDetailAPIUrl(distinctId)
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("DELETE", urlStr, nil)
	if err != nil {
		return err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	//
	err = u.client.parseApiResponse(httpResponse, nil)
	if err != nil {
		return err
	}
	return nil
}

type UserBulkDeletePayload struct {
	DistinctIds []string `json:"distinct_ids"`
}

// payload: {"distinct_ids": ["id1", "id2"]}
func (u *usersService) BulkDelete(ctx context.Context, payload UserBulkDeletePayload) error {
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("DELETE", u._bulkUrl, payload)
	if err != nil {
		return err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	//
	err = u.client.parseApiResponse(httpResponse, nil)
	if err != nil {
		return err
	}
	return nil
}

func (u *usersService) GetObjectsSubscribedTo(ctx context.Context, distinctId string, opts *CursorListApiOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%ssubscribed_to/object/", u.userDetailAPIUrl(distinctId)), opts.BuildQuery())
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &CursorListApiResponse{}
	err = u.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) GetListsSubscribedTo(ctx context.Context, distinctId string, opts *CursorListApiOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(fmt.Sprintf("%ssubscribed_to/list/", u.userDetailAPIUrl(distinctId)), opts.BuildQuery())
	// prepare http.Request object
	request, err := u.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	//
	resp := &CursorListApiResponse{}
	err = u.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *usersService) GetEditInstance(distinctId string) UserEdit {
	return newUserEdit(u.client, distinctId)
}

func (u *usersService) GetBulkEditInstance() BulkUsersEdit {
	return newBulkUsersEdit(u.client)
}

// Deprecated: this method will be removed in near future. Use GetEditInstance instead.
func (u *usersService) GetInstance(distinctId string) Subscriber {
	return newSubscriber(u.client, distinctId)
}

type UserPreferencesResponse struct {
	Sections           []any `json:"sections"`
	ChannelPreferences []any `json:"channel_preferences"`
}

// GetUserPreferences fetches the current notification preferences for the user across all categories and channels.
func (u *usersService) GetUserPreferences(ctx context.Context, distinctId string, opts *UserPreferencesOptions) (*UserPreferencesResponse, error) {
	if strings.TrimSpace(distinctId) == "" {
		return nil, &Error{Message: "distinct_id is required"}
	}
	urlStr := fmt.Sprintf("%s%s/preference/", u._url, url.PathEscape(strings.TrimSpace(distinctId)))
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.ShowOptOutChannels != nil {
			query.Set("show_opt_out_channels", fmt.Sprintf("%v", *opts.ShowOptOutChannels))
		}
		if opts.Tags != "" {
			query.Set("tags", opts.Tags)
		}
	}
	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}
	request, err := u.client.prepareHttpRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := u.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &UserPreferencesResponse{}
	err = u.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type UserPreferencesOptions struct {
	TenantId           string
	ShowOptOutChannels *bool
	Tags               string // can be a simple tag or a JSON string for advanced filtering
}
