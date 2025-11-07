package suprsend

import (
	"fmt"
	"slices"
	"strings"
)

// ---------
type objectEditHelper struct {
	setDict       map[string]any
	setOnceDict   map[string]any
	incrementDict map[string]any
	appendDict    map[string]any
	removeDict    map[string]any
	//
	unsetList []string
	//
	_errors []string
	_info   []string
}

func newObjectEditHelper() *objectEditHelper {
	return &objectEditHelper{
		setDict:       map[string]any{},
		setOnceDict:   map[string]any{},
		incrementDict: map[string]any{},
		appendDict:    map[string]any{},
		removeDict:    map[string]any{},
		unsetList:     []string{},
		//
		_errors: []string{},
		_info:   []string{},
	}
}

func (o *objectEditHelper) reset() {
	o.setDict = map[string]any{}
	o.setOnceDict = map[string]any{}
	o.incrementDict = map[string]any{}
	o.appendDict = map[string]any{}
	o.removeDict = map[string]any{}
	o.unsetList = []string{}
	//
	o._errors = []string{}
	o._info = []string{}
}

type getObjectEditOperationResult struct {
	errors []string
	info   []string
	//
	operation map[string]any
}

func (o *objectEditHelper) getOperationResult() *getObjectEditOperationResult {
	operation := o._formOperation()
	retVal := &getObjectEditOperationResult{
		errors:    o._errors,
		info:      o._info,
		operation: operation,
	}
	o.reset()
	return retVal
}

func (o *objectEditHelper) _formOperation() map[string]any {
	payload := map[string]any{}
	if len(o.setDict) > 0 {
		payload["$set"] = o.setDict
	}
	if len(o.setOnceDict) > 0 {
		payload["$set_once"] = o.setOnceDict
	}
	if len(o.incrementDict) > 0 {
		payload["$add"] = o.incrementDict
	}
	if len(o.appendDict) > 0 {
		payload["$append"] = o.appendDict
	}
	if len(o.removeDict) > 0 {
		payload["$remove"] = o.removeDict
	}
	if len(o.unsetList) > 0 {
		payload["$unset"] = o.unsetList
	}
	return payload
}

func (o *objectEditHelper) _validateKeyBasic(key, caller string) (string, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		o._info = append(o._info, fmt.Sprintf("[%s] skipping key: empty string", caller))
		return key, false
	}
	return key, true
}

func (o *objectEditHelper) _isIdentityKey(key string) bool {
	return slices.Contains(IDENT_KEYS_ALL, key)
}

func (o *objectEditHelper) appendKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if o._isIdentityKey(key) {
		o.addIdentity(key, val, kvMap, caller)
	} else {
		o.appendDict[key] = val
	}
}

func (o *objectEditHelper) removeKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if o._isIdentityKey(key) {
		o.removeIdentity(key, val, kvMap, caller)
	} else {
		o.removeDict[key] = val
	}
}

func (o *objectEditHelper) setKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		o.setDict[key] = val
	}
}

func (o *objectEditHelper) setOnceKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		o.setOnceDict[key] = val
	}
}

func (o *objectEditHelper) incrementKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		o.incrementDict[key] = val
	}
}

func (o *objectEditHelper) unsetKey(key string, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	o.unsetList = append(o.unsetList, key)
}

func (o *objectEditHelper) setPreferredLanguage(langCode string, caller string) {
	o.setDict[KEY_PREFERRED_LANGUAGE] = langCode
}

func (o *objectEditHelper) setLocale(localeCode string, caller string) {
	o.setDict[KEY_LOCALE] = localeCode
}

func (o *objectEditHelper) setTimezone(timezone string, caller string) {
	o.setDict[KEY_TIMEZONE] = timezone
}

