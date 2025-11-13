package models

import "encoding/json"

// --- 3X-UI MODELS ---

// LoginPayload is the JSON structure for the login request.
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// APIResponse represents a generic response from the API.
type APIResponse struct {
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Obj     json.RawMessage `json:"obj"`
}

// Inbound represents the data structure for an inbound received from the API.
type Inbound struct {
	ID             int    `json:"id"`
	Remark         string `json:"remark"`
	Protocol       string `json:"protocol"`
	Port           int    `json:"port"`
	Settings       string `json:"settings"`       // This field is a JSON string.
	StreamSettings string `json:"streamSettings"` // Transmission settings.
	Sniffing       string `json:"sniffing"`       // Sniffing settings (External Proxy).
	Enable         bool   `json:"enable"`
	Tag            string `json:"tag"`
	ExpiryTime     int64  `json:"expiryTime"`
	Total          int64  `json:"total"`
	Listen         string `json:"listen"`
}

// InboundSettings is the inner structure of the 'settings' field.
type InboundSettings struct {
	Clients []ClientSetting `json:"clients"`
}

// ClientSetting represents a client within the settings section.
type ClientSetting struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Enable     bool   `json:"enable"`
	TotalGB    int64  `json:"totalGB"`    // Total traffic in bytes.
	ExpiryTime int64  `json:"expiryTime"` // As a timestamp.
	LimitIP    int    `json:"limitIp"`
	Flow       string `json:"flow"`
	SubID      string `json:"subId"`
	TgID       string `json:"tgId"`
	Reset      int    `json:"reset"`
}

// ClientTraffic represents user traffic data (upload and download).
type ClientTraffic struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

// ClientDetails is the final client information structure for nesting within InboundData.
type ClientDetails struct {
	ClientEmail         string  `json:"client_email"`
	ClientID            string  `json:"client_id"`
	ClientEnable        bool    `json:"client_enable"`
	ClientLimitIP       int     `json:"-"`
	ClientTotalGB       int64   `json:"client_total_gb"` // in bytes
	ClientExpiryTime    int64   `json:"client_expiry_time"`
	ClientFlow          string  `json:"-"`
	ClientSubID         string  `json:"client_sub_id"`
	ClientTgID          string  `json:"-"`
	ClientReset         int     `json:"-"`
	TrafficUsed         int64   `json:"traffic_used"`      // in bytes
	TrafficRemaining    int64   `json:"traffic_remaining"` // in bytes (-1 for unlimited)
	TrafficUsagePercent float64 `json:"-"`
}

// ThreeXUIUser represents a simplified user structure for 3X-UI export (similar to PasarGuardUser)
type ThreeXUIUser struct {
	ID               int    `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	UUID             string `json:"uuid"`
	Enable           bool   `json:"enable"`
	TotalGB          int64  `json:"totalGB"`    // Total traffic in bytes
	ExpiryTime       int64  `json:"expiryTime"` // As a timestamp
	LimitIP          int    `json:"limitIp"`
	UsedTraffic      int64  `json:"usedTraffic"`      // Used traffic in bytes
	RemainingTraffic int64  `json:"remainingTraffic"` // Remaining traffic in bytes
	Protocol         string `json:"protocol"`
	Port             int    `json:"port"`
	Remark           string `json:"remark"`
}

// ThreeXUIUsersExportFile is the structure for exporting 3X-UI users (similar to PasarGuardUsersExportFile)
type ThreeXUIUsersExportFile struct {
	ExportDate string         `json:"export_date"`
	PanelType  string         `json:"panel_type"`
	TotalUsers int            `json:"total_users"`
	Users      []ThreeXUIUser `json:"users"`
}

// InboundData is the final structure for an inbound in the output JSON file.
type InboundData struct {
	ID               int             `json:"id"`
	Remark           string          `json:"remark"`
	Protocol         string          `json:"protocol"`
	Port             int             `json:"port"`
	Enable           bool            `json:"enable"`
	Tag              string          `json:"tag"`
	Listen           string          `json:"listen"`
	InboundExpiry    int64           `json:"inbound_expiry_time"`
	InboundTotalGB   int64           `json:"inbound_total_gb_bytes"`
	Transmission     string          `json:"transmission"`
	ExternalProxy    string          `json:"external_proxy"`
	OriginalSettings string          `json:"original_settings"` // Save raw settings for WireGuard
	Clients          []ClientDetails `json:"clients"`
}

// AddInboundPayload is the structure for the JSON body when adding an inbound.
type AddInboundPayload struct {
	Remark         string `json:"remark"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	Settings       string `json:"settings"`
	StreamSettings string `json:"streamSettings"`
	Sniffing       string `json:"sniffing"`
	Enable         bool   `json:"enable"`
	Listen         string `json:"listen"`
	Total          int64  `json:"total"`
	ExpiryTime     int64  `json:"expiryTime"`
}

