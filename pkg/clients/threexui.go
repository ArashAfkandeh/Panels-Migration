package clients

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"panels_user_manager/pkg/models"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// ThreeXUIClient is the client to manage communication with the 3x-ui panel.
type ThreeXUIClient struct {
	BaseURL    string
	Username   string
	Password   string
	HttpClient *http.Client
}

// NewThreeXUIClient creates a new client for the 3x-ui panel.
func NewThreeXUIClient(baseURL, username, password string) *ThreeXUIClient {
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
	return &ThreeXUIClient{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Username:   username,
		Password:   password,
		HttpClient: client,
	}
}

// Login logs into the panel.
func (c *ThreeXUIClient) Login() error {
	url := fmt.Sprintf("%s/login", c.BaseURL)
	payload := models.LoginPayload{Username: c.Username, Password: c.Password}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error building login JSON: %v", err)
	}
	resp, err := c.HttpClient.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error connecting to the panel: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("error parsing login response: %v", err)
	}
	if !apiResp.Success {
		return fmt.Errorf("login failed: %s", apiResp.Msg)
	}
	fmt.Println(" │ ✅ Authentication successful")
	return nil
}

// GetAllInbounds fetches all inbounds from the panel.
func (c *ThreeXUIClient) GetAllInbounds() ([]models.Inbound, error) {
	url := fmt.Sprintf("%s/panel/api/inbounds/list", c.BaseURL)
	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching inbound list: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error parsing inbound list: %v", err)
	}
	if !apiResp.Success {
		return nil, fmt.Errorf("API error getting inbounds: %s", apiResp.Msg)
	}
	var inbounds []models.Inbound
	if err := json.Unmarshal(apiResp.Obj, &inbounds); err != nil {
		return nil, fmt.Errorf("error parsing inbound data: %v", err)
	}
	return inbounds, nil
}

// GetClientTraffic fetches traffic stats for a specific client by email.
func (c *ThreeXUIClient) GetClientTraffic(email string) (*models.ClientTraffic, error) {
	url := fmt.Sprintf("%s/panel/api/inbounds/getClientTraffics/%s", c.BaseURL, url.PathEscape(email))
	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if !apiResp.Success {
		return nil, fmt.Errorf("api error: %s", apiResp.Msg)
	}
	if string(apiResp.Obj) == "null" || string(apiResp.Obj) == "" {
		return &models.ClientTraffic{Up: 0, Down: 0}, nil
	}
	var traffic models.ClientTraffic
	if err := json.Unmarshal(apiResp.Obj, &traffic); err != nil {
		return nil, err
	}
	return &traffic, nil
}

