package suprsend

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/xeipuuv/gojsonschema"
)

func CurrentTimeFormatted() string {
	t := time.Now().UTC()
	return t.Format(HEADER_DATE_FMT)
}

func validateWorkflowBodySchema(body map[string]any) (map[string]any, error) {
	// In case data is not provided, set it to empty dict
	if d, found := body["data"]; !found || d == nil {
		body["data"] = map[string]any{}
	}
	schema, err := GetSchema("workflow")
	if err != nil {
		return body, err
	}
	// validate body
	loadedBody := gojsonschema.NewGoLoader(body)
	result, err := schema.Validate(loadedBody)
	if err != nil {
		return body, &Error{Err: err}
	}
	if !result.Valid() {
		errList := []string{}
		for _, err := range result.Errors() {
			errList = append(errList, fmt.Sprintf(" - %v", err))
		}
		return body, &Error{Message: fmt.Sprintf("SuprsendValidationError: error in workflow body \n%v", strings.Join(errList, "\n"))}
	}
	return body, nil
}

func validateWorkflowTriggerBodySchema(body map[string]any) (map[string]any, error) {
	// In case data is not provided, set it to empty dict
	if d, found := body["data"]; !found || d == nil {
		body["data"] = map[string]any{}
	}
	schema, err := GetSchema("workflow_trigger")
	if err != nil {
		return body, err
	}
	// validate body
	loadedBody := gojsonschema.NewGoLoader(body)
	result, err := schema.Validate(loadedBody)
	if err != nil {
		return body, &Error{Err: err}
	}
	if !result.Valid() {
		errList := []string{}
		for _, err := range result.Errors() {
			errList = append(errList, fmt.Sprintf(" - %v", err))
		}
		return body, &Error{Message: fmt.Sprintf("SuprsendValidationError: error in workflow body \n%v", strings.Join(errList, "\n"))}
	}
	return body, nil
}

func validateTrackEventSchema(body map[string]any) (map[string]any, error) {
	// In case props is not provided, set it to empty dict
	if d, found := body["properties"]; !found || d == nil {
		body["properties"] = map[string]any{}
	}
	schema, err := GetSchema("event")
	if err != nil {
		return body, err
	}
	// validate body
	loadedBody := gojsonschema.NewGoLoader(body)
	result, err := schema.Validate(loadedBody)
	if err != nil {
		return body, &Error{Err: err}
	}
	if !result.Valid() {
		errList := []string{}
		for _, err := range result.Errors() {
			errList = append(errList, fmt.Sprintf(" - %v", err))
		}
		return body, &Error{Message: fmt.Sprintf("SuprsendValidationError: \n%v", strings.Join(errList, "\n"))}
	}
	return body, nil
}

func validateListBroadcastBodySchema(body map[string]any) (map[string]any, error) {
	// In case props is not provided, set it to empty dict
	if d, found := body["data"]; !found || d == nil {
		body["data"] = map[string]any{}
	}
	schema, err := GetSchema("list_broadcast")
	if err != nil {
		return body, err
	}
	// validate body
	loadedBody := gojsonschema.NewGoLoader(body)
	result, err := schema.Validate(loadedBody)
	if err != nil {
		return body, &Error{Err: err}
	}
	if !result.Valid() {
		errList := []string{}
		for _, err := range result.Errors() {
			errList = append(errList, fmt.Sprintf(" - %v", err))
		}
		return body, &Error{Message: fmt.Sprintf("SuprsendValidationError: \n%v", strings.Join(errList, "\n"))}
	}
	return body, nil
}

func getAttachmentCountInWorkflowBody(body map[string]any) int {
	numAttachments := 0
	if d, dfound := body["data"]; dfound {
		dm := d.(map[string]any)
		if a, afound := dm["$attachments"]; afound {
			am := a.([]map[string]any)
			numAttachments = len(am)
		}
	}
	return numAttachments
}

