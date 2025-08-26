package suprsend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WorkflowsService interface {
	Trigger(*WorkflowTriggerRequest) (*Response, error)
	BulkTriggerInstance() BulkWorkflowsTrigger
}

type workflowsService struct {
	client *Client
}

var _ WorkflowsService = &workflowsService{}

func newWorkflowService(client *Client) *workflowsService {
	ws := &workflowsService{
		client: client,
	}
	return ws
}

func (w *workflowsService) Trigger(workflow *WorkflowTriggerRequest) (*Response, error) {
	wfBody, _, err := workflow.getFinalJson(w.client, false)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%strigger/", w.client.baseUrl)
	// prepare http.Request object
	request, err := w.client.prepareHttpRequest("POST", url, wfBody)
	if err != nil {
		return nil, err
	}
	httpResponse, err := w.client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	suprResponse, err := w.formatAPIResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return suprResponse, nil
}

func (w *workflowsService) formatAPIResponse(httpRes *http.Response) (*Response, error) {
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, &Error{Err: err}
	}
	if httpRes.StatusCode >= 400 {
		var errorResponse map[string]interface{}
		if jsonErr := json.Unmarshal(respBody, &errorResponse); jsonErr == nil {
			if errorObj, exists := errorResponse["error"]; exists {
				if errorMap, ok := errorObj.(map[string]interface{}); ok {
					if errorMessage, exists := errorMap["message"]; exists {
						if msg, ok := errorMessage.(string); ok {
							return nil, &Error{Code: httpRes.StatusCode, Message: msg}
						}
					}
				}
			}
			if errorMessage, exists := errorResponse["message"]; exists {
				if msg, ok := errorMessage.(string); ok {
					return nil, &Error{Code: httpRes.StatusCode, Message: msg}
				}
			}
		}
		return nil, &Error{Code: httpRes.StatusCode, Message: string(respBody)}
	}
	var responseData map[string]interface{}
	var messageID string
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		if msgID, exists := responseData["message_id"]; exists {
			if msgIDtr, ok := msgID.(string); ok {
				messageID = msgIDtr
			}
		}
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: messageID, RawResponse: responseData}, nil
}

func (w *workflowsService) BulkTriggerInstance() BulkWorkflowsTrigger {
	return &bulkWorkflowsTrigger{
		client:   w.client,
		response: &BulkResponse{},
	}
}
