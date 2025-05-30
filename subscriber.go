package suprsend

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscriber interface {
	Save() (*Response, error)
	//
	AppendKV(string, any)
	Append(map[string]any)
	SetKV(string, any)
	Set(map[string]any)
	SetOnceKV(string, any)
	SetOnce(map[string]any)
	IncrementKV(string, any)
	Increment(map[string]any)
	RemoveKV(string, any)
	Remove(map[string]any)
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
	AddWebpush(value map[string]any, provider string)
	RemoveWebpush(value map[string]any, provider string)
	//
	AddSlack(value map[string]any)
	RemoveSlack(value map[string]any)
	//
	AddMSTeams(value map[string]any)
	RemoveMSTeams(value map[string]any)
}

var _ Subscriber = &subscriber{}

type subscriber struct {
	client     *Client
	distinctId string
	_url       string
	//
	_errors        []string
	_warnings      []string
	userOperations []map[string]any
	//
	_helper *subscriberHelper
	//
	_warningsList []string
}

func newSubscriber(client *Client, distinctId string) Subscriber {
	s := &subscriber{
		client:     client,
		distinctId: distinctId,
	}
	// events url
	s._url = fmt.Sprintf("%sevent/", client.baseUrl)
	s._helper = newSubscriberHelper()
	return s
}

func (s *subscriber) getEvent() map[string]any {
	return map[string]any{
		"$schema":          "2",
		"$insert_id":       uuid.New().String(),
		"$time":            time.Now().UnixMilli(),
		"env":              s.client.getWsIdentifierValue(),
		"distinct_id":      s.distinctId,
		"$user_operations": s.userOperations,
		"properties":       map[string]any{"$ss_sdk_version": s.client.userAgent},
	}
}

func (s *subscriber) asJson() map[string]any {
	return map[string]any{
		"distinct_id":      s.distinctId,
		"$user_operations": s.userOperations,
		"warnings":         s._warningsList,
	}
}

func (s *subscriber) validateEventSize(event map[string]any) (map[string]any, int, error) {
	apparentSize, err := getApparentIdentityEventSize(event)
	if err != nil {
		return event, 0, err
	}
	if apparentSize > IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("User Event size too big - %d Bytes, must not cross %s", apparentSize,
			IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, &Error{Code: 413, Message: errStr}
	}
	return event, apparentSize, nil
}

func (s *subscriber) validateBody(isPartOfBulk bool) ([]string, error) {
	if s.distinctId == "" {
		s._errors = append(s._errors, "missing distinct_id")
	}
	s._warningsList = []string{}
	if len(s._warnings) > 0 {
		msg := fmt.Sprintf("[distinct_id: %s] %s", s.distinctId, strings.Join(s._warnings, "\n"))
		s._warningsList = append(s._warningsList, msg)
		// print on console as well
		log.Println("WARNING:", msg)
	}
	if len(s._errors) > 0 {
		msg := fmt.Sprintf("[distinct_id: %s] %s", s.distinctId, strings.Join(s._errors, "\n"))
		s._warningsList = append(s._warningsList, msg)
		errMsg := "ERROR: " + msg
		if isPartOfBulk {
			// print on console in case of bulk-api
			log.Println(errMsg)
		} else {
			// raise error in case of single api.
			// return nil, &Error{Message: errMsg} // Removed exception throwing, let backend handle it
			log.Println(errMsg)
		}
	}
	return s._warningsList, nil
}

