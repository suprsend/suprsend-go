package suprsend

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	// read request body for logging
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		return nil, err
	}
	// prepare request log
	logString := "DEBUG: HTTP Request ------------------\n" +
		"METHOD:\t%v\nURL:\t%v\nHEADER\t%v\nBODY:\t%v\n" +
		"------------------\n"
	log.Printf(logString, req.Method, req.URL, req.Header, string(body))

	// Set new body
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	//
	res, e = l.Proxied.RoundTrip(req)
	return
}
