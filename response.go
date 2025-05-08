package suprsend

import "fmt"

type Response struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
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
		if b.Status == "success" {
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
}

func emptyChunkSuccessResponse() *chunkResponse {
	return &chunkResponse{
		status:        "success",
		statusCode:    200,
		total:         0,
		success:       0,
		failure:       0,
		failedRecords: []map[string]any{},
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
	}
}
