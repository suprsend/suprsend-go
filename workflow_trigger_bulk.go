package suprsend

import (
	"fmt"
	"log"

	"github.com/jinzhu/copier"
)

type BulkWorkflowsTrigger interface {
	Append(...*WorkflowTriggerRequest)
	Trigger() (*BulkResponse, error)
}

var _ BulkWorkflowsTrigger = &bulkWorkflowsTrigger{}

type bulkWorkflowsTrigger struct {
	client *Client
	//
	_workflows      []WorkflowTriggerRequest
	_pendingRecords []pendingWorkflowTriggerRecord
	chunks          []*bulkWorkflowsRequestChunk
	//
	response *BulkResponse
	// invalid_record json: {"record": event-json, "error": error_str, "code": 500}
	_invalidRecords []map[string]any
}

type pendingWorkflowTriggerRecord struct {
	record     map[string]any
	recordSize int
}

func (b *bulkWorkflowsTrigger) _validateWorkflows() {
	for _, wf := range b._workflows {
		wfJson, bodySize, err := wf.getFinalJson(b.client, true)
		if err != nil {
			invRec := invalidRecordJson(wf.asJson(), err)
			b._invalidRecords = append(b._invalidRecords, invRec)
		} else {
			b._pendingRecords = append(
				b._pendingRecords,
				pendingWorkflowTriggerRecord{
					record:     wfJson,
					recordSize: bodySize,
				},
			)
		}
	}
}

func (b *bulkWorkflowsTrigger) _chunkify(startIdx int) {
	currChunk := newBulkWorkflowsRequestChunk(b.client)
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

func (b *bulkWorkflowsTrigger) Append(workflows ...*WorkflowTriggerRequest) {
	for _, wf := range workflows {
		if wf == nil {
			continue
		}
		wfCopy := WorkflowTriggerRequest{}
		copier.CopyWithOption(&wfCopy, wf, copier.Option{DeepCopy: true})
		b._workflows = append(b._workflows, wfCopy)
	}
}

func (b *bulkWorkflowsTrigger) Trigger() (*BulkResponse, error) {
	b._validateWorkflows()
	if len(b._invalidRecords) > 0 {
		chResponse := invalidRecordsChunkResponse(b._invalidRecords)
		b.response.mergeChunkResponse(chResponse)
	}
	if len(b._pendingRecords) > 0 {
		b._chunkify(0)
		for cIdx, ch := range b.chunks {
			if b.client.debug {
				log.Printf("DEBUG: triggering api call for chunk: %d\n", cIdx)
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

type bulkWorkflowsRequestChunk struct {
	_chunkApparentSizeInBytes int
	_maxRecordsInChunk        int
	//
	client *Client
	_url   string
	//
	_chunk         []map[string]any
	_runningSize   int
	_runningLength int
	response       *chunkResponse
}

func newBulkWorkflowsRequestChunk(client *Client) *bulkWorkflowsRequestChunk {
	bwc := &bulkWorkflowsRequestChunk{
		_chunkApparentSizeInBytes: BODY_MAX_APPARENT_SIZE_IN_BYTES,
		_maxRecordsInChunk:        MAX_WORKFLOWS_IN_BULK_API,
		//
		client: client,
		_url:   fmt.Sprintf("%strigger/", client.baseUrl),
		_chunk: []map[string]any{},
	}
	return bwc
}

func (b *bulkWorkflowsRequestChunk) _addBodyToChunk(body map[string]any, bodySize int) {
	// First add size, then event to reduce effects of race condition
	b._runningSize += bodySize
	b._chunk = append(b._chunk, body)
	b._runningLength += 1
}

func (b *bulkWorkflowsRequestChunk) _checkLimitReached() bool {
	return b._runningLength >= b._maxRecordsInChunk || b._runningSize >= b._chunkApparentSizeInBytes
}

/*
returns whether passed body was able to get added to this chunk or not,
if true, body gets added to chunk
*/
func (b *bulkWorkflowsRequestChunk) tryToAddIntoChunk(body map[string]any, bodySize int) bool {
	if body == nil {
		return true
	}
	if b._checkLimitReached() {
		return false
	}
	// if apparent_size of body crosses limit
	if (b._runningSize + bodySize) > b._chunkApparentSizeInBytes {
		return false
	}
	if !ALLOW_ATTACHMENTS_IN_BULK_API {
		delete(body["data"].(map[string]any), "$attachments")
	}
	// Add workflow to chunk
	b._addBodyToChunk(body, bodySize)
	//
	return true
}

func (b *bulkWorkflowsRequestChunk) trigger() {
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
