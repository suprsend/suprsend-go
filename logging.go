package suprsend

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
)

const redactedHeaderValue = "[REDACTED]"

var sensitiveHeaders = map[string]struct{}{
	"Authorization":       {},
	"Proxy-Authorization": {},
	"Cookie":              {},
	"Set-Cookie":          {},
}

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			log.Printf("DEBUG: error reading request body: %v", err)
			return nil, err
		}
		// Set new body
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	logString := "DEBUG: HTTP Request ------------------\n" +
		"METHOD:\t%v\nURL:\t%v\nHEADER\t%v\nBODY:\t%v\n" +
		"------------------\n"
	log.Printf(logString, req.Method, req.URL, sanitizedHeaders(req.Header), string(body))

	res, e = l.Proxied.RoundTrip(req)
	return
}

// sanitizedHeaders returns a copy of h with sensitive header values redacted.
// The original header is not modified.
func sanitizedHeaders(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	out := h.Clone()
	for k, vals := range out {
		canonical := http.CanonicalHeaderKey(k)
		if _, ok := sensitiveHeaders[canonical]; !ok {
			continue
		}
		redacted := make([]string, len(vals))
		for i, v := range vals {
			redacted[i] = redactHeaderValue(canonical, v)
		}
		out[k] = redacted
	}
	return out
}

func redactHeaderValue(canonicalKey, value string) string {
	if slices.Contains([]string{"Authorization", "Proxy-Authorization"}, canonicalKey) {
		const bearerPrefix = "Bearer "
		if len(value) >= len(bearerPrefix) && strings.EqualFold(value[:len(bearerPrefix)], bearerPrefix) {
			return value[:len(bearerPrefix)] + redactedHeaderValue
		}
	}
	return redactedHeaderValue
}
