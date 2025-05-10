package suprsend

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jinzhu/copier"
)

type BulkUsersEdit interface {
	Append(users ...UserEdit)
	Save() (*BulkResponse, error)
}

var _ BulkUsersEdit = &bulkUsersEdit{}

type bulkUsersEdit struct {
	client *Client
	//
	_users          []userEdit
	_pendingRecords []pendingIdentityEventRecord2
	chunks          []*bulkUsersEditChunk
	//
	response *BulkResponse
	// invalid_record json: {"record": event-json, "error": error_str, "code": 500}
	_invalidRecords []map[string]any
}

func newBulkUsersEdit(client *Client) BulkUsersEdit {
	u := &bulkUsersEdit{
		client:   client,
		response: &BulkResponse{},
	}
	return u
}

type pendingIdentityEventRecord2 struct {
	record     map[string]any
	recordSize int
}

func (b *bulkUsersEdit) _validateUsers() {
	for _, u := range b._users {
		// -- check if there is any error/warning, if so add it to warnings list of BulkResponse
		warningsList := u.validateBody()
		if len(warningsList) > 0 {
			b.response.Warnings = append(b.response.Warnings, warningsList...)
		}
		//
		pl := u.GetAsyncPayload()
		plJson, plSize, err := u.validatePayloadSize(pl)
		if err != nil {
			invRec := invalidRecordJson(u.asJsonAsync(), err)
			b._invalidRecords = append(b._invalidRecords, invRec)
		} else {
			b._pendingRecords = append(
				b._pendingRecords,
				pendingIdentityEventRecord2{
					record:     plJson,
					recordSize: plSize,
				},
			)
		}
	}
}

func (b *bulkUsersEdit) _chunkify(startIdx int) {
	currChunk := newBulkUsersEditChunk(b.client)
	b.chunks = append(b.chunks, currChunk)
	for relIdx, rec := range b._pendingRecords[startIdx:] {
		isAdded := currChunk.tryToAddIntoChunk(rec.record, rec.recordSize)
		if !isAdded {
			// create chunks from remaining records
			b._chunkify(startIdx + relIdx)
			// Don't forget to break. As current loop must not continue further
			break
		}
	}
}

func (b *bulkUsersEdit) Append(users ...UserEdit) {
	for _, u := range users {
		if u == nil {
			continue
		}
		if ue, ok := u.(*userEdit); ok {
			ueCopy := userEdit{}
			copier.CopyWithOption(&ueCopy, ue, copier.Option{DeepCopy: true})
			b._users = append(b._users, ueCopy)
		}
	}
}

func (b *bulkUsersEdit) Save() (*BulkResponse, error) {
	b._validateUsers()
	if len(b._invalidRecords) > 0 {
		chResponse := invalidRecordsChunkResponse(b._invalidRecords)
		b.response.mergeChunkResponse(chResponse)
	}
	if len(b._pendingRecords) > 0 {
		b._chunkify(0)
		for cIdx, ch := range b.chunks {
			if b.client.debug {
				log.Printf("DEBUG: triggering api call for chunk: %d", cIdx)
			}
			// do api call
			ch.trigger()
			// merge response
			b.response.mergeChunkResponse(ch.response)
		}
	} else {
		if len(b._invalidRecords) == 0 {
			b.response.mergeChunkResponse(emptyChunkSuccessResponse())
		}
	}
	return b.response, nil
}

// ==========================================================

type bulkUsersEditChunk struct {
	_chunk_apparent_size_in_bytes int
	_max_records_in_chunk         int
	//
	client *Client
	_url   string
	//
	_chunk         []map[string]any
	_runningSize   int
	_runningLength int
	response       *chunkResponse
}

func newBulkUsersEditChunk(client *Client) *bulkUsersEditChunk {
	bsc := &bulkUsersEditChunk{
		_chunk_apparent_size_in_bytes: BODY_MAX_APPARENT_SIZE_IN_BYTES,
		_max_records_in_chunk:         MAX_IDENTITY_EVENTS_IN_BULK_API,
		//
		client: client,
		_url:   fmt.Sprintf("%sevent/", client.baseUrl),
		_chunk: []map[string]any{},
	}
	return bsc
}

func (b *bulkUsersEditChunk) _addEventToChunk(event map[string]any, eventSize int) {
	// First add size, then event to reduce effects of race condition
	b._runningSize += eventSize
	b._chunk = append(b._chunk, event)
	b._runningLength += 1
}

func (b *bulkUsersEditChunk) _checkLimitReached() bool {
	return b._runningLength >= b._max_records_in_chunk || b._runningSize >= b._chunk_apparent_size_in_bytes
}

/*
returns whether passed event was able to get added to this chunk or not,
if true, event gets added to chunk
*/
func (b *bulkUsersEditChunk) tryToAddIntoChunk(event map[string]any, eventSize int) bool {
	if event == nil {
		return true
	}
	if b._checkLimitReached() {
		return false
	}
	// if apparent_size of event crosses limit
	if (b._runningSize + eventSize) > b._chunk_apparent_size_in_bytes {
		return false
	}
	// Add Event to chunk
	b._addEventToChunk(event, eventSize)
	return true
}

func (b *bulkUsersEditChunk) trigger() {
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("POST", b._url, b._chunk)
	if err != nil {
		suprResponse := b.formatAPIResponse(nil, err)
		b.response = suprResponse
	}
	//
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		suprResponse := b.formatAPIResponse(nil, err)
		b.response = suprResponse

	} else {
		defer httpResponse.Body.Close()
		suprResponse := b.formatAPIResponse(httpResponse, nil)
		b.response = suprResponse
	}
}

func (b *bulkUsersEditChunk) formatAPIResponse(httpRes *http.Response, err error) *chunkResponse {
	//
	bulkRespFunc := func(statusCode int, errMsg string) *chunkResponse {
		failedRecords := []map[string]any{}
		if statusCode >= 400 {
			for _, c := range b._chunk {
				failedRecords = append(failedRecords,
					map[string]any{
						"record": c,
						"error":  errMsg,
						"code":   statusCode,
					})
			}
			return &chunkResponse{
				status: "fail", statusCode: statusCode,
				total: len(b._chunk), success: 0, failure: len(b._chunk),
				failedRecords: failedRecords,
			}
		} else {
			return &chunkResponse{
				status: "success", statusCode: statusCode,
				total: len(b._chunk), success: len(b._chunk), failure: 0,
				failedRecords: failedRecords,
			}
		}
	}
	if err != nil {
		return bulkRespFunc(500, err.Error())

	} else if httpRes != nil {
		respBody, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return bulkRespFunc(500, err.Error())
		}
		//
		return bulkRespFunc(httpRes.StatusCode, string(respBody))
	}
	return nil
}
