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

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func validateWorkflowBodySchema(body map[string]interface{}) (map[string]interface{}, error) {
	// In case data is not provided, set it to empty dict
	if d, found := body["data"]; !found || d == nil {
		body["data"] = map[string]interface{}{}
	}
	schema, err := GetSchema("workflow")
	if err != nil {
		return body, err
	}
	// validate body
	loadedBody := gojsonschema.NewGoLoader(body)
	result, err := schema.Validate(loadedBody)
	if err != nil {
		return body, err
	}
	if !result.Valid() {
		errList := []string{}
		for _, err := range result.Errors() {
			errList = append(errList, fmt.Sprintf(" - %v", err))
		}
		return body, fmt.Errorf("SuprsendValidationError: error in workflow body \n%v", strings.Join(errList, "\n"))
	}
	return body, nil
}

func validateTrackEventSchema(body map[string]interface{}) (map[string]interface{}, error) {
	// In case props is not provided, set it to empty dict
	if d, found := body["properties"]; !found || d == nil {
		body["properties"] = map[string]interface{}{}
	}
	schema, err := GetSchema("event")
	if err != nil {
		return body, err
	}
	// validate body
	loadedBody := gojsonschema.NewGoLoader(body)
	result, err := schema.Validate(loadedBody)
	if err != nil {
		return body, err
	}
	if !result.Valid() {
		errList := []string{}
		for _, err := range result.Errors() {
			errList = append(errList, fmt.Sprintf(" - %v", err))
		}
		return body, fmt.Errorf("SuprsendValidationError: \n%v", strings.Join(errList, "\n"))
	}
	return body, nil
}

func getAttachmentCountInWorkflowBody(body map[string]interface{}) int {
	numAttachments := 0
	if d, dfound := body["data"]; dfound {
		dm := d.(map[string]interface{})
		if a, afound := dm["$attachments"]; afound {
			am := a.([]map[string]interface{})
			numAttachments = len(am)
		}
	}
	return numAttachments
}

func getApparentWorkflowBodySize(body map[string]interface{}, isPartOfBulk bool) (int, error) {
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
					bodyCopy := map[string]interface{}{}
					copier.CopyWithOption(&bodyCopy, &body, copier.Option{DeepCopy: true})
					attachs := bodyCopy["data"].(map[string]interface{})["$attachments"].([]map[string]interface{})
					for _, a := range attachs {
						delete(a, "data")
					}
					apparentBody = bodyCopy

				} else {
					// if auto upload is not enabled, attachment data will be passed as it is.
				}
			} else {
				// If attachment not allowed, then remove data->$attachments before calculating size
				bodyCopy := map[string]interface{}{}
				copier.CopyWithOption(&bodyCopy, &body, copier.Option{DeepCopy: true})

				delete(bodyCopy["data"].(map[string]interface{}), "$attachments")
				apparentBody = bodyCopy
			}
		} else {
			if ATTACHMENT_UPLOAD_ENABLED {
				// if auto upload enabled, to calculate size, replace attachment size with equivalent url size
				extraBytes += numAttachments * ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES
				// -- remove attachments->data key to calculate data size
				bodyCopy := map[string]interface{}{}
				copier.CopyWithOption(&bodyCopy, &body, copier.Option{DeepCopy: true})
				attachs := bodyCopy["data"].(map[string]interface{})["$attachments"].([]map[string]interface{})
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
		return 0, err
	}
	apparentSize := len(bodyBytes) + extraBytes
	// ------
	return apparentSize, nil
}

func getAttachmentCountInEventProperties(event map[string]interface{}) int {
	numAttachments := 0
	if d, dfound := event["properties"]; dfound {
		dm := d.(map[string]interface{})
		if a, afound := dm["$attachments"]; afound {
			am := a.([]map[string]interface{})
			numAttachments = len(am)
		}
	}
	return numAttachments
}

func getApparentEventSize(event map[string]interface{}, isPartOfBulk bool) (int, error) {
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
					eventCopy := map[string]interface{}{}
					copier.CopyWithOption(&eventCopy, &event, copier.Option{DeepCopy: true})
					attachs := eventCopy["properties"].(map[string]interface{})["$attachments"].([]map[string]interface{})
					for _, a := range attachs {
						delete(a, "data")
					}
					apparentBody = eventCopy
				} else {
					// if auto upload is not enabled, attachment data will be passed as it is.
				}
			} else {
				// If attachment not allowed, then remove data->$attachments before calculating size
				eventCopy := map[string]interface{}{}
				copier.CopyWithOption(&eventCopy, &event, copier.Option{DeepCopy: true})

				delete(eventCopy["properties"].(map[string]interface{}), "$attachments")
				apparentBody = eventCopy
			}
		} else {
			if ATTACHMENT_UPLOAD_ENABLED {
				// if auto upload enabled, to calculate size, replace attachment size with equivalent url size
				extraBytes += numAttachments * ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES
				// -- remove attachments->data key to calculate data size
				eventCopy := map[string]interface{}{}
				copier.CopyWithOption(&eventCopy, &event, copier.Option{DeepCopy: true})
				attachs := eventCopy["properties"].(map[string]interface{})["$attachments"].([]map[string]interface{})
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
		return 0, err
	}
	apparentSize := len(bodyBytes) + extraBytes
	// ------
	return apparentSize, nil
}

func getApparentIdentityEventSize(event map[string]interface{}) (int, error) {
	bodyBytes, err := json.Marshal(event)
	if err != nil {
		return 0, err
	}
	return len(bodyBytes), nil
}
