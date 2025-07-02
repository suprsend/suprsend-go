package main

import (
	"context"
	"log"

	suprsend "github.com/suprsend/suprsend-go"
)

func preferencesApiExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()

	// ----- Users
	/// ----- get full preferences
	preferences, err := suprClient.Users.GetUserPreferences(ctx, "__distinct_id1__", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences)

	/// ----- get categories preferences
	preferences_cat, err := suprClient.Users.GetCategoriesPreferences(ctx, "__distinct_id1__", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cat)

	// ----- get categories preferences
	preferences_cat_1, err := suprClient.Users.GetGlobalChannelPreferences(ctx, "__distinct_id1__", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cat_1)

	// ----- get category preferences
	preferences_cat_2, err := suprClient.Users.GetCategoryPreference(ctx, "__distinct_id1__", "akhil-system", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cat_2)

	// --- update category preference within channel
	pref := "opt_in"
	ch1 := "email"
	ch2 := "sms"
	body := &suprsend.UserUpdateCategoryPreferenceBody{
		Preference:     &pref,
		OptOutChannels: []*string{&ch1, &ch2},
	}

	opts := &suprsend.UserCategoryUpdatePreferenceOptions{
		TenantId: "__tenant_id1__",
	}

	preferences_cat_ch, err := suprClient.Users.UpdateCategoryPreference(ctx, "__distinct_id1__", "__category_slug__", *body, opts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cat_ch)

	// ----- update global channel preferences
	channel_preferences := []suprsend.UserChannelPreferenceIn{
		{Channel: "email", IsRestricted: true},
		{Channel: "inbox", IsRestricted: true},
	}
	body_ch := suprsend.UserGlobalChannelPreferenceUpdateBody{
		ChannelPreferences: channel_preferences,
	}

	opts_ch := &suprsend.UserGlobalPreferenceOptions{
		TenantId: "__tenant_id1__",
	}

	preferences_ch, err := suprClient.Users.UpdateGlobalChannelPreferences(ctx, "__tenant_id1__", body_ch, opts_ch)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_ch)

	// ----- bulk update
	channel_preferences_b := []*suprsend.UserChannelPreferenceIn{
		{Channel: "email", IsRestricted: true},
		{Channel: "inbox", IsRestricted: false},
	}

	categories := []*suprsend.UserCategoryPreferenceIn{
		{Category: "__category_slug__", Preference: "opt_out", OptOutChannels: []*string{&ch1, &ch2}},
	}

	body_b := suprsend.UserBulkPreferenceUpdateBody{
		DistinctIDs:        []string{"__distinct_id1__"},
		ChannelPreferences: channel_preferences_b,
		Categories:         categories,
	}

	preferences_b, err := suprClient.Users.BulkUpdatePreferences(ctx, body_b, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_b)

	// -- bulk reset
	resetChannels := true
	resetCategories := true
	body_br := suprsend.UserBulkResetPreferenceBody{
		DistinctIDs:             []string{"__distinct_id1__"},
		ResetChannelPreferences: &resetChannels,
		ResetCategories:         &resetCategories,
	}

	opts_br := &suprsend.UserPreferenceResetOptions{
		TenantId: "__tenant_id1__",
	}

	preferences_br, err := suprClient.Users.ResetPreferences(ctx, body_br, opts_br)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_br)

	// ----- Tenants
	/// Get All Categories Preference for a tenant
	preferences_gt, err := suprClient.Tenants.GetAllCategoriesPreference(ctx, "__distinct_id1__")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_gt)

	//  Update Categories preference for a tenant
	body_c_t := suprsend.TenantPreferenceCategoryUpdateBody{
		Preference:          "opt_in",
		VisibleToSubscriber: true,
		MandatoryChannels:   []string{"email", "sms", "inbox"},
		BlockedChannels:     []string{"slack"},
	}
	preferences_c_t, err := suprClient.Tenants.UpdateCategoryPreference(ctx, "__tenant_id1__", "__category_slug__", body_c_t, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_c_t)

	// ----- Objects
	/// Get All Preferences for an object
	obj := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}
	preferences_ao, err := suprClient.Objects.GetPreference(ctx, obj, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_ao)

	// Get Preference Category for an object
	obj_co := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}
	preferences_co, err := suprClient.Objects.GetAllCategoriesPreference(ctx, obj_co, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_co)

	// Get Global Channels Preferences
	obj_cho := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}

	preferences_cho, err := suprClient.Objects.GetGlobalChannelsPreferences(ctx, obj_cho, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cho)

	// Get a single category preference for an object
	obj_cato := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}

	preferences_cato, err := suprClient.Objects.GetCategoryPreference(ctx, obj_cato, "__category_slug__", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cato)

	// Update category preferences
	obj_cpo := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}

	body_cpo := suprsend.ObjectUpdateCategoryPreferenceBody{
		Preference:     "opt_in",
		OptOutChannels: []string{"iospush", "slack"},
	}

	preferences_cpo, err := suprClient.Objects.UpdateCategoryPreference(ctx, obj_cpo, "__category_slug__", body_cpo, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(preferences_cpo)

	// Update global channel preferences
	obj_gco := suprsend.ObjectIdentifier{
		ObjectType: "__object_type__",
		Id:         "__object_id__",
	}

	channel_preferences_gco := []suprsend.UserChannelPreferenceIn{
		{Channel: "email", IsRestricted: true},
		{Channel: "inbox", IsRestricted: false},
	}

	body_gco := suprsend.ObjectGlobalChannelPreferenceUpdateBody{
		ChannelPreferences: channel_preferences_gco,
	}

	preferences_gco, err := suprClient.Objects.UpdateGlobalChannelsPreferences(ctx, obj_gco, body_gco, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(preferences_gco)
}
