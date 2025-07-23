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