// OutputFile is the structure of the output JSON file.
type OutputFile struct {
	ExportDate    string        `json:"export_date"`
	TotalInbounds int           `json:"total_inbounds"`
	TotalUsers    int           `json:"total_users"`
	Inbounds      []InboundData `json:"inbounds"`
}

// --- PASARGUARD MODELS ---

// PasarGuardTokenResponse represents the token response from PasarGuard API.
type PasarGuardTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// PasarGuardUser represents a user in PasarGuard panel.
type PasarGuardUser struct {
	ID               int                    `json:"id"`
	Username         string                 `json:"username"`
	Email            string                 `json:"email"`
	UUID             string                 `json:"uuid"`
	Enable           bool                   `json:"enable"`
	TotalGB          int64                  `json:"totalGB"`    // Total traffic in bytes
	ExpiryTime       int64                  `json:"expiryTime"` // As a timestamp
	LimitIP          int                    `json:"limitIp"`
	UsedTraffic      int64                  `json:"usedTraffic"`      // Used traffic in bytes
	RemainingTraffic int64                  `json:"remainingTraffic"` // Remaining traffic in bytes
	Protocol         string                 `json:"protocol"`
	Port             int                    `json:"port"`
	Remark           string                 `json:"remark"`
	SubscriptionURL  string                 `json:"subscription_url"` // Subscription URL
	Note             string                 `json:"note"`             // Note field (can contain email or other info)
	ProxySettings    map[string]interface{} `json:"proxy_settings"`   // All protocol UUIDs
	GroupIDs         []int                  `json:"group_ids"`
}

// PasarGuardUserListResponse represents the response from getting users list.
type PasarGuardUserListResponse struct {
	Success bool             `json:"success"`
	Msg     string           `json:"msg"`
	Obj     []PasarGuardUser `json:"obj"`
}

// PasarGuardUsersResponse represents the actual API response structure from /api/users
type PasarGuardUsersResponse struct {
	Users []PasarGuardUserAPI `json:"users"`
	Total int                 `json:"total"`
}

// PasarGuardUserAPI represents the actual user structure from PasarGuard API
type PasarGuardUserAPI struct {
	ID                     int                    `json:"id"`
	Username               string                 `json:"username"`
	Status                 string                 `json:"status"`
	Expire                 string                 `json:"expire"`     // ISO 8601 format
	DataLimit              int64                  `json:"data_limit"` // in bytes
	UsedTraffic            int64                  `json:"used_traffic"`
	LifetimeUsedTraffic    int64                  `json:"lifetime_used_traffic"`
	DataLimitResetStrategy string                 `json:"data_limit_reset_strategy"`
	Note                   *string                `json:"note"`
	OnHoldExpireDuration   *int                   `json:"on_hold_expire_duration"`
	OnHoldTimeout          *int                   `json:"on_hold_timeout"`
	GroupIDs               []int                  `json:"group_ids"`
	AutoDeleteInDays       *int                   `json:"auto_delete_in_days"`
	NextPlan               *string                `json:"next_plan"`
	CreatedAt              string                 `json:"created_at"`
	EditAt                 string                 `json:"edit_at"`
	OnlineAt               string                 `json:"online_at"`
	SubscriptionURL        string                 `json:"subscription_url"`
	ProxySettings          map[string]interface{} `json:"proxy_settings"`
	Admin                  struct {
		Username string `json:"username"`
	} `json:"admin"`
}

// PasarGuardAddUserPayload is the structure for adding a user.
type PasarGuardAddUserPayload struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	UUID       string `json:"uuid"`
	Enable     bool   `json:"enable"`
	TotalGB    int64  `json:"totalGB"`
	ExpiryTime int64  `json:"expiryTime"`
	LimitIP    int    `json:"limitIp"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	Remark     string `json:"remark"`
}

// PasarGuardUsersExportFile is the structure for exporting PasarGuard users.
type PasarGuardUsersExportFile struct {
	ExportDate string           `json:"export_date"`
	PanelType  string           `json:"panel_type"`
	TotalUsers int              `json:"total_users"`
	Users      []PasarGuardUser `json:"users"`
}

// PasarGuardGroup represents a group object in PasarGuard panel
type PasarGuardGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// OpenAPISchema represents the OpenAPI schema structure
type OpenAPISchema struct {
	Paths map[string]map[string]interface{} `json:"paths"`
}
