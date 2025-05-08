package suprsend

const (
	//
	VERSION = "0.7.0"
	//
	DEFAULT_URL = "https://hub.suprsend.com/"

	// a API call should not have apparent body size of more than 800KB
	BODY_MAX_APPARENT_SIZE_IN_BYTES          = 800 * 1024
	BODY_MAX_APPARENT_SIZE_IN_BYTES_READABLE = "800KB"

	// in general url-size wont exceed 2048 chars or 2048 utf-8 bytes
	ATTACHMENT_URL_POTENTIAL_SIZE_IN_BYTES = 2100

	// few keys added in-flight, amounting to almost 200 bytes increase per workflow-body
	WORKFLOW_RUNTIME_KEYS_POTENTIAL_SIZE_IN_BYTES = 200

	// max workflow-records in one bulk api call.
	MAX_WORKFLOWS_IN_BULK_API = 100

	// max event-records in one bulk api call
	MAX_EVENTS_IN_BULK_API = 100

	ALLOW_ATTACHMENTS_IN_BULK_API = true

	ATTACHMENT_UPLOAD_ENABLED = false

	// single Identity event limit
	IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES          = 10 * 1024
	IDENTITY_SINGLE_EVENT_MAX_APPARENT_SIZE_IN_BYTES_READABLE = "10KB"

	MAX_IDENTITY_EVENTS_IN_BULK_API = 400

	// time.RFC1123
	HEADER_DATE_FMT = "Mon, 02 Jan 2006 15:04:05 GMT"
)
