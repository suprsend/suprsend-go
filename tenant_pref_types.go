package suprsend

// Tenant Category Digest Schedule in Response
type TenantCategoryDigestScheduleOut struct {
	OriginalOptions []TenantCategoryDigestScheduleOptionOut `json:"original_options"`
	Options         []TenantCategoryDigestScheduleOptionOut `json:"options"`
}

type TenantCategoryDigestScheduleOptionOut struct {
	Id        string                                    `json:"id"`
	Label     string                                    `json:"label"`
	Frequency string                                    `json:"frequency"`
	Interval  int                                       `json:"interval,omitempty"`
	IsDefault bool                                      `json:"is_default"`
	Time      *TenantCategoryDigestScheduleTimeOut      `json:"time,omitempty"`
	Weekdays  *TenantCategoryDigestScheduleWeekdaysOut  `json:"weekdays,omitempty"`
	Monthdays *TenantCategoryDigestScheduleMonthdaysOut `json:"monthdays,omitempty"`
}

type TenantCategoryDigestScheduleTimeOut struct {
	DefaultValue string `json:"default_value"`
	EditPolicy   string `json:"edit_policy"`
}

type TenantCategoryDigestScheduleWeekdaysOut struct {
	DefaultValue []string `json:"default_value"`
	EditPolicy   string   `json:"edit_policy"`
}

type DigestScheduleMonthdayComponent struct {
	Pos int    `json:"pos"`
	Day string `json:"day,omitempty"`
}

type TenantCategoryDigestScheduleMonthdaysOut struct {
	DefaultValue []DigestScheduleMonthdayComponent `json:"default_value"`
	EditPolicy   string                            `json:"edit_policy"`
}

type PreferenceCategoryPropertyOut struct {
	Key               string                                `json:"key"`
	Label             string                                `json:"label"`
	ValueType         string                                `json:"value_type"`
	DefaultValue      any                                   `json:"default_value"`
	Choices           []PreferenceCategoryPropertyChoiceOut `json:"choices,omitempty"`
	DynamicChoicesKey string                                `json:"dynamic_choices_key,omitempty"`
	//
	EditPolicy string `json:"edit_policy"`
	IsOptional bool   `json:"is_optional"`
	//
	IsOverridden bool `json:"is_overridden"`
	Value        any  `json:"value"`
}

type PreferenceCategoryPropertyChoiceOut struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ------------

// Tenant Category Digest Schedule in Input
type TenantCategoryDigestScheduleIn struct {
	Options []*TenantCategoryDigestScheduleOptionIn `json:"options"`
}

// Tenant Category Property in Input
type PreferenceCategoryPropertyIn struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type TenantCategoryDigestScheduleOptionIn struct {
	Id string `json:"id"`
	//
	Time      *TenantCategoryDigestScheduleOptionFieldConfig[string]                            `json:"time,omitempty"`
	Weekdays  *TenantCategoryDigestScheduleOptionFieldConfig[[]string]                          `json:"weekdays,omitempty"`
	Monthdays *TenantCategoryDigestScheduleOptionFieldConfig[[]DigestScheduleMonthdayComponent] `json:"monthdays,omitempty"`
	//
	Dtstart *TenantCategoryDigestScheduleOptionFieldConfig[string] `json:"dtstart,omitempty"`
	//
	IsDefault bool `json:"is_default"`
}

type TenantCategoryDigestScheduleOptionFieldConfig[T any] struct {
	DefaultValue T `json:"default_value" validate:"required"`
}
