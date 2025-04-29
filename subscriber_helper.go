package suprsend

import (
	"fmt"
	"regexp"
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

var OTHER_RESERVED_KEYS = []string{
	"$messenger", "$inbox",
	KEY_ID_PROVIDER, "$device_id",
	"$insert_id", "$time",
	"$set", "$set_once", "$add", "$append", "$remove", "$unset",
	"$identify", "$anon_id", "$identified_id",
	KEY_PREFERRED_LANGUAGE, KEY_TIMEZONE,
	"$notification_delivered", "$notification_dismiss", "$notification_clicked",
}

var SUPER_PROPERTY_KEYS = []string{
	"$app_version_string", "$app_build_number", "$brand", "$carrier", "$manufacturer", "$model", "$os",
	"$ss_sdk_version", "$insert_id", "$time",
}

var ALL_RESERVED_KEYS = append(append(SUPER_PROPERTY_KEYS, OTHER_RESERVED_KEYS...), IDENT_KEYS_ALL...)

// ---------
const MOBILE_REGEX = "^\\+[0-9\\s]+"

var mobileRegexCompiled = regexp.MustCompile(MOBILE_REGEX)

const EMAIL_REGEX = "^\\S+@\\S+\\.\\S+$"

var emailRegexCompiled = regexp.MustCompile(EMAIL_REGEX)

// ---------

type subscriberHelper struct {
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

func newSubscriberHelper() *subscriberHelper {
	return &subscriberHelper{
		setDict:       map[string]interface{}{},
		setOnceDict:   map[string]interface{}{},
		incrementDict: map[string]interface{}{},
		appendDict:    map[string]interface{}{},
		removeDict:    map[string]interface{}{},
		unsetList:     []string{},
		_errors:       []string{},
		_info:         []string{},
	}
}

func (s *subscriberHelper) reset() {
	s.setDict = map[string]interface{}{}
	s.setOnceDict = map[string]interface{}{}
	s.incrementDict = map[string]interface{}{}
	s.appendDict = map[string]interface{}{}
	s.removeDict = map[string]interface{}{}
	s.unsetList = []string{}
	s._errors, s._info = []string{}, []string{}
}

type getIdentityEventResp struct {
	errors []string
	info   []string
	//
	event map[string]interface{}
}

func (s *subscriberHelper) getIdentityEvent() *getIdentityEventResp {
	evt := s._formEvent()
	retVal := &getIdentityEventResp{
		errors: s._errors,
		info:   s._info,
		event:  evt,
	}
	s.reset()
	return retVal
}

func (s *subscriberHelper) _formEvent() map[string]interface{} {
	event := map[string]interface{}{}
	if len(s.setDict) > 0 {
		event["$set"] = s.setDict
	}
	if len(s.setOnceDict) > 0 {
		event["$set_once"] = s.setOnceDict
	}
	if len(s.incrementDict) > 0 {
		event["$add"] = s.incrementDict
	}
	if len(s.appendDict) > 0 {
		event["$append"] = s.appendDict
	}
	if len(s.removeDict) > 0 {
		event["$remove"] = s.removeDict
	}
	if len(s.unsetList) > 0 {
		event["$unset"] = s.unsetList
	}
	return event
}

func (s *subscriberHelper) _validateKeyBasic(key, caller string) (string, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		s._info = append(s._info, fmt.Sprintf("[%s] skipping key: empty string", caller))
		return key, false
	}
	return key, true
}

func (s *subscriberHelper) _validateKeyPrefix(key, caller string) bool {
	if !slices.Contains(ALL_RESERVED_KEYS, key) {
		keyLower := strings.ToLower(key)
		if strings.HasPrefix(keyLower, "$") || strings.HasPrefix(keyLower, "ss_") {
			s._info = append(s._info, fmt.Sprintf("[%s] skipping key: %s. key starting with [$,ss_] are reserved", caller, key))
			return false
		}
	}
	return true
}

func (s *subscriberHelper) _isIdentityKey(key string) bool {
	return slices.Contains(IDENT_KEYS_ALL, key)
}

// -------------------------

func (s *subscriberHelper) appendKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if s._isIdentityKey(key) {
		s.addIdentity(key, val, kvMap, caller)
	} else {
		isKeyValid := s._validateKeyPrefix(key, caller)
		if isKeyValid {
			s.appendDict[key] = val
		}
	}
}

