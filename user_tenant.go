package suprsend

// UserTenantUpsertBody is the request body for UpsertTenant.
type UserTenantUpsertBody struct {
	UserProperties     map[string]any   `json:"user_properties,omitempty"`
	RestrictedChannels []string         `json:"restricted_channels,omitempty"`
	Locale             string           `json:"$locale,omitempty"`
	Timezone           string           `json:"$timezone,omitempty"`
	// Inline channel identities
	Email       any `json:"$email,omitempty"`
	Sms         any `json:"$sms,omitempty"`
	Whatsapp    any `json:"$whatsapp,omitempty"`
	Androidpush any `json:"$androidpush,omitempty"`
	Iospush     any `json:"$iospush,omitempty"`
	Webpush     any `json:"$webpush,omitempty"`
	Slack       any `json:"$slack,omitempty"`
	MsTeams     any `json:"$ms_teams,omitempty"`
	// Operation-based identity management
	Append     map[string]any   `json:"$append,omitempty"`
	Remove     map[string]any   `json:"$remove,omitempty"`
	Unset      []string         `json:"$unset,omitempty"`
	Operations []map[string]any `json:"operations,omitempty"`
}

// UserIdentityItem represents a single subscriber identity (email, sms, push token, etc.).
type UserIdentityItem struct {
	Value              string         `json:"value,omitempty"`
	ValueJson          map[string]any `json:"value_json,omitempty"`
	IdProvider         string         `json:"id_provider,omitempty"`
	DeviceId           string         `json:"device_id,omitempty"`
	Status             string         `json:"status"`
	PermaStatus        string         `json:"perma_status,omitempty"`
	PermaStatusCode    string         `json:"perma_status_code,omitempty"`
	PermaStatusReason  string         `json:"perma_status_reason,omitempty"`
	PermaStatusTimeout string         `json:"perma_status_timeout,omitempty"`
}

// UserTenantItem represents a single user-tenant mapping.
type UserTenantItem struct {
	TenantId           string                        `json:"tenant_id"`
	TenantName         string                        `json:"tenant_name"`
	UserProperties     map[string]any                `json:"user_properties"`
	RestrictedChannels []string                      `json:"restricted_channels"`
	Locale             string                        `json:"locale"`
	Timezone           string                        `json:"timezone"`
	Identities         map[string][]UserIdentityItem `json:"identities,omitempty"`
}

// UserTenantsListMeta is pagination metadata for GetTenants.
type UserTenantsListMeta struct {
	Count   int     `json:"count"`
	Limit   int     `json:"limit"`
	HasPrev bool    `json:"has_prev"`
	Before  *string `json:"before"`
	HasNext bool    `json:"has_next"`
	After   *string `json:"after"`
}

// UserTenantsListResult is the inner result object for GetTenants.
type UserTenantsListResult struct {
	DistinctId  string              `json:"distinct_id"`
	Properties  map[string]any      `json:"properties"`
	CreatedAt   string              `json:"created_at,omitempty"`
	UpdatedAt   string              `json:"updated_at,omitempty"`
	Email       []UserIdentityItem  `json:"$email,omitempty"`
	Sms         []UserIdentityItem  `json:"$sms,omitempty"`
	Whatsapp    []UserIdentityItem  `json:"$whatsapp,omitempty"`
	Androidpush []UserIdentityItem  `json:"$androidpush,omitempty"`
	Iospush     []UserIdentityItem  `json:"$iospush,omitempty"`
	Webpush     []UserIdentityItem  `json:"$webpush,omitempty"`
	Inbox       []UserIdentityItem  `json:"$inbox,omitempty"`
	Slack       []UserIdentityItem  `json:"$slack,omitempty"`
	MsTeams     []UserIdentityItem  `json:"$ms_teams,omitempty"`
	Tenants     []UserTenantItem    `json:"tenants"`
}

// UserTenantsListResponse is the response for GetTenants.
type UserTenantsListResponse struct {
	Results *UserTenantsListResult `json:"results"`
	Meta    *UserTenantsListMeta   `json:"meta"`
}

// UserTenantDetailResponse is the response for GetTenantDetail and UpsertTenant.
type UserTenantDetailResponse struct {
	DistinctId  string              `json:"distinct_id"`
	Properties  map[string]any      `json:"properties"`
	CreatedAt   string              `json:"created_at,omitempty"`
	UpdatedAt   string              `json:"updated_at,omitempty"`
	Email       []UserIdentityItem  `json:"$email,omitempty"`
	Sms         []UserIdentityItem  `json:"$sms,omitempty"`
	Whatsapp    []UserIdentityItem  `json:"$whatsapp,omitempty"`
	Androidpush []UserIdentityItem  `json:"$androidpush,omitempty"`
	Iospush     []UserIdentityItem  `json:"$iospush,omitempty"`
	Webpush     []UserIdentityItem  `json:"$webpush,omitempty"`
	Inbox       []UserIdentityItem  `json:"$inbox,omitempty"`
	Slack       []UserIdentityItem  `json:"$slack,omitempty"`
	MsTeams     []UserIdentityItem  `json:"$ms_teams,omitempty"`
	Tenant      *UserTenantItem     `json:"tenant"`
}