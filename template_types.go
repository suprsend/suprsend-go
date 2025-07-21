package suprsend

import (
	"net/url"
	"strconv"
	"time"
)

type TemplateOptions struct {
	HasTagIDsAny   string `url:"has_tag_ids_any,omitempty"`  // Comma-separated tag IDs
	HasChannelsAny string `url:"has_channels_any,omitempty"` // Comma-separated channels (e.g., "email,sms")
	IsActive       *bool  `url:"is_active,omitempty"`        // true = published, false = draft
	IsArchived     *bool  `url:"is_archived,omitempty"`      // true = archived, default is false
}

func (opts *TemplateOptions) BuildQuery() string {
	query := url.Values{}
	if opts != nil {
		if opts.HasTagIDsAny != "" {
			query.Set("has_tag_ids_any", opts.HasTagIDsAny)

		}
		if opts.HasChannelsAny != "" {
			query.Set("has_channels_any", opts.HasChannelsAny)
		}
		if strconv.FormatBool(*opts.IsActive) != "" {
			query.Set("is_active", strconv.FormatBool(*opts.IsActive))
		}
		if strconv.FormatBool(*opts.IsArchived) != "" {
			query.Set("is_archived", strconv.FormatBool(*opts.IsArchived))
		}
	}

	return query.Encode()
}

type ChannelTemplateResponse struct {
	ID      int `json:"id"`
	Channel struct {
		Name                     string `json:"name"`
		Slug                     string `json:"slug"`
		IsTemplateApprovalNeeded bool   `json:"is_template_approval_needed"`
	} `json:"channel"`
	IsActive          bool      `json:"is_active"`
	IsEnabled         bool      `json:"is_enabled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DisabledLanguages []string  `json:"disabled_languages"`
	Versions          []struct {
		ID        int `json:"id"`
		Templates []struct {
			ID       int `json:"id"`
			Language struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"language"`
			IsEnabled      bool      `json:"is_enabled"`
			ApprovalStatus string    `json:"approval_status"`
			Content        any       `json:"content"` // Use appropriate struct if known
			CreatedAt      time.Time `json:"created_at"`
			UpdatedAt      time.Time `json:"updated_at"`
			UpdatedBy      struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"updated_by"`
			ApprovalCycle           any  `json:"approval_cycle"` // Use appropriate struct if needed
			IsApprovalNeeded        bool `json:"is_approval_needed"`
			IsClonedFromLastVersion bool `json:"is_cloned_from_last_version"`
		} `json:"templates"`
		Status     string    `json:"status"`
		VersionTag string    `json:"version_tag"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		UpdatedBy  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"updated_by"`
		VersionTagUser             string         `json:"version_tag_user"`
		PublishedLanguages         []string       `json:"published_languages"`
		ApparentPublishedLanguages []string       `json:"apparent_published_languages"`
		SystemApprovalInfo         map[string]any `json:"system_approval_info"` // Assuming it's an object
	} `json:"versions"`
}
