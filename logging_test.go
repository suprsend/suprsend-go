package suprsend

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSanitizedHeadersRedactsSensitiveValues(t *testing.T) {
	tests := []struct {
		name   string
		header string
		value  string
		want   string
	}{
		{"bearer", "Authorization", "Bearer ss_api_key_xyz", "Bearer [REDACTED]"},
		{"bearer lowercase", "Authorization", "bearer ss_api_key_xyz", "bearer [REDACTED]"},
		{"bearer uppercase", "Authorization", "BEARER ss_api_key_xyz", "BEARER [REDACTED]"},
		{"ws key secret", "Authorization", "ws_key_abc:signature", "[REDACTED]"},
		{"proxy bearer", "Proxy-Authorization", "Bearer proxy-token", "Bearer [REDACTED]"},
		{"proxy basic", "Proxy-Authorization", "Basic abcdef", "[REDACTED]"},
		{"cookie", "Cookie", "session=abc", "[REDACTED]"},
		{"set-cookie", "Set-Cookie", "session=abc", "[REDACTED]"},
		{"content-type", "Content-Type", "application/json", "application/json"},
		{"user-agent", "User-Agent", "suprsend-go", "suprsend-go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			h.Set(tt.header, tt.value)
			got := sanitizedHeaders(h)
			if got.Get(tt.header) != tt.want {
				t.Errorf("%s = %q, want %q", tt.header, got.Get(tt.header), tt.want)
			}
			if h.Get(tt.header) != tt.value {
				t.Errorf("original %s mutated: got %q, want %q", tt.header, h.Get(tt.header), tt.value)
			}
		})
	}
}

func TestSanitizedHeadersNilAndEmpty(t *testing.T) {
	if got := sanitizedHeaders(nil); got != nil {
		t.Errorf("sanitizedHeaders(nil) = %v, want nil", got)
	}
	got := sanitizedHeaders(http.Header{})
	if len(got) != 0 {
		t.Errorf("sanitizedHeaders(empty) = %v, want empty", got)
	}
}

func TestSanitizedHeadersRedactsAllValues(t *testing.T) {
	h := http.Header{}
	h.Add("Authorization", "Bearer first")
	h.Add("Authorization", "Bearer second")
	got := sanitizedHeaders(h)
	vals := got.Values("Authorization")
	if len(vals) != 2 {
		t.Fatalf("got %d values, want 2", len(vals))
	}
	for i, v := range vals {
		if v != "Bearer [REDACTED]" {
			t.Errorf("Authorization[%d] = %q, want %q", i, v, "Bearer [REDACTED]")
		}
	}
}

func TestLoggingRoundTripperRedactsAuthorizationInLogs(t *testing.T) {
	var logged bytes.Buffer
	log.SetOutput(&logged)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	})

	const secret = "ss_api_key_xyz"
	var proxiedAuth string
	rt := LoggingRoundTripper{
		Proxied: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			proxiedAuth = req.Header.Get("Authorization")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	req, err := http.NewRequest(http.MethodPost, "https://hub.suprsend.com/v1/user/", strings.NewReader(`{"a":1}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Content-Type", "application/json")

	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if proxiedAuth != "Bearer "+secret {
		t.Errorf("proxied Authorization = %q, want original value", proxiedAuth)
	}

	out := logged.String()
	if strings.Contains(out, secret) {
		t.Errorf("log leaked secret: %s", out)
	}
	if !strings.Contains(out, "Bearer [REDACTED]") {
		t.Errorf("log missing redacted bearer: %s", out)
	}
	if !strings.Contains(out, "application/json") {
		t.Errorf("log missing non-sensitive header: %s", out)
	}
}

func TestLoggingRoundTripperNilBody(t *testing.T) {
	var logged bytes.Buffer
	log.SetOutput(&logged)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	})

	var proxiedBody io.ReadCloser
	rt := LoggingRoundTripper{
		Proxied: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			proxiedBody = req.Body
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	req, err := http.NewRequest(http.MethodGet, "https://hub.suprsend.com/v1/user/", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proxiedBody != nil {
		t.Error("proxied request Body was replaced; want nil")
	}
}
