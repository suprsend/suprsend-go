package suprsend

import (
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
		return nil, err
	}
	if httpRes.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpRes.StatusCode, string(respBody))
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}

func (w *workflowsService) BulkTriggerInstance() BulkWorkflowsTrigger {
	return &bulkWorkflowsTrigger{
		client:   w.client,
		response: &BulkResponse{},
	}
}
