package suprsend

import (
	"fmt"
	"slices"
	"strings"
)

// ---------- Identity keys
const (
	IDENT_KEY_EMAIL       = "$email"
	IDENT_KEY_SMS         = "$sms"
	IDENT_KEY_ANDROIDPUSH = "$androidpush"
	IDENT_KEY_IOSPUSH     = "$iospush"
	IDENT_KEY_WHATSAPP    = "$whatsapp"
	IDENT_KEY_WEBPUSH     = "$webpush"
	IDENT_KEY_SLACK       = "$slack"
	IDENT_KEY_MS_TEAMS    = "$ms_teams"
)

var IDENT_KEYS_ALL = []string{IDENT_KEY_EMAIL, IDENT_KEY_SMS, IDENT_KEY_ANDROIDPUSH, IDENT_KEY_IOSPUSH,
	IDENT_KEY_WHATSAPP, IDENT_KEY_WEBPUSH, IDENT_KEY_SLACK, IDENT_KEY_MS_TEAMS}

const (
	KEY_ID_PROVIDER        = "$id_provider"
	KEY_PREFERRED_LANGUAGE = "$preferred_language"
	KEY_TIMEZONE           = "$timezone"
)

// ---------
type userEditHelper struct {
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

func newUserEditHelper() *userEditHelper {
	return &userEditHelper{
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

func (u *userEditHelper) reset() {
	u.setDict = map[string]any{}
	u.setOnceDict = map[string]any{}
	u.incrementDict = map[string]any{}
	u.appendDict = map[string]any{}
	u.removeDict = map[string]any{}
	u.unsetList = []string{}
	//
	u._errors = []string{}
	u._info = []string{}
}

type getUserEditOperationResult struct {
	errors []string
	info   []string
	//
	operation map[string]any
}

func (u *userEditHelper) getOperationResult() *getUserEditOperationResult {
	operation := u._formOperation()
	retVal := &getUserEditOperationResult{
		errors:    u._errors,
		info:      u._info,
		operation: operation,
	}
	u.reset()
	return retVal
}

func (u *userEditHelper) _formOperation() map[string]any {
	payload := map[string]any{}
	if len(u.setDict) > 0 {
		payload["$set"] = u.setDict
	}
	if len(u.setOnceDict) > 0 {
		payload["$set_once"] = u.setOnceDict
	}
	if len(u.incrementDict) > 0 {
		payload["$add"] = u.incrementDict
	}
	if len(u.appendDict) > 0 {
		payload["$append"] = u.appendDict
	}
	if len(u.removeDict) > 0 {
		payload["$remove"] = u.removeDict
	}
	if len(u.unsetList) > 0 {
		payload["$unset"] = u.unsetList
	}
	return payload
}

func (u *userEditHelper) _validateKeyBasic(key, caller string) (string, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		u._info = append(u._info, fmt.Sprintf("[%s] skipping key: empty string", caller))
		return key, false
	}
	return key, true
}

func (u *userEditHelper) _isIdentityKey(key string) bool {
	return slices.Contains(IDENT_KEYS_ALL, key)
}

func (u *userEditHelper) appendKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := u._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if u._isIdentityKey(key) {
		u.addIdentity(key, val, kvMap, caller)
	} else {
		u.appendDict[key] = val
	}
}

func (u *userEditHelper) removeKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := u._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if u._isIdentityKey(key) {
		u.removeIdentity(key, val, kvMap, caller)
	} else {
		u.removeDict[key] = val
	}
}

func (u *userEditHelper) setKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := u._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		u.setDict[key] = val
	}
}

func (u *userEditHelper) setOnceKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := u._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		u.setOnceDict[key] = val
	}
}

func (u *userEditHelper) incrementKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := u._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		u.incrementDict[key] = val
	}
}

func (u *userEditHelper) unsetKey(key string, caller string) {
	key, isKeyValid := u._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	u.unsetList = append(u.unsetList, key)
}

func (u *userEditHelper) setPreferredLanguage(langCode string, caller string) {
	u.setDict[KEY_PREFERRED_LANGUAGE] = langCode
}

