package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func objectApisExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	_ = ctx

	// List objects for a given object type
	resp, err := suprClient.Objects.List(ctx, "office_locations", &suprsend.CursorListApiOptions{Limit: 10})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	// Create/Upsert object
	o1, err := suprClient.Objects.Upsert(ctx, suprsend.ObjectIdentifier{ObjectType: "office_locations", Id: "nyc"},
		map[string]any{
			"prop1": "val1",
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(o1)

	// ============= Get object
	o2, err := suprClient.Objects.Get(ctx, suprsend.ObjectIdentifier{ObjectType: "office_locations", Id: "nyc"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(o2)

	// ============= edit object
	o3, err := suprClient.Objects.Edit(ctx,
		suprsend.ObjectEditRequest{
			Identifier: &suprsend.ObjectIdentifier{ObjectType: "office_locations", Id: "nyc"},
			Payload: map[string]any{
				"operations": []map[string]any{
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
			}})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(o3)

	// -------- delete [create and delete]
	oid4 := suprsend.ObjectIdentifier{ObjectType: "office_locations", Id: "la"}
	o4, err := suprClient.Objects.Upsert(ctx, oid4, map[string]any{
		"prop1": "val1",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(o4)

	err = suprClient.Objects.Delete(ctx, oid4)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("deleted %s/%s\n", oid4.ObjectType, oid4.Id)

	// ----- bulk delete
	err = suprClient.Objects.BulkDelete(ctx, "office_locations", suprsend.ObjectBulkDeletePayload{
		ObjectIds: []string{"loc1", "loc2"},
	})
	if err != nil {
		log.Fatalln("Error Bulk delete: ", err)
	}

	// ====================== subscriptions
	// Create subscriptions
	oid5 := suprsend.ObjectIdentifier{ObjectType: "office_locations", Id: "nyc"}
	subs, err := suprClient.Objects.CreateSubscriptions(ctx, oid5, map[string]any{
		"recipients": []string{"distinct_id_1", "distinct_id_2"},
		"properties": map[string]any{
			"type": "admin",
		},
	})
	if err != nil {
		log.Fatalln("Error creating subscriptions: ", err)
	}
	log.Println(subs)

	// Fetch subscriptions
	subsList, err := suprClient.Objects.GetSubscriptions(ctx, oid5, &suprsend.CursorListApiOptions{Limit: 10})
	if err != nil {
		log.Fatalln("Error fetching subscriptions: ", err)
	}
	log.Println(subsList)

	// Delete subscriptions
	err = suprClient.Objects.DeleteSubscriptions(ctx, oid5, map[string]any{
		"recipients": []string{"distinct_id_1"},
	})
	if err != nil {
		log.Fatalln("Error deleting subscriptions: ", err)
	}
	log.Println("subscriptions deleted")

	// List objects, this object is subscribed to
	res, err := suprClient.Objects.GetObjectsSubscribedTo(ctx, suprsend.ObjectIdentifier{ObjectType: "office_locations", Id: "nyc"}, &suprsend.CursorListApiOptions{Limit: 10})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)
}

func objectEditApiExample() {
	// Instantiate Client
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	obj := suprsend.ObjectIdentifier{
		ObjectType: "office_locations",
		Id:         "nyc",
	}
	//
	o := suprClient.Objects.GetEditInstance(obj)
	// Add email channel
	o.AddEmail("user@example.com")
	// add sms channel
	o.AddSms("+1444455555")
	// Add whatsapp channel
	o.AddWhatsapp("+1444455555")
	// Add androidpush token, token providers: fcm/xiaomi
	o.AddAndroidpush("__fcm_push_token__", "fcm")
	// Add iospush token, token providers: apns
	o.AddIospush("__ios_push_token__", "apns")
	// Add webpush token (vapid)
	o.AddWebpush(map[string]any{
		"keys": map[string]any{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")

	// add slack using email
	o.AddSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// add slack using user_id
	o.AddSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// add slack using channel_id
	o.AddSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// add slack using incoming-webhook
	o.AddSlack(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// DM on Team's channel using conversation id
	o.AddMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// add teams via DM user using team user id
	o.AddMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// add teams using incoming webhook
	o.AddMSTeams(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})
	//
	// remove email channel
	o.RemoveEmail("user@example.com")
	// remove sms channel
	o.RemoveSms("+1444455555")
	// remove whatsapp channel
	o.RemoveWhatsapp("+1444455555")
	// remove androidpush token
	o.RemoveAndroidpush("__fcm_push_token__", "fcm")
	// remove iospush token
	o.RemoveIospush("__ios_push_token__", "apns")

	// remove webpush token
	o.RemoveWebpush(map[string]any{
		"keys": map[string]any{
			"auth":   "",
			"p256dh": "",
		},
		"endpoint": "",
	}, "vapid")
	// remove slack using email
	o.RemoveSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "email": "user@example.com"})
	// remove slack using user_id
	o.RemoveSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "user_id": "UXXXXXXXX"})
	// remove slack using channel_id
	o.RemoveSlack(map[string]any{"access_token": "xoxb-xxxxxxxxxxxxxxxxx", "channel_id": "CXXXXXXXX"})
	// remove slack using incoming-webhook
	o.RemoveSlack(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://hooks.slack.com/services/TXXXXXX/BXXXXX/XXXXXXXXXXX"},
	})
	// remove teams via DM on Team's channel using conversation id
	o.RemoveMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "conversation_id": "XXXXXXXXXXXX",
	})
	// remove teams via DM user using team user id
	o.RemoveMSTeams(map[string]any{
		"tenant_id": "XXXXXXX", "service_url": "https://smba.trafficmanager.net/XXXXXXXXXX", "user_id": "XXXXXXXXXXXX",
	})
	// remove teams using incoming webhook
	o.RemoveMSTeams(map[string]any{
		"incoming_webhook": map[string]any{"url": "https://XXXXX.webhook.office.com/webhookb2/XXXXXXXXXX@XXXXXXXXXX/IncomingWebhook/XXXXXXXXXX/XXXXXXXXXX"},
	})

	// Set object locale. localeCode must be either in [ISO 639-1 2-letter] format OR lang-region format (e.g. en-US, fr-FR, etc.)
	// Locale ISO codes combine a language code (ISO 639-1 2-letter) and a country code (ISO 3166-1 alpha-2), separated by a hyphen
	o.SetLocale("en-US")
	// set timezone property at subscriber level based on IANA timezone info
	o.SetTimezone("America/Los_Angeles")
	// If you need to remove all emails for this user, call user.Unset(["$email"])
	o.Unset([]string{"$email"})
	// set a user property using a map
	o.Set(map[string]any{"prop1": "val1", "prop2": "val2"})
	// set a user property using a key, value pair
	o.SetKV("prop", "value")
	// set a user property once using map
	o.SetOnce(map[string]any{"prop3": "val3"})
	// set a user property once using a key value pair
	o.SetOnceKV("prop4", "val4")
	// increment an already existing user property using key value pair
	o.IncrementKV("increment_prop", 2)
	// increment an already existing property using map
	o.Increment(map[string]any{"increment_prop1": 5})

	ctx := context.Background()
	_, _ = ctx, o
	// Save user
	resp, err := suprClient.Objects.Edit(ctx, suprsend.ObjectEditRequest{
		EditInstance: o,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}
