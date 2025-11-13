package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"panels_user_manager/pkg/models"
	"panels_user_manager/pkg/utils"
)

// AddUser creates a new user on the PasarGuard panel.
func (c *PasarGuardClient) AddUser(user models.PasarGuardUser) error {
	if c.Token == "" {
		return fmt.Errorf("not authenticated. Please login first")
	}

	endpoints := []string{
		"/api/user",
		"/api/users",
		"/api/v1/user",
		"/api/v1/users",
		"/api/admin/users",
	}

	expireStr := ""
	if user.ExpiryTime > 0 {
		expireStr = time.Unix(user.ExpiryTime, 0).Format(time.RFC3339)
	}

	proxySettings := make(map[string]interface{})
	switch user.Protocol {
	case "vmess":
		proxySettings["vmess"] = map[string]interface{}{
			"id": user.UUID,
		}
	case "vless":
		proxySettings["vless"] = map[string]interface{}{
			"id":   user.UUID,
			"flow": "",
		}
	case "trojan":
		proxySettings["trojan"] = map[string]interface{}{
			"password": user.UUID,
		}
	case "shadowsocks":
		proxySettings["shadowsocks"] = map[string]interface{}{
			"password": user.UUID,
			"method":   "chacha20-ietf-poly1305",
		}
	}

	status := "active"
	if !user.Enable {
		status = "disabled"
	}

	payload := map[string]interface{}{
		"username":       user.Username,
		"proxy_settings": proxySettings,
		"status":         status,
	}

	if user.TotalGB > 0 {
		payload["data_limit"] = user.TotalGB
	} else {
		payload["data_limit"] = 0
	}

	if expireStr != "" {
		payload["expire"] = expireStr
	}

	if user.UsedTraffic >= 0 {
		payload["used_traffic"] = user.UsedTraffic
		payload["lifetime_used_traffic"] = user.UsedTraffic
	}

	if user.Note != "" {
		payload["note"] = user.Note
	}
	if user.LimitIP > 0 {
		payload["limit_ip"] = user.LimitIP
	}
	payload["group_ids"] = user.GroupIDs

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling add user payload: %v", err)
	}

	utils.VerboseLog("AddUser payload for %s: %s", user.Username, string(payloadBytes))

	var lastErr error
	for _, endpoint := range endpoints {
		reqURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
		req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			lastErr = fmt.Errorf("error creating request for %s: %v", reqURL, err)
			continue
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request to %s: %v", reqURL, err)
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		utils.VerboseLog("AddUser response from %s: status=%d, body=%s", reqURL, resp.StatusCode, string(bodyBytes))

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var userResp map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &userResp); err == nil {
				if username, ok := userResp["username"].(string); ok && username != "" {
					utils.VerboseLog("AddUser successful (user object) from %s, username=%s", reqURL, username)
					return nil
				}
			}

			var apiResp models.APIResponse
			if err := json.Unmarshal(bodyBytes, &apiResp); err == nil {
				utils.VerboseLog("AddUser parsed as APIResponse, Success=%v, Msg=%s", apiResp.Success, apiResp.Msg)
				if apiResp.Success {
					utils.VerboseLog("AddUser successful (APIResponse) from %s", reqURL)
					return nil
				}
			}

			utils.VerboseLog("AddUser successful (status %d) from %s, assuming success", resp.StatusCode, reqURL)

			if user.UsedTraffic > 0 {
				if err := c.setUserTraffic(user.Username, user.UsedTraffic); err != nil {
					utils.VerboseLog("Failed to set traffic for %s: %v (this is non-critical)", user.Username, err)
				}
			}

			return nil
		}

		if resp.StatusCode == http.StatusNotFound {
			lastErr = fmt.Errorf("endpoint not found: %s", reqURL)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed. Token may have expired. Please login again")
		}

		if resp.StatusCode == http.StatusConflict {
			utils.VerboseLog("AddUser got 409 Conflict from %s, user already exists", reqURL)
			return fmt.Errorf("user already exists")
		}

		responseBody := strings.ToLower(string(bodyBytes))
		if strings.Contains(responseBody, "already exists") || strings.Contains(responseBody, "user already exists") {
			utils.VerboseLog("AddUser detected 'already exists' in response from %s", reqURL)
			return fmt.Errorf("user already exists")
		}

		if resp.StatusCode == http.StatusMethodNotAllowed {
			lastErr = fmt.Errorf("method not allowed: %s", reqURL)
			continue
		}

		lastErr = fmt.Errorf("server returned status %d from %s. Response: %s", resp.StatusCode, reqURL, string(bodyBytes))
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusConflict {
			break
		}
	}

	return fmt.Errorf("failed to add user after trying multiple endpoints. Last error: %v", lastErr)
}

// UpdateUser updates an existing user on the PasarGuard panel.
func (c *PasarGuardClient) UpdateUser(user models.PasarGuardUser) error {
	return c.UpdateUserByIdentifier(user.Username, user)
}

