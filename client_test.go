package suprsend

import (
	"errors"
	"io"
	"strings"
	"testing"
)

const (
	// A placeholder Workspace UID, from SuprSend dashboard -> Settings ->
	// General -> Workspace UID. It is 20 characters, because the "env" field of
	// request_json/event.json sets minLength 20.
	testWorkspaceUid = "abcd1234EFGH5678ijkl"
	testHttpApiKey   = "ss_api_key_xyz"
)

func TestNewClientWithWorkspaceAPIKeySetsApiKeyAuthMethod(t *testing.T) {
	c, err := NewClientWithWorkspaceAPIKey(testWorkspaceUid, testHttpApiKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.AuthMethod != AuthMethod_ApiKey {
		t.Errorf("AuthMethod = %q, want %q", c.AuthMethod, AuthMethod_ApiKey)
	}
	if c.WorkspaceUid != testWorkspaceUid {
		t.Errorf("WorkspaceUid = %q, want %q", c.WorkspaceUid, testWorkspaceUid)
	}
	if c.HttpApiKey != testHttpApiKey {
		t.Errorf("HttpApiKey = %q, want %q", c.HttpApiKey, testHttpApiKey)
	}
}

func TestNewClientWithWorkspaceAPIKeyRejectsEmptyWorkspaceUid(t *testing.T) {
	_, err := NewClientWithWorkspaceAPIKey("", testHttpApiKey)
	if !errors.Is(err, ErrMissingWorkspaceUid) {
		t.Errorf("err = %v, want ErrMissingWorkspaceUid", err)
	}
}

func TestNewClientWithWorkspaceAPIKeyRejectsEmptyApiKey(t *testing.T) {
	_, err := NewClientWithWorkspaceAPIKey(testWorkspaceUid, "")
	if !errors.Is(err, ErrMissingAPIKey) {
		t.Errorf("err = %v, want ErrMissingAPIKey", err)
	}
}

func apiKeyTestClient(t *testing.T) *Client {
	t.Helper()
	c, err := NewClientWithWorkspaceAPIKey(testWorkspaceUid, testHttpApiKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return c
}

func TestApiKeyAuthSetsBearerAuthorizationHeader(t *testing.T) {
	c := apiKeyTestClient(t)
	req, err := c.prepareHttpRequest(t.Context(), "GET", "https://hub.suprsend.com/v1/user/", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := req.Header.Get("Authorization"), "Bearer "+testHttpApiKey; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

func TestApiKeyAuthDoesNotSendDateHeader(t *testing.T) {
	c := apiKeyTestClient(t)
	req, err := c.prepareHttpRequest(t.Context(), "GET", "https://hub.suprsend.com/v1/user/", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := req.Header.Get("Date"); got != "" {
		t.Errorf("Date = %q, want empty", got)
	}
}

func TestApiKeyAuthSendsJsonEncodedBodyOnPost(t *testing.T) {
	c := apiKeyTestClient(t)
	body := map[string]any{"distinct_id": "user-1"}
	req, err := c.prepareHttpRequest(t.Context(), "POST", "https://hub.suprsend.com/v1/user/", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sent, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := string(sent), `{"distinct_id":"user-1"}`; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestApiKeyAuthSendsEmptyBodyOnGet(t *testing.T) {
	c := apiKeyTestClient(t)
	req, err := c.prepareHttpRequest(t.Context(), "GET", "https://hub.suprsend.com/v1/user/", map[string]any{"a": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sent, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sent) != 0 {
		t.Errorf("body = %q, want empty", string(sent))
	}
}

func TestWsKeySecretAuthStillSignsRequest(t *testing.T) {
	c, err := NewClient("ws_key_abc", "ws_secret_xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req, err := c.prepareHttpRequest(t.Context(), "GET", "https://hub.suprsend.com/v1/user/", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "ws_key_abc:") {
		t.Errorf("Authorization = %q, want prefix %q", auth, "ws_key_abc:")
	}
	if req.Header.Get("Date") == "" {
		t.Error("Date header is empty, want an RFC1123 value")
	}
}

func TestApiKeyAuthUsesWorkspaceUidAsWorkspaceIdentifier(t *testing.T) {
	c := apiKeyTestClient(t)
	if got, want := c.getWsIdentifierValue(), testWorkspaceUid; got != want {
		t.Errorf("getWsIdentifierValue() = %q, want %q", got, want)
	}
}

func TestApiKeyAuthPutsWorkspaceUidInBroadcastUrl(t *testing.T) {
	c := apiKeyTestClient(t)
	want := DEFAULT_URL + testWorkspaceUid + "/broadcast/"
	if got := c.SubscriberLists._broadcastUrl; got != want {
		t.Errorf("broadcast url = %q, want %q", got, want)
	}
}

func TestApiKeyAuthPutsWorkspaceUidInEventEnv(t *testing.T) {
	c := apiKeyTestClient(t)
	ev := &Event{DistinctId: "user-1", EventName: "product_purchased", Properties: map[string]any{}}
	body, _, err := ev.getFinalJson(c, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := body["env"], testWorkspaceUid; got != want {
		t.Errorf("env = %v, want %q", got, want)
	}
}
