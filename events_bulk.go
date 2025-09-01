package suprsend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jinzhu/copier"
)

type bulkEventsService struct {
	client *Client
}

func (b *bulkEventsService) NewInstance() BulkEvents {
	return &bulkEvents{
		client:   b.client,
		response: &BulkResponse{},
	}
}

type BulkEvents interface {
	Append(...*Event)
	Trigger() (*BulkResponse, error)
}

var _ BulkEvents = &bulkEvents{}

type bulkEvents struct {
	client *Client
	//
	_events         []Event
	_pendingRecords []pendingEventRecord
	chunks          []*bulkEventsChunk
	//
	response *BulkResponse
	// invalid_record json: {"record": event-json, "error": error_str, "code": 500}
	_invalidRecords []map[string]any
}

type pendingEventRecord struct {
	record     map[string]any
	recordSize int
}

func (b *bulkEvents) _validateEvents() {
	for _, ev := range b._events {
		evJson, bodySize, err := ev.getFinalJson(b.client, true)
		if err != nil {
			invRec := invalidRecordJson(ev.asJson(), err)
			b._invalidRecords = append(b._invalidRecords, invRec)
		} else {
			b._pendingRecords = append(
				b._pendingRecords,
				pendingEventRecord{
					record:     evJson,
					recordSize: bodySize,
				},
			)
		}
	}
}

func (b *bulkEvents) _chunkify(startIdx int) {
	currChunk := newBulkEventsChunk(b.client)
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

func (b *bulkEvents) Append(events ...*Event) {
	for _, ev := range events {
		if ev == nil {
			continue
		}
		eventCopy := Event{}
		copier.CopyWithOption(&eventCopy, ev, copier.Option{DeepCopy: true})
		b._events = append(b._events, eventCopy)
	}
}

func (b *bulkEvents) Trigger() (*BulkResponse, error) {
	b._validateEvents()
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

type bulkEventsChunk struct {
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

func newBulkEventsChunk(client *Client) *bulkEventsChunk {
	bec := &bulkEventsChunk{
		_chunk_apparent_size_in_bytes: BODY_MAX_APPARENT_SIZE_IN_BYTES,
		_max_records_in_chunk:         MAX_EVENTS_IN_BULK_API,
		//
		client: client,
		_url:   fmt.Sprintf("%sv2/bulk/event/", client.baseUrl),
		_chunk: []map[string]any{},
	}
	return bec
}

func (b *bulkEventsChunk) _addEventToChunk(event map[string]any, eventSize int) {
	// First add size, then event to reduce effects of race condition
	b._runningSize += eventSize
	b._chunk = append(b._chunk, event)
	b._runningLength += 1
}

func (b *bulkEventsChunk) _checkLimitReached() bool {
	return b._runningLength >= b._max_records_in_chunk || b._runningSize >= b._chunk_apparent_size_in_bytes
}

/*
returns whether passed event was able to get added to this chunk or not,
if true, event gets added to chunk
*/
func (b *bulkEventsChunk) tryToAddIntoChunk(event map[string]any, eventSize int) bool {
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
	if !ALLOW_ATTACHMENTS_IN_BULK_API {
		delete(event["properties"].(map[string]any), "$attachments")
	}
	// Add Event to chunk
	b._addEventToChunk(event, eventSize)
	return true
}

func (b *bulkEventsChunk) trigger() {
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("POST", b._url, b._chunk)
	if err != nil {
		suprResponse := parseV2BulkEventResponse(nil, err, b._chunk)
		b.response = suprResponse
	}
	httpResponse, err := b.client.httpClient.Do(request)
	if err != nil {
		suprResponse := parseV2BulkEventResponse(nil, err, b._chunk)
		b.response = suprResponse

	} else {
		defer httpResponse.Body.Close()
		suprResponse := parseV2BulkEventResponse(httpResponse, nil, b._chunk)
		b.response = suprResponse
	}
}

// Used by bulk apis: /v2/bulk/event/ and /trigger/ endpoints
func parseV2BulkEventResponse(httpRes *http.Response, err error, _chunk []map[string]any) *chunkResponse {
	/*
		"string"
		OR
		{"status": "error", "error": {"message": "string", "type": "string"}}
		{"status": "success", records: [
			{"status": "success", "message_id": "string", "status_code": "string"},
			{"status": "error", "error": {"message": "string", "type": "string"}, "status_code": "string"}
		]}
	*/
	bulkRespFunc := func(statusCode int, errMsg string, respPtr *v2EventBulkResponse) *chunkResponse {
		failedRecords := []map[string]any{}
		if statusCode >= 400 {
			// pick error message from response pointer if present
			if respPtr != nil && respPtr.Error != nil {
				errMsg = respPtr.Error.Message
			}
			for _, c := range _chunk {
				failedRecords = append(failedRecords,
					map[string]any{
						"record": c,
						"error":  errMsg,
						"code":   statusCode,
					})
			}
			return &chunkResponse{
				status: "fail", statusCode: statusCode,
				total: len(_chunk), success: 0, failure: len(_chunk),
				failedRecords: failedRecords,
			}
		} else {
			// multi-status 207 response. Filter failed records
			for ri, r := range respPtr.Records {
				if r.Status == "error" {
					failedR := map[string]any{
						"record": nil,
						"error":  r.Error.Message,
						"code":   r.StatusCode,
					}
					if ri < len(_chunk) {
						failedR["record"] = _chunk[ri]
					}
					failedRecords = append(failedRecords, failedR)
				}
			}
			// set derived fields
			respPtr.setDerivedFields()
			return &chunkResponse{
				status:     respPtr.dStatus,
				statusCode: statusCode,
				total:      respPtr.dTotal, success: respPtr.dSuccess, failure: respPtr.dFailure,
				failedRecords: failedRecords,
			}
		}
	}
	// error during http request
	if err != nil { //
		return bulkRespFunc(500, err.Error(), nil)
	}
	// try to parse
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return bulkRespFunc(500, err.Error(), nil)
	}
	// First try to unmarshal to map. If fails, response is likely "string"
	var tempMap map[string]any
	var respPtr *v2EventBulkResponse
	var isOldResp bool
	if err := json.Unmarshal(respBody, &tempMap); err != nil {
		isOldResp = true
	} else {
		// If unmarshal to map succeeds, it's new response format
		respPtr = &v2EventBulkResponse{}
		if err := mapstructure.WeakDecode(tempMap, respPtr); err != nil || respPtr.Status == "" {
			// this should never happen, but just in case
			isOldResp = true
		}
	}
	if isOldResp {
		return bulkRespFunc(httpRes.StatusCode, string(respBody), nil)
	} else {
		res := bulkRespFunc(httpRes.StatusCode, "", respPtr)
		res.rawResponse = tempMap
		return res
	}
}
