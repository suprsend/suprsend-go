package suprsend

import (
	"fmt"
	"net/url"
	"strconv"
)

type UserFullPreferencesOptions struct {
	TenantId           string
	ShowOptOutChannels *bool
	Tags               string // can be a simple tag or a JSON string for advanced filtering
}

func (opts *UserFullPreferencesOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.ShowOptOutChannels != nil {
			query.Set("show_opt_out_channels", fmt.Sprintf("%v", *opts.ShowOptOutChannels))
		}
		if opts.Tags != "" {
			query.Set("tags", opts.Tags)
		}
	}
	return query.Encode()
}

type UserFullPreferenceResponse struct {
	Sections []struct {
		Name          *string                  `json:"name"`
		Subcategories []UserCategoryPreference `json:"subcategories"`
	} `json:"sections"`
	ChannelPreferences []UserGlobalChannelPreference `json:"channel_preferences"`
}

// ------------------------------------------------------------

type UserGlobalChannelPreference struct {
	Channel      string `json:"channel"`
	IsRestricted bool   `json:"is_restricted"`
}

type UserGlobalChannelsPreferenceOptions struct {
	TenantId string
}

func (opts *UserGlobalChannelsPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}
	return query.Encode()
}

type UserGlobalChannelsPreferenceResponse struct {
	ChannelPreferences []UserGlobalChannelPreference `json:"channel_preferences"`
}

type UserGlobalChannelsPreferenceUpdateBody struct {
	ChannelPreferences []UserGlobalChannelPreference `json:"channel_preferences"`
}

// ------------------------------------------------------------

type UserCategoriesPreferenceOptions struct {
	Limit  int
	Offset int
	//
	TenantId           string
	ShowOptOutChannels *bool
	Tags               string // can be a simple tag or a JSON string for advanced filtering
}

func (opts *UserCategoriesPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			query.Set("offset", strconv.Itoa(opts.Offset))
		}
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.ShowOptOutChannels != nil {
			query.Set("show_opt_out_channels", fmt.Sprintf("%v", *opts.ShowOptOutChannels))
		}
		if opts.Tags != "" {
			query.Set("tags", opts.Tags)
		}
	}
	return query.Encode()
}

type UserCategoriesPreferenceResponse struct {
	Meta    *ListApiMetaInfo         `json:"meta"`
	Results []UserCategoryPreference `json:"results"`
}

type UserCategoryPreference struct {
	Name               string  `json:"name"`
	Category           string  `json:"category"`
	Description        string  `json:"description"`
	OriginalPreference *string `json:"original_preference"`
	Preference         string  `json:"preference"`
	IsEditable         bool    `json:"is_editable"`
	Channels           []struct {
		Channel    string `json:"channel"`
		Preference string `json:"preference"`
		IsEditable bool   `json:"is_editable"`
	} `json:"channels"`
	Tags          []string `json:"tags"`
	EffectiveTags []string `json:"effective_tags"`
	//
	DigestScheduleOptions *UserCategoryDigestScheduleOptionsOut `json:"digest_schedule_options,omitzero"`
	DigestSchedule        *UserCategoryDigestScheduleOut        `json:"digest_schedule,omitzero"`
	Properties            []PreferenceCategoryPropertyOut       `json:"properties,omitzero"`
}

type UserCategoryDigestScheduleOptionsOut struct {
	Options []UserCategoryDigestScheduleOptionOut `json:"options"`
}

type UserCategoryDigestScheduleOptionOut struct {
	Id        string                                  `json:"id"`
	Label     string                                  `json:"label"`
	Frequency string                                  `json:"frequency"`
	Interval  int                                     `json:"interval,omitempty"`
	IsDefault bool                                    `json:"is_default"`
	Time      *UserCategoryDigestScheduleTimeOut      `json:"time,omitempty"`
	Weekdays  *UserCategoryDigestScheduleWeekdaysOut  `json:"weekdays,omitempty"`
	Monthdays *UserCategoryDigestScheduleMonthdaysOut `json:"monthdays,omitempty"`
	//
	IsUserSelected bool `json:"is_user_selected"`
}

