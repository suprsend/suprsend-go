package signature

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

func GetRequestSignature(urlStr string, httpVerb string, content any,
	headers map[string]string, secret string,
) ([]byte, string, error) {
	//
	var contentBody []byte
	var contentMd5 string
	// possible methods: POST/GET/PUT
	if httpVerb == "GET" || SafeCheckNil(content) {
		contentBody, contentMd5 = []byte(""), ""
	} else {
		cBytes, err := json.Marshal(content)
		if err != nil {
			return contentBody, contentMd5, fmt.Errorf("failed to marshal content: %w", err)
		}
		contentBody = cBytes
		// MD5 of the content
		md5Hash := md5.Sum(cBytes)
		contentMd5 = hex.EncodeToString(md5Hash[:])
	}
	// Proper Url encoding
	requestUrlPath, err := getUrlPath(urlStr)
	if err != nil {
		return contentBody, contentMd5, err
	}
	// Create string to sign
	stringToSign := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		httpVerb,
		contentMd5,
		headers["Content-Type"],
		headers["Date"],
		requestUrlPath,
	)
	// fmt.Printf("stringToSign:\n%s", stringToSign)
	// ----- HMAC-SHA-256
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(stringToSign))
	// signature
	sig := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	//
	return contentBody, sig, nil
}

func getUrlPath(urlStr string) (string, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return "", err
	}
	uri := u.Path
	if u.RawQuery != "" {
		uri = fmt.Sprintf("%s?%s", uri, u.Query().Encode())
	}
	return uri, nil
}

func SafeCheckNil(input any) bool {
	if input == nil {
		return true
	}
	v := reflect.ValueOf(input)
	// Only call IsNil on types that can be nil
	return (v.Kind() == reflect.Ptr || v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map || v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Interface || v.Kind() == reflect.Func) &&
		v.IsNil()
}
