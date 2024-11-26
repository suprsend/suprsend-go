package suprsend

import (
	"fmt"
	"strings"
)

// ---------
type objectHelper struct {
	setDict       map[string]interface{}
	setOnceDict   map[string]interface{}
	incrementDict map[string]interface{}
	appendDict    map[string]interface{}
	removeDict    map[string]interface{}
	//
	unsetList []string
	//
	_errors []string
	_info   []string
}

type getObjectIdentityEventResp struct {
	errors []string
	info   []string
	//
	payload map[string]interface{}
}

func newObjectHelper() *objectHelper {
	return &objectHelper{
		setDict:       map[string]interface{}{},
		setOnceDict:   map[string]interface{}{},
		incrementDict: map[string]interface{}{},
		appendDict:    map[string]interface{}{},
		removeDict:    map[string]interface{}{},
		unsetList:     []string{},
		//
		_errors: []string{},
		_info:   []string{},
	}
}

func (o *objectHelper) getIdentityEvent() *getObjectIdentityEventResp {
	payload := o._formPayload()
	retVal := &getObjectIdentityEventResp{
		errors:  o._errors,
		info:    o._info,
		payload: payload,
	}
	o.reset()
	return retVal
}

func (o *objectHelper) reset() {
	o.setDict = map[string]interface{}{}
	o.setOnceDict = map[string]interface{}{}
	o.incrementDict = map[string]interface{}{}
	o.appendDict = map[string]interface{}{}
	o.removeDict = map[string]interface{}{}
	o.unsetList = []string{}
	//
	o._errors = []string{}
	o._info = []string{}
}

func (o *objectHelper) _formPayload() map[string]interface{} {
	payload := map[string]interface{}{}
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

func (o *objectHelper) _validateKeyBasic(key, caller string) (string, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		o._info = append(o._info, fmt.Sprintf("[%s] skipping key: empty string", caller))
		return key, false
	}
	return key, true
}

func (o *objectHelper) _validateKeyPrefix(key, caller string) bool {
	if !Contains(ALL_RESERVED_KEYS, key) {
		keyLower := strings.ToLower(key)
		if strings.HasPrefix(keyLower, "$") || strings.HasPrefix(keyLower, "ss_") {
			o._info = append(o._info, fmt.Sprintf("[%s] skipping key: %s. key starting with [$,ss_] are reserved", caller, key))
			return false
		}
	}
	return true
}

func (o *objectHelper) _isIdentityKey(key string) bool {
	return Contains(IDENT_KEYS_ALL, key)
}

func (o *objectHelper) appendKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if o._isIdentityKey(key) {
		o.addIdentity(key, val, kvMap, caller)
	} else {
		isKeyValid := o._validateKeyPrefix(key, caller)
		if isKeyValid {
			o.appendDict[key] = val
		}
	}
}

func (o *objectHelper) setKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		isKeyValid := o._validateKeyPrefix(key, caller)
		if isKeyValid {
			o.setDict[key] = val
		}
	}
}

func (o *objectHelper) setOnceKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		isKeyValid := o._validateKeyPrefix(key, caller)
		if isKeyValid {
			o.setOnceDict[key] = val
		}
	}
}

func (o *objectHelper) incrementKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		isKeyValid := o._validateKeyPrefix(key, caller)
		if isKeyValid {
			o.incrementDict[key] = val
		}
	}
}

func (o *objectHelper) removeKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if o._isIdentityKey(key) {
		o.removeIdentity(key, val, kvMap, caller)
	} else {
		isKeyValid := o._validateKeyPrefix(key, caller)
		if isKeyValid {
			o.removeDict[key] = val
		}
	}
}

func (o *objectHelper) unsetKey(key string, caller string) {
	key, isKeyValid := o._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	o.unsetList = append(o.unsetList, key)
}

func (o *objectHelper) setPreferredLanguage(langCode string, caller string) {
	// Check language code is in the list
	if !Contains(ALL_LANG_CODES, langCode) {
		o._info = append(o._info, fmt.Sprintf("[%s] invalid value %s", caller, langCode))
		return
	}
	o.setDict[KEY_PREFERRED_LANGUAGE] = langCode
}

