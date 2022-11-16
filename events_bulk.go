package suprsend

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
}

type pendingEventRecord struct {
	record     map[string]interface{}
	recordSize int
}

func (b *bulkEvents) _validateEvents() error {
	if len(b._events) == 0 {
		return fmt.Errorf("events list is empty in bulk request")
	}
	for _, ev := range b._events {
		evJson, bodySize, err := ev.getFinalJson(b.client, true)
		if err != nil {
			return err
		}
		b._pendingRecords = append(
			b._pendingRecords,
			pendingEventRecord{
				record:     evJson,
				recordSize: bodySize,
			},
		)
	}
	return nil
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
	err := b._validateEvents()
	if err != nil {
		return nil, err
	}
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
	_chunk         []map[string]interface{}
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
		_url:   fmt.Sprintf("%sevent/", client.baseUrl),
		_chunk: []map[string]interface{}{},
	}
	return bec
}

func (b *bulkEventsChunk) _addEventToChunk(event map[string]interface{}, eventSize int) {
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
func (b *bulkEventsChunk) tryToAddIntoChunk(event map[string]interface{}, eventSize int) bool {
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
		delete(event["properties"].(map[string]interface{}), "$attachments")
	}
	// Add Event to chunk
	b._addEventToChunk(event, eventSize)
	return true
}

func (b *bulkEventsChunk) trigger() error {
	// prepare http.Request object
	request, err := b.client.prepareHttpRequest("POST", b._url, b._chunk)
	if err != nil {
		return err
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
	return nil
}

func (b *bulkEventsChunk) formatAPIResponse(httpRes *http.Response, err error) *chunkResponse {
	//
	bulkRespFunc := func(statusCode int, errMsg string) *chunkResponse {
		failedRecords := []map[string]interface{}{}
		if statusCode >= 400 {
			for _, c := range b._chunk {
				failedRecords = append(failedRecords,
					map[string]interface{}{
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
