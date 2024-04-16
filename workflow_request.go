package suprsend

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type WorkflowRequest struct {
	Body            map[string]interface{}
	IdempotencyKey  string
	TenantId        string
	CancellationKey string
}

func (w *WorkflowRequest) AddAttachment(filePath string, ao *AttachmentOption) error {
	if d, found := w.Body["data"]; !found || d == nil {
		w.Body["data"] = map[string]interface{}{}
	}
	attachment, err := GetAttachmentJson(filePath, ao)
	if err != nil {
		return err
	}
	if attachment == nil {
		return nil
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

func (w *WorkflowRequest) getFinalJson(client *Client, isPartOfBulk bool) (map[string]interface{}, int, error) {
	// Add idempotency_key if present
	if w.IdempotencyKey != "" {
		w.Body["$idempotency_key"] = w.IdempotencyKey
	}
	// Add tenant_id if present
	if w.TenantId != "" {
		w.Body["tenant_id"] = w.TenantId
	}
	if w.CancellationKey != "" {
		w.Body["cancellation_key"] = w.CancellationKey
	}
	body, err := validateWorkflowTriggerBodySchema(w.Body)
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

func (w *WorkflowRequest) asJson() map[string]interface{} {
	body := map[string]interface{}{}
	copier.CopyWithOption(&body, w.Body, copier.Option{DeepCopy: true})

	// Add idempotency_key if present
	if w.IdempotencyKey != "" {
		body["$idempotency_key"] = w.IdempotencyKey
	}
	// Add tenant_id if present
	if w.TenantId != "" {
		body["tenant_id"] = w.TenantId
	}
	if w.CancellationKey != "" {
		body["cancellation_key"] = w.CancellationKey
	}
	return body
}
