package suprsend

import (
	"net/url"
	"strconv"
)

type ObjectGlobalChannelPreference struct {
	Channel      string `json:"channel"`
	IsRestricted bool   `json:"is_restricted"`
}

type ObjectPreferenceOptions struct {
	TenantId           string `json:"tenant_id"`
	ShowOptOutChannels bool   `json:"show_opt_out_channels"`
	Tags               string `json:"tags"`
}

func (opts *ObjectPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
		if opts.Tags != "" {
			query.Set("tags", opts.Tags)
		}
		if strconv.FormatBool(opts.ShowOptOutChannels) != "" {
			query.Set("show_opt_out_channels", strconv.FormatBool(opts.ShowOptOutChannels))
		}
	}

	return query.Encode()
}

type ObjectGlobalPreferenceOptions struct {
	TenantId string `json:"tenant_id"`
}

func (opts *ObjectGlobalPreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	return query.Encode()
}

type ObjectCategoryPreferenceOptions struct {
	TenantId           string `json:"tenant_id"`
	ShowOptOutChannels bool   `json:"show_opt_out_channels"`
}

func (opts *ObjectCategoryPreferenceOptions) BuildQuery() string {
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

type ObjectCategoryUpdatePreferenceOptions struct {
	TenantId string `json:"tenant_id"`
}

func (opts *ObjectCategoryUpdatePreferenceOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.TenantId != "" {
			query.Set("tenant_id", opts.TenantId)
		}
	}

	return query.Encode()
}

type ObjectUpdateCategoryPreferenceBody struct {
	Preference     string   `json:"preference"`
	OptOutChannels []string `json:"opt_out_channels"`
}

type ObjectGlobalChannelPreferenceUpdateBody struct {
	ChannelPreferences []ObjectGlobalChannelPreference `json:"channel_preferences"`
}

type ObjectPreferenceResponse struct {
	Sections           []any `json:"sections"`
	ChannelPreferences []any `json:"channel_preferences"`
}

type ObjectGlobalChannelPreferencesResponse struct {
	ChannelPreferences []any `json:"channel_preferences"`
}

type ObjectCategoryPreferenceResponse struct {
	Name               string `json:"name"`
	Category           string `json:"category"`
	Description        string `json:"description"`
	OriginalPreference string `json:"original_preference"`
	Preference         string `json:"preference"`
	IsEditable         bool   `json:"is_editable"`
	Channels           []any  `json:"channels"`
}

type ObjectCategoriesPreferenceResponse struct {
	Meta    *ListApiMetaInfo                   `json:"meta"`
	Results []ObjectCategoryPreferenceResponse `json:"results"`
}
