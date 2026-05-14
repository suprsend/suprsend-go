# suprsend-go
SuprSend Go SDK

## Installation
```bash
go get github.com/suprsend/suprsend-go
```

## Usage
Initialize the SuprSend SDK
```go
import (
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func main() {
    opts := []suprsend.ClientOption{
		// suprsend.WithDebug(true),
	}
    suprClient, err := suprsend.NewClient("__api_key__", "__api_secret__", opts...)
	if err != nil {
		log.Println(err)
	}
}

```

### Trigger Workflow
```go
package main

import (
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func main() {
	// Instantiate Client
	suprClient, err := suprsend.NewClient("__api_key__", "__api_secret__")
	if err != nil {
		log.Println(err)
		return
	}
	// Create WorkflowTriggerRequest body
	wfReqBody := map[string]any{
		"workflow": "workflow-slug",
		"recipients": []map[string]any{
			{
				"distinct_id": "0f988f74-6982-41c5-8752-facb6911fb08",
				// if $channels is present, communication will be tried on mentioned channels only (for this request).
				// "$channels": []string{"email"},
				"$email": []string{"user@example.com"},
				"$androidpush": []map[string]any{
					{"token": "__android_push_token__", "provider": "fcm", "device_id": ""},
				},
			},
		},
		// data can be any json / serializable map
		"data": map[string]any{
			"first_name":   "User",
			"spend_amount": "$10",
			"nested_key_example": map[string]any{
				"nested_key1": "some_value_1",
				"nested_key2": map[string]any{
					"nested_key3": "some_value_3",
				},
			},
		},
	}

	wf := &suprsend.WorkflowTriggerRequest{
		Body:           wfReqBody,
		IdempotencyKey: "",
		TenantId:        "",
	}
    // Call TriggerWorkflow to send request to Suprsend
	resp, err = suprClient.Workflows.Trigger(wf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}


```
Check SuprSend docs here https://docs.suprsend.com/docs

Check examples directory to understand how to use different functionalities.

### Messages API

#### List Messages
Fetch a paginated list of messages for your workspace. All fields in `MessageListOptions` are optional.

```go
import (
    "context"
    "log"

    suprsend "github.com/suprsend/suprsend-go"
)

// Basic call — returns first page with default limit
resp, err := suprClient.Messages.List(context.Background(), nil)
if err != nil {
    log.Fatalln(err)
}
log.Println(resp)

// With filters
isCampaign := false
resp, err = suprClient.Messages.List(context.Background(), &suprsend.MessageListOptions{
    // Pagination
    Limit:  20,             // records per page (default: 1000, max: 1000)
    After:  "__cursor__",   // cursor for next page (from response Meta.After)
    Before: "__cursor__",   // cursor for previous page (from response Meta.Before)

    // Message filters
    MessageID:      "__message_id__",
    IdempotencyKey: "__idempotency_key__",

    // Recipient filters
    RecipientID: []string{"user1", "user2"}, // recipient_id[]
    TenantID:    "default",

    // Object recipient filters (both required together)
    ObjectType: "__object_type__",
    ObjectID:   "__object_id__",

    // Workflow / execution filters
    WorkflowSlug: "purchase-made",
    ExecutionID:  "__execution_id__",

    // Channel filter
    // valid: email, sms, whatsapp, androidpush, iospush, webpush, slack, ms_teams
    Channel: "email",

    // status[] — valid: triggered, delivered, delivery_failed, seen, clicked, dismissed, read, archived, unread
    Status: []string{"delivered", "seen"},

    // category[]
    Category: []string{"transactional"},

    IsCampaign: &isCampaign,

    // Date range filters (RFC3339 format)
    CreatedAtGte: "2026-01-01T00:00:00Z",
    CreatedAtLte: "2026-12-31T23:59:59Z",
})
if err != nil {
    log.Fatalln(err)
}
log.Println(resp)
```

Response structure (`*CursorListApiResponse`):
```json
{
  "meta": {
    "count": 150,
    "limit": 20,
    "has_prev": true,
    "has_next": true,
    "before": null,
    "after": null
  },
  "results": [
    {
      "message_id": "01KQVGPW9ZJKH6T5TSxxxxxxx",
      "created_at": "2025-08-27T15:24:38.14Z",
      "updated_at": "2025-08-27T15:24:41.00Z",
      "triggered_at": "2025-08-27T15:24:38.29Z",
      "delivered_at": "2025-08-27T15:24:41.037Z",
      "seen_at": "2025-08-27T15:24:45.65Z",
      "clicked_at": null,
      "dismissed_at": null,
      "read_at": null,
      "unread_at": null,
      "archived_at": null,
      "unarchived_at": null,
      "is_read": false,
      "is_archived": false,
      "status": "seen",
      "channel": "email",
      "category": "transactional",
      "idempotency_key": "8087c3e7-6612-4d16-9660-xxxxxxxx",
      "failure_reason": "",
      "recipient": { "$type": "user", "distinct_id": "user_123" },
      "parent_entity_id": "__object:TEAMS:teams_1",
      "parent_entity_type": "object",
      "vendor": { "name": "amazon_ses", "nickname": "AWS SES" },
      "execution_id": "dsl_w1_id3741_xxxxxxxx_0_1",
      "parent_execution_id": "dsl_w1_id3741_xxxxxxxx_0",
      "is_campaign": false,
      "tenant_id": "default",
      "workflow": {
        "slug": "purchase-made",
        "version_id": "wf_v_01KQVGxxxxxxx_chkp",
        "name": "Purchase Workflow",
        "node_ref": ""
      },
      "template": { "name": "Purchase Template", "slug": "amazon_ses", "version_no": 1 },
      "channel_identity": { "email": "user@example.com" }
    }
  ]
}
```

#### Bulk Update Message Status
Update the status of one or more messages in a single call.
Valid actions: `seen`, `clicked`, `dismissed`, `read`, `unread`, `archived`, `unarchived`.

```go
messages := []suprsend.MessageUpdateItem{
    {MessageID: "__message_id_1__", Action: "read"},
    {MessageID: "__message_id_2__", Action: "archived"},
}
resp, err := suprClient.Messages.BulkUpdate(context.Background(), messages)
if err != nil {
    log.Fatalln(err)
}
log.Println(resp)
```

Response structure (`*MessageBulkUpdateResponse`):
```json
{
  "records": [
    {
      "message_id": "__message_id_1__",
      "status_code": 202,
      "error": null
    },
    {
      "message_id": "__message_id_2__",
      "status_code": 404,
      "error": {
        "type": "not_found",
        "message": "message not found"
      }
    }
  ]
}
```

Error `status_code` values: `202` success · `404` not found · `422` action not supported · `500` internal error.
