package suprsend

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

type subscribersService struct {
	client *Client
}

func (s *subscribersService) GetInstance(distinctId string) Subscriber {
	distinctId = strings.TrimSpace(distinctId)
	return newSubscriber(s.client, distinctId)
}

type Subscriber interface {
	Save() (*Response, error)
	//
	AppendKV(string, interface{})
	Append(map[string]interface{})
	RemoveKV(string, interface{})
	Remove(map[string]interface{})
	Unset([]string)
	//
	SetPreferredLanguage(string)
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
	AddSlackEmail(value string)
	RemoveSlackEmail(value string)
	//
	AddSlackUserid(value string)
	RemoveSlackUserid(value string)
}

var _ Subscriber = &subscriber{}

type subscriber struct {
	client      *Client
	distinctId  string
	_url        string
	_superProps map[string]interface{}
	//
	_errors     []string
	_warnings   []string
	appendCount int
	removeCount int
	setCount    int
	unsetCount  int
	//
	_events []map[string]interface{}
	_helper *subscriberHelper
}

func newSubscriber(client *Client, distinctId string) Subscriber {
	s := &subscriber{
		client:     client,
		distinctId: distinctId,
	}
	// events url
	s._url = fmt.Sprintf("%sevent/", client.baseUrl)
	s._superProps = map[string]interface{}{"$ss_sdk_version": client.userAgent}
	s._helper = newSubscriberHelper(distinctId, client.ApiKey)
	return s
}

func (s *subscriber) validateEventSize(event map[string]interface{}) (map[string]interface{}, int, error) {
	apparentSize, err := getApparentIdentityEventSize(event)
	if err != nil {
		return event, 0, err
	}
	if apparentSize > IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("User Event size too big - %d Bytes, must not cross %s", apparentSize,
			IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, errors.New(errStr)
	}
	return event, apparentSize, nil
}

func (s *subscriber) getEvents() []map[string]interface{} {
	// TOOD: Don't mutate the original array. Make a copy of _events
	allEvents := s._events[:]
	for _, e := range allEvents {
		e["properties"] = s._superProps
	}
	//  # Add $identify event by default, if new properties get added
	if len(allEvents) == 0 || s.setCount > 0 || s.appendCount > 0 {
		// Don't add $anon_id to properties,
		identifyEventProps := map[string]interface{}{"$identified_id": s.distinctId}
		maps.Copy(identifyEventProps, s._superProps)

		userIdentifyEvent := map[string]interface{}{
			"$insert_id": uuid.New().String(),
			"$time":      time.Now().UnixMilli(),
			"env":        s.client.ApiKey,
			"event":      "$identify",
			"properties": identifyEventProps,
		}
		// # Add $identify event at the 0th index, so that $identify runs before $append/$remove/$reset
		allEvents = append([]map[string]interface{}{userIdentifyEvent}, allEvents...)

	}
	return allEvents
}

func (s *subscriber) validateBody(isPartOfBulk bool) ([]string, error) {
	if s.distinctId == "" {
		s._errors = append([]string{"missing distinct_id"}, s._errors...)
	}
	warningsList := []string{}
	if len(s._warnings) > 0 {
		msg := fmt.Sprintf("[distinct_id: %s] %s", s.distinctId, strings.Join(s._warnings, "\n"))
		warningsList = append(warningsList, msg)
		// print on console as well
		log.Println("WARNING:", msg)
	}
	if len(s._errors) > 0 {
		msg := fmt.Sprintf("[distinct_id: %s] %s", s.distinctId, strings.Join(s._errors, "\n"))
		warningsList = append(warningsList, msg)
		errMsg := "ERROR: " + msg
		if isPartOfBulk {
			// print on console in case of bulk-api
			log.Println(errMsg)
		} else {
			// raise error in case of single api
			return nil, errors.New(errMsg)
		}
	}
	return warningsList, nil
}

func (s *subscriber) Save() (*Response, error) {
	if _, err := s.validateBody(false); err != nil {
		return nil, err
	}
	//
	events := s.getEvents()
	for _, ev := range events {
		if _, _, err := s.validateEventSize(ev); err != nil {
			return nil, err
		}
	}
	// prepare http.Request object
	request, err := s.client.prepareHttpRequest("POST", s._url, events)
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
		return nil, err
	}
	if httpRes.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpRes.StatusCode, string(respBody))
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}

