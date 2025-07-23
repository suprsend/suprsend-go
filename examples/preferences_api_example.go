package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func preferencesApiExample() {
	userPreferencesApiExample()
	objectPreferencesApiExample()
	tenantPreferencesApiExample()
}

func userPreferencesApiExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()

	// ----- get full preferences
	preferences, err := suprClient.Users.GetFullPreference(ctx, "__distinct_id1__", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences)

	// ------ get global channels preferences
	globalPrefs, err := suprClient.Users.GetGlobalChannelsPreference(ctx, "__distinct_id1__",
		&suprsend.UserGlobalChannelsPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(globalPrefs)

	// ----- update global channel preferences
	globalChannelsPrefbody := suprsend.UserGlobalChannelsPreferenceUpdateBody{
		ChannelPreferences: []suprsend.UserGlobalChannelPreference{
			{Channel: "email", IsRestricted: true},
			{Channel: "inbox", IsRestricted: true},
		},
	}
	globalPreferences, err := suprClient.Users.UpdateGlobalChannelsPreference(ctx, "__distinct_id1__", globalChannelsPrefbody,
		&suprsend.UserGlobalChannelsPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(globalPreferences)

	// ----- get categories preferences
	allCatsPreferences, err := suprClient.Users.GetAllCategoriesPreference(ctx, "__distinct_id1__",
		&suprsend.UserCategoriesPreferenceOptions{
			Limit:    10,
			Offset:   0,
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(allCatsPreferences)

	// ----- get category preference
	singleCatPreference, err := suprClient.Users.GetCategoryPreference(ctx, "__distinct_id1__", "__category_slug__",
		&suprsend.UserCategoryPreferenceOptions{
			TenantId:           "__tenant_id1__",
			ShowOptOutChannels: suprsend.Bool(true),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(singleCatPreference)

	// --- update category preference
	updateSingleCatPreference, err := suprClient.Users.UpdateCategoryPreference(ctx, "__distinct_id1__", "__category_slug__",
		suprsend.UserUpdateCategoryPreferenceBody{
			Preference:     "opt_in",
			OptOutChannels: []string{"email", "sms"},
		},
		&suprsend.UserCategoryPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(updateSingleCatPreference)

	// ----- bulk update
	bulkUpdateBody := suprsend.UserBulkPreferenceUpdateBody{
		DistinctIDs: []string{"__distinct_id1__"},
		ChannelPreferences: []*suprsend.UserGlobalChannelPreference{
			{Channel: "email", IsRestricted: true},
			{Channel: "inbox", IsRestricted: false},
		},
		Categories: []*suprsend.UserCategoryPreferenceIn{
			{Category: "__category_slug__", Preference: "opt_out", OptOutChannels: []string{"email", "sms"}},
		},
	}
	bulkPreferenceUpdateResponse, err := suprClient.Users.BulkUpdatePreferences(ctx, bulkUpdateBody,
		&suprsend.UserBulkPreferenceUpdateOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bulkPreferenceUpdateResponse)

	// -- bulk reset
	bulkPrefResetBody := suprsend.UserBulkPreferenceResetBody{
		DistinctIDs:             []string{"__distinct_id1__"},
		ResetChannelPreferences: true,
		ResetCategories:         true,
	}
	bulkPrefResetResp, err := suprClient.Users.ResetPreferences(ctx, bulkPrefResetBody,
		&suprsend.UserBulkPreferenceUpdateOptions{
			TenantId: "__tenant_id1__",
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bulkPrefResetResp)
}

func objectPreferencesApiExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()

	// ----- get full preferences
	obj := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}
	preferences, err := suprClient.Objects.GetFullPreference(ctx, obj,
		&suprsend.ObjectFullPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences)

	// ------ get global channels preferences
	globalPrefs, err := suprClient.Objects.GetGlobalChannelsPreference(ctx, obj,
		&suprsend.ObjectGlobalChannelsPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(globalPrefs)

	// ----- update global channel preferences
	globalChannelsPrefbody := suprsend.ObjectGlobalChannelsPreferenceUpdateBody{
		ChannelPreferences: []suprsend.ObjectGlobalChannelPreference{
			{Channel: "email", IsRestricted: true},
			{Channel: "inbox", IsRestricted: false},
		},
	}
	globalPreferences, err := suprClient.Objects.UpdateGlobalChannelsPreference(ctx, obj, globalChannelsPrefbody,
		&suprsend.ObjectGlobalChannelsPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(globalPreferences)

	// ----- get categories preferences
	allCatsPreferences, err := suprClient.Objects.GetAllCategoriesPreference(ctx, obj,
		&suprsend.ObjectCategoriesPreferenceOptions{
			Limit:    10,
			Offset:   0,
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(allCatsPreferences)

	// ----- Get a single category preference
	singleCatPreference, err := suprClient.Objects.GetCategoryPreference(ctx, obj, "__category_slug__",
		&suprsend.ObjectCategoryPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(singleCatPreference)

	// --- update category preference
	updateSingleCatPreference, err := suprClient.Objects.UpdateCategoryPreference(ctx, obj, "__category_slug__",
		suprsend.ObjectUpdateCategoryPreferenceBody{
			Preference:     "opt_in",
			OptOutChannels: []string{"iospush", "slack"},
		},
		&suprsend.ObjectCategoryPreferenceOptions{
			TenantId: "__tenant_id1__",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(updateSingleCatPreference)
}

func tenantPreferencesApiExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	// ----- Tenants
	// Get All Categories Preference for a tenant
	preferences_gt, err := suprClient.Tenants.GetAllCategoriesPreference(ctx, "__tenant_id__",
		&suprsend.TenantCategoriesPreferenceOptions{
			Limit:  10,
			Offset: 0,
			// Tags:   "tag1",
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_gt)

	//  Update Categories preference for a tenant
	body_c_t := suprsend.TenantCategoryPreferenceUpdateBody{
		Preference:          "opt_in",
		VisibleToSubscriber: suprsend.Bool(true),
		MandatoryChannels:   []string{"email", "sms", "inbox"},
		BlockedChannels:     []string{"slack"},
	}
	preferences_c_t, err := suprClient.Tenants.UpdateCategoryPreference(ctx, "__tenant_id__", "__category_slug__", body_c_t)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_c_t)
}
