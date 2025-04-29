package suprsend

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

var RESERVED_EVENT_NAMES = []string{
	"$identify",
	"$notification_delivered", "$notification_dismiss", "$notification_clicked",
	"$app_launched", "$user_login", "$user_logout",
}

type Event struct {
	DistinctId     string
	EventName      string
	Properties     map[string]interface{}
	IdempotencyKey string
	TenantId       string
	// Brand has been renamed to Tenant. Brand is kept for backward-compatibilty.
	// Use Tenant instead of Brand
	BrandId string
}

func (e *Event) validateDistinctId() error {
	e.DistinctId = strings.TrimSpace(e.DistinctId)
	if e.DistinctId == "" {
		return errors.New("distinct_id missing")
	}
	return nil
}

func (e *Event) checkProperties() {
	if e.Properties == nil {
		e.Properties = map[string]interface{}{}
	}
}

func (e *Event) checkEventPrefix() error {
	if !slices.Contains(RESERVED_EVENT_NAMES, e.EventName) {
		if strings.HasPrefix(e.EventName, "$") || strings.HasPrefix(e.EventName, "ss_") ||
			strings.HasPrefix(e.EventName, "SS_") {
			return errors.New("event_names starting with [$,ss_] are reserved by SuprSend")
		}
	}
	return nil
}

func (e *Event) validateEventName() error {
	e.EventName = strings.TrimSpace(e.EventName)
	if e.EventName == "" {
		return errors.New("event_name missing")
	}
	err := e.checkEventPrefix()
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) AddAttachment(filePath string, ao *AttachmentOption) error {
	e.checkProperties()
	attachment, err := GetAttachmentJson(filePath, ao)
	if err != nil {
		return err
	}
	if attachment == nil {
		return nil
	}
	// add the attachment to properties->$attachments
	if a, found := e.Properties["$attachments"]; !found || a == nil {
		e.Properties["$attachments"] = []map[string]interface{}{}
	}
	allAttachments := e.Properties["$attachments"].([]map[string]interface{})
	allAttachments = append(allAttachments, attachment)
	e.Properties["$attachments"] = allAttachments
	//
	return nil
}

func (e *Event) getFinalJson(client *Client, isPartOfBulk bool) (map[string]interface{}, int, error) {
	var err error
	err = e.validateDistinctId()
	if err != nil {
		return nil, 0, err
	}
	err = e.validateEventName()
	if err != nil {
		return nil, 0, err
	}
	e.checkProperties()
	//
	suprProps := map[string]interface{}{"$ss_sdk_version": client.userAgent}
	// props
	maps.Copy(e.Properties, suprProps)
	//
	eventMap := map[string]interface{}{
		"$insert_id":  uuid.New().String(),
		"$time":       time.Now().UnixMilli(),
		"event":       e.EventName,
		"env":         client.ApiKey,
		"distinct_id": e.DistinctId,
		"properties":  e.Properties,
	}
	// Add idempotency_key if present
	if e.IdempotencyKey != "" {
		eventMap["$idempotency_key"] = e.IdempotencyKey
	}
	// Add tenant_id if present
	if e.TenantId != "" {
		eventMap["tenant_id"] = e.TenantId
	}
	if e.BrandId != "" {
		eventMap["brand_id"] = e.BrandId
	}
	eventMap, err = validateTrackEventSchema(eventMap)
	if err != nil {
		return nil, 0, err
	}
	// Check request size
	apparentSize, err := getApparentEventSize(eventMap, isPartOfBulk)
	if err != nil {
		return nil, 0, err
	}
	if apparentSize > SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("event size too big - %d Bytes, must not cross %s", apparentSize,
			SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, errors.New(errStr)
	}
	return eventMap, apparentSize, nil
}

func (e *Event) asJson() map[string]interface{} {
	eventMap := map[string]interface{}{
		"event":       e.EventName,
		"distinct_id": e.DistinctId,
		"properties":  e.Properties,
	}
	// Add idempotency_key if present
	if e.IdempotencyKey != "" {
		eventMap["$idempotency_key"] = e.IdempotencyKey
	}
	// Add tenant_id if present
	if e.TenantId != "" {
		eventMap["tenant_id"] = e.TenantId
	}
	if e.BrandId != "" {
		eventMap["brand_id"] = e.BrandId
	}
	return eventMap
}

type eventsCollector struct {
	client *Client
	_url   string
}

func newEventCollectorInstance(client *Client) *eventsCollector {
	ec := &eventsCollector{
		client: client,
		// events url
		_url: fmt.Sprintf("%sevent/", client.baseUrl),
	}
	return ec
}

func (e *eventsCollector) Collect(event *Event) (*Response, error) {
	eventMap, _, err := event.getFinalJson(e.client, false)
	if err != nil {
		return nil, err
	}
	suprResp, err := e.send(eventMap)
	if err != nil {
		return nil, err
	}
	return suprResp, nil
}

func (e *eventsCollector) send(eventMap map[string]interface{}) (*Response, error) {
	// prepare http.Request object
	request, err := e.client.prepareHttpRequest("POST", e._url, eventMap)
	if err != nil {
		return nil, err
	}
	//
	httpResponse, err := e.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	suprResponse, err := e.formatAPIResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return suprResponse, nil
}

func (e *eventsCollector) formatAPIResponse(httpRes *http.Response) (*Response, error) {
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpRes.StatusCode, string(respBody))

	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}