// UpdateUserByIdentifier updates an existing user using a provided identifier.
func (c *PasarGuardClient) UpdateUserByIdentifier(identifier string, user models.PasarGuardUser) error {
	if c.Token == "" {
		return fmt.Errorf("not authenticated. Please login first")
	}

	endpoints := []struct {
		path   string
		method string
	}{
		{fmt.Sprintf("/api/user/%s", identifier), "PUT"},
		{fmt.Sprintf("/api/user/%s", identifier), "PATCH"},
		{fmt.Sprintf("/api/users/%s", identifier), "PUT"},
		{fmt.Sprintf("/api/users/%s", identifier), "PATCH"},
		{fmt.Sprintf("/api/v1/user/%s", identifier), "PUT"},
		{fmt.Sprintf("/api/v1/user/%s", identifier), "PATCH"},
	}

	expireStr := ""
	if user.ExpiryTime > 0 {
		expireStr = time.Unix(user.ExpiryTime, 0).Format(time.RFC3339)
	}

	proxySettings := make(map[string]interface{})
	switch user.Protocol {
	case "vmess":
		proxySettings["vmess"] = map[string]interface{}{"id": user.UUID}
	case "vless":
		proxySettings["vless"] = map[string]interface{}{"id": user.UUID, "flow": ""}
	case "trojan":
		proxySettings["trojan"] = map[string]interface{}{"password": user.UUID}
	case "shadowsocks":
		proxySettings["shadowsocks"] = map[string]interface{}{"password": user.UUID, "method": "chacha20-ietf-poly1305"}
	}

	status := "active"
	if !user.Enable {
		status = "disabled"
	}

	payload := map[string]interface{}{
		"username":       user.Username,
		"proxy_settings": proxySettings,
		"status":         status,
	}
	if user.TotalGB > 0 {
		payload["data_limit"] = user.TotalGB
	} else {
		payload["data_limit"] = 0
	}
	if expireStr != "" {
		payload["expire"] = expireStr
	}
	if user.UsedTraffic >= 0 {
		payload["used_traffic"] = user.UsedTraffic
		payload["lifetime_used_traffic"] = user.UsedTraffic
	}
	if user.Note != "" {
		payload["note"] = user.Note
	}
	if user.LimitIP > 0 {
		payload["limit_ip"] = user.LimitIP
	}
	if len(user.GroupIDs) > 0 {
		payload["group_ids"] = user.GroupIDs
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling update user payload: %v", err)
	}

	utils.VerboseLog("UpdateUserByIdentifier payload for %s (ident=%s): %s", user.Username, identifier, string(payloadBytes))

	var lastErr error
	for _, ep := range endpoints {
		reqURL := fmt.Sprintf("%s%s", c.BaseURL, ep.path)
		req, err := http.NewRequest(ep.method, reqURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			lastErr = fmt.Errorf("error creating request for %s %s: %v", ep.method, reqURL, err)
			continue
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request to %s %s: %v", ep.method, reqURL, err)
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		utils.VerboseLog("UpdateUserByIdentifier response from %s %s: status=%d, body=%s", ep.method, reqURL, resp.StatusCode, string(bodyBytes))

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var userResp map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &userResp); err == nil {
				if username, ok := userResp["username"].(string); ok && username != "" {
					utils.VerboseLog("UpdateUserByIdentifier successful for %s (got user object, username=%s)", user.Username, username)
					return nil
				}
			}
			var apiResp models.APIResponse
			if err := json.Unmarshal(bodyBytes, &apiResp); err == nil {
				if apiResp.Success {
					utils.VerboseLog("UpdateUserByIdentifier successful for %s (APIResponse)", user.Username)
					return nil
				}
				lastErr = fmt.Errorf("API error from %s %s: %s", ep.method, reqURL, apiResp.Msg)
				continue
			}

			if user.UsedTraffic > 0 {
				if err := c.setUserTraffic(user.Username, user.UsedTraffic); err != nil {
					utils.VerboseLog("Failed to set traffic for %s: %v (non-critical)", user.Username, err)
				}
			}
			return nil
		}
		if resp.StatusCode == http.StatusNotFound {
			lastErr = fmt.Errorf("endpoint not found: %s %s", ep.method, reqURL)
			continue
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed. Token may have expired. Please login again")
		}
		if resp.StatusCode == http.StatusMethodNotAllowed {
			lastErr = fmt.Errorf("method %s not allowed for %s", ep.method, reqURL)
			continue
		}
		lastErr = fmt.Errorf("server returned status %d from %s %s. Response: %s", resp.StatusCode, ep.method, reqURL, string(bodyBytes))
	}

	return fmt.Errorf("failed to update user after trying multiple endpoints and methods. Last error: %v", lastErr)
}

