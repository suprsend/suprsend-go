package suprsend

import (
	"net/url"
	"strconv"
)

type ListApiMetaInfo struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CursorListApiMetaInfo struct {
	Limit   int    `json:"limit"`
	Count   int    `json:"count"`
	Before  string `json:"before"`
	After   string `json:"after"`
	HasPrev bool   `json:"has_prev"`
	HasNext bool   `json:"has_next"`
}

type CursorListApiResponse struct {
	Meta    *CursorListApiMetaInfo `json:"meta"`
	Results []map[string]any       `json:"results"`
}

type CursorListApiOptions struct {
	Limit  int
	Before string
	After  string
	// add more filters like this: {"key": "val1"}
	ExtraParams map[string]string
	// For multivalue params: {"key2[]": ["val2", "val3"]}
	MultiValueParams map[string][]string
}

func (o *CursorListApiOptions) BuildQuery() string {
	if o == nil {
		return ""
	}
	params := url.Values{}
	if o.Limit > 0 {
		params.Add("limit", strconv.Itoa(o.Limit))
	}
	if o.Before != "" {
		params.Add("before", o.Before)
	}
	if o.After != "" {
		params.Add("after", o.After)
	}
	for k, v := range o.ExtraParams {
		params.Add(k, v)
	}
	for k, v := range o.MultiValueParams {
		for _, vv := range v {
			params.Add(k, vv)
		}
	}
	return params.Encode()
}
