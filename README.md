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
	wfReqBody := map[string]interface{}{
		"workflow": "workflow-slug",
		"recipients": []map[string]interface{}{
			{
				"distinct_id": "0f988f74-6982-41c5-8752-facb6911fb08",
				// if $channels is present, communication will be tried on mentioned channels only (for this request).
				// "$channels": []string{"email"},
				"$email": []string{"user@example.com"},
				"$androidpush": []map[string]interface{}{
					{"token": "__android_push_token__", "provider": "fcm", "device_id": ""},
				},
			},
		},
		// data can be any json / serializable map
		"data": map[string]interface{}{
			"first_name":   "User",
			"spend_amount": "$10",
			"nested_key_example": map[string]interface{}{
				"nested_key1": "some_value_1",
				"nested_key2": map[string]interface{}{
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
