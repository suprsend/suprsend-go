package suprsend

// UserTenantUpsertBody is the request body for UpsertTenant.
type UserTenantUpsertBody struct {
	UserProperties     map[string]any `json:"user_properties,omitempty"`
	RestrictedChannels []string       `json:"restricted_channels,omitempty"`
	Locale             string         `json:"$locale,omitempty"`
	Timezone           string         `json:"$timezone,omitempty"`
}

// UserTenantItem represents a single user-tenant mapping.
type UserTenantItem struct {
	TenantId           string         `json:"tenant_id"`
	TenantName         string         `json:"tenant_name"`
	UserProperties     map[string]any `json:"user_properties"`
	RestrictedChannels []string       `json:"restricted_channels"`
	Locale             string         `json:"locale"`
	Timezone           string         `json:"timezone"`
}

// UserTenantsListResponse is the response for GetTenants.
type UserTenantsListResponse struct {
	Tenants []UserTenantItem `json:"tenants"`
}

// UserTenantDetailResponse is the response for UpsertTenant.
type UserTenantDetailResponse struct {
	Tenant *UserTenantItem `json:"tenant"`
}