package suprsend

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SubscriberList struct {
	ListId          string `json:"list_id,omitempty"`
	ListName        string `json:"list_name,omitempty"`
	ListDescription string `json:"list_description,omitempty"`
	ListType        string `json:"list_type,omitempty"`
	//
	SubscribersCount int    `json:"subscribers_count,omitempty"`
	Source           string `json:"source,omitempty"`
	IsReadonly       bool   `json:"is_readonly,omitempty"`
	Status           string `json:"status,omitempty"`
	//
	TrackUserEntry bool `json:"track_user_entry,omitempty"`
	TrackUserExit  bool `json:"track_user_exit,omitempty"`
	//
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	// version_id will be present its a draft version
	VersionId string `json:"version_id,omitempty"`
	// drafts will be present if there are any drafts started from this list
	Drafts []*SubscriberListVersion `json:"drafts,omitempty"`
}

type SubscriberListVersion struct {
	VersionId        string `json:"version_id,omitempty"`
	SubscribersCount int    `json:"subscribers_count,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// GetAll response
type SubscriberListAll struct {
	Meta    *ListApiMetaInfo  `json:"meta"`
	Results []*SubscriberList `json:"results"`
}

// GetAll options
type SubscriberListAllOptions struct {
	Limit  int
	Offset int
}

func (b *SubscriberListAllOptions) cleanParams() {
	// limit must be 0 < x <= 1000
	if b.Limit <= 0 || b.Limit > 1000 {
		b.Limit = 20
	}
	if b.Offset < 0 {
		b.Offset = 0
	}
}

// Create subscriberlist request input
type SubscriberListCreateInput struct {
	ListId          string `json:"list_id,omitempty"`
	ListName        string `json:"list_name,omitempty"`
	ListDescription string `json:"list_description,omitempty"`
	//
	TrackUserEntry *bool `json:"track_user_entry,omitempty"`
	TrackUserExit  *bool `json:"track_user_exit,omitempty"`
	// list_type enums: query_based, static_list
	ListType *string `json:"list_type,omitempty"`
	// Query: applicable when list_type='query_based'
	Query *string `json:"query,omitempty"`
}

// Broadcast request params on SubscriberList
type SubscriberListBroadcast struct {
	Body           map[string]any
	IdempotencyKey string
	TenantId       string
	// Brand has been renamed to Tenant. BrandId is kept for backward-compatibilty.
	// Use TenantId instead of BrandId
	BrandId string
}

func (s *SubscriberListBroadcast) AddAttachment(filePath string, ao *AttachmentOption) error {
	if d, found := s.Body["data"]; !found || d == nil {
		s.Body["data"] = map[string]any{}
	}
	attachment, err := GetAttachmentJson(filePath, ao)
	if err != nil {
		return err
	}
	if attachment == nil {
		return nil
	}
	data := s.Body["data"].(map[string]any)
	if a, found := data["$attachments"]; !found || a == nil {
		data["$attachments"] = []map[string]any{}
	}
	allAttachments := data["$attachments"].([]map[string]any)
	allAttachments = append(allAttachments, attachment)
	data["$attachments"] = allAttachments
	return nil
}

func (s *SubscriberListBroadcast) getFinalJson(client *Client) (map[string]any, int, error) {
	s.Body["$insert_id"] = uuid.New().String()
	s.Body["$time"] = time.Now().UnixMilli()
	if s.IdempotencyKey != "" {
		s.Body["$idempotency_key"] = s.IdempotencyKey
	}
	if s.TenantId != "" {
		s.Body["tenant_id"] = s.TenantId
	}
	if s.BrandId != "" {
		s.Body["brand_id"] = s.BrandId
	}
	body, err := validateListBroadcastBodySchema(s.Body)
	if err != nil {
		return nil, 0, err
	}
	s.Body = body
	// Check request size
	apparentSize, err := getApparentListBroadcastBodySize(body)
	if err != nil {
		return nil, 0, err
	}
	if apparentSize > BODY_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("SubscriberListBroadcast body too big - %d Bytes, must not cross %s", apparentSize,
			BODY_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, &Error{Code: 413, Message: errStr}
	}
	return s.Body, apparentSize, nil
}
