package suprsend

import (
	"fmt"
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
	suprResponse, err := parseV2EventResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	return suprResponse, nil
}

func (w *workflowsService) BulkTriggerInstance() BulkWorkflowsTrigger {
	return &bulkWorkflowsTrigger{
		client:   w.client,
		response: &BulkResponse{},
	}
}
