package suprsend

type Response struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type BulkResponse struct {
	Status        string
	FailedRecords []map[string]interface{}
	Total         int
	Success       int
	Failure       int
	Warnings      []string
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
	failedRecords []map[string]interface{}
}