func (o *objectHelper) setTimezone(timezone string, caller string) {
	o.setDict[KEY_TIMEZONE] = timezone
}

func (o *objectHelper) addIdentity(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			o.addEmail(val.(string), newCaller)
		}
	case IDENT_KEY_SMS:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			o.addSms(val.(string), newCaller)
		}

	case IDENT_KEY_WHATSAPP:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			o.addWhatsapp(val.(string), newCaller)
		}

	case IDENT_KEY_ANDROIDPUSH:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			pushvendor := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					pushvendor = pvStr
				}
			}
			o.addAndroidpush(val.(string), pushvendor, newCaller)
		}

	case IDENT_KEY_IOSPUSH:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			pushvendor := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					pushvendor = pvStr
				}
			}
			o.addIospush(val.(string), pushvendor, newCaller)
		}
	case IDENT_KEY_WEBPUSH:
		val, isValid := o._checkIdentValDict(val, newCaller)
		if isValid {
			pushvendor := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					pushvendor = pvStr
				}
			}
			o.addWebpush(val.(map[string]interface{}), pushvendor, newCaller)
		}

	case IDENT_KEY_SLACK:
		val, isValid := o._checkIdentValDict(val, newCaller)
		if isValid {
			o.addSlack(val.(map[string]interface{}), newCaller)
		}

	case IDENT_KEY_MS_TEAMS:
		val, isValid := o._checkIdentValDict(val, newCaller)
		if isValid {
			o.addMSTeams(val.(map[string]interface{}), newCaller)
		}
	}
}

func (o *objectHelper) removeIdentity(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			o.removeEmail(val.(string), newCaller)
		}

	case IDENT_KEY_SMS:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			o.removeSms(val.(string), newCaller)
		}

	case IDENT_KEY_WHATSAPP:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			o.removeWhatsapp(val.(string), newCaller)
		}

	case IDENT_KEY_ANDROIDPUSH:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			pushvendor := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					pushvendor = pvStr
				}
			}
			o.removeAndroidpush(val.(string), pushvendor, newCaller)
		}

	case IDENT_KEY_IOSPUSH:
		val, isValid := o._checkIdentValString(val, newCaller)
		if isValid {
			pushvendor := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					pushvendor = pvStr
				}
			}
			o.removeIospush(val.(string), pushvendor, newCaller)
		}
	case IDENT_KEY_WEBPUSH:
		val, isValid := o._checkIdentValDict(val, newCaller)
		if isValid {
			pushvendor := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					pushvendor = pvStr
				}
			}
			o.removeWebpush(val.(map[string]interface{}), pushvendor, newCaller)
		}

	case IDENT_KEY_SLACK:
		val, isValid := o._checkIdentValDict(val, newCaller)
		if isValid {
			o.removeSlack(val.(map[string]interface{}), newCaller)
		}

	case IDENT_KEY_MS_TEAMS:
		val, isValid := o._checkIdentValDict(val, newCaller)
		if isValid {
			o.removeMSTeams(val.(map[string]interface{}), newCaller)
		}
	}
}

// ------------------------

func (o *objectHelper) _checkIdentValString(value interface{}, caller string) (interface{}, bool) {
	msg := "value must be a string with proper value"
	if vstring, ok := value.(string); !ok {
		o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		vstring = strings.TrimSpace(vstring)
		if vstring == "" {
			o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
			return value, false
		}
		return vstring, true
	}
}

func (o *objectHelper) _checkIdentValDict(value interface{}, caller string) (interface{}, bool) {
	msg := "value must be a valid dict/map"
	if valMap, ok := value.(map[string]interface{}); !ok {
		o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		if len(valMap) == 0 {
			o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
			return value, false
		}
		return valMap, true
	}
}

// ------------------------ Email

func (o *objectHelper) addEmail(value string, caller string) {
	val, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_EMAIL] = val
}

func (o *objectHelper) removeEmail(value string, caller string) {
	val, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_EMAIL] = val
}

// ------------------------ SMS

func (o *objectHelper) addSms(value string, caller string) {
	val, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_SMS] = val
}

func (o *objectHelper) removeSms(value string, caller string) {
	val, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_SMS] = val
}

