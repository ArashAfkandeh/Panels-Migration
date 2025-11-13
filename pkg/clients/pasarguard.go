package clients

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"panels_user_manager/pkg/models"
	"panels_user_manager/pkg/utils"
)

// PasarGuardClient is the client to manage communication with the PasarGuard panel.
type PasarGuardClient struct {
	BaseURL    string
	Username   string
	Password   string
	HttpClient *http.Client
	Token      string // OAuth2 access token
}

// NewPasarGuardClient creates a new client for the PasarGuard panel.
func NewPasarGuardClient(baseURL, username, password string) *PasarGuardClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("Error creating cookie jar: %v", err))
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Jar:       jar,
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	return &PasarGuardClient{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Username:   username,
		Password:   password,
		HttpClient: client,
	}
}

// Login logs into the PasarGuard panel using OAuth2 token endpoint.
func (c *PasarGuardClient) Login() error {
	reqURL := fmt.Sprintf("%s/api/admin/token", c.BaseURL)
	data := fmt.Sprintf("grant_type=password&username=%s&password=%s", c.Username, c.Password)

	resp, err := c.HttpClient.Post(reqURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("error connecting to panel: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp models.PasarGuardTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("error parsing token response: %v", err)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("no access token received in response")
	}

	c.Token = tokenResp.AccessToken
	fmt.Println(" │ ✅ Authentication successful")
	return nil
}

// extractAllUUIDsFromProxySettings extracts ALL UUID-like identifiers from all protocols
func extractAllUUIDsFromProxySettings(proxySettings map[string]interface{}) []string {
	var uuids []string

	if proxySettings == nil {
		return uuids
	}

	protocolKeys := []string{"vmess", "vless", "trojan", "shadowsocks", "hysteria", "ss", "hy2"}
	fieldNames := []string{"id", "uuid", "password"}

	for _, protocol := range protocolKeys {
		if settings, ok := proxySettings[protocol].(map[string]interface{}); ok {
			for _, fieldName := range fieldNames {
				if value, ok := settings[fieldName].(string); ok && value != "" {
					normalized := strings.ToLower(strings.TrimSpace(value))
					alreadyExists := false
					for _, existing := range uuids {
						if existing == normalized {
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						uuids = append(uuids, normalized)
					}
				}
			}
		}
	}

	return uuids
}

// GetAllUsers fetches all users from the PasarGuard panel.
func (c *PasarGuardClient) GetAllUsers() ([]models.PasarGuardUser, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("not authenticated. Please login first")
	}

	if endpoints, err := c.discoverEndpoints(); err == nil && len(endpoints) > 0 {
		for _, ep := range endpoints {
			if users, err := c.tryGetUsers(ep.path, ep.method); err == nil {
				return users, nil
			}
		}
	}

	endpoints := []struct {
		path   string
		method string
	}{
		{"/api/users", "GET"},
		{"/api/users", "POST"},
		{"/api/admin/users", "GET"},
		{"/api/admin/users", "POST"},
		{"/api/admin/user", "GET"},
		{"/api/admin/user", "POST"},
		{"/api/v1/admin/users", "GET"},
		{"/api/v1/admin/users", "POST"},
		{"/api/v1/admin/user", "GET"},
		{"/api/v1/admin/user", "POST"},
		{"/api/v1/users", "GET"},
		{"/api/v1/users", "POST"},
		{"/api/user/list", "GET"},
		{"/api/user/list", "POST"},
		{"/api/admin/user/list", "GET"},
		{"/api/admin/user/list", "POST"},
		{"/api/v1/user/list", "GET"},
		{"/api/v1/user/list", "POST"},
		{"/api/v1/admin/user/list", "GET"},
		{"/api/v1/admin/user/list", "POST"},
	}

	var lastErr error
	for _, ep := range endpoints {
		users, err := c.tryGetUsers(ep.path, ep.method)
		if err == nil {
			return users, nil
		}
		lastErr = err
	}

	fmt.Println(" │ ⚠️  Trying to discover endpoints from API documentation...")
	if discoveredEndpoints, err := c.discoverEndpoints(); err == nil && len(discoveredEndpoints) > 0 {
		fmt.Printf(" │ ℹ️  Found %d user-related endpoints in API docs, trying them...\n", len(discoveredEndpoints))
		for _, ep := range discoveredEndpoints {
			if users, err := c.tryGetUsers(ep.path, ep.method); err == nil {
				return users, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to fetch users after trying multiple endpoints. Last error: %v", lastErr)
}

// GetAllGroups fetches available groups from the panel
func (c *PasarGuardClient) GetAllGroups() ([]models.PasarGuardGroup, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("not authenticated. Please login first")
	}

	endpoints := []string{
		"/api/groups",
		"/api/admin/groups",
		"/api/v1/groups",
		"/api/group",
		"/api/admin/group",
		"/api/v1/group",
	}

	var lastErr error
	for _, ep := range endpoints {
		reqURL := fmt.Sprintf("%s%s", c.BaseURL, ep)
		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			lastErr = fmt.Errorf("error creating request for %s: %v", reqURL, err)
			continue
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error requesting %s: %v", reqURL, err)
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("non-200 status %d from %s: %s", resp.StatusCode, reqURL, string(bodyBytes))
			continue
		}

		var groups []models.PasarGuardGroup
		if err := json.Unmarshal(bodyBytes, &groups); err == nil && len(groups) > 0 {
			return groups, nil
		}

		var apiResp models.APIResponse
		if err := json.Unmarshal(bodyBytes, &apiResp); err == nil && apiResp.Success {
			var raw []map[string]interface{}
			if err := json.Unmarshal(apiResp.Obj, &raw); err == nil {
				out := make([]models.PasarGuardGroup, 0, len(raw))
				for _, item := range raw {
					id := 0
					if v, ok := item["id"].(float64); ok {
						id = int(v)
					}
					name := ""
					if v, ok := item["name"].(string); ok {
						name = v
					} else if v, ok := item["title"].(string); ok {
						name = v
					}
					out = append(out, models.PasarGuardGroup{ID: id, Name: name})
				}
				if len(out) > 0 {
					return out, nil
				}
			}
		}

		var wrapper map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &wrapper); err == nil {
			if arr, ok := wrapper["groups"].([]interface{}); ok && len(arr) > 0 {
				out := make([]models.PasarGuardGroup, 0, len(arr))
				for _, a := range arr {
					if m, ok := a.(map[string]interface{}); ok {
						id := 0
						if v, ok := m["id"].(float64); ok {
							id = int(v)
						}
						name := ""
						if v, ok := m["name"].(string); ok {
							name = v
						} else if v, ok := m["title"].(string); ok {
							name = v
						}
						out = append(out, models.PasarGuardGroup{ID: id, Name: name})
					}
				}
				if len(out) > 0 {
					return out, nil
				}
			}
		}

		lastErr = fmt.Errorf("no groups found in response from %s", reqURL)
	}

	return nil, fmt.Errorf("failed to fetch groups: %v", lastErr)
}

// tryGetUsers attempts to get users from a specific endpoint
func (c *PasarGuardClient) tryGetUsers(path, method string) ([]models.PasarGuardUser, error) {
	reqURL := fmt.Sprintf("%s%s", c.BaseURL, path)

	var req *http.Request
	var err error

	if method == "POST" {
		req, err = http.NewRequest("POST", reqURL, bytes.NewBuffer([]byte("{}")))
	} else {
		req, err = http.NewRequest("GET", reqURL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request for %s %s: %v", method, reqURL, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching users from %s %s: %v", method, reqURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		var debugResp map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &debugResp); err == nil {
			if usersArray, ok := debugResp["users"].([]interface{}); ok && len(usersArray) > 0 {
				if firstUser, ok := usersArray[0].(map[string]interface{}); ok {
					utils.VerboseLog("First user raw JSON - used_traffic: %v, lifetime_used_traffic: %v",
						firstUser["used_traffic"], firstUser["lifetime_used_traffic"])
				}
			}
		}

		var usersResp models.PasarGuardUsersResponse
		if err := json.Unmarshal(bodyBytes, &usersResp); err == nil && len(usersResp.Users) > 0 {
			users := make([]models.PasarGuardUser, 0, len(usersResp.Users))
			for i, apiUser := range usersResp.Users {
				utils.VerboseLog("User %s - UsedTraffic: %d, LifetimeUsedTraffic: %d", apiUser.Username, apiUser.UsedTraffic, apiUser.LifetimeUsedTraffic)

				if apiUser.UsedTraffic == 0 && apiUser.LifetimeUsedTraffic == 0 {
					var debugResp map[string]interface{}
					if err := json.Unmarshal(bodyBytes, &debugResp); err == nil {
						if usersArray, ok := debugResp["users"].([]interface{}); ok && i < len(usersArray) {
							if userMap, ok := usersArray[i].(map[string]interface{}); ok {
								if usedTrafficVal, ok := userMap["used_traffic"]; ok {
									switch v := usedTrafficVal.(type) {
									case float64:
										apiUser.UsedTraffic = int64(v)
										utils.VerboseLog("Extracted used_traffic from raw JSON: %d", apiUser.UsedTraffic)
									case int64:
										apiUser.UsedTraffic = v
										utils.VerboseLog("Extracted used_traffic from raw JSON: %d", apiUser.UsedTraffic)
									}
								}
								if lifetimeVal, ok := userMap["lifetime_used_traffic"]; ok {
									switch v := lifetimeVal.(type) {
									case float64:
										apiUser.LifetimeUsedTraffic = int64(v)
										utils.VerboseLog("Extracted lifetime_used_traffic from raw JSON: %d", apiUser.LifetimeUsedTraffic)
									case int64:
										apiUser.LifetimeUsedTraffic = v
										utils.VerboseLog("Extracted lifetime_used_traffic from raw JSON: %d", apiUser.LifetimeUsedTraffic)
									}
								}
							}
						}
					}
				}

				var expiryTime int64
				if apiUser.Expire != "" {
					if t, err := time.Parse(time.RFC3339, apiUser.Expire); err == nil {
						expiryTime = t.Unix()
					}
				}

				protocol := ""
				uuid := ""
				port := 0
				if apiUser.ProxySettings != nil {
					if vmess, ok := apiUser.ProxySettings["vmess"].(map[string]interface{}); ok {
						protocol = "vmess"
						if id, ok := vmess["id"].(string); ok {
							uuid = id
						}
					} else if vless, ok := apiUser.ProxySettings["vless"].(map[string]interface{}); ok {
						protocol = "vless"
						if id, ok := vless["id"].(string); ok {
							uuid = id
						}
					} else if trojan, ok := apiUser.ProxySettings["trojan"].(map[string]interface{}); ok {
						protocol = "trojan"
						if password, ok := trojan["password"].(string); ok {
							uuid = password
						}
					} else if shadowsocks, ok := apiUser.ProxySettings["shadowsocks"].(map[string]interface{}); ok {
						protocol = "shadowsocks"
						if password, ok := shadowsocks["password"].(string); ok {
							uuid = password
						}
					}
				}

				email := ""
				if apiUser.Note != nil && *apiUser.Note != "" {
					note := *apiUser.Note
					if strings.Contains(note, "@") {
						email = note
					}
				}

				usedTraffic := apiUser.UsedTraffic
				if usedTraffic == 0 && apiUser.LifetimeUsedTraffic > 0 {
					usedTraffic = apiUser.LifetimeUsedTraffic
				}

				user := models.PasarGuardUser{
					ID:              apiUser.ID,
					Username:        apiUser.Username,
					Email:           email,
					UUID:            uuid,
					Enable:          apiUser.Status == "active",
					TotalGB:         apiUser.DataLimit,
					ExpiryTime:      expiryTime,
					LimitIP:         0,
					UsedTraffic:     usedTraffic,
					Protocol:        protocol,
					Port:            port,
					Remark:          "",
					SubscriptionURL: apiUser.SubscriptionURL,
					ProxySettings:   apiUser.ProxySettings,
					GroupIDs:        apiUser.GroupIDs,
				}

				if apiUser.Note != nil && *apiUser.Note != "" {
					user.Note = *apiUser.Note
					user.Remark = *apiUser.Note
				}

				if user.TotalGB > 0 {
					user.RemainingTraffic = user.TotalGB - user.UsedTraffic
					if user.RemainingTraffic < 0 {
						user.RemainingTraffic = 0
					}
				} else {
					user.RemainingTraffic = -1
				}

				users = append(users, user)
			}
			return users, nil
		}

		var apiResp models.APIResponse
		if err := json.Unmarshal(bodyBytes, &apiResp); err == nil && apiResp.Success {
			var users []models.PasarGuardUser
			if err := json.Unmarshal(apiResp.Obj, &users); err == nil {
				for i := range users {
					if users[i].TotalGB > 0 {
						users[i].RemainingTraffic = users[i].TotalGB - users[i].UsedTraffic
						if users[i].RemainingTraffic < 0 {
							users[i].RemainingTraffic = 0
						}
					} else {
						users[i].RemainingTraffic = -1
					}
				}
				return users, nil
			}
		}

		var users []models.PasarGuardUser
		if err := json.Unmarshal(bodyBytes, &users); err == nil && len(users) > 0 {
			for i := range users {
				if users[i].TotalGB > 0 {
					users[i].RemainingTraffic = users[i].TotalGB - users[i].UsedTraffic
					if users[i].RemainingTraffic < 0 {
						users[i].RemainingTraffic = 0
					}
				} else {
					users[i].RemainingTraffic = -1
				}
			}
			return users, nil
		}

		return nil, fmt.Errorf("error parsing response from %s %s: invalid format. Response: %s", method, reqURL, string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("endpoint not found: %s %s", method, reqURL)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed. Token may have expired. Please login again")
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, fmt.Errorf("method %s not allowed for %s", method, reqURL)
	}

	return nil, fmt.Errorf("server returned status %d from %s %s. Response: %s", resp.StatusCode, method, reqURL, string(bodyBytes))
}

// discoverEndpoints tries to discover API endpoints from OpenAPI/Swagger docs
func (c *PasarGuardClient) discoverEndpoints() ([]struct{ path, method string }, error) {
	docsEndpoints := []string{
		"/openapi.json",
		"/api/openapi.json",
		"/docs",
		"/api/docs",
		"/swagger.json",
		"/api/swagger.json",
	}

	var endpoints []struct{ path, method string }

	for _, docPath := range docsEndpoints {
		reqURL := fmt.Sprintf("%s%s", c.BaseURL, docPath)
		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		resp, err := c.HttpClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var schema models.OpenAPISchema
		if err := json.Unmarshal(bodyBytes, &schema); err == nil && schema.Paths != nil {
			for path, methods := range schema.Paths {
				if strings.Contains(strings.ToLower(path), "user") {
					for method := range methods {
						methodUpper := strings.ToUpper(method)
						if methodUpper == "GET" || methodUpper == "POST" {
							endpoints = append(endpoints, struct{ path, method string }{path, methodUpper})
						}
					}
				}
			}
			if len(endpoints) > 0 {
				return endpoints, nil
			}
		}
	}

	return endpoints, nil
}
