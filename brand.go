package suprsend

type Brand struct {
	BrandId        *string                `json:"brand_id,omitempty"`
	BrandName      *string                `json:"brand_name,omitempty"`
	Logo           *string                `json:"logo,omitempty"`
	PrimaryColor   *string                `json:"primary_color,omitempty"`
	SecondaryColor *string                `json:"secondary_color,omitempty"`
	TertiaryColor  *string                `json:"tertiary_color,omitempty"`
	SocialLinks    *BrandSocialLinks      `json:"social_links,omitempty"`
	Properties     map[string]interface{} `json:"properties,omitempty"`
}

type BrandSocialLinks struct {
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

type BrandList struct {
	Meta    *ListApiMetaInfo `json:"meta"`
	Results []*Brand         `json:"results"`
}

type ListApiMetaInfo struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type BrandListOptions struct {
	Limit  int
	Offset int
}

func (b *BrandListOptions) cleanParams() {
	// limit must be 0 < x <= 1000
	if b.Limit <= 0 || b.Limit > 1000 {
		b.Limit = 20
	}
	if b.Offset < 0 {
		b.Offset = 0
	}
}
