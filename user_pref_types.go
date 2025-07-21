package suprsend

import (
	"fmt"
	"net/url"
	"strconv"
)

type UserGlobalChannelPreference struct {
	Channel      string `json:"channel"`
	IsRestricted bool   `json:"is_restricted"`
}

type UserPreferencesOptions struct {
	TenantId           string
	ShowOptOutChannels *bool
	Tags               string // can be a simple tag or a JSON string for advanced filtering
}

func (opts *UserPreferencesOptions) BuildQuery() string {
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

type UserGlobalPreferenceOptions struct {
	TenantId string `json:"tenant_id"`
}

func (opts *UserGlobalPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	return query.Encode()
}

type UserCategoryPreferenceOptions struct {
	TenantId           string `json:"tenant_id"`
	ShowOptOutChannels bool   `json:"show_opt_out_channels"`
}

func (opts *UserCategoryPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}

		if strconv.FormatBool(opts.ShowOptOutChannels) != "" {
			query.Set("show_opt_out_channels", strconv.FormatBool(opts.ShowOptOutChannels))
		}
	}

	return query.Encode()
}

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

type UserPreferenceResetOptions struct {
	TenantId string `json:"tenant_id"`
}

func (opts *UserPreferenceResetOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	return query.Encode()
}

type UserUpdateCategoryPreferenceBody struct {
	Preference     *string   `json:"preference"`
	OptOutChannels []*string `json:"opt_out_channels"`
}

type UserGlobalChannelPreferenceUpdateBody struct {
	ChannelPreferences []UserGlobalChannelPreference `json:"channel_preferences"`
}

type UserCategoryPreferenceIn struct {
	Category       string    `json:"category"`
	Preference     string    `json:"preference"`
	OptOutChannels []*string `json:"opt_out_channels,omitempty"`
}

type UserBulkPreferenceUpdateBody struct {
	DistinctIDs        []string                       `json:"distinct_ids,omitempty"`
	ChannelPreferences []*UserGlobalChannelPreference `json:"channel_preferences,omitempty"`
	Categories         []*UserCategoryPreferenceIn    `json:"categories,omitempty"`
}

type UserBulkResetPreferenceBody struct {
	DistinctIDs             []string `json:"distinct_ids"`
	ResetChannelPreferences *bool    `json:"reset_channel_preferences"`
	ResetCategories         *bool    `json:"reset_categories"`
}

type UserPreferencesResponse struct {
	Sections []struct {
		Name          *string                          `json:"name"`
		Subcategories []UserCategoryPreferenceResponse `json:"subcategories"`
	} `json:"sections"`
	ChannelPreferences []UserGlobalChannelPreferencesResponse `json:"channel_preferences"`
}

type UserGlobalChannelPreferencesResponse struct {
	ChannelPreferences []any `json:"channel_preferences"`
}

type UserCategoryPreferenceResponse struct {
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
}

type UserBulkPreferenceResponse struct {
	Success bool `json:"success"`
}

type UserCategoriesPreferenceResponse struct {
	Meta    *ListApiMetaInfo                 `json:"meta"`
	Results []UserCategoryPreferenceResponse `json:"results"`
}