// AddInbound creates a new inbound on the panel using data from a file.
func (c *ThreeXUIClient) AddInbound(inboundData models.InboundData) error {
	var finalSettings string
	// Sanitize streamSettings and sniffing to prevent "malformed JSON" errors.
	finalStreamSettings := inboundData.Transmission
	if strings.TrimSpace(finalStreamSettings) == "" {
		finalStreamSettings = "{}"
	}
	finalSniffing := inboundData.ExternalProxy
	if strings.TrimSpace(finalSniffing) == "" {
		finalSniffing = "{}"
	}
	// --- THE MAIN FIX: REBUILD THE SETTINGS OBJECT FROM SCRATCH ---
	if inboundData.Protocol == "wireguard" {
		fmt.Println("→ Detected WireGuard protocol. Generating new cryptographic keys...")
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return fmt.Errorf("failed to generate wireguard private key: %v", err)
		}
		publicKey := privateKey.PublicKey()
		var originalWgSettings map[string]interface{}
		if err := json.Unmarshal([]byte(inboundData.OriginalSettings), &originalWgSettings); err != nil {
			originalWgSettings = make(map[string]interface{})
			fmt.Printf(" Warning: Could not parse original WireGuard settings for '%s'. Peers will not be migrated.\n", inboundData.Remark)
		}
		newWgSettings := map[string]interface{}{
			"privateKey": privateKey.String(),
			"publicKey":  publicKey.String(),
			"peers":      originalWgSettings["peers"],
			"mtu":        originalWgSettings["mtu"],
			"listenPort": originalWgSettings["listenPort"],
		}
		settingsBytes, err := json.Marshal(newWgSettings)
		if err != nil {
			return fmt.Errorf("failed to marshal new wireguard settings: %v", err)
		}
		finalSettings = string(settingsBytes)
		fmt.Println(" ✓ New keys generated and settings object rebuilt.")
	} else {
		// For client-based protocols (VLESS, Trojan, Vmess, Shadowsocks):
		clientSettings := []models.ClientSetting{}
		for _, cd := range inboundData.Clients {
			cs := models.ClientSetting{
				ID:         cd.ClientID,
				Email:      cd.ClientEmail,
				Enable:     cd.ClientEnable,
				TotalGB:    cd.ClientTotalGB,
				ExpiryTime: cd.ClientExpiryTime,
				LimitIP:    cd.ClientLimitIP,
				Flow:       cd.ClientFlow,
				SubID:      cd.ClientSubID,
				TgID:       cd.ClientTgID,
				Reset:      cd.ClientReset,
			}
			clientSettings = append(clientSettings, cs)
		}
		// Try to preserve other fields from the original settings
		settingsMap := make(map[string]interface{})
		json.Unmarshal([]byte(inboundData.OriginalSettings), &settingsMap)
		settingsMap["clients"] = clientSettings
		if inboundData.Protocol == "vless" {
			settingsMap["decryption"] = "none"
		}
		settingsBytes, err := json.Marshal(settingsMap)
		if err != nil {
			return fmt.Errorf("failed to marshal rebuilt settings for '%s': %v", inboundData.Remark, err)
		}
		finalSettings = string(settingsBytes)
		fmt.Printf(" ✓ Rebuilt settings for %d clients for inbound '%s'.\n", len(clientSettings), inboundData.Remark)
	}
	// Create and send the final payload
	payload := models.AddInboundPayload{
		Remark:         inboundData.Remark,
		Port:           inboundData.Port,
		Protocol:       inboundData.Protocol,
		Enable:         inboundData.Enable,
		Listen:         inboundData.Listen,
		Total:          inboundData.InboundTotalGB,
		ExpiryTime:     inboundData.InboundExpiry,
		Settings:       finalSettings,
		StreamSettings: finalStreamSettings,
		Sniffing:       finalSniffing,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling add inbound payload: %v", err)
	}
	requestURL := fmt.Sprintf("%s/panel/api/inbounds/add", c.BaseURL)
	resp, err := c.HttpClient.Post(requestURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()
	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("error parsing API response: %v", err)
	}
	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Msg)
	}
	return nil
}

