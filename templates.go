package suprsend

import "time"

type Template struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	IsActive        bool   `json:"is_active"`
	DefaultLanguage struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"default_language"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"updated_by"`
	LastTriggeredAt        time.Time `json:"last_triggered_at"`
	IsAutoTranslateEnabled bool      `json:"is_auto_translate_enabled"`
	EnabledLanguages       []string  `json:"enabled_languages"`
	Channels               []struct {
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
		Versions          []any     `json:"versions"` // Change if you have a known type
	} `json:"channels"`
	Tags []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

type TemplateList struct {
	Meta    ListApiMetaInfo `json:"meta"`
	Results []Template      `json:"results"`
}

type TemplateListOptions struct {
	Limit  int
	Offset int
}

func (t *TemplateListOptions) cleanParams() {
	// limit must be 0 < x <= 1000
	if t.Limit <= 0 || t.Limit > 1000 {
		t.Limit = 20
	}
	if t.Offset < 0 {
		t.Offset = 0
	}
}