func (s *subscriberHelper) setKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		isKeyValid := s._validateKeyPrefix(key, caller)
		if isKeyValid {
			s.setDict[key] = val
		}
	}
}

func (s *subscriberHelper) setOnceKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		isKeyValid := s._validateKeyPrefix(key, caller)
		if isKeyValid {
			s.setOnceDict[key] = val
		}
	}
}

func (s *subscriberHelper) incrementKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		isKeyValid := s._validateKeyPrefix(key, caller)
		if isKeyValid {
			s.incrementDict[key] = val
		}
	}
}

func (s *subscriberHelper) removeKV(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if s._isIdentityKey(key) {
		s.removeIdentity(key, val, kvMap, caller)
	} else {
		isKeyValid := s._validateKeyPrefix(key, caller)
		if isKeyValid {
			s.removeDict[key] = val
		}
	}
}

func (s *subscriberHelper) unsetKey(key string, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	s.unsetList = append(s.unsetList, key)
}

func (s *subscriberHelper) setPreferredLanguage(langCode string, caller string) {
	s.setDict[KEY_PREFERRED_LANGUAGE] = langCode
}

func (s *subscriberHelper) setTimezone(timezone string, caller string) {
	s.setDict[KEY_TIMEZONE] = timezone
}

func (s *subscriberHelper) addIdentity(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			s.addEmail(val.(string), newCaller)
		}
	case IDENT_KEY_SMS:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			s.addSms(val.(string), newCaller)
		}

	case IDENT_KEY_WHATSAPP:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			s.addWhatsapp(val.(string), newCaller)
		}

	case IDENT_KEY_ANDROIDPUSH:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			idProvider := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					idProvider = pvStr
				}
			}
			s.addAndroidpush(val.(string), idProvider, newCaller)
		}

	case IDENT_KEY_IOSPUSH:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			idProvider := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					idProvider = pvStr
				}
			}
			s.addIospush(val.(string), idProvider, newCaller)
		}
	case IDENT_KEY_WEBPUSH:
		val, isValid := s._checkIdentValDict(val, newCaller)
		if isValid {
			idProvider := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					idProvider = pvStr
				}
			}
			s.addWebpush(val.(map[string]interface{}), idProvider, newCaller)
		}

	case IDENT_KEY_SLACK:
		val, isValid := s._checkIdentValDict(val, newCaller)
		if isValid {
			s.addSlack(val.(map[string]interface{}), newCaller)
		}

	case IDENT_KEY_MS_TEAMS:
		val, isValid := s._checkIdentValDict(val, newCaller)
		if isValid {
			s.addMSTeams(val.(map[string]interface{}), newCaller)
		}
	}
}

func (s *subscriberHelper) removeIdentity(key string, val interface{}, kvMap map[string]interface{}, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			s.removeEmail(val.(string), newCaller)
		}

	case IDENT_KEY_SMS:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			s.removeSms(val.(string), newCaller)
		}

	case IDENT_KEY_WHATSAPP:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			s.removeWhatsapp(val.(string), newCaller)
		}

	case IDENT_KEY_ANDROIDPUSH:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			idProvider := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					idProvider = pvStr
				}
			}
			s.removeAndroidpush(val.(string), idProvider, newCaller)
		}

	case IDENT_KEY_IOSPUSH:
		val, isValid := s._checkIdentValString(val, newCaller)
		if isValid {
			idProvider := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					idProvider = pvStr
				}
			}
			s.removeIospush(val.(string), idProvider, newCaller)
		}
	case IDENT_KEY_WEBPUSH:
		val, isValid := s._checkIdentValDict(val, newCaller)
		if isValid {
			idProvider := ""
			if pv, found := kvMap[KEY_ID_PROVIDER]; found {
				if pvStr, ok := pv.(string); ok {
					idProvider = pvStr
				}
			}
			s.removeWebpush(val.(map[string]interface{}), idProvider, newCaller)
		}

	case IDENT_KEY_SLACK:
		val, isValid := s._checkIdentValDict(val, newCaller)
		if isValid {
			s.removeSlack(val.(map[string]interface{}), newCaller)
		}

	case IDENT_KEY_MS_TEAMS:
		val, isValid := s._checkIdentValDict(val, newCaller)
		if isValid {
			s.removeMSTeams(val.(map[string]interface{}), newCaller)
		}
	}
}

