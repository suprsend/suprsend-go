package suprsend

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserEdit interface {
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

var _ UserEdit = &userEdit{}

type userEdit struct {
	client     *Client
	distinctId string
	//
	_errors    []string
	_infos     []string
	operations []map[string]any
	//
	_helper       *userEditHelper
	_warningsList []string
}

func newUserEdit(client *Client, distinctId string) UserEdit {
	u := &userEdit{
		client:     client,
		distinctId: distinctId,
		_helper:    newUserEditHelper(),
	}
	return u
}

func (u *userEdit) GetPayload() map[string]any {
	return map[string]any{
		"operations": u.operations,
	}
}

func (u *userEdit) GetAsyncPayload() map[string]any {
	return map[string]any{
		"$schema":          "2",
		"$insert_id":       uuid.New().String(),
		"$time":            time.Now().UnixMilli(),
		"env":              u.client.getWsIdentifierValue(),
		"distinct_id":      u.distinctId,
		"$user_operations": u.operations,
		"properties":       map[string]any{"$ss_sdk_version": u.client.userAgent},
	}
}

func (u *userEdit) asJsonAsync() map[string]any {
	return map[string]any{
		"distinct_id":      u.distinctId,
		"$user_operations": u.operations,
		"warnings":         u._warningsList,
	}
}

func (u *userEdit) validatePayloadSize(payload map[string]any) (map[string]any, int, error) {
	apparentSize, err := getApparentIdentityEventSize(payload)
	if err != nil {
		return payload, 0, err
	}
	if apparentSize > IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("User Event size too big - %d Bytes, must not cross %s", apparentSize,
			IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, &Error{Code: 413, Message: errStr}
	}
	return payload, apparentSize, nil
}

func (u *userEdit) validateBody() []string {
	u._warningsList = []string{}
	if len(u._infos) > 0 {
		msg := fmt.Sprintf("[distinct_id: %s] %s", u.distinctId, strings.Join(u._infos, "\n"))
		u._warningsList = append(u._warningsList, msg)
		// print on console as well
		log.Println("WARNING:", msg)
	}
	if len(u._errors) > 0 {
		msg := fmt.Sprintf("[distinct_id: %s] %s", u.distinctId, strings.Join(u._errors, "\n"))
		u._warningsList = append(u._warningsList, msg)
		// print on console as well
		log.Println("ERROR:", msg)
	}
	return u._warningsList
}

func (u *userEdit) _collectOperation() {
	resp := u._helper.getOperationResult()
	if len(resp.errors) > 0 {
		u._errors = append(u._errors, resp.errors...)
	}
	if len(resp.info) > 0 {
		u._infos = append(u._infos, resp.info...)
	}
	if len(resp.operation) > 0 {
		u.operations = append(u.operations, resp.operation)
	}
}

/*
Usage:
 1. append(k, v)
 2. append({k1: v1, k2: v2})
*/
func (u *userEdit) AppendKV(k string, v any) {
	caller := "appendKV"
	u._helper.appendKV(k, v, map[string]any{}, caller)
	u._collectOperation()
}

func (u *userEdit) Append(kvMap map[string]any) {
	caller := "append"
	for k, v := range kvMap {
		u._helper.appendKV(k, v, kvMap, caller)
	}
	u._collectOperation()
}

/*
Usage:
 1. SetKV(k, v)
 2. Set({k1: v1, k2: v2})
*/
func (u *userEdit) SetKV(k string, v any) {
	caller := "setKV"
	u._helper.setKV(k, v, map[string]any{}, caller)
	u._collectOperation()
}

func (u *userEdit) Set(kvMap map[string]any) {
	caller := "set"
	for k, v := range kvMap {
		u._helper.setKV(k, v, kvMap, caller)
	}
	u._collectOperation()
}

/*
Usage:
 1. SetOnceKV(k, v)
 2. SetOnce({k1: v1, k2: v2})
*/
func (u *userEdit) SetOnceKV(k string, v any) {
	caller := "set_onceKV"
	u._helper.setOnceKV(k, v, map[string]any{}, caller)
	u._collectOperation()
}

func (u *userEdit) SetOnce(kvMap map[string]any) {
	caller := "set_once"
	for k, v := range kvMap {
		u._helper.setOnceKV(k, v, kvMap, caller)
	}
	u._collectOperation()
}

/*
Usage:
 1. IncrementKV(k, v)
 2. Increment({k1: v1, k2: v2})
*/
func (u *userEdit) IncrementKV(k string, v any) {
	caller := "incrementKV"
	u._helper.incrementKV(k, v, map[string]any{}, caller)
	u._collectOperation()
}

func (u *userEdit) Increment(kvMap map[string]any) {
	caller := "increment"
	for k, v := range kvMap {
		u._helper.incrementKV(k, v, kvMap, caller)
	}
	u._collectOperation()
}

/*
Usage:
 1. RemoveKV(k, v)
 2. Remove({k1: v1, k2: v2})
*/
func (u *userEdit) RemoveKV(k string, v any) {
	caller := "removeKV"
	u._helper.removeKV(k, v, map[string]any{}, caller)
	u._collectOperation()
}

func (u *userEdit) Remove(kvMap map[string]any) {
	caller := "remove"
	for k, v := range kvMap {
		u._helper.removeKV(k, v, kvMap, caller)
	}
	u._collectOperation()
}

/*
Usage:
 1. unset([k1, k2])
*/
func (u *userEdit) Unset(keys []string) {
	caller := "unset"
	for _, key := range keys {
		u._helper.unsetKey(key, caller)
	}
	u._collectOperation()
}

// ----------------------- Preferred language

func (u *userEdit) SetPreferredLanguage(langCode string) {
	caller := "set_preferred_language"
	u._helper.setPreferredLanguage(langCode, caller)
	u._collectOperation()
}

// SetTimezone : set IANA supported timezone
func (u *userEdit) SetTimezone(timezone string) {
	caller := "set_timezone"
	u._helper.setTimezone(timezone, caller)
	u._collectOperation()
}

// ------------------------ Email

func (u *userEdit) AddEmail(value string) {
	caller := "add_email"
	u._helper.addEmail(value, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveEmail(value string) {
	caller := "remove_email"
	u._helper.removeEmail(value, caller)
	u._collectOperation()
}

// ------------------------ SMS

func (u *userEdit) AddSms(value string) {
	caller := "add_sms"
	u._helper.addSms(value, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveSms(value string) {
	caller := "remove_sms"
	u._helper.removeSms(value, caller)
	u._collectOperation()
}

// ------------------------ Whatsapp

func (u *userEdit) AddWhatsapp(value string) {
	caller := "add_whatsapp"
	u._helper.addWhatsapp(value, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveWhatsapp(value string) {
	caller := "remove_whatsapp"
	u._helper.removeWhatsapp(value, caller)
	u._collectOperation()
}

// ------------------------ Androidpush [providers: fcm]

func (u *userEdit) AddAndroidpush(value, provider string) {
	caller := "add_androidpush"
	u._helper.addAndroidpush(value, provider, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveAndroidpush(value, provider string) {
	caller := "remove_androidpush"
	u._helper.removeAndroidpush(value, provider, caller)
	u._collectOperation()
}

// ------------------------ Iospush [providers: apns]

func (u *userEdit) AddIospush(value, provider string) {
	caller := "add_iospush"
	u._helper.addIospush(value, provider, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveIospush(value, provider string) {
	caller := "remove_iospush"
	u._helper.removeIospush(value, provider, caller)
	u._collectOperation()
}

// ------------------------ Webpush [providers: vapid]

func (u *userEdit) AddWebpush(value map[string]any, provider string) {
	caller := "add_webpush"
	u._helper.addWebpush(value, provider, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveWebpush(value map[string]any, provider string) {
	caller := "remove_webpush"
	u._helper.removeWebpush(value, provider, caller)
	u._collectOperation()
}

// ------------------------ Slack

func (u *userEdit) AddSlack(value map[string]any) {
	caller := "add_slack"
	u._helper.addSlack(value, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveSlack(value map[string]any) {
	caller := "remove_slack"
	u._helper.removeSlack(value, caller)
	u._collectOperation()
}

// ------------------------ MS Teams

func (u *userEdit) AddMSTeams(value map[string]any) {
	caller := "add_ms_teams"
	u._helper.addMSTeams(value, caller)
	u._collectOperation()
}

func (u *userEdit) RemoveMSTeams(value map[string]any) {
	caller := "remove_ms_teams"
	u._helper.removeMSTeams(value, caller)
	u._collectOperation()
}
