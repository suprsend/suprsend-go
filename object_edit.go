package suprsend

import (
	"fmt"
	"log"
	"strings"
)

type ObjectEdit interface {
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
	SetLocale(string)
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

var _ ObjectEdit = &objectEdit{}

type objectEdit struct {
	client     *Client
	objectType string
	objectId   string
	//
	_errors    []string
	_infos     []string
	operations []map[string]any
	//
	_helper *objectEditHelper
}

func newObjectEdit(client *Client, obj ObjectIdentifier) ObjectEdit {
	o := &objectEdit{
		client:     client,
		objectType: obj.ObjectType,
		objectId:   obj.Id,
		_helper:    newObjectEditHelper(),
	}
	return o
}

func (o *objectEdit) GetPayload() map[string]any {
	return map[string]any{
		"operations": o.operations,
	}
}

func (o *objectEdit) validateBody() {
	if len(o._infos) > 0 {
		msg := fmt.Sprintf("[object: %s/%s] %s", o.objectType, o.objectId, strings.Join(o._infos, "\n"))
		log.Println("WARNING:", msg)
	}
	if len(o._errors) > 0 {
		msg := fmt.Sprintf("[object: %s/%s] %s", o.objectType, o.objectId, strings.Join(o._errors, "\n"))
		log.Println("ERROR:", msg)
	}
}

func (o *objectEdit) _collectOperation() {
	resp := o._helper.getOperationResult()
	if len(resp.errors) > 0 {
		o._errors = append(o._errors, resp.errors...)
	}
	if len(resp.info) > 0 {
		o._infos = append(o._infos, resp.info...)
	}
	if len(resp.operation) > 0 {
		o.operations = append(o.operations, resp.operation)
	}
}

/*
Usage:
 1. append(k, v)
 2. append({k1: v1, k2: v2})
*/
func (o *objectEdit) AppendKV(k string, v any) {
	caller := "appendKV"
	o._helper.appendKV(k, v, map[string]any{}, caller)
	o._collectOperation()
}

func (o *objectEdit) Append(kvMap map[string]any) {
	caller := "append"
	for k, v := range kvMap {
		o._helper.appendKV(k, v, kvMap, caller)
	}
	o._collectOperation()
}

/*
Usage:
 1. SetKV(k, v)
 2. Set({k1: v1, k2: v2})
*/
func (o *objectEdit) SetKV(k string, v any) {
	caller := "setKV"
	o._helper.setKV(k, v, map[string]any{}, caller)
	o._collectOperation()
}

func (o *objectEdit) Set(kvMap map[string]any) {
	caller := "set"
	for k, v := range kvMap {
		o._helper.setKV(k, v, kvMap, caller)
	}
	o._collectOperation()
}

/*
Usage:
 1. SetOnceKV(k, v)
 2. SetOnce({k1: v1, k2: v2})
*/
func (o *objectEdit) SetOnceKV(k string, v any) {
	caller := "set_onceKV"
	o._helper.setOnceKV(k, v, map[string]any{}, caller)
	o._collectOperation()
}

func (o *objectEdit) SetOnce(kvMap map[string]any) {
	caller := "set_once"
	for k, v := range kvMap {
		o._helper.setOnceKV(k, v, kvMap, caller)
	}
	o._collectOperation()
}

/*
Usage:
 1. IncrementKV(k, v)
 2. Increment({k1: v1, k2: v2})
*/
func (o *objectEdit) IncrementKV(k string, v any) {
	caller := "incrementKV"
	o._helper.incrementKV(k, v, map[string]any{}, caller)
	o._collectOperation()
}

func (o *objectEdit) Increment(kvMap map[string]any) {
	caller := "increment"
	for k, v := range kvMap {
		o._helper.incrementKV(k, v, kvMap, caller)
	}
	o._collectOperation()
}

/*
Usage:
 1. RemoveKV(k, v)
 2. Remove({k1: v1, k2: v2})
*/
func (o *objectEdit) RemoveKV(k string, v any) {
	caller := "removeKV"
	o._helper.removeKV(k, v, map[string]any{}, caller)
	o._collectOperation()
}

func (o *objectEdit) Remove(kvMap map[string]any) {
	caller := "remove"
	for k, v := range kvMap {
		o._helper.removeKV(k, v, kvMap, caller)
	}
	o._collectOperation()
}

/*
Usage:
 1. unset([k1, k2])
*/
func (o *objectEdit) Unset(keys []string) {
	caller := "unset"
	for _, key := range keys {
		o._helper.unsetKey(key, caller)
	}
	o._collectOperation()
}

// ----------------------- Preferred language

func (o *objectEdit) SetPreferredLanguage(langCode string) {
	caller := "set_preferred_language"
	o._helper.setPreferredLanguage(langCode, caller)
	o._collectOperation()
}

func (o *objectEdit) SetLocale(localeCode string) {
	caller := "set_locale"
	o._helper.setLocale(localeCode, caller)
	o._collectOperation()
}

// SetTimezone : set IANA supported timezone
func (o *objectEdit) SetTimezone(timezone string) {
	caller := "set_timezone"
	o._helper.setTimezone(timezone, caller)
	o._collectOperation()
}

// ------------------------ Email

func (o *objectEdit) AddEmail(value string) {
	caller := "add_email"
	o._helper.addEmail(value, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveEmail(value string) {
	caller := "remove_email"
	o._helper.removeEmail(value, caller)
	o._collectOperation()
}

// ------------------------ SMS

func (o *objectEdit) AddSms(value string) {
	caller := "add_sms"
	o._helper.addSms(value, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveSms(value string) {
	caller := "remove_sms"
	o._helper.removeSms(value, caller)
	o._collectOperation()
}

// ------------------------ Whatsapp

func (o *objectEdit) AddWhatsapp(value string) {
	caller := "add_whatsapp"
	o._helper.addWhatsapp(value, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveWhatsapp(value string) {
	caller := "remove_whatsapp"
	o._helper.removeWhatsapp(value, caller)
	o._collectOperation()
}

// ------------------------ Androidpush [providers: fcm]

func (o *objectEdit) AddAndroidpush(value, provider string) {
	caller := "add_androidpush"
	o._helper.addAndroidpush(value, provider, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveAndroidpush(value, provider string) {
	caller := "remove_androidpush"
	o._helper.removeAndroidpush(value, provider, caller)
	o._collectOperation()
}

// ------------------------ Iospush [providers: apns]

func (o *objectEdit) AddIospush(value, provider string) {
	caller := "add_iospush"
	o._helper.addIospush(value, provider, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveIospush(value, provider string) {
	caller := "remove_iospush"
	o._helper.removeIospush(value, provider, caller)
	o._collectOperation()
}

// ------------------------ Webpush [providers: vapid]

func (o *objectEdit) AddWebpush(value map[string]any, provider string) {
	caller := "add_webpush"
	o._helper.addWebpush(value, provider, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveWebpush(value map[string]any, provider string) {
	caller := "remove_webpush"
	o._helper.removeWebpush(value, provider, caller)
	o._collectOperation()
}

// ------------------------ Slack

func (o *objectEdit) AddSlack(value map[string]any) {
	caller := "add_slack"
	o._helper.addSlack(value, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveSlack(value map[string]any) {
	caller := "remove_slack"
	o._helper.removeSlack(value, caller)
	o._collectOperation()
}

// ------------------------ MS Teams

func (o *objectEdit) AddMSTeams(value map[string]any) {
	caller := "add_ms_teams"
	o._helper.addMSTeams(value, caller)
	o._collectOperation()
}

func (o *objectEdit) RemoveMSTeams(value map[string]any) {
	caller := "remove_ms_teams"
	o._helper.removeMSTeams(value, caller)
	o._collectOperation()
}
