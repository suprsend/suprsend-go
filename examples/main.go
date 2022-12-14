package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func main() {
	triggerWorkflowExample()
	sendEventExample()
	updateUserProfileExample()
	//
	bulkWorkflowsExample()
	bulkEventsExample()
	bulkUserProfileUpdateExample()
	//
	brandExample()
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

func triggerWorkflowExample() {
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
		BrandId:        "",
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
	}
	// Add attachment (If needed) by calling .AddAttachment
	err = ev.AddAttachment("~/Downloads/attachment.pdf", &suprsend.AttachmentOption{FileName: "My Attachment.pdf"})
	if err != nil {
		log.Println(err)
	}
	// Send event to Suprsend by calling .TrackEvent
	_, err = suprClient.TrackEvent(ev)
	if err != nil {
		log.Fatalln(err)
	}
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

	// Set user preferred language. languageCode must be in [ISO 639-1 2-letter] format
	user.SetPreferredLanguage("en")

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

	// Save user
	_, err = user.Save()
	if err != nil {
		log.Fatalln(err)
	}
}

func bulkWorkflowsExample() {
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
		BrandId:        "",
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
		BrandId:        "default",
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
	user1.AddEmail("user2@example.com")
	user1.AddWhatsapp("+2909090900")

	// Append all users to bulk instance
	bulkIns.Append(user1, user2)

	// Call save
	bulkResponse, err := bulkIns.Save()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bulkResponse)
}

func brandExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// ================= Fetch existing brand by ID
	brand1, err := suprClient.Brands.Get(context.Background(), "__brand_id__")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(brand1)

	// ================= Fetch all brands
	brandsList, err := suprClient.Brands.List(context.Background(), &suprsend.BrandListOptions{Limit: 10})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(brandsList)

	// ================= Update/Insert a brand
	brandPayload := &suprsend.Brand{
		BrandName:      suprsend.String("Brand Name"),
		Logo:           suprsend.String("Brand logo url"),
		PrimaryColor:   suprsend.String("#FFFFFF"),
		SecondaryColor: suprsend.String("#000000"),
		TertiaryColor:  nil,
		SocialLinks: &suprsend.BrandSocialLinks{
			Facebook: suprsend.String("https://facebook.com/brand"),
		},
		Properties: map[string]interface{}{
			"k1": "brand settings 1",
			"k2": "brand settings 2",
		},
	}
	res, err := suprClient.Brands.Upsert(context.Background(), "__brand_id__", brandPayload)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)
}