// setUserTraffic attempts to set the used traffic for a user via a separate API call
func (c *PasarGuardClient) setUserTraffic(username string, usedTraffic int64) error {
	if c.Token == "" {
		return fmt.Errorf("not authenticated. Please login first")
	}

	endpoints := []struct {
		path   string
		method string
	}{
		{fmt.Sprintf("/api/user/%s/traffic", username), "PUT"},
		{fmt.Sprintf("/api/user/%s/traffic", username), "POST"},
		{fmt.Sprintf("/api/user/%s/set_traffic", username), "PUT"},
		{fmt.Sprintf("/api/user/%s/set_traffic", username), "POST"},
		{fmt.Sprintf("/api/users/%s/traffic", username), "PUT"},
		{fmt.Sprintf("/api/users/%s/traffic", username), "POST"},
		{fmt.Sprintf("/api/admin/user/%s/traffic", username), "PUT"},
		{fmt.Sprintf("/api/admin/user/%s/traffic", username), "POST"},
	}

	payload := map[string]interface{}{
		"used_traffic":          usedTraffic,
		"lifetime_used_traffic": usedTraffic,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling traffic payload: %v", err)
	}

	var lastErr error
	for _, ep := range endpoints {
		reqURL := fmt.Sprintf("%s%s", c.BaseURL, ep.path)

		req, err := http.NewRequest(ep.method, reqURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			lastErr = fmt.Errorf("error creating request for %s %s: %v", ep.method, reqURL, err)
			continue
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request to %s %s: %v", ep.method, reqURL, err)
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			utils.VerboseLog("setUserTraffic successful for %s via %s %s", username, ep.method, reqURL)
			return nil
		}

		if resp.StatusCode == http.StatusNotFound {
			lastErr = fmt.Errorf("endpoint not found: %s %s", ep.method, reqURL)
			continue
		}

		lastErr = fmt.Errorf("server returned status %d from %s %s. Response: %s", resp.StatusCode, ep.method, reqURL, string(bodyBytes))
	}

	return fmt.Errorf("failed to set traffic after trying multiple endpoints. Last error: %v", lastErr)
}

// ClearUserGroups removes all groups from a user using the bulk remove endpoint
func (c *PasarGuardClient) ClearUserGroups(username string) error {
	if c.Token == "" {
		return fmt.Errorf("not authenticated. Please login first")
	}

	allGroups, err := c.GetAllGroups()
	if err != nil {
		utils.VerboseLog("Could not fetch groups list for clearing: %v", err)
		return fmt.Errorf("failed to fetch groups for clearing: %v", err)
	}

	if len(allGroups) == 0 {
		utils.VerboseLog("No groups exist on panel, nothing to clear")
		return nil
	}

	groupIDs := make([]int, len(allGroups))
	for i, g := range allGroups {
		groupIDs[i] = g.ID
	}

	utils.VerboseLog("clearUserGroups will remove these group IDs %v from user %s", groupIDs, username)

	reqURL := fmt.Sprintf("%s/api/user/%s", c.BaseURL, username)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error creating GET request for user %s: %v", username, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching user %s: %v", username, err)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch user %s: status %d", username, resp.StatusCode)
	}

	var userResp models.PasarGuardUser
	if err := json.Unmarshal(bodyBytes, &userResp); err != nil {
		return fmt.Errorf("error parsing user response: %v", err)
	}

	utils.VerboseLog("Got user ID %d for username %s", userResp.ID, username)

	bulkPayload := map[string]interface{}{
		"group_ids": groupIDs,
		"users":     []int{userResp.ID},
	}

	bulkPayloadBytes, err := json.Marshal(bulkPayload)
	if err != nil {
		return fmt.Errorf("error marshalling bulk remove payload: %v", err)
	}

	utils.VerboseLog("clearUserGroups bulk remove payload: %s", string(bulkPayloadBytes))

	bulkURL := fmt.Sprintf("%s/api/groups/bulk/remove", c.BaseURL)
	bulkReq, err := http.NewRequest("POST", bulkURL, bytes.NewBuffer(bulkPayloadBytes))
	if err != nil {
		return fmt.Errorf("error creating bulk remove request: %v", err)
	}
	bulkReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	bulkReq.Header.Set("Content-Type", "application/json")

	bulkResp, err := c.HttpClient.Do(bulkReq)
	if err != nil {
		return fmt.Errorf("error calling bulk remove endpoint: %v", err)
	}
	bulkRespBytes, _ := io.ReadAll(bulkResp.Body)
	bulkResp.Body.Close()

	if bulkResp.StatusCode == http.StatusOK {
		utils.VerboseLog("clearUserGroups successful for user %s via bulk remove endpoint", username)
		utils.VerboseLog("clearUserGroups bulk response: %s", string(bulkRespBytes))
		return nil
	}

	return fmt.Errorf("bulk remove groups failed for user %s: status %d, response: %s", username, bulkResp.StatusCode, string(bulkRespBytes))
}