func (o *objectEditHelper) addIdentity(key string, val any, kvMap map[string]any, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		o.addEmail(val, newCaller)

	case IDENT_KEY_SMS:
		o.addSms(val, newCaller)

	case IDENT_KEY_WHATSAPP:
		o.addWhatsapp(val, newCaller)

	case IDENT_KEY_ANDROIDPUSH:
		o.addAndroidpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_IOSPUSH:
		o.addIospush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_WEBPUSH:
		o.addWebpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_SLACK:
		o.addSlack(val, newCaller)

	case IDENT_KEY_MS_TEAMS:
		o.addMSTeams(val, newCaller)
	}
}

func (o *objectEditHelper) removeIdentity(key string, val any, kvMap map[string]any, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		o.removeEmail(val, newCaller)

	case IDENT_KEY_SMS:
		o.removeSms(val, newCaller)

	case IDENT_KEY_WHATSAPP:
		o.removeWhatsapp(val, newCaller)

	case IDENT_KEY_ANDROIDPUSH:
		o.removeAndroidpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_IOSPUSH:
		o.removeIospush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_WEBPUSH:
		o.removeWebpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_SLACK:
		o.removeSlack(val, newCaller)

	case IDENT_KEY_MS_TEAMS:
		o.removeMSTeams(val, newCaller)
	}
}

// ------------------------ Email

func (o *objectEditHelper) addEmail(value any, caller string) {
	o.appendDict[IDENT_KEY_EMAIL] = value
}

func (o *objectEditHelper) removeEmail(value any, caller string) {
	o.removeDict[IDENT_KEY_EMAIL] = value
}

// ------------------------ SMS

func (o *objectEditHelper) addSms(value any, caller string) {
	o.appendDict[IDENT_KEY_SMS] = value
}

func (o *objectEditHelper) removeSms(value any, caller string) {
	o.removeDict[IDENT_KEY_SMS] = value
}

// ------------------------ Whatsapp

func (o *objectEditHelper) addWhatsapp(value any, caller string) {
	o.appendDict[IDENT_KEY_WHATSAPP] = value
}

func (o *objectEditHelper) removeWhatsapp(value any, caller string) {
	o.removeDict[IDENT_KEY_WHATSAPP] = value
}

// ------------------------ Androidpush

func (o *objectEditHelper) addAndroidpush(value any, provider any, caller string) {
	o.appendDict[IDENT_KEY_ANDROIDPUSH] = value
	o.appendDict[KEY_ID_PROVIDER] = provider
}

func (o *objectEditHelper) removeAndroidpush(value any, provider any, caller string) {
	o.removeDict[IDENT_KEY_ANDROIDPUSH] = value
	o.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Iospush

func (o *objectEditHelper) addIospush(value any, provider any, caller string) {
	o.appendDict[IDENT_KEY_IOSPUSH] = value
	o.appendDict[KEY_ID_PROVIDER] = provider
}

func (o *objectEditHelper) removeIospush(value any, provider any, caller string) {
	o.removeDict[IDENT_KEY_IOSPUSH] = value
	o.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Webpush [providers: vapid]

func (o *objectEditHelper) addWebpush(value any, provider any, caller string) {
	o.appendDict[IDENT_KEY_WEBPUSH] = value
	o.appendDict[KEY_ID_PROVIDER] = provider
}

func (o *objectEditHelper) removeWebpush(value any, provider any, caller string) {
	o.removeDict[IDENT_KEY_WEBPUSH] = value
	o.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Slack

func (o *objectEditHelper) addSlack(value any, caller string) {
	o.appendDict[IDENT_KEY_SLACK] = value
}

func (o *objectEditHelper) removeSlack(value any, caller string) {
	o.removeDict[IDENT_KEY_SLACK] = value
}

// ------------------------ MS Teams

func (o *objectEditHelper) addMSTeams(value any, caller string) {
	o.appendDict[IDENT_KEY_MS_TEAMS] = value
}

func (o *objectEditHelper) removeMSTeams(value any, caller string) {
	o.removeDict[IDENT_KEY_MS_TEAMS] = value
}
