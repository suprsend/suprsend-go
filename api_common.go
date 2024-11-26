package suprsend

type ListApiMetaInfo struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CursorPaginationMetaInfo struct {
	Limit   int    `json:"limit"`
	Count   int    `json:"count"`
	Before  string `json:"before"`
	After   string `json:"after"`
	HasPrev bool   `json:"has_prev"`
	HasNext bool   `json:"has_next"`
}

type CursorPaginationList struct {
	Meta    *CursorPaginationMetaInfo `json:"meta"`
	Results []*map[string]any         `json:"results"`
}

type CursorPaginationListOptions struct {
	Limit  int
	Before string
	After  string
}

func (c *CursorPaginationListOptions) cleanParams() {
	if c.Limit <= 0 {
		c.Limit = 20
	}
}
