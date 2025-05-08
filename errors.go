package suprsend

import (
	"encoding/json"
)

var (
	ErrInvalidAuthMethod   = &Error{Code: 400, Message: "suprsend: invalid auth_method"}
	ErrMissingAPIKey       = &Error{Code: 400, Message: "suprsend: missing api_key"}
	ErrMissingAPISecret    = &Error{Code: 400, Message: "suprsend: missing api_secret"}
	ErrMissingServiceToken = &Error{Code: 400, Message: "suprsend: missing service_token"}
	ErrMissingWorkspaceUid = &Error{Code: 400, Message: "suprsend: missing workspace_uid"}
	ErrMissingBaseUrl      = &Error{Code: 400, Message: "suprsend: missing base_url"}
)

type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Detail  json.RawMessage `json:"detail"`
	//
	Err error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	} else if e.Detail != nil {
		return string(e.Detail)
	} else if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func (e *Error) Unwrap() error {
	return e.Err
}
