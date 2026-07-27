package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func messagesApisExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	opts := &suprsend.MessageListOptions{
		Limit: 10,
	}
	// Fetch Messages list
	resp, err := suprClient.Messages.List(ctx, opts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	// -- Update messages
	messages := []suprsend.MessageUpdateItem{
		{MessageID: "xx1", Action: "read"},
		{MessageID: "xx2", Action: "read"},
	}
	updateResp, err := suprClient.Messages.BulkUpdate(ctx, messages)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(updateResp)
}
