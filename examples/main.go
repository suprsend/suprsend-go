package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func main() {
	triggerWorkflowAPIExample()
	bulkWorkflowTriggerAPIExample()
	//
	triggerDynamicWorkflowExample()
	sendEventExample()
	updateUserProfileExample()
	//
	bulkDynamicWorkflowsExample()
	bulkEventsExample()
	bulkUserProfileUpdateExample()
	//
	tenantExample()
	//
	subscriberListExample()
	subscriberListVersioningExample()
	//
	objectCrudOperationsExample()
	updateObjectPropertiesExample()
}

func getSuprsendClient() (*suprsend.Client, error) {
	opts := []suprsend.ClientOption{
		suprsend.WithDebug(true),
	}
	suprClient, err := suprsend.NewClient("__api_key__", "__api_secret__", opts...)
	if err != nil {
		return nil, err
	}
	return suprClient, nil
}

func triggerWorkflowAPIExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// Create WorkflowTriggerRequest body
	wfReqBody := map[string]interface{}{
		"workflow": "workflow-slug", // mandatory
		// "actor":    "actor-distinct-id", // optional
		// recipients: an array. each element is either string/dict
		// in case of string, the value must the distinct_id of recipient/user
		// in case of dict, along with distinct_id, value can contain user profile info like channels etc..
		// e.g ["distinct_id1", "distinct_id1"]
		// or [{"distinct_id": "__distinct_id_1__", "$email": ["a@example.com"], "prop1": "v1"}]
		"recipients": []map[string]interface{}{
			{
				"distinct_id": "__distinct_id1__",
				// if $channels is present, communication will be tried on mentioned channels only (for this request).
				// "$channels": []string{"email"},
				"$email": []string{"user@example.com"},
				"$androidpush": []map[string]interface{}{
					{"token": "__android_push_token__", "provider": "fcm", "device_id": ""},
				},
				"name": "Recipient 1",
			},
		},
		// # data can be any json / serializable map
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
		TenantId:       "",
	}
	// Add attachment by calling .AddAttachment
	err = wf.AddAttachment("https://attachment-url", &suprsend.AttachmentOption{})
	if err != nil {
		log.Fatalln(err)
	}
	// Call Workflows.Trigger to send request to Suprsend
	resp, err := suprClient.Workflows.Trigger(wf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}

func bulkWorkflowTriggerAPIExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// WorkflowTriggerRequest: 1
	wf1 := &suprsend.WorkflowTriggerRequest{
		Body: map[string]interface{}{
			"workflow": "workflow-slug",
			// "actor":    "actor-distinct-id", // optional
			// recipients: an array. each element is either string/dict
			// in case of string, the value must the distinct_id of recipient/user
			// in case of dict, along with distinct_id, value can contain user profile info like channels etc..
			// e.g ["distinct_id1", "distinct_id1"]
			// or [{"distinct_id": "__distinct_id_1__", "$email": ["a@example.com"], "prop1": "v1"}]
			"recipients": []map[string]interface{}{
				{
					"distinct_id": "__distinct_id1__",
					// if $channels is present, communication will be tried on mentioned channels only (for this request).
					// "$channels": []string{"email"},
					"$email": []string{"user@example.com"},
					"$androidpush": []map[string]interface{}{
						{"token": "__android_push_token__", "provider": "fcm", "device_id": ""},
					},
					"name": "Recipient 1",
				},
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
		},
		IdempotencyKey: "",
		TenantId:       "",
	}

	// WorkflowTriggerRequest: 2
	wf2 := &suprsend.WorkflowTriggerRequest{
		Body: map[string]interface{}{
			"workflow": "workflow-slug",
			// "actor":    "actor-distinct-id", // optional
			// recipients: an array. each element is either string/dict
			// in case of string, the value must the distinct_id of recipient/user
			// in case of dict, along with distinct_id, value can contain user profile info like channels etc..
			// e.g ["distinct_id1", "distinct_id1"]
			// or [{"distinct_id": "__distinct_id_1__", "$email": ["a@example.com"], "prop1": "v1"}]
			"recipients": []map[string]interface{}{
				{
					"distinct_id": "__distinct_id1__",
					// if $channels is present, communication will be tried on mentioned channels only (for this request).
					// "$channels": []string{"email"},
					"$email": []string{"user@example.com"},
					"$androidpush": []map[string]interface{}{
						{"token": "__android_push_token__", "provider": "fcm", "device_id": ""},
					},
					"name": "Recipient 1",
				},
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
		},
		IdempotencyKey: "123456",
		TenantId:       "default",
	}
	// ...... Add as many Workflow records as required.

	// Create workflows bulk instance
	bulkIns := suprClient.Workflows.BulkTriggerInstance()
	// add all your workflows to bulkInstance
	bulkIns.Append(wf1, wf2)
	// Trigger
	bulkResponse, err := bulkIns.Trigger()
	if err != nil {
		log.Println(err)
		//
	}
	log.Println(bulkResponse)
}

// Deprecated: dynamic workflows will be deprecated in near future
func triggerDynamicWorkflowExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
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
		TenantId:       "",
	}
	// Add attachment by calling .AddAttachment
	err = wf.AddAttachment("https://attachment-url", &suprsend.AttachmentOption{})
	if err != nil {
		log.Fatalln(err)
	}
	// Call TriggerWorkflow to send request to Suprsend
	_, err = suprClient.TriggerWorkflow(wf)
	if err != nil {
		log.Fatalln(err)
	}
}

func sendEventExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	//
	ev := &suprsend.Event{
		EventName:  "__event_name__",
		DistinctId: "__distinct_id__",
		Properties: map[string]interface{}{
			"k1": "v1",
		},
		// IdempotencyKey: "",
		// TenantId: "",
	}
	// Add attachment (If needed) by calling .AddAttachment
	err = ev.AddAttachment("~/Downloads/attachment.pdf", &suprsend.AttachmentOption{FileName: "My Attachment.pdf"})
	if err != nil {
		log.Println(err)
	}
	// Send event to Suprsend by calling .TrackEvent
	resp, err := suprClient.TrackEvent(ev)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}

func updateUserProfileExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	//
	user := suprClient.Users.GetInstance("__distinct_id__")
	// Add email channel
	user.AddEmail("user@example.com")
	// add sms channel
	user.AddSms("+1444455555")
	// Add whatsapp channel
	user.AddWhatsapp("+1444455555")
	// Add androidpush token, token providers: fcm/xiaomi
	user.AddAndroidpush("__fcm_push_token__", "fcm")
	// Add iospush token, token providers: apns
	user.AddIospush("__ios_push_token__", "apns")
	// Add webpush token (vapid)
	user.AddWebpush(map[string]interface{}{
		"keys": map[string]interface{}{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")

	// add slack using email
	user.AddSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// add slack using user_id
	user.AddSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// add slack using channel_id
	user.AddSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// add slack using incoming-webhook
	user.AddSlack(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// DM on Team's channel using conversation id
	user.AddMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// add teams via DM user using team user id
	user.AddMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// add teams using incoming webhook
	user.AddMSTeams(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})
	//
	// remove email channel
	user.RemoveEmail("user@example.com")
	// remove sms channel
	user.RemoveSms("+1444455555")
	// remove whatsapp channel
	user.RemoveWhatsapp("+1444455555")
	// remove androidpush token
	user.RemoveAndroidpush("__fcm_push_token__", "fcm")
	// remove iospush token
	user.RemoveIospush("__ios_push_token__", "apns")

	// remove webpush token
	user.RemoveWebpush(map[string]interface{}{
		"keys": map[string]interface{}{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")
	// remove slack using email
	user.RemoveSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// remove slack using user_id
	user.RemoveSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// remove slack using channel_id
	user.RemoveSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// remove slack using incoming-webhook
	user.RemoveSlack(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// remove teams via DM on Team's channel using conversation id
	user.RemoveMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// remove teams via DM user using team user id
	user.RemoveMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// remove teams using incoming webhook
	user.RemoveMSTeams(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})

	// Set user preferred language. languageCode must be in [ISO 639-1 2-letter] format
	user.SetPreferredLanguage("en")
	// set timezone property at subscriber level based on IANA timezone info
	user.SetTimezone("America/Los_Angeles")
	// If you need to remove all emails for this user, call user.Unset(["$email"])
	user.Unset([]string{"$email"})
	// # what value to pass to unset channels
	// # for email:                $email
	// # for whatsapp:             $whatsapp
	// # for SMS:                  $sms
	// # for androidpush tokens:   $androidpush
	// # for iospush tokens:       $iospush
	// # for webpush tokens:       $webpush
	// # for slack:                $slack
	// # for ms teams:             $ms_teams

	// set a user property using a map
	user.Set(map[string]interface{}{"prop1": "val1", "prop2": "val2"})
	// set a user property using a key, value pair
	user.SetKV("prop", "value")
	// set a user property once using map
	user.SetOnce(map[string]interface{}{"prop3": "val3"})
	// set a user property once using a key value pair
	user.SetOnceKV("prop4", "val4")
	// increment an already existing user property using key value pair
	user.IncrementKV("increment_prop", 2)
	// increment an already existing property using map
	user.Increment(map[string]interface{}{"increment_prop1": 5})

	// Save user
	_, err = user.Save()
	if err != nil {
		log.Fatalln(err)
	}
}

// Deprecated: dynamic workflows will be deprecated in near future
func bulkDynamicWorkflowsExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// Workflow: 1
	wf1 := &suprsend.Workflow{
		Body: map[string]interface{}{
			"name":                  "__workflow_name__",
			"template":              "__template_slug__",
			"notification_category": "__category__", // system/transactional/promotional
			"users": []map[string]interface{}{
				{
					"distinct_id": "__distinct_id__",
					"$email":      []string{"__email__"},
				},
			},
		},
		IdempotencyKey: "",
		TenantId:       "",
	}

	// Workflow: 2
	wf2 := &suprsend.Workflow{
		Body: map[string]interface{}{
			"name":                  "__workflow_name__",
			"template":              "__template_slug__",
			"notification_category": "__category__", // system/transactional/promotional
			"users": []map[string]interface{}{
				{
					"distinct_id": "__distinct_id__",
					"$email":      []string{"__email__"},
				},
			},
		},
		IdempotencyKey: "123456",
		TenantId:       "default",
	}
	// ...... Add as many Workflow records as required.

	// Create bulk workflows instance
	bulkIns := suprClient.BulkWorkflows.NewInstance()
	// add all your workflows to bulkInstance
	bulkIns.Append(wf1, wf2)
	// Trigger
	bulkResponse, err := bulkIns.Trigger()
	if err != nil {
		log.Println(err)
		//
	}
	log.Println(bulkResponse)
}

func bulkEventsExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	//
	ev1 := &suprsend.Event{
		EventName:  "__event_name__",
		DistinctId: "__distinct_id__",
		Properties: map[string]interface{}{
			"k1": "v1",
		},
	}
	ev2 := &suprsend.Event{
		EventName:  "__event_name__",
		DistinctId: "__distinct_id__",
	}
	// Create bulkEvents instance
	bulkIns := suprClient.BulkEvents.NewInstance()
	// Add all events to bulk Instance
	bulkIns.Append(ev1, ev2)
	// call trigger to send all these events to SuprSend
	bulkResponse, err := bulkIns.Trigger()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bulkResponse)
}

func bulkUserProfileUpdateExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// create bulkUsers instance
	bulkIns := suprClient.BulkUsers.NewInstance()

	// Prepare user1
	user1 := suprClient.Users.GetInstance("sanjeev1")
	user1.AddEmail("user1@example.com")
	user1.AddWhatsapp("+1909090900")

	// prepare user 2
	user2 := suprClient.Users.GetInstance("sanjeev1")
	user2.AddEmail("user2@example.com")
	user2.AddWhatsapp("+2909090900")

	// Append all users to bulk instance
	bulkIns.Append(user1, user2)

	// Call save
	bulkResponse, err := bulkIns.Save()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bulkResponse)
}

func tenantExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// ================= Fetch existing tenant by ID
	tenant1, err := suprClient.Tenants.Get(context.Background(), "__tenant_id__")
	if err != nil {
		log.Println(err)
	}
	log.Println(tenant1)

	// ================= Fetch all tenants
	tenantsList, err := suprClient.Tenants.List(context.Background(), &suprsend.TenantListOptions{Limit: 10})
	if err != nil {
		log.Println(err)
	}
	log.Println(tenantsList)

	// ================= Update/Insert a tenant
	tenantPayload := &suprsend.Tenant{
		TenantName: suprsend.String("Tenant Name"),
		Logo:       suprsend.String("Tenant logo url"),
		Timezone:   suprsend.String("America/Los_Angeles"),
		// BlockedChannels: []string{},
		// EmbeddedPreferenceUrl:  suprsend.String("https://company-url.com/preferences"),
		// HostedPreferenceDomain: suprsend.String("preferences.suprsend.com"),
		PrimaryColor:   suprsend.String("#FFFFFF"),
		SecondaryColor: suprsend.String("#000000"),
		TertiaryColor:  nil,
		SocialLinks: &suprsend.TenantSocialLinks{
			Facebook: suprsend.String("https://facebook.com/tenant"),
		},
		Properties: map[string]interface{}{
			"k1": "tenant settings 1",
			"k2": "tenant settings 2",
		},
	}
	res, err := suprClient.Tenants.Upsert(context.Background(), "__tenant_id__", tenantPayload)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)
	// -- Delete tenant
	err = suprClient.Tenants.Delete(context.Background(), "__tenant_id__")
	if err != nil {
		log.Fatalln(err)
	}
}

func subscriberListExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	// ================= Fetch all subscriber lists
	allSubscriberList, err := suprClient.SubscriberLists.GetAll(ctx, &suprsend.SubscriberListAllOptions{Limit: 10, Offset: 0})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(allSubscriberList)

	// ================= Create subscriber list
	subscriberListCreated, err := suprClient.SubscriberLists.Create(ctx, &suprsend.SubscriberListCreateInput{
		ListId:          "users-with-prepaid-vouchers-1", // max length 64 characters
		ListName:        "Users With Prepaid Vouchers above $250",
		ListDescription: "Users With Prepaid Vouchers above $250",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(subscriberListCreated)

	// ================= Fetch existing subscriber-list
	existingSubsList, err := suprClient.SubscriberLists.Get(ctx, "users-with-prepaid-vouchers-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(existingSubsList)

	// ================= Add users to a list
	addDistinctIds := []string{"distinct_id_1", "distinct_id_2"}
	addResponse, err := suprClient.SubscriberLists.Add(ctx, "users-with-prepaid-vouchers-1", addDistinctIds)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(addResponse)

	// ================= remove users from a list
	removeDistinctIds := []string{"distinct_id_1", "distinct_id_2"}
	removeResponse, err := suprClient.SubscriberLists.Remove(ctx, "users-with-prepaid-vouchers-1", removeDistinctIds)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(removeResponse)

	// ================= broadcast to a list
	broadcastIns := &suprsend.SubscriberListBroadcast{
		Body: map[string]interface{}{
			"list_id":               "users-with-prepaid-vouchers-1",
			"template":              "template slug",
			"notification_category": "category",
			// broadcast channels.
			// if empty: broadcast will be tried on all available channels
			// if present: broadcast will be tried on passed channels only
			"channels": []string{"email"},
			"delay":    "1m", // check docs for delay format
			// "trigger_at":            "", // check docs for trigger_at format
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
		},
		IdempotencyKey: "",
		TenantId:       "",
	}
	// If need to add attachment
	// err = broadcastIns.AddAttachment("https://attachment-url", &suprsend.AttachmentOption{IgnoreIfError: true})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	res, err := suprClient.SubscriberLists.Broadcast(ctx, broadcastIns)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)

	// delete active list
	deleteListResp, err := suprClient.SubscriberLists.Delete(ctx, "users-with-prepaid-vouchers-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("delete list resp: ", deleteListResp)
}

func subscriberListVersioningExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()

	// ================= Create subscriber list
	subscriberListCreated, err := suprClient.SubscriberLists.Create(ctx, &suprsend.SubscriberListCreateInput{
		ListId:          "users-with-prepaid-vouchers-1", // max length 64 characters
		ListName:        "Users With Prepaid Vouchers above $250",
		ListDescription: "Users With Prepaid Vouchers above $250",
		TrackUserEntry:  suprsend.Bool(false),
		TrackUserExit:   suprsend.Bool(false),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list create resp: ", subscriberListCreated)

	// ================= Fetch existing subscriber-list
	existingSubsList, err := suprClient.SubscriberLists.Get(ctx, "users-with-prepaid-vouchers-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("fetch list resp: ", existingSubsList)

	// start sync
	newVersion, err := suprClient.SubscriberLists.StartSync(ctx, "users-with-prepaid-vouchers-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("start sync resp: ", newVersion)

	versionId := newVersion.VersionId

	// =============================== fetch active and draft lists after start sync
	subsList, err := suprClient.SubscriberLists.Get(ctx, "users-with-prepaid-vouchers-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("fetch list resp: ", subsList)

	subsListVersion, err := suprClient.SubscriberLists.GetVersion(ctx, "users-with-prepaid-vouchers-1", versionId)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("fetch list version resp: ", subsListVersion)

	// ================= Add users to a draft list (with versionId)
	addDistinctIds := []string{"id-399999", "id-399998"}
	addResponse, err := suprClient.SubscriberLists.AddToVersion(ctx, "users-with-prepaid-vouchers-1", versionId, addDistinctIds)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("add to version resp:", addResponse)

	// ================= remove users from a list
	removeDistinctIds := []string{"distinct_id_1", "distinct_id_2"}
	removeResponse, err := suprClient.SubscriberLists.RemoveFromVersion(ctx, "users-with-prepaid-vouchers-1", versionId, removeDistinctIds)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("remove from version resp: ", removeResponse)

	// finish sync
	finishSyncResp, err := suprClient.SubscriberLists.FinishSync(ctx, "users-with-prepaid-vouchers-1", versionId)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("finish sync resp: ", finishSyncResp)

	// ******************************************* delete version **************************************//
	// create a new version to be deleted later
	tempListVersion, err := suprClient.SubscriberLists.StartSync(ctx, "users-with-prepaid-vouchers-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(tempListVersion)

	// delete versioned list
	deleteVersionResp, err := suprClient.SubscriberLists.DeleteVersion(ctx, "users-with-prepaid-vouchers-1", tempListVersion.VersionId)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("delete version resp: ", deleteVersionResp)
}

func objectCrudOperationsExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	objectType := "sdk"
	objectId := "go"
	//
	// ================= Create Object
	object, _ := suprClient.Objects.Upsert(context.Background(), objectType, objectId, map[string]interface{}{
		"prop1": "val1",
	})
	log.Println(object)

	// ============= get object

	object, _ = suprClient.Objects.Get(context.Background(), objectType, objectId)
	log.Println(object)

	// =============== list object

	objects, _ := suprClient.Objects.List(context.Background(), objectType, nil)
	log.Println(objects)

	// ============= edit object

	object, _ = suprClient.Objects.Edit(context.Background(), objectType, objectId, map[string]interface{}{
		"prop1": "val1",
		"operations": []map[string]interface{}{
			{
				"$set": map[string]any{
					"k1": "v1",
					"k2": "v2",
				},
			},
			{
				"$remove": map[string]any{
					"k1": "v1",
				},
			},
		},
	})
	log.Println(object)

	// -------- delete [create and delete]
	object, _ = suprClient.Objects.Upsert(context.Background(), "delete_obj", "delete_obj_id", map[string]interface{}{
		"prop1": "val1",
	})
	log.Println(object)

	object, _ = suprClient.Objects.Delete(context.Background(), "delete_obj", "delete_obj_id")
	log.Println(object)

	// ----- bulk delete
	object, _ = suprClient.Objects.Upsert(context.Background(), "ot1", "oid1", map[string]interface{}{
		"prop1": "val1",
	})
	log.Println(object)

	object, _ = suprClient.Objects.Upsert(context.Background(), "ot1", "oid2", map[string]interface{}{
		"prop1": "val1",
	})
	log.Println(object)

	object, _ = suprClient.Objects.BulkDelete(context.Background(), "ot1", map[string]any{
		"object_ids": []string{"oid1", "oid2"},
	})
	log.Println(object)

	// ====================== subscriptions
	object, _ = suprClient.Objects.CreateSubscriptions(context.Background(), objectType, objectId, map[string]any{
		"recipients": []string{"praveen@suprsend.com"},
	})
	log.Println(object)

	object, _ = suprClient.Objects.GetSubscriptions(context.Background(), objectType, objectId, nil)
	log.Println(object)

	object, _ = suprClient.Objects.DeleteSubscriptions(context.Background(), objectType, objectId, map[string]any{
		"recipients": []string{"praveen@suprsend.com"},
	})
	log.Println(object)
}

func updateObjectPropertiesExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	objectType := "sdk"
	objectId := "go"
	//
	object, _ := suprClient.Objects.GetInstance(objectType, objectId)
	// Add email channel
	object.AddEmail("user@example.com")
	// add sms channel
	object.AddSms("+1444455555")
	// Add whatsapp channel
	object.AddWhatsapp("+1444455555")
	// Add androidpush token, token providers: fcm/xiaomi
	object.AddAndroidpush("__fcm_push_token__", "fcm")
	// Add iospush token, token providers: apns
	object.AddIospush("__ios_push_token__", "apns")
	// Add webpush token (vapid)
	object.AddWebpush(map[string]interface{}{
		"keys": map[string]interface{}{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")

	// add slack using email
	object.AddSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// add slack using user_id
	object.AddSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// add slack using channel_id
	object.AddSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// add slack using incoming-webhook
	object.AddSlack(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// DM on Team's channel using conversation id
	object.AddMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// add teams via DM user using team user id
	object.AddMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// add teams using incoming webhook
	object.AddMSTeams(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})
	//
	// remove email channel
	object.RemoveEmail("user@example.com")
	// remove sms channel
	object.RemoveSms("+1444455555")
	// remove whatsapp channel
	object.RemoveWhatsapp("+1444455555")
	// remove androidpush token
	object.RemoveAndroidpush("__fcm_push_token__", "fcm")
	// remove iospush token
	object.RemoveIospush("__ios_push_token__", "apns")

	// remove webpush token
	object.RemoveWebpush(map[string]interface{}{
		"keys": map[string]interface{}{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")
	// remove slack using email
	object.RemoveSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// remove slack using user_id
	object.RemoveSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// remove slack using channel_id
	object.RemoveSlack(map[string]interface{}{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// remove slack using incoming-webhook
	object.RemoveSlack(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// remove teams via DM on Team's channel using conversation id
	object.RemoveMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// remove teams via DM user using team user id
	object.RemoveMSTeams(map[string]interface{}{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// remove teams using incoming webhook
	object.RemoveMSTeams(map[string]interface{}{
		"incoming_webhook": map[string]interface{}{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})

	// Set user preferred language. languageCode must be in [ISO 639-1 2-letter] format
	object.SetPreferredLanguage("en")
	// set timezone property at subscriber level based on IANA timezone info
	object.SetTimezone("America/Los_Angeles")
	// If you need to remove all emails for this user, call user.Unset(["$email"])
	object.Unset([]string{"$email"})
	// set a user property using a map
	object.Set(map[string]interface{}{"prop1": "val1", "prop2": "val2"})
	// set a user property using a key, value pair
	object.SetKV("prop", "value")
	// set a user property once using map
	object.SetOnce(map[string]interface{}{"prop3": "val3"})
	// set a user property once using a key value pair
	object.SetOnceKV("prop4", "val4")
	// increment an already existing user property using key value pair
	object.IncrementKV("increment_prop", 2)
	// increment an already existing property using map
	object.Increment(map[string]interface{}{"increment_prop1": 5})

	// Save user
	resp, err := object.Save()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}