type UserCategoryDigestScheduleTimeOut struct {
	EditPolicy   string `json:"edit_policy"`
	DefaultValue string `json:"default_value"`
	Value        string `json:"value,omitempty"`
}

type UserCategoryDigestScheduleWeekdaysOut struct {
	EditPolicy   string   `json:"edit_policy"`
	DefaultValue []string `json:"default_value"`
	Value        []string `json:"value,omitempty"`
}

type UserCategoryDigestScheduleMonthdaysOut struct {
	EditPolicy   string                            `json:"edit_policy"`
	DefaultValue []DigestScheduleMonthdayComponent `json:"default_value"`
	Value        []DigestScheduleMonthdayComponent `json:"value,omitempty"`
}

type UserCategoryDigestScheduleOut struct {
	Id        string                            `json:"id"`
	Label     string                            `json:"label"`
	Frequency string                            `json:"frequency"`
	Interval  int                               `json:"interval,omitempty"`
	IsDefault bool                              `json:"is_default"`
	Time      string                            `json:"time,omitempty"`
	Weekdays  []string                          `json:"weekdays,omitempty"`
	Monthdays []DigestScheduleMonthdayComponent `json:"monthdays,omitempty"`
	//
	IsUserSelected bool `json:"is_user_selected"`
}

type UserCategoryDigestScheduleIn struct {
	Id        string                            `json:"id"`
	Time      string                            `json:"time,omitempty"`
	Weekdays  []string                          `json:"weekdays,omitempty"`
	Monthdays []DigestScheduleMonthdayComponent `json:"monthdays,omitempty"`
}

// ------------------------------------------------------------

type UserCategoryPreferenceOptions struct {
	TenantId           string
	ShowOptOutChannels *bool
}

func (opts *UserCategoryPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.ShowOptOutChannels != nil {
			query.Set("show_opt_out_channels", fmt.Sprintf("%v", *opts.ShowOptOutChannels))
		}
	}
	return query.Encode()
}

type UserUpdateCategoryPreferenceBody struct {
	Preference     string   `json:"preference"`
	OptOutChannels []string `json:"opt_out_channels"`
	//
	DigestSchedule Nullable[UserCategoryDigestScheduleIn]   `json:"digest_schedule,omitzero"`
	Properties     Nullable[[]PreferenceCategoryPropertyIn] `json:"properties,omitzero"`
}

// ------------------------------------------------------------

type UserBulkPreferenceUpdateOptions struct {
	TenantId string `json:"tenant_id"`
}

func (opts *UserBulkPreferenceUpdateOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}
	return query.Encode()
}

type UserBulkPreferenceUpdateBody struct {
	DistinctIDs        []string                       `json:"distinct_ids"`
	ChannelPreferences []*UserGlobalChannelPreference `json:"channel_preferences,omitempty"`
	Categories         []*UserCategoryPreferenceIn    `json:"categories,omitempty"`
}

type UserCategoryPreferenceIn struct {
	Category       string   `json:"category"`
	Preference     string   `json:"preference"`
	OptOutChannels []string `json:"opt_out_channels"`
	//
	DigestSchedule Nullable[UserCategoryDigestScheduleIn]   `json:"digest_schedule,omitzero"`
	Properties     Nullable[[]PreferenceCategoryPropertyIn] `json:"properties,omitzero"`
}

type UserBulkPreferenceUpdateResponse struct {
	Success bool `json:"success"`
}

// ------------------------------------------------------------

type UserBulkPreferenceResetBody struct {
	DistinctIDs             []string `json:"distinct_ids"`
	ResetChannelPreferences bool     `json:"reset_channel_preferences"`
	ResetCategories         bool     `json:"reset_categories"`
}