// ------------------------

func (s *subscriberHelper) _checkIdentValString(value interface{}, caller string) (interface{}, bool) {
	msg := "value must be a string with proper value"
	if vstring, ok := value.(string); !ok {
		s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		vstring = strings.TrimSpace(vstring)
		if vstring == "" {
			s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
			return value, false
		}
		return vstring, true
	}
}

func (s *subscriberHelper) _checkIdentValDict(value interface{}, caller string) (interface{}, bool) {
	msg := "value must be a valid dict/map"
	if valMap, ok := value.(map[string]interface{}); !ok {
		s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		if len(valMap) == 0 {
			s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
			return value, false
		}
		return valMap, true
	}
}

// ------------------------ Email

func (s *subscriberHelper) _validateEmail(email string, caller string) (string, bool) {
	iEmail, isValid := s._checkIdentValString(email, caller)
	if !isValid {
		return email, false
	}
	email = iEmail.(string)
	// --- validate basic regex
	msg := "value in email format required. e.g. user@example.com"
	minLength, maxLength := 6, 127
	// ---
	if !emailRegexCompiled.MatchString(email) {
		s._errors = append(s._errors, fmt.Sprintf("[%s] invalid value %s. %s", caller, email, msg))
		return email, false
	}
	if len(email) < minLength || len(email) > maxLength {
		s._errors = append(s._errors, fmt.Sprintf("[%s] invalid value %s. must be 6 <= len(email) <= 127", caller, email))
		return email, false
	}
	// ----
	return email, true
}

func (s *subscriberHelper) addEmail(value string, caller string) {
	val, isValid := s._validateEmail(value, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_EMAIL] = val
}

func (s *subscriberHelper) removeEmail(value string, caller string) {
	val, isValid := s._validateEmail(value, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_EMAIL] = val
}

// ------------------------ Mobile no

func (s *subscriberHelper) _validateMobileNo(mobileNo string, caller string) (string, bool) {
	iMobileNo, isValid := s._checkIdentValString(mobileNo, caller)
	if !isValid {
		return mobileNo, false
	}
	mobileNo = iMobileNo.(string)
	// --- validate basic regex
	msg := "number must start with + and must contain country code. e.g. +41446681800"
	minLength := 8
	// ---
	if !mobileRegexCompiled.MatchString(mobileNo) {
		s._errors = append(s._errors, fmt.Sprintf("[%s] invalid value %s. %s", caller, mobileNo, msg))
		return mobileNo, false
	}
	if len(mobileNo) < minLength {
		s._errors = append(s._errors, fmt.Sprintf("[%s] invalid value %s. len(mobile_no) must be >= 8", caller, mobileNo))
		return mobileNo, false
	}
	// ----
	return mobileNo, true
}

// ------------------------ SMS

func (s *subscriberHelper) addSms(value string, caller string) {
	val, isValid := s._validateMobileNo(value, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_SMS] = val
}

func (s *subscriberHelper) removeSms(value string, caller string) {
	val, isValid := s._validateMobileNo(value, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_SMS] = val
}

// ------------------------ Whatsapp

func (s *subscriberHelper) addWhatsapp(value string, caller string) {
	val, isValid := s._validateMobileNo(value, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_WHATSAPP] = val
}

func (s *subscriberHelper) removeWhatsapp(value string, caller string) {
	val, isValid := s._validateMobileNo(value, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_WHATSAPP] = val
}

// ------------------------ Androidpush [providers: fcm / xiaomi / oppo]

func (s *subscriberHelper) _checkAndroidpushValue(value string, provider string, caller string) (string, string, bool) {
	iValue, isValid := s._checkIdentValString(value, caller)
	if !isValid {
		return value, provider, false
	}
	value = iValue.(string)
	// -- validate provider
	if provider == "" {
		provider = "fcm"
	}
	// convert to lowercase to make it case-insensitive
	provider = strings.ToLower(provider)
	if !slices.Contains([]string{"fcm", "xiaomi", "oppo"}, provider) {
		s._errors = append(s._errors, fmt.Sprintf("[%s] unsupported androidpush provider %s", caller, provider))
		return value, provider, false
	}
	return value, provider, true
}