// ------------------------ Whatsapp

func (o *objectHelper) addWhatsapp(value string, caller string) {
	val, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_WHATSAPP] = val
}

func (o *objectHelper) removeWhatsapp(value string, caller string) {
	val, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_WHATSAPP] = val
}

// ------------------------ Androidpush

func (o *objectHelper) _checkAndroidpushValue(value string, provider string, caller string) (string, string, bool) {
	iValue, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return value, provider, false
	}
	value = iValue.(string)
	// convert to lowercase to make it case-insensitive
	provider = strings.ToLower(provider)
	return value, provider, true
}

func (o *objectHelper) addAndroidpush(value string, provider string, caller string) {
	value, provider, isValid := o._checkAndroidpushValue(value, provider, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_ANDROIDPUSH] = value
	o.appendDict[KEY_ID_PROVIDER] = provider
}

func (o *objectHelper) removeAndroidpush(value string, provider string, caller string) {
	value, provider, isValid := o._checkAndroidpushValue(value, provider, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_ANDROIDPUSH] = value
	o.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Iospush

func (o *objectHelper) _checkIospushValue(value string, provider string, caller string) (string, string, bool) {
	iValue, isValid := o._checkIdentValString(value, caller)
	if !isValid {
		return value, provider, false
	}
	value = iValue.(string)
	// convert provider to lowercase to make it case-insensitive
	provider = strings.ToLower(provider)
	return value, provider, true
}

func (o *objectHelper) addIospush(value string, provider string, caller string) {
	value, provider, isValid := o._checkIospushValue(value, provider, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_IOSPUSH] = value
	o.appendDict[KEY_ID_PROVIDER] = provider
}

func (o *objectHelper) removeIospush(value string, provider string, caller string) {
	value, provider, isValid := o._checkIospushValue(value, provider, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_IOSPUSH] = value
	o.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Webpush [providers: vapid]

func (o *objectHelper) _checkWebpushDict(value map[string]interface{}, provider string, caller string) (interface{}, string, bool) {
	msg := "value must be a valid dict/map representing webpush-token"
	if len(value) == 0 {
		o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, provider, false
	}
	// convert provider to lowercase to make it case-insensitive
	provider = strings.ToLower(provider)
	return value, provider, true
}

func (o *objectHelper) addWebpush(value map[string]interface{}, provider string, caller string) {
	iValue, provider, isValid := o._checkWebpushDict(value, provider, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_WEBPUSH] = iValue.(map[string]interface{})
	o.appendDict[KEY_ID_PROVIDER] = provider
}

func (o *objectHelper) removeWebpush(value map[string]interface{}, provider string, caller string) {
	iValue, provider, isValid := o._checkWebpushDict(value, provider, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_WEBPUSH] = iValue.(map[string]interface{})
	o.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Slack

func (o *objectHelper) _checkSlackDict(value map[string]interface{}, caller string) (map[string]interface{}, bool) {
	msg := "value must be a valid dict/map with proper keys"
	if len(value) == 0 {
		o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		return value, true
	}
}

func (o *objectHelper) addSlack(value map[string]interface{}, caller string) {
	value, isValid := o._checkSlackDict(value, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_SLACK] = value
}

func (o *objectHelper) removeSlack(value map[string]interface{}, caller string) {
	value, isValid := o._checkSlackDict(value, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_SLACK] = value
}

// ------------------------ MS Teams

func (o *objectHelper) _checkMSTeamsDict(value map[string]interface{}, caller string) (map[string]interface{}, bool) {
	msg := "value must be a valid dict/map with proper keys"
	if len(value) == 0 {
		o._errors = append(o._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		return value, true
	}
}

func (o *objectHelper) addMSTeams(value map[string]interface{}, caller string) {
	value, isValid := o._checkMSTeamsDict(value, caller)
	if !isValid {
		return
	}
	o.appendDict[IDENT_KEY_MS_TEAMS] = value
}

func (o *objectHelper) removeMSTeams(value map[string]interface{}, caller string) {
	value, isValid := o._checkMSTeamsDict(value, caller)
	if !isValid {
		return
	}
	o.removeDict[IDENT_KEY_MS_TEAMS] = value
}