func (s *subscriber) Save() (*Response, error) {
	if _, err := s.validateBody(false); err != nil {
		return nil, err
	}
	//
	event := s.getEvent()
	if _, _, err := s.validateEventSize(event); err != nil {
		return nil, err
	}
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("POST", s._url, event)
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

func (s *subscriber) formatAPIResponse(httpRes *http.Response) (*Response, error) {
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, &Error{Err: err}
	}
	if httpRes.StatusCode >= 400 {
		return nil, &Error{Code: httpRes.StatusCode, Message: string(respBody)}
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}

func (s *subscriber) _collectEvent() {
	resp := s._helper.getIdentityEvent()
	if len(resp.errors) > 0 {
		s._errors = append(s._errors, resp.errors...)
	}
	if len(resp.info) > 0 {
		s._warnings = append(s._warnings, resp.info...)
	}
	if len(resp.event) > 0 {
		s.userOperations = append(s.userOperations, resp.event)
	}
}

/*
Usage:
 1. append(k, v)
 2. append({k1: v1, k2: v2})
*/
func (s *subscriber) AppendKV(k string, v any) {
	caller := "appendKV"
	s._helper.appendKV(k, v, map[string]any{}, caller)
	s._collectEvent()
}

func (s *subscriber) Append(kvMap map[string]any) {
	caller := "append"
	for k, v := range kvMap {
		s._helper.appendKV(k, v, kvMap, caller)
	}
	s._collectEvent()
}

/*
Usage:
 1. SetKV(k, v)
 2. Set({k1: v1, k2: v2})
*/
func (s *subscriber) SetKV(k string, v any) {
	caller := "setKV"
	s._helper.setKV(k, v, map[string]any{}, caller)
	s._collectEvent()
}

func (s *subscriber) Set(kvMap map[string]any) {
	caller := "set"
	for k, v := range kvMap {
		s._helper.setKV(k, v, kvMap, caller)
	}
	s._collectEvent()
}

/*
Usage:
 1. SetOnceKV(k, v)
 2. SetOnce({k1: v1, k2: v2})
*/
func (s *subscriber) SetOnceKV(k string, v any) {
	caller := "set_onceKV"
	s._helper.setOnceKV(k, v, map[string]any{}, caller)
	s._collectEvent()
}

func (s *subscriber) SetOnce(kvMap map[string]any) {
	caller := "set_once"
	for k, v := range kvMap {
		s._helper.setOnceKV(k, v, kvMap, caller)
	}
	s._collectEvent()
}

/*
Usage:
 1. IncrementKV(k, v)
 2. Increment({k1: v1, k2: v2})
*/
func (s *subscriber) IncrementKV(k string, v any) {
	caller := "incrementKV"
	s._helper.incrementKV(k, v, map[string]any{}, caller)
	s._collectEvent()
}

func (s *subscriber) Increment(kvMap map[string]any) {
	caller := "increment"
	for k, v := range kvMap {
		s._helper.incrementKV(k, v, kvMap, caller)
	}
	s._collectEvent()
}

/*
Usage:
 1. RemoveKV(k, v)
 2. Remove({k1: v1, k2: v2})
*/
func (s *subscriber) RemoveKV(k string, v any) {
	caller := "removeKV"
	s._helper.removeKV(k, v, map[string]any{}, caller)
	s._collectEvent()
}

func (s *subscriber) Remove(kvMap map[string]any) {
	caller := "remove"
	for k, v := range kvMap {
		s._helper.removeKV(k, v, kvMap, caller)
	}
	s._collectEvent()
}

/*
Usage:
 1. unset([k1, k2])
*/
func (s *subscriber) Unset(keys []string) {
	caller := "unset"
	for _, key := range keys {
		s._helper.unsetKey(key, caller)
	}
	s._collectEvent()
}

// ----------------------- Preferred language
func (s *subscriber) SetPreferredLanguage(langCode string) {
	caller := "set_preferred_language"
	s._helper.setPreferredLanguage(langCode, caller)
	s._collectEvent()
}

// SetTimezone : set IANA supported timezone as subscriber property
func (s *subscriber) SetTimezone(timezone string) {
	caller := "set_timezone"
	s._helper.setTimezone(timezone, caller)
	s._collectEvent()
}

// ------------------------ Email
func (s *subscriber) AddEmail(value string) {
	caller := "add_email"
	s._helper.addEmail(value, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveEmail(value string) {
	caller := "remove_email"
	s._helper.removeEmail(value, caller)
	s._collectEvent()
}

// ------------------------ SMS
func (s *subscriber) AddSms(value string) {
	caller := "add_sms"
	s._helper.addSms(value, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveSms(value string) {
	caller := "remove_sms"
	s._helper.removeSms(value, caller)
	s._collectEvent()
}

// ------------------------ Whatsapp
func (s *subscriber) AddWhatsapp(value string) {
	caller := "add_whatsapp"
	s._helper.addWhatsapp(value, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveWhatsapp(value string) {
	caller := "remove_whatsapp"
	s._helper.removeWhatsapp(value, caller)
	s._collectEvent()
}

// ------------------------ Androidpush [providers: fcm]

func (s *subscriber) AddAndroidpush(value, provider string) {
	caller := "add_androidpush"
	s._helper.addAndroidpush(value, provider, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveAndroidpush(value, provider string) {
	caller := "remove_androidpush"
	s._helper.removeAndroidpush(value, provider, caller)
	s._collectEvent()
}

// ------------------------ Iospush [providers: apns]

func (s *subscriber) AddIospush(value, provider string) {
	caller := "add_iospush"
	s._helper.addIospush(value, provider, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveIospush(value, provider string) {
	caller := "remove_iospush"
	s._helper.removeIospush(value, provider, caller)
	s._collectEvent()
}

// ------------------------ Webpush [providers: vapid]

func (s *subscriber) AddWebpush(value map[string]any, provider string) {
	caller := "add_webpush"
	s._helper.addWebpush(value, provider, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveWebpush(value map[string]any, provider string) {
	caller := "remove_webpush"
	s._helper.removeWebpush(value, provider, caller)
	s._collectEvent()
}

// ------------------------ Slack

func (s *subscriber) AddSlack(value map[string]any) {
	caller := "add_slack"
	s._helper.addSlack(value, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveSlack(value map[string]any) {
	caller := "remove_slack"
	s._helper.removeSlack(value, caller)
	s._collectEvent()
}

// ------------------------ MS Teams

func (s *subscriber) AddMSTeams(value map[string]any) {
	caller := "add_ms_teams"
	s._helper.addMSTeams(value, caller)
	s._collectEvent()
}

func (s *subscriber) RemoveMSTeams(value map[string]any) {
	caller := "remove_ms_teams"
	s._helper.removeMSTeams(value, caller)
	s._collectEvent()
}