func (s *subscriber) _collectEvent(discardIfError bool) {
	resp := s._helper.getIdentityEvent()
	if len(resp.errors) > 0 {
		s._errors = append(s._errors, resp.errors...)
	}
	if len(resp.info) > 0 {
		s._warnings = append(s._warnings, resp.info...)
	}
	if len(resp.event) > 0 {
		s._events = append(s._events, resp.event)
		s.setCount += resp.set
		s.appendCount += resp.append
		s.removeCount += resp.remove
		s.unsetCount += resp.unset
	}
}

/*
Usage:
 1. append(k, v)
 2. append({k1: v1, k2, v2})
*/
func (s *subscriber) AppendKV(k string, v interface{}) {
	caller := "appendKV"
	s._helper.appendKV(k, v, map[string]interface{}{}, caller)
	s._collectEvent(true)
}

func (s *subscriber) Append(kvMap map[string]interface{}) {
	caller := "append"
	for k, v := range kvMap {
		s._helper.appendKV(k, v, kvMap, caller)

	}
	s._collectEvent(false)
}

/*
Usage:
 1. remove(k, v)
 2. remove({k1: v1, k2, v2})
*/
func (s *subscriber) RemoveKV(k string, v interface{}) {
	caller := "removeKV"
	s._helper.removeKV(k, v, map[string]interface{}{}, caller)
	s._collectEvent(true)
}

func (s *subscriber) Remove(kvMap map[string]interface{}) {
	caller := "remove"
	for k, v := range kvMap {
		s._helper.removeKV(k, v, kvMap, caller)

	}
	s._collectEvent(false)
}

/*
Usage:
 1. unset(k)
 2. unset([k1, k2])
*/
func (s *subscriber) Unset(keys []string) {
	caller := "unset"
	for _, key := range keys {
		s._helper.unsetKey(key, caller)
	}
	s._collectEvent(false)
}

// ----------------------- Preferred language
func (s *subscriber) SetPreferredLanguage(langCode string) {
	caller := "set_preferred_language"
	s._helper.setPreferredLanguage(langCode, caller)
	s._collectEvent(true)
}

// ------------------------ Email
func (s *subscriber) AddEmail(value string) {
	caller := "add_email"
	s._helper.addEmail(value, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveEmail(value string) {
	caller := "remove_email"
	s._helper.removeEmail(value, caller)
	s._collectEvent(true)
}

// ------------------------ SMS
func (s *subscriber) AddSms(value string) {
	caller := "add_sms"
	s._helper.addSms(value, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveSms(value string) {
	caller := "remove_sms"
	s._helper.removeSms(value, caller)
	s._collectEvent(true)
}

// ------------------------ Whatsapp
func (s *subscriber) AddWhatsapp(value string) {
	caller := "add_whatsapp"
	s._helper.addWhatsapp(value, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveWhatsapp(value string) {
	caller := "remove_whatsapp"
	s._helper.removeWhatsapp(value, caller)
	s._collectEvent(true)
}

// ------------------------ Androidpush [providers: fcm / xiaomi / oppo]

func (s *subscriber) AddAndroidpush(value, provider string) {
	caller := "add_androidpush"
	s._helper.addAndroidpush(value, provider, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveAndroidpush(value, provider string) {
	caller := "remove_androidpush"
	s._helper.removeAndroidpush(value, provider, caller)
	s._collectEvent(true)
}

// ------------------------ Iospush [providers: apns]

func (s *subscriber) AddIospush(value, provider string) {
	caller := "add_iospush"
	s._helper.addIospush(value, provider, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveIospush(value, provider string) {
	caller := "remove_iospush"
	s._helper.removeIospush(value, provider, caller)
	s._collectEvent(true)
}

// ------------------------ Webpush [providers: vapid]

func (s *subscriber) AddWebpush(value map[string]interface{}, provider string) {
	caller := "add_webpush"
	s._helper.addWebpush(value, provider, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveWebpush(value map[string]interface{}, provider string) {
	caller := "remove_webpush"
	s._helper.removeWebpush(value, provider, caller)
	s._collectEvent(true)
}

// ------------------------ Slack

func (s *subscriber) AddSlackEmail(value string) {
	caller := "add_slack_email"
	s._helper.addSlack(map[string]string{"email": value}, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveSlackEmail(value string) {
	caller := "remove_slack_email"
	s._helper.removeSlack(map[string]string{"email": value}, caller)
	s._collectEvent(true)
}

func (s *subscriber) AddSlackUserid(value string) {
	caller := "add_slack_userid"
	s._helper.addSlack(map[string]string{"user_id": value}, caller)
	s._collectEvent(true)
}

func (s *subscriber) RemoveSlackUserid(value string) {
	caller := "remove_slack_userid"
	s._helper.removeSlack(map[string]string{"user_id": value}, caller)
	s._collectEvent(true)
}
