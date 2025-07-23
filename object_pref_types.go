package suprsend

import (
	"fmt"
	"net/url"
	"strconv"
)

type ObjectFullPreferenceOptions struct {
	TenantId           string
	ShowOptOutChannels *bool
	Tags               string
}

func (opts *ObjectFullPreferenceOptions) BuildQuery() string {
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

type ObjectFullPreferenceResponse struct {
	Sections []struct {
		Name          *string                  `json:"name"`
		Subcategories []UserCategoryPreference `json:"subcategories"`
	} `json:"sections"`
	ChannelPreferences []ObjectGlobalChannelPreference `json:"channel_preferences"`
}

// ------------------------------------------------------------

type ObjectGlobalChannelPreference struct {
	Channel      string `json:"channel"`
	IsRestricted bool   `json:"is_restricted"`
}

type ObjectGlobalChannelsPreferenceOptions struct {
	TenantId string
}

func (opts *ObjectGlobalChannelsPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}
	return query.Encode()
}

type ObjectGlobalChannelsPreferenceUpdateBody struct {
	ChannelPreferences []ObjectGlobalChannelPreference `json:"channel_preferences"`
}

type ObjectGlobalChannelsPreferenceResponse struct {
	ChannelPreferences []ObjectGlobalChannelPreference `json:"channel_preferences"`
}

// ------------------------------------------------------------

type ObjectCategoriesPreferenceOptions struct {
	Limit  int
	Offset int
	//
	TenantId           string
	ShowOptOutChannels *bool
	Tags               string // can be a simple tag or a JSON string for advanced filtering
}

func (opts *ObjectCategoriesPreferenceOptions) BuildQuery() string {
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

type ObjectCategoriesPreferenceResponse struct {
	Meta    *ListApiMetaInfo           `json:"meta"`
	Results []ObjectCategoryPreference `json:"results"`
}

type ObjectCategoryPreference struct {
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

type ObjectCategoryPreferenceOptions struct {
	TenantId           string
	ShowOptOutChannels *bool
}

func (opts *ObjectCategoryPreferenceOptions) BuildQuery() string {
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

type ObjectUpdateCategoryPreferenceBody struct {
	Preference     string   `json:"preference"`
	OptOutChannels []string `json:"opt_out_channels"`
}
