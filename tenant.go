package suprsend

type Tenant struct {
	TenantId   *string `json:"tenant_id,omitempty"`
	TenantName *string `json:"tenant_name,omitempty"`
	Logo       *string `json:"logo,omitempty"`
	Timezone   *string `json:"timezone,omitempty"`
	//
	BlockedChannels        []string `json:"blocked_channels"`
	EmbeddedPreferenceUrl  *string  `json:"embedded_preference_url,omitempty"`
	HostedPreferenceDomain *string  `json:"hosted_preference_domain,omitempty"`
	//
	PrimaryColor   *string                `json:"primary_color,omitempty"`
	SecondaryColor *string                `json:"secondary_color,omitempty"`
	TertiaryColor  *string                `json:"tertiary_color,omitempty"`
	SocialLinks    *TenantSocialLinks     `json:"social_links,omitempty"`
	Properties     map[string]interface{} `json:"properties,omitempty"`
}

type TenantSocialLinks struct {
	Website   *string `json:"website,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Linkedin  *string `json:"linkedin,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Medium    *string `json:"medium,omitempty"`
	Discord   *string `json:"discord,omitempty"`
	Telegram  *string `json:"telegram,omitempty"`
	Youtube   *string `json:"youtube,omitempty"`
}

type TenantList struct {
	Meta    *ListApiMetaInfo `json:"meta"`
	Results []*Tenant        `json:"results"`
}

type TenantListOptions struct {
	Limit  int
	Offset int
}

func (t *TenantListOptions) cleanParams() {
	// limit must be 0 < x <= 1000
	if t.Limit <= 0 || t.Limit > 1000 {
		t.Limit = 20
	}
	if t.Offset < 0 {
		t.Offset = 0
	}
}