func (s *subscriberHelper) addAndroidpush(value string, provider string, caller string) {
	value, provider, isValid := s._checkAndroidpushValue(value, provider, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_ANDROIDPUSH] = value
	s.appendDict[KEY_ID_PROVIDER] = provider
}

func (s *subscriberHelper) removeAndroidpush(value string, provider string, caller string) {
	value, provider, isValid := s._checkAndroidpushValue(value, provider, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_ANDROIDPUSH] = value
	s.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Iospush [providers: apns]

func (s *subscriberHelper) _checkIospushValue(value string, provider string, caller string) (string, string, bool) {
	iValue, isValid := s._checkIdentValString(value, caller)
	if !isValid {
		return value, provider, false
	}
	value = iValue.(string)
	// -- validate provider
	if provider == "" {
		provider = "apns"
	}
	// convert to lowercase to make it case-insensitive
	provider = strings.ToLower(provider)
	if !slices.Contains([]string{"apns"}, provider) {
		s._errors = append(s._errors, fmt.Sprintf("[%s] unsupported iospush provider %s", caller, provider))
		return value, provider, false
	}
	return value, provider, true
}

func (s *subscriberHelper) addIospush(value string, provider string, caller string) {
	value, provider, isValid := s._checkIospushValue(value, provider, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_IOSPUSH] = value
	s.appendDict[KEY_ID_PROVIDER] = provider
}

func (s *subscriberHelper) removeIospush(value string, provider string, caller string) {
	value, provider, isValid := s._checkIospushValue(value, provider, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_IOSPUSH] = value
	s.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Webpush [providers: vapid]

func (s *subscriberHelper) _checkWebpushDict(value map[string]interface{}, provider string, caller string) (interface{}, string, bool) {
	msg := "value must be a valid dict/map representing webpush-token"
	if len(value) == 0 {
		s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, provider, false
	}
	// -- validate provider
	if provider == "" {
		provider = "vapid"
	}
	// convert to lowercase to make it case-insensitive
	provider = strings.ToLower(provider)
	if !slices.Contains([]string{"vapid"}, provider) {
		s._errors = append(s._errors, fmt.Sprintf("[%s] unsupported webpush provider %s", caller, provider))
		return value, provider, false
	}
	return value, provider, true
}

func (s *subscriberHelper) addWebpush(value map[string]interface{}, provider string, caller string) {
	iValue, provider, isValid := s._checkWebpushDict(value, provider, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_WEBPUSH] = iValue.(map[string]interface{})
	s.appendDict[KEY_ID_PROVIDER] = provider
}

func (s *subscriberHelper) removeWebpush(value map[string]interface{}, provider string, caller string) {
	iValue, provider, isValid := s._checkWebpushDict(value, provider, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_WEBPUSH] = iValue.(map[string]interface{})
	s.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Slack

func (s *subscriberHelper) _checkSlackDict(value map[string]interface{}, caller string) (map[string]interface{}, bool) {
	msg := "value must be a valid dict/map with proper keys"
	if len(value) == 0 {
		s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		return value, true
	}
}

func (s *subscriberHelper) addSlack(value map[string]interface{}, caller string) {
	value, isValid := s._checkSlackDict(value, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_SLACK] = value
}

func (s *subscriberHelper) removeSlack(value map[string]interface{}, caller string) {
	value, isValid := s._checkSlackDict(value, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_SLACK] = value
}

// ------------------------ MS Teams

func (s *subscriberHelper) _checkMSTeamsDict(value map[string]interface{}, caller string) (map[string]interface{}, bool) {
	msg := "value must be a valid dict/map with proper keys"
	if len(value) == 0 {
		s._errors = append(s._errors, fmt.Sprintf("[%s] %s", caller, msg))
		return value, false
	} else {
		return value, true
	}
}

func (s *subscriberHelper) addMSTeams(value map[string]interface{}, caller string) {
	value, isValid := s._checkMSTeamsDict(value, caller)
	if !isValid {
		return
	}
	s.appendDict[IDENT_KEY_MS_TEAMS] = value
}

func (s *subscriberHelper) removeMSTeams(value map[string]interface{}, caller string) {
	value, isValid := s._checkMSTeamsDict(value, caller)
	if !isValid {
		return
	}
	s.removeDict[IDENT_KEY_MS_TEAMS] = value
}
