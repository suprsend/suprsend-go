package suprsend

import (
	"fmt"
	"slices"
	"strings"
)

type subscriberHelper struct {
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

func newSubscriberHelper() *subscriberHelper {
	return &subscriberHelper{
		setDict:       map[string]any{},
		setOnceDict:   map[string]any{},
		incrementDict: map[string]any{},
		appendDict:    map[string]any{},
		removeDict:    map[string]any{},
		unsetList:     []string{},
		_errors:       []string{},
		_info:         []string{},
	}
}

func (s *subscriberHelper) reset() {
	s.setDict = map[string]any{}
	s.setOnceDict = map[string]any{}
	s.incrementDict = map[string]any{}
	s.appendDict = map[string]any{}
	s.removeDict = map[string]any{}
	s.unsetList = []string{}
	s._errors, s._info = []string{}, []string{}
}

type getIdentityEventResp struct {
	errors []string
	info   []string
	//
	event map[string]any
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

func (s *subscriberHelper) _formEvent() map[string]any {
	event := map[string]any{}
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

func (s *subscriberHelper) _isIdentityKey(key string) bool {
	return slices.Contains(IDENT_KEYS_ALL, key)
}

// -------------------------

func (s *subscriberHelper) appendKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if s._isIdentityKey(key) {
		s.addIdentity(key, val, kvMap, caller)
	} else {
		s.appendDict[key] = val
	}
}

func (s *subscriberHelper) setKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		s.setDict[key] = val
	}
}

func (s *subscriberHelper) setOnceKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		s.setOnceDict[key] = val
	}
}

func (s *subscriberHelper) incrementKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	} else {
		s.incrementDict[key] = val
	}
}

func (s *subscriberHelper) removeKV(key string, val any, kvMap map[string]any, caller string) {
	key, isKeyValid := s._validateKeyBasic(key, caller)
	if !isKeyValid {
		return
	}
	if s._isIdentityKey(key) {
		s.removeIdentity(key, val, kvMap, caller)
	} else {
		s.removeDict[key] = val
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

func (s *subscriberHelper) setLocale(localeCode string, caller string) {
	s.setDict[KEY_LOCALE] = localeCode
}

func (s *subscriberHelper) setTimezone(timezone string, caller string) {
	s.setDict[KEY_TIMEZONE] = timezone
}

func (s *subscriberHelper) addIdentity(key string, val any, kvMap map[string]any, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		s.addEmail(val, newCaller)

	case IDENT_KEY_SMS:
		s.addSms(val, newCaller)

	case IDENT_KEY_WHATSAPP:
		s.addWhatsapp(val, newCaller)

	case IDENT_KEY_ANDROIDPUSH:
		s.addAndroidpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_IOSPUSH:
		s.addIospush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_WEBPUSH:
		s.addWebpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_SLACK:
		s.addSlack(val, newCaller)

	case IDENT_KEY_MS_TEAMS:
		s.addMSTeams(val, newCaller)
	}
}

func (s *subscriberHelper) removeIdentity(key string, val any, kvMap map[string]any, caller string) {
	newCaller := fmt.Sprintf("%s:%s", caller, key)
	switch key {
	case IDENT_KEY_EMAIL:
		s.removeEmail(val, newCaller)

	case IDENT_KEY_SMS:
		s.removeSms(val, newCaller)

	case IDENT_KEY_WHATSAPP:
		s.removeWhatsapp(val, newCaller)

	case IDENT_KEY_ANDROIDPUSH:
		s.removeAndroidpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_IOSPUSH:
		s.removeIospush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_WEBPUSH:
		s.removeWebpush(val, kvMap[KEY_ID_PROVIDER], newCaller)

	case IDENT_KEY_SLACK:
		s.removeSlack(val, newCaller)

	case IDENT_KEY_MS_TEAMS:
		s.removeMSTeams(val, newCaller)
	}
}

// ------------------------ Email

func (s *subscriberHelper) addEmail(value any, caller string) {
	s.appendDict[IDENT_KEY_EMAIL] = value
}

func (s *subscriberHelper) removeEmail(value any, caller string) {
	s.removeDict[IDENT_KEY_EMAIL] = value
}

// ------------------------ SMS

func (s *subscriberHelper) addSms(value any, caller string) {
	s.appendDict[IDENT_KEY_SMS] = value
}

func (s *subscriberHelper) removeSms(value any, caller string) {
	s.removeDict[IDENT_KEY_SMS] = value
}

// ------------------------ Whatsapp

func (s *subscriberHelper) addWhatsapp(value any, caller string) {
	s.appendDict[IDENT_KEY_WHATSAPP] = value
}

func (s *subscriberHelper) removeWhatsapp(value any, caller string) {
	s.removeDict[IDENT_KEY_WHATSAPP] = value
}

// ------------------------ Androidpush

func (s *subscriberHelper) addAndroidpush(value any, provider any, caller string) {
	s.appendDict[IDENT_KEY_ANDROIDPUSH] = value
	s.appendDict[KEY_ID_PROVIDER] = provider
}

func (s *subscriberHelper) removeAndroidpush(value any, provider any, caller string) {
	s.removeDict[IDENT_KEY_ANDROIDPUSH] = value
	s.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Iospush

func (s *subscriberHelper) addIospush(value any, provider any, caller string) {
	s.appendDict[IDENT_KEY_IOSPUSH] = value
	s.appendDict[KEY_ID_PROVIDER] = provider
}

func (s *subscriberHelper) removeIospush(value any, provider any, caller string) {
	s.removeDict[IDENT_KEY_IOSPUSH] = value
	s.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Webpush [providers: vapid]

func (s *subscriberHelper) addWebpush(value any, provider any, caller string) {
	s.appendDict[IDENT_KEY_WEBPUSH] = value
	s.appendDict[KEY_ID_PROVIDER] = provider
}

func (s *subscriberHelper) removeWebpush(value any, provider any, caller string) {
	s.removeDict[IDENT_KEY_WEBPUSH] = value
	s.removeDict[KEY_ID_PROVIDER] = provider
}

// ------------------------ Slack

func (s *subscriberHelper) addSlack(value any, caller string) {
	s.appendDict[IDENT_KEY_SLACK] = value
}

func (s *subscriberHelper) removeSlack(value any, caller string) {
	s.removeDict[IDENT_KEY_SLACK] = value
}

// ------------------------ MS Teams

func (s *subscriberHelper) addMSTeams(value any, caller string) {
	s.appendDict[IDENT_KEY_MS_TEAMS] = value
}

func (s *subscriberHelper) removeMSTeams(value any, caller string) {
	s.removeDict[IDENT_KEY_MS_TEAMS] = value
}
