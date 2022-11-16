package suprsend

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Workflow struct {
	Body           map[string]interface{}
	IdempotencyKey string
	BrandId        string
}

func (w *Workflow) AddAttachment(filePath string, ao *AttachmentOption) error {
	if d, found := w.Body["data"]; !found || d == nil {
		w.Body["data"] = map[string]interface{}{}
	}
	attachment, err := GetAttachmentJson(filePath, ao)
	if err != nil {
		return err
	}
	data := w.Body["data"].(map[string]interface{})
	if a, found := data["$attachments"]; !found || a == nil {
		data["$attachments"] = []map[string]interface{}{}
	}
	allAttachments := data["$attachments"].([]map[string]interface{})
	allAttachments = append(allAttachments, attachment)
	data["$attachments"] = allAttachments
	return nil
}

func (w *Workflow) getFinalJson(client *Client, isPartOfBulk bool) (map[string]interface{}, int, error) {
	// Add idempotency_key if present
	if w.IdempotencyKey != "" {
		w.Body["$idempotency_key"] = w.IdempotencyKey
	}
	// Add brand_id if present
	if w.BrandId != "" {
		w.Body["brand_id"] = w.BrandId
	}
	body, err := validateWorkflowBodySchema(w.Body)
	if err != nil {
		return nil, 0, err
	}
	w.Body = body
	// Check request size
	apparentSize, err := getApparentWorkflowBodySize(body, isPartOfBulk)
	if err != nil {
		return nil, 0, err
	}
	if apparentSize > SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES {
		errStr := fmt.Sprintf("workflow body too big - %d Bytes, must not cross %s", apparentSize,
			SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE)
		return nil, 0, errors.New(errStr)
	}
	return w.Body, apparentSize, nil
}

type workflowTrigger struct {
	client *Client
	_url   string
}

func newWorkflowTriggerInstance(client *Client) *workflowTrigger {
	wt := &workflowTrigger{
		client: client,
		// workflow trigger url
		_url: fmt.Sprintf("%s%s/trigger/", client.baseUrl, client.ApiKey),
	}
	return wt
}

func (w *workflowTrigger) Trigger(workflow *Workflow) (*Response, error) {
	wfBody, _, err := workflow.getFinalJson(w.client, false)
	if err != nil {
		return nil, err
	}
	suprResp, err := w.send(wfBody)
	if err != nil {
		return nil, err
	}
	return suprResp, nil
}

func (w *workflowTrigger) send(wfBody map[string]interface{}) (*Response, error) {
	// prepare http.Request object
	request, err := w.client.prepareHttpRequest("POST", w._url, wfBody)
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

func (w *workflowTrigger) formatAPIResponse(httpRes *http.Response) (*Response, error) {
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode >= 400 {
		return nil, fmt.Errorf("code: %v. message: %v", httpRes.StatusCode, string(respBody))
	}
	return &Response{Success: true, StatusCode: httpRes.StatusCode, Message: string(respBody)}, nil
}