func getApparentWorkflowBodySize(body map[string]any, isPartOfBulk bool) (int, error) {
	extraBytes := WORKFLOW_RUNTIME_KEYS_POTENTIAL_SIZE_IN_BYTES
	apparentBody := body
	numAttachments := getAttachmentCountInWorkflowBody(body)
	if numAttachments > 0 {
		if isPartOfBulk {
			if ALLOW_ATTACHMENTS_IN_BULK_API {
				// if attachment is allowed in bulk api, then calculate size based on whether auto Upload is enabled
				if ATTACHMENT_UPLOAD_ENABLED {
					// If auto upload enabled, To calculate size, replace attachment size with equivalent url size
					extraBytes += numAttachments * ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES
					// -- remove attachments->data key to calculate data size
					bodyCopy := map[string]any{}
					copier.CopyWithOption(&bodyCopy, &body, copier.Option{DeepCopy: true})
					attachs := bodyCopy["data"].(map[string]any)["$attachments"].([]map[string]any)
					for _, a := range attachs {
						delete(a, "data")
					}
					apparentBody = bodyCopy

				} else {
					// if auto upload is not enabled, attachment data will be passed as it is.
				}
			} else {
				// If attachment not allowed, then remove data->$attachments before calculating size
				bodyCopy := map[string]any{}
				copier.CopyWithOption(&bodyCopy, &body, copier.Option{DeepCopy: true})

				delete(bodyCopy["data"].(map[string]any), "$attachments")
				apparentBody = bodyCopy
			}
		} else {
			if ATTACHMENT_UPLOAD_ENABLED {
				// if auto upload enabled, to calculate size, replace attachment size with equivalent url size
				extraBytes += numAttachments * ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES
				// -- remove attachments->data key to calculate data size
				bodyCopy := map[string]any{}
				copier.CopyWithOption(&bodyCopy, &body, copier.Option{DeepCopy: true})
				attachs := bodyCopy["data"].(map[string]any)["$attachments"].([]map[string]any)
				for _, a := range attachs {
					delete(a, "data")
				}
				apparentBody = bodyCopy

			} else {
				// if auto upload is not enabled, attachment data will be passed as it is.
			}
		}
	}
	// ------
	bodyBytes, err := json.Marshal(apparentBody)
	if err != nil {
		return 0, &Error{Err: err}
	}
	apparentSize := len(bodyBytes) + extraBytes
	// ------
	return apparentSize, nil
}

func getAttachmentCountInEventProperties(event map[string]any) int {
	numAttachments := 0
	if d, dfound := event["properties"]; dfound {
		dm := d.(map[string]any)
		if a, afound := dm["$attachments"]; afound {
			am := a.([]map[string]any)
			numAttachments = len(am)
		}
	}
	return numAttachments
}

func getApparentEventSize(event map[string]any, isPartOfBulk bool) (int, error) {
	extraBytes := 0
	apparentBody := event
	numAttachments := getAttachmentCountInEventProperties(event)
	if numAttachments > 0 {
		if isPartOfBulk {
			if ALLOW_ATTACHMENTS_IN_BULK_API {
				// if attachment is allowed in bulk api, then calculate size based on whether auto Upload is enabled
				if ATTACHMENT_UPLOAD_ENABLED {
					// If auto upload enabled, To calculate size, replace attachment size with equivalent url size
					extraBytes += numAttachments * ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES
					// -- remove attachments->data key to calculate data size
					eventCopy := map[string]any{}
					copier.CopyWithOption(&eventCopy, &event, copier.Option{DeepCopy: true})
					attachs := eventCopy["properties"].(map[string]any)["$attachments"].([]map[string]any)
					for _, a := range attachs {
						delete(a, "data")
					}
					apparentBody = eventCopy
				} else {
					// if auto upload is not enabled, attachment data will be passed as it is.
				}
			} else {
				// If attachment not allowed, then remove data->$attachments before calculating size
				eventCopy := map[string]any{}
				copier.CopyWithOption(&eventCopy, &event, copier.Option{DeepCopy: true})

				delete(eventCopy["properties"].(map[string]any), "$attachments")
				apparentBody = eventCopy
			}
		} else {
			if ATTACHMENT_UPLOAD_ENABLED {
				// if auto upload enabled, to calculate size, replace attachment size with equivalent url size
				extraBytes += numAttachments * ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES
				// -- remove attachments->data key to calculate data size
				eventCopy := map[string]any{}
				copier.CopyWithOption(&eventCopy, &event, copier.Option{DeepCopy: true})
				attachs := eventCopy["properties"].(map[string]any)["$attachments"].([]map[string]any)
				for _, a := range attachs {
					delete(a, "data")
				}
				apparentBody = eventCopy
			} else {
				// if auto upload is not enabled, attachment data will be passed as it is.
			}
		}
	}
	// ------
	bodyBytes, err := json.Marshal(apparentBody)
	if err != nil {
		return 0, &Error{Err: err}
	}
	apparentSize := len(bodyBytes) + extraBytes
	// ------
	return apparentSize, nil
}

func getApparentIdentityEventSize(event map[string]any) (int, error) {
	bodyBytes, err := json.Marshal(event)
	if err != nil {
		return 0, &Error{Err: err}
	}
	return len(bodyBytes), nil
}

func getApparentListBroadcastBodySize(body map[string]any) (int, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, &Error{Err: err}
	}
	return len(bodyBytes), nil
}

func invalidRecordJson(failedRecord map[string]any, err error) map[string]any {
	return map[string]any{
		"record": failedRecord,
		"error":  err.Error(),
		"code":   500,
	}
}

func appendQueryParamPart(url string, qp string) string {
	if qp == "" {
		return url
	} else {
		return fmt.Sprintf("%s?%s", url, qp)
	}
}
