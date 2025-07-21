package main

import (
	"context"
	"log"
)

func templatesApisExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	_ = ctx

	// --- List all the templates
	resp, err := suprClient.Templates.GetAllTemplates(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	// --- get details of a single template
	resp_2, err := suprClient.Templates.GetDetails(ctx, "__template_slug__")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp_2)

	// -- get details of template used in a channel
	resp_3, err := suprClient.Templates.GetChannelContent(ctx, "__template_slug__", "__channel_slug__")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp_3)
}
