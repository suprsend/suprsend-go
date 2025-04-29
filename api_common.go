package suprsend

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
	// other filter params can added to more
	More map[string]string
}