// UpdateInbound updates an existing inbound on the panel.
// Used for conflict resolution when import detects existing inbound.
func (c *ThreeXUIClient) UpdateInbound(inboundID int, inboundData models.InboundData) error {
	var finalSettings string
	// Same logic as AddInbound
	finalStreamSettings := inboundData.Transmission
	if strings.TrimSpace(finalStreamSettings) == "" {
		finalStreamSettings = "{}"
	}
	finalSniffing := inboundData.ExternalProxy
	if strings.TrimSpace(finalSniffing) == "" {
		finalSniffing = "{}"
	}

	if inboundData.Protocol == "wireguard" {
		if strings.TrimSpace(inboundData.OriginalSettings) != "" {
			finalSettings = inboundData.OriginalSettings
		} else {
			finalSettings = "{\"peers\":[]}"
		}
	} else {
		var settingsMap map[string]interface{}
		if strings.TrimSpace(inboundData.OriginalSettings) != "" {
			json.Unmarshal([]byte(inboundData.OriginalSettings), &settingsMap)
		}
		if settingsMap == nil {
			settingsMap = make(map[string]interface{})
		}

		var clientSettings []map[string]interface{}
		for _, clientDetail := range inboundData.Clients {
			clientSetting := map[string]interface{}{
				"id":         clientDetail.ClientID,
				"email":      clientDetail.ClientEmail,
				"enable":     clientDetail.ClientEnable,
				"totalGB":    clientDetail.ClientTotalGB,
				"expiryTime": clientDetail.ClientExpiryTime,
				"limitIp":    clientDetail.ClientLimitIP,
				"flow":       clientDetail.ClientFlow,
				"subId":      clientDetail.ClientSubID,
				"tgId":       clientDetail.ClientTgID,
				"reset":      clientDetail.ClientReset,
			}
			clientSettings = append(clientSettings, clientSetting)
		}

		settingsMap["clients"] = clientSettings
		if inboundData.Protocol == "shadowsocks" {
			settingsMap["method"] = "aes-128-gcm"
			settingsMap["password"] = "password"
		} else if inboundData.Protocol == "socks" {
			settingsMap["accounts"] = clientSettings
			delete(settingsMap, "clients")
		} else {
			settingsMap["decryption"] = "none"
		}
		settingsBytes, err := json.Marshal(settingsMap)
		if err != nil {
			return fmt.Errorf("failed to marshal rebuilt settings for '%s': %v", inboundData.Remark, err)
		}
		finalSettings = string(settingsBytes)
	}

	payload := models.AddInboundPayload{
		Remark:         inboundData.Remark,
		Port:           inboundData.Port,
		Protocol:       inboundData.Protocol,
		Enable:         inboundData.Enable,
		Listen:         inboundData.Listen,
		Total:          inboundData.InboundTotalGB,
		ExpiryTime:     inboundData.InboundExpiry,
		Settings:       finalSettings,
		StreamSettings: finalStreamSettings,
		Sniffing:       finalSniffing,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling update inbound payload: %v", err)
	}
	// درست endpoint برای update: /panel/api/inbounds/update/{id}
	requestURL := fmt.Sprintf("%s/panel/api/inbounds/update/%d", c.BaseURL, inboundID)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating update request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	// اگر status کامیاب ہو تو OK ہے
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

// ExtractClientsFromInbounds processes the raw inbound list and enriches it with client data.
func (c *ThreeXUIClient) ExtractClientsFromInbounds(inbounds []models.Inbound) ([]models.InboundData, int, error) {
	var allInboundData []models.InboundData
	totalUserCount := 0
	for _, inbound := range inbounds {
		inboundData := models.InboundData{
			ID:               inbound.ID,
			Remark:           inbound.Remark,
			Protocol:         inbound.Protocol,
			Port:             inbound.Port,
			Enable:           inbound.Enable,
			Tag:              inbound.Tag,
			Listen:           inbound.Listen,
			InboundExpiry:    inbound.ExpiryTime,
			InboundTotalGB:   inbound.Total,
			Transmission:     inbound.StreamSettings,
			ExternalProxy:    inbound.Sniffing,
			OriginalSettings: inbound.Settings,
			Clients:          []models.ClientDetails{},
		}
		var settings models.InboundSettings
		if strings.TrimSpace(inbound.Settings) == "" || inbound.Settings == "{}" {
			fmt.Printf("→ Processing inbound: %s (%s:%d) - No clients (empty settings)\n", inbound.Remark, inbound.Protocol, inbound.Port)
			allInboundData = append(allInboundData, inboundData)
			continue
		}
		if err := json.Unmarshal([]byte(inbound.Settings), &settings); err != nil {
			fmt.Printf("→ Processing inbound: %s (%s:%d) - Could not parse clients (likely not a client-based protocol like WireGuard). This is OK.\n", inbound.Remark, inbound.Protocol, inbound.Port)
			allInboundData = append(allInboundData, inboundData)
			continue
		}
		fmt.Printf("\n→ Processing inbound: %s (%s:%d)\n", inbound.Remark, inbound.Protocol, inbound.Port)
		fmt.Printf(" Number of clients: %d\n", len(settings.Clients))
		totalUserCount += len(settings.Clients)
		for _, client := range settings.Clients {
			clientDetails := models.ClientDetails{
				ClientEmail:      client.Email,
				ClientID:         client.ID,
				ClientEnable:     client.Enable,
				ClientLimitIP:    client.LimitIP,
				ClientTotalGB:    client.TotalGB,
				ClientExpiryTime: client.ExpiryTime,
				ClientFlow:       client.Flow,
				ClientSubID:      client.SubID,
				ClientTgID:       client.TgID,
				ClientReset:      client.Reset,
			}
			if client.Email != "" {
				traffic, err := c.GetClientTraffic(client.Email)
				if err == nil && traffic != nil {
					trafficUsed := traffic.Up + traffic.Down
					clientDetails.TrafficUsed = trafficUsed
					if client.TotalGB > 0 {
						remaining := client.TotalGB - trafficUsed
						if remaining < 0 {
							remaining = 0
						}
						clientDetails.TrafficRemaining = remaining
						clientDetails.TrafficUsagePercent = (float64(trafficUsed) / float64(client.TotalGB)) * 100
					} else {
						clientDetails.TrafficRemaining = -1
						clientDetails.TrafficUsagePercent = 0
					}
				}
			}
			inboundData.Clients = append(inboundData.Clients, clientDetails)
		}
		allInboundData = append(allInboundData, inboundData)
	}
	return allInboundData, totalUserCount, nil
}
