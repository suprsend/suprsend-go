package suprsend

import (
	"errors"
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
	Source           string `json:"source,omitempty,omitempty"`
	IsReadonly       bool   `json:"is_readonly,omitempty"`
	Status           string `json:"status,omitempty"`
	//
	TrackUserEntry bool `json:"track_user_entry,omitempty"`
	TrackUserExit  bool `json:"track_user_exit,omitempty"`
	//
	RequestedForDelete bool   `json:"requested_for_delete"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	//
	Drafts []SubscriberListVersion `json:"drafts,omitempty"`
}

type SubscriberListVersion struct {
	ListId          string `json:"list_id,omitempty"`
	ListName        string `json:"list_name,omitempty"`
	ListDescription string `json:"list_description,omitempty"`
	ListType        string `json:"list_type,omitempty"`
	//
	SubscribersCount int    `json:"subscribers_count,omitempty"`
	Source           string `json:"source,omitempty,omitempty"`
	IsReadonly       bool   `json:"is_readonly,omitempty"`
	Status           string `json:"status,omitempty"`
	//
	TrackUserEntry bool `json:"track_user_entry,omitempty"`
	TrackUserExit  bool `json:"track_user_exit,omitempty"`
	//
	RequestedForDelete bool   `json:"requested_for_delete,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	//
	VersionId string `json:"version_id,omitempty"`
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
	ListId          string  `json:"list_id,omitempty"`
	ListName        string  `json:"list_name,omitempty"`
	ListDescription string  `json:"list_description,omitempty"`
	ListType        *string `json:"list_type"`
	Query           string  `json:"query"`
	Source          string  `json:"source"`
	TrackUserEntry  bool    `json:"track_user_entry"`
	TrackUserExit   bool    `json:"track_user_exit"`
}

// Broadcast request params on SubscriberList
type SubscriberListBroadcast struct {
	Body           map[string]interface{}
	IdempotencyKey string
	BrandId        string
}

func (s *SubscriberListBroadcast) AddAttachment(filePath string, ao *AttachmentOption) error {
	if d, found := s.Body["data"]; !found || d == nil {
		s.Body["data"] = map[string]interface{}{}
	}
	attachment, err := GetAttachmentJson(filePath, ao)
	if err != nil {
		return err
	}
	data := s.Body["data"].(map[string]interface{})
	if a, found := data["$attachments"]; !found || a == nil {
		data["$attachments"] = []map[string]interface{}{}
	}
	allAttachments := data["$attachments"].([]map[string]interface{})
	allAttachments = append(allAttachments, attachment)
	data["$attachments"] = allAttachments
	return nil
}

func (s *SubscriberListBroadcast) getFinalJson(client *Client) (map[string]interface{}, int, error) {
	s.Body["$insert_id"] = uuid.New().String()
	s.Body["$time"] = time.Now().UnixMilli()
	if s.IdempotencyKey != "" {
		s.Body["$idempotency_key"] = s.IdempotencyKey
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
	if apparentSize > SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("SubscriberListBroadcast body too big - %d Bytes, must not cross %s", apparentSize,
			SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, errors.New(errStr)
	}
	return s.Body, apparentSize, nil
}