func (u *userEditHelper) setTimezone(timezone string, caller string) {
	u.setDict[KEY_TIMEZONE] = timezone
}

func (u *userEditHelper) addIdentity(key string, val any, kvMap map[string]any, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		u.addEmail(val, newCaller)

	case IDENT_KEY_SMS:
		u.addSms(val, newCaller)

	case IDENT_KEY_WHATSAPP:
		u.addWhatsapp(val, newCaller)

	case IDENT_KEY_ANDROIDPUSH:
		u.addAndroidpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_IOSPUSH:
		u.addIospush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_WEBPUSH:
		u.addWebpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_SLACK:
		u.addSlack(val, newCaller)

	case IDENT_KEY_MS_TEAMS:
		u.addMSTeams(val, newCaller)
	}
}

func (u *userEditHelper) removeIdentity(key string, val any, kvMap map[string]any, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		u.removeEmail(val, newCaller)

	case IDENT_KEY_SMS:
		u.removeSms(val, newCaller)

	case IDENT_KEY_WHATSAPP:
		u.removeWhatsapp(val, newCaller)

	case IDENT_KEY_ANDROIDPUSH:
		u.removeAndroidpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_IOSPUSH:
		u.removeIospush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_WEBPUSH:
		u.removeWebpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_SLACK:
		u.removeSlack(val, newCaller)

	case IDENT_KEY_MS_TEAMS:
		u.removeMSTeams(val, newCaller)
	}
}

// ------------------------ Email

func (u *userEditHelper) addEmail(value any, caller string) {
	u.appendDict[IDENT_KEY_EMAIL] = value
}

func (u *userEditHelper) removeEmail(value any, caller string) {
	u.removeDict[IDENT_KEY_EMAIL] = value
}

// ------------------------ SMS

func (u *userEditHelper) addSms(value any, caller string) {
	u.appendDict[IDENT_KEY_SMS] = value
}

func (u *userEditHelper) removeSms(value any, caller string) {
	u.removeDict[IDENT_KEY_SMS] = value
}

// ------------------------ Whatsapp

func (u *userEditHelper) addWhatsapp(value any, caller string) {
	u.appendDict[IDENT_KEY_WHATSAPP] = value
}

func (u *userEditHelper) removeWhatsapp(value any, caller string) {
	u.removeDict[IDENT_KEY_WHATSAPP] = value
}

// ------------------------ Androidpush

func (u *userEditHelper) addAndroidpush(value any, provider any, caller string) {
	u.appendDict[IDENT_KEY_ANDROIDPUSH] = value
	u.appendDict[KEY_ID_PROVIDER] = provider
}

func (u *userEditHelper) removeAndroidpush(value any, provider any, caller string) {
	u.removeDict[IDENT_KEY_ANDROIDPUSH] = value
	u.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Iospush

func (u *userEditHelper) addIospush(value any, provider any, caller string) {
	u.appendDict[IDENT_KEY_IOSPUSH] = value
	u.appendDict[KEY_ID_PROVIDER] = provider
}

func (u *userEditHelper) removeIospush(value any, provider any, caller string) {
	u.removeDict[IDENT_KEY_IOSPUSH] = value
	u.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Webpush [providers: vapid]

func (u *userEditHelper) addWebpush(value any, provider any, caller string) {
	u.appendDict[IDENT_KEY_WEBPUSH] = value
	u.appendDict[KEY_ID_PROVIDER] = provider
}

func (u *userEditHelper) removeWebpush(value any, provider any, caller string) {
	u.removeDict[IDENT_KEY_WEBPUSH] = value
	u.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Slack

func (u *userEditHelper) addSlack(value any, caller string) {
	u.appendDict[IDENT_KEY_SLACK] = value
}

func (u *userEditHelper) removeSlack(value any, caller string) {
	u.removeDict[IDENT_KEY_SLACK] = value
}

// ------------------------ MS Teams

func (u *userEditHelper) addMSTeams(value any, caller string) {
	u.appendDict[IDENT_KEY_MS_TEAMS] = value
}

func (u *userEditHelper) removeMSTeams(value any, caller string) {
	u.removeDict[IDENT_KEY_MS_TEAMS] = value
}
