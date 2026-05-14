package suprsend

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// AppInfo describes the user's application integrating the SDK.
// It is optional and shown in the User-Agent header when set.
type AppInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// userAgentInfo is the structured runtime fingerprint serialized as JSON
// into the X-Suprsend-Client-User-Agent header.
type userAgentInfo struct {
	SDK         string `json:"sdk"`
	SDKVersion  string `json:"sdk_version"`
	Lang        string `json:"lang"`
	LangVersion string `json:"lang_version,omitempty"`
	// platform: server / android / ios / react-native / flutter / browser / macos
	Platform string `json:"platform"`
	// environment: web/mobile/desktop
	Environment *string `json:"environment,omitempty"`
	// linux / darwin / windows / android / ios
	OS        *string  `json:"os,omitempty"`
	OSVersion *string  `json:"os_version,omitempty"`
	AppInfo   *AppInfo `json:"app_info,omitempty"`
	// for mobile only (e.g "Pixel 8")
	DeviceModel *string `json:"device_model,omitempty"`
	// For browser (js-sdk only)
	Browser        *string `json:"browser,omitempty"` // chrome/firebox/safari
	BrowserVersion *string `json:"browser_version,omitempty"`
}

func formatAppInfo(info *AppInfo) string {
	if info == nil || strings.TrimSpace(info.Name) == "" {
		return ""
	}
	s := strings.TrimSpace(info.Name)
	if strings.TrimSpace(info.Version) != "" {
		s += "/" + strings.TrimSpace(info.Version)
	}
	return s
}

// buildUserAgent returns the human-readable User-Agent header value and
// the JSON-encoded X-Suprsend-Client-User-Agent value.
func buildUserAgent(appInfo *AppInfo) (userAgent string, clientUserAgent string) {
	info := userAgentInfo{
		SDK:         "suprsend-go",
		SDKVersion:  VERSION,
		Lang:        "go",
		LangVersion: runtime.Version(),
		Platform:    "server",
		OS:          String(runtime.GOOS),
		OSVersion:   String(""),
		AppInfo:     appInfo,
	}
	if buf, err := json.Marshal(info); err == nil {
		clientUserAgent = string(buf)
	}
	userAgent = fmt.Sprintf("%s/%s (go/%s; %s)", info.SDK, info.SDKVersion, info.LangVersion, *info.OS)
	if app := formatAppInfo(appInfo); app != "" {
		userAgent += " (" + app + ")"
	}
	return userAgent, clientUserAgent
}
