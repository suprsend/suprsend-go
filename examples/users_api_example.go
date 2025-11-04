package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func userApisExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	opts := &suprsend.CursorListApiOptions{
		Limit: 10,
	}
	// Fetch Users list
	users, err := suprClient.Users.List(ctx, opts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(users)

	// -- Upsert user
	userProps := map[string]any{
		"prop1": "val1",
	}
	user2, err := suprClient.Users.Upsert(ctx, "__distinct_id1__", userProps)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(user2)

	// Get user by distinct_id
	user, err := suprClient.Users.Get(ctx, "__distinct_id1__")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(user)

	// -- Edit user (by preparing exact payload)
	editReq := suprsend.UserEditRequest{
		DistinctId: "__distinct_id2__",
		Payload: map[string]any{
			"operations": []map[string]any{
				{
					"$set": map[string]any{
						"prop1": "val1",
					},
				},
			},
		},
	}
	resp, err := suprClient.Users.Edit(ctx, editReq)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	// --- Merge API
	res, err := suprClient.Users.Merge(ctx, "__distinct_id1__", suprsend.UserMergeRequest{
		FromUserId: "__distinct_id2__",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)

	// Delete user
	user3, err := suprClient.Users.Upsert(ctx, "__distinct_id2__", map[string]any{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(user3)
	err = suprClient.Users.Delete(ctx, "__distinct_id2__")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("User deleted")

	// Bulk delete users
	err = suprClient.Users.BulkDelete(ctx, suprsend.UserBulkDeletePayload{DistinctIds: []string{"123"}})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("User deleted")

	// Get objects the user is subscribed to
	resp2, err := suprClient.Users.GetObjectsSubscribedTo(ctx, "__distinct_id__", &suprsend.CursorListApiOptions{Limit: 10})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp2)

	// Get Subscriber-ists, the user is subscribed to
	resp3, err := suprClient.Users.GetListsSubscribedTo(ctx, "__distinct_id__", &suprsend.CursorListApiOptions{Limit: 10})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp3)
}

func userEditApiExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	//
	user := suprClient.Users.GetEditInstance("__distinct_id__")
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
	user.AddWebpush(map[string]any{
		"keys": map[string]any{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")

	// add slack using email
	user.AddSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// add slack using user_id
	user.AddSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// add slack using channel_id
	user.AddSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// add slack using incoming-webhook
	user.AddSlack(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// DM on Team's channel using conversation id
	user.AddMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// add teams via DM user using team user id
	user.AddMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// add teams using incoming webhook
	user.AddMSTeams(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
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
	user.RemoveWebpush(map[string]any{
		"keys": map[string]any{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")
	// remove slack using email
	user.RemoveSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// remove slack using user_id
	user.RemoveSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// remove slack using channel_id
	user.RemoveSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// remove slack using incoming-webhook
	user.RemoveSlack(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// remove teams via DM on Team's channel using conversation id
	user.RemoveMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// remove teams via DM user using team user id
	user.RemoveMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// remove teams using incoming webhook
	user.RemoveMSTeams(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})

	// Set user preferred locale. languageCode must be in [ISO 639-1 2-letter] format
	user.SetLocale("en")
	// deprecated: use SetLocale instead
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
	user.Set(map[string]any{"prop1": "val1", "prop2": "val2"})
	// set a user property using a key, value pair
	user.SetKV("prop", "value")
	// set a user property once using map
	user.SetOnce(map[string]any{"prop3": "val3"})
	// set a user property once using a key value pair
	user.SetOnceKV("prop4", "val4")
	// increment an already existing user property using key value pair
	user.IncrementKV("increment_prop", 2)
	// increment an already existing property using map
	user.Increment(map[string]any{"increment_prop1": 5})

	ctx := context.Background()
	_, _ = ctx, user
	// Sync API: Update user
	resp, err := suprClient.Users.Edit(ctx, suprsend.UserEditRequest{EditInstance: user})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	// Async Edit API: Update user
	resp2, err := suprClient.Users.AsyncEdit(ctx, user)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp2)
}

func userEditBulkExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// create bulkUsers instance
	bulkIns := suprClient.Users.GetBulkEditInstance()

	// Prepare user1
	user1 := suprClient.Users.GetEditInstance("distinct_id_1")
	user1.AddEmail("user1@example.com")
	user1.AddWhatsapp("+1909090900")

	// prepare user 2
	user2 := suprClient.Users.GetEditInstance("distinct_id_2")
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
