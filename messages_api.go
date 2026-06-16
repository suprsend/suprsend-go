package suprsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type MessagesService interface {
	List(context.Context, *MessageListOptions) (*CursorListApiResponse, error)
	BulkUpdate(context.Context, []MessageUpdateItem) (*MessageBulkUpdateResponse, error)
}

type messagesService struct {
	client   *Client
	_url     string
	_bulkUrl string
}

var _ MessagesService = &messagesService{}

func newMessagesService(client *Client) *messagesService {
	return &messagesService{
		client:   client,
		_url:     fmt.Sprintf("%sv1/message/", client.baseUrl),
		_bulkUrl: fmt.Sprintf("%sv1/bulk/message/", client.baseUrl),
	}
}

// MessageListOptions holds filter and pagination params for the message list API.
// Multi-value fields (RecipientID, Status, Category) accept one or more values.
type MessageListOptions struct {
	Limit          int
	Before         string
	After          string
	RecipientID    []string
	Status         []string
	MessageID      string
	IdempotencyKey string
	TenantID       string
	WorkflowSlug   string
	Channel        string
	ExecutionID    string
	CreatedAtGte   string // RFC3339 timestamp
	CreatedAtLte   string // RFC3339 timestamp
	ObjectID       string // must be paired with ObjectType
	ObjectType     string // must be paired with ObjectID
	IsCampaign     *bool
	Category       []string
}

func (o *MessageListOptions) buildQuery() string {
	if o == nil {
		return ""
	}
	params := url.Values{}
	if o.Limit > 0 {
		params.Add("limit", strconv.Itoa(o.Limit))
	}
	if o.Before != "" {
		params.Add("before", o.Before)
	}
	if o.After != "" {
		params.Add("after", o.After)
	}
	for _, v := range o.RecipientID {
		params.Add("recipient_id[]", v)
	}
	for _, v := range o.Status {
		params.Add("status[]", v)
	}
	if o.MessageID != "" {
		params.Add("message_id", o.MessageID)
	}
	if o.IdempotencyKey != "" {
		params.Add("idempotency_key", o.IdempotencyKey)
	}
	if o.TenantID != "" {
		params.Add("tenant_id", o.TenantID)
	}
	if o.WorkflowSlug != "" {
		params.Add("workflow_slug", o.WorkflowSlug)
	}
	if o.Channel != "" {
		params.Add("channel", o.Channel)
	}
	if o.ExecutionID != "" {
		params.Add("execution_id", o.ExecutionID)
	}
	if o.CreatedAtGte != "" {
		params.Add("created_at_gte", o.CreatedAtGte)
	}
	if o.CreatedAtLte != "" {
		params.Add("created_at_lte", o.CreatedAtLte)
	}
	if o.ObjectID != "" {
		params.Add("object_id", o.ObjectID)
	}
	if o.ObjectType != "" {
		params.Add("object_type", o.ObjectType)
	}
	if o.IsCampaign != nil {
		params.Add("is_campaign", strconv.FormatBool(*o.IsCampaign))
	}
	for _, v := range o.Category {
		params.Add("category[]", v)
	}
	return params.Encode()
}

func (m *messagesService) List(ctx context.Context, opts *MessageListOptions) (*CursorListApiResponse, error) {
	urlStr := appendQueryParamPart(m._url, opts.buildQuery())
	request, err := m.client.prepareHttpRequest(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := m.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &CursorListApiResponse{}
	err = m.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MessageUpdateItem is a single message status update in a bulk update request.
type MessageUpdateItem struct {
	MessageID string `json:"message_id"`
	Action    string `json:"action"`
}

// MessageUpdateRecord is the per-item result from a bulk update.
type MessageUpdateRecord struct {
	MessageID  string `json:"message_id"`
	StatusCode int    `json:"status_code"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// MessageBulkUpdateResponse is the response from the bulk update API.
type MessageBulkUpdateResponse struct {
	Records []MessageUpdateRecord `json:"records"`
}

func (m *messagesService) BulkUpdate(ctx context.Context, messages []MessageUpdateItem) (*MessageBulkUpdateResponse, error) {
	payload := map[string]any{"messages": messages}
	request, err := m.client.prepareHttpRequest(ctx, "PATCH", m._bulkUrl, payload)
	if err != nil {
		return nil, err
	}
	httpResponse, err := m.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	resp := &MessageBulkUpdateResponse{}
	err = m.client.parseApiResponse(httpResponse, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// // MessageContent holds the rendered content fields for a message.
// // Which fields are populated depends on the channel.
// type MessageContent struct {
// 	Subject  string `json:"subject"`
// 	Body     string `json:"body"`
// 	Title    string `json:"title"`    // push channels
// 	Subtitle string `json:"subtitle"` // push channels
// 	Data     string `json:"data"`     // push channels, custom JSON string
// }

// // MessageContentResponse is the response from the message content API.
// type MessageContentResponse struct {
// 	NotificationID string          `json:"notification_id"`
// 	Channel        string          `json:"channel"`
// 	RenderedAt     *time.Time      `json:"rendered_at"`
// 	Content        *MessageContent `json:"content"`
// }

// func (m *messagesService) messageContentURL(messageID string) string {
// 	return fmt.Sprintf("%s/%s/content", m._url, url.PathEscape(strings.TrimSpace(messageID)))
// }

// func (m *messagesService) GetContent(ctx context.Context, messageID string) (*MessageContentResponse, error) {
// 	urlStr := m.messageContentURL(messageID)
// 	request, err := m.client.prepareHttpRequest(ctx, "GET", urlStr, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	httpResponse, err := m.client.httpClient.Do(request)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer httpResponse.Body.Close()
// 	resp := &MessageContentResponse{}
// 	err = m.client.parseApiResponse(httpResponse, resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }
