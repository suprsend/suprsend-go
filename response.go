package suprsend

import (
	"fmt"
)

type Response struct {
	Success     bool           `json:"success"`
	StatusCode  int            `json:"status_code"`
	Message     string         `json:"message"`
	RawResponse map[string]any `json:"raw_response"`
}

func (r *Response) String() string {
	return fmt.Sprintf("Response{Success: %v, StatusCode: %v, Message: %v}",
		r.Success, r.StatusCode, r.Message)
}

type BulkResponse struct {
	Status        string
	FailedRecords []map[string]any
	Total         int
	Success       int
	Failure       int
	Warnings      []string
}

func (b *BulkResponse) String() string {
	return fmt.Sprintf("BulkResponse{Status: %v, Total: %v, Success: %v, Failure: %v, Warnings: %v}",
		b.Status, b.Total, b.Success, b.Failure, len(b.Warnings))
}

func (b *BulkResponse) mergeChunkResponse(chResponse *chunkResponse) {
	if chResponse == nil {
		return
	}
	// possible status: success/partial/fail
	if b.Status == "" {
		b.Status = chResponse.status
	} else {
		if chResponse.status == "partial" {
			b.Status = "partial"
		} else if b.Status == "success" {
			if chResponse.status == "fail" {
				b.Status = "partial"
			}
		} else if b.Status == "fail" {
			if chResponse.status == "success" {
				b.Status = "partial"
			}
		}
	}
	b.Total += chResponse.total
	b.Success += chResponse.success
	b.Failure += chResponse.failure
	b.FailedRecords = append(b.FailedRecords, chResponse.failedRecords...)
}

type chunkResponse struct {
	status        string
	statusCode    int
	total         int
	success       int
	failure       int
	failedRecords []map[string]any
	rawResponse   map[string]any
}

func emptyChunkSuccessResponse() *chunkResponse {
	return &chunkResponse{
		status:        "success",
		statusCode:    200,
		total:         0,
		success:       0,
		failure:       0,
		failedRecords: []map[string]any{},
		rawResponse:   nil,
	}
}

func invalidRecordsChunkResponse(invalidRecords []map[string]any) *chunkResponse {
	return &chunkResponse{
		status:        "fail",
		statusCode:    500,
		total:         len(invalidRecords),
		success:       0,
		failure:       len(invalidRecords),
		failedRecords: invalidRecords,
		rawResponse:   nil,
	}
}

type v2EventSingleResponse struct {
	Status    string `json:"status" mapstructure:"status"`
	MessageId string `json:"message_id" mapstructure:"message_id"`
	Error     *struct {
		Message string `json:"message" mapstructure:"message"`
		Type    string `json:"type" mapstructure:"type"`
	} `json:"error" mapstructure:"error"`
	// in case of bulk, each record has its own status-code
	StatusCode int `json:"status_code" mapstructure:"status_code"`
}

type v2EventBulkResponse struct {
	Status  string                  `json:"status" mapstructure:"status"`
	Records []v2EventSingleResponse `json:"records" mapstructure:"records"`
	Error   *struct {
		Message string `json:"message" mapstructure:"message"`
		Type    string `json:"type" mapstructure:"type"`
	} `json:"error" mapstructure:"error"`
	// derived fields
	dStatus  string `json:"-"`
	dTotal   int    `json:"-"`
	dSuccess int    `json:"-"`
	dFailure int    `json:"-"`
}

func (b *v2EventBulkResponse) setDerivedFields() {
	b.dTotal = len(b.Records)
	for _, r := range b.Records {
		if r.Status == "success" {
			b.dSuccess++
		}
	}
	b.dFailure = b.dTotal - b.dSuccess
	if b.dFailure > 0 {
		if b.dSuccess > 0 {
			b.dStatus = "partial"
		} else {
			b.dStatus = "fail"
		}
	} else {
		b.dStatus = "success"
	}
}
