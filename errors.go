package suprsend

import "errors"

var (
	ErrMissingAPIKey    = errors.New("suprsend: missing api_key")
	ErrMissingAPISecret = errors.New("suprsend: missing api_secret")
	ErrMissingBaseUrl   = errors.New("suprsend: missing base_url")
)
