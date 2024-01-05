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
	// Create workflow body
	wfBody := map[string]interface{}{
		"name":                  "Workflow Name",
		"template":              "template slug",
		"notification_category": "category",
		// "delay":                 "15m", // Chek duration format in documentation
		"users": []map[string]interface{}{
			{
				"distinct_id": "0f988f74-6982-41c5-8752-facb6911fb08",
				// if $channels is present, communication will be tried on mentioned channels only.
				// "$channels": []string{"email"},
				"$email": []string{"user@example.com"},
				"$androidpush": []map[string]interface{}{
					{"token": "__android_push_token__", "provider": "fcm", "device_id": ""},
				},
			},
		},
		// delivery instruction. how should notifications be sent, and whats the success metric
		"delivery": map[string]interface{}{
			"smart":   false,
			"success": "seen",
		},
		// # data can be any json / serializable python-dictionary
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

	wf := &suprsend.Workflow{
		Body:           wfBody,
		IdempotencyKey: "",
		TenantId:        "",
	}
    // Call TriggerWorkflow to send request to Suprsend
	_, err = suprClient.TriggerWorkflow(wf)
	if err != nil {
		log.Fatalln(err)
	}
}


```
Check SuprSend docs here https://docs.suprsend.com/docs

Check examples directory to understand how to use different functionalities.
