package exporters

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"panels_user_manager/pkg/models"
	"panels_user_manager/pkg/utils"
)

// SaveToJSON saves the extracted 3X-UI data (inbounds + users) to a JSON file and prints stats.
func SaveToJSON(inboundsData []models.InboundData, totalUsers int, filename string) error {
	output := models.OutputFile{
		ExportDate:    time.Now().Format(time.RFC3339),
		TotalInbounds: len(inboundsData),
		TotalUsers:    totalUsers,
		Inbounds:      inboundsData,
	}

	fileData, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return fmt.Errorf("error creating output JSON: %v", err)
	}
	if err := os.WriteFile(filename, fileData, 0644); err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	// --- Display stats ---
	fmt.Println("\n" + utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"✅ EXPORT COMPLETED SUCCESSFULLY"+utils.ColorReset, 70) + utils.ColorBrightGreen + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)

	activeUsers := 0
	var totalTrafficUsed, totalTrafficLimit, totalTrafficRemaining int64

	// Calculate stats from inbounds data
	for _, inbound := range inboundsData {
		for _, client := range inbound.Clients {
			if client.ClientEnable {
				activeUsers++
			}
			totalTrafficUsed += client.TrafficUsed
			if client.ClientTotalGB > 0 {
				totalTrafficLimit += client.ClientTotalGB
			}
			if client.TrafficRemaining > 0 {
				totalTrafficRemaining += client.TrafficRemaining
			}
		}
	}

	fmt.Println("\n " + utils.ColorBrightBlue + "┌─ GENERAL INFORMATION" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorGreen+"📁 Export File: "+utils.ColorReset+"%s\n", filename)
	fmt.Printf(" │ "+utils.ColorCyan+"🕐 Export Date: "+utils.ColorReset+"%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(" " + utils.ColorBrightBlue + "└" + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightGreen + "┌─ INBOUND STATISTICS" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorYellow+"📍 Total Inbounds: "+utils.ColorReset+"%d\n", len(inboundsData))
	enabledInbounds := 0
	disabledInbounds := 0
	for _, ib := range inboundsData {
		if ib.Enable {
			enabledInbounds++
		} else {
			disabledInbounds++
		}
	}
	fmt.Printf(" │ "+utils.ColorGreen+"✅ Enabled: "+utils.ColorReset+"%d | "+utils.ColorRed+"❌ Disabled: "+utils.ColorReset+"%d\n", enabledInbounds, disabledInbounds)
	fmt.Println(" " + utils.ColorBrightGreen + "└" + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightMagenta + "┌─ USER STATISTICS" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorMagenta+"👥 Total Users: "+utils.ColorReset+"%d\n", totalUsers)
	fmt.Printf(" │ "+utils.ColorGreen+"✅ Active Users: "+utils.ColorReset+"%d\n", activeUsers)
	fmt.Printf(" │ "+utils.ColorRed+"❌ Inactive Users: "+utils.ColorReset+"%d\n", totalUsers-activeUsers)
	fmt.Println(" " + utils.ColorBrightMagenta + "└" + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightYellow + "┌─ TRAFFIC STATISTICS" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorCyan+"📤 Total Traffic Used: "+utils.ColorReset+"%s\n", utils.FormatBytes(totalTrafficUsed))
	fmt.Printf(" │ "+utils.ColorBlue+"📊 Total Traffic Allocated: "+utils.ColorReset+"%s\n", utils.FormatBytes(totalTrafficLimit))
	fmt.Printf(" │ "+utils.ColorGreen+"📥 Total Traffic Remaining: "+utils.ColorReset+"%s\n", utils.FormatBytes(totalTrafficRemaining))
	if totalTrafficLimit > 0 {
		usagePercent := (float64(totalTrafficUsed) / float64(totalTrafficLimit)) * 100
		bar := utils.GenerateProgressBar(int(usagePercent))
		fmt.Printf(" │ "+utils.ColorYellow+"📈 Overall Usage: "+utils.ColorReset+"%s (%.1f%%)\n", bar, usagePercent)
	}
	fmt.Println(" " + utils.ColorBrightYellow + "└" + utils.ColorReset)
	fmt.Println("\n" + utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset + "\n")
	return nil
}

// SaveThreeXUIUsersToJSON saves 3X-UI users in PasarGuard format for compatibility
// تمام کاربران را به ساختار PasarGuard convert می‌کند تا با فایل‌های PasarGuard compatible باشند
func SaveThreeXUIUsersToJSON(inboundsData []models.InboundData, filename string) error {
	var users []models.PasarGuardUser
	userID := 1

	for _, inbound := range inboundsData {
		for _, client := range inbound.Clients {
			user := models.PasarGuardUser{
				ID:               userID,
				Username:         client.ClientEmail,
				Email:            client.ClientEmail,
				UUID:             client.ClientID,
				Enable:           client.ClientEnable,
				TotalGB:          client.ClientTotalGB,
				ExpiryTime:       client.ClientExpiryTime,
				LimitIP:          client.ClientLimitIP,
				UsedTraffic:      client.TrafficUsed,
				RemainingTraffic: client.TrafficRemaining,
				Protocol:         inbound.Protocol,
				Port:             inbound.Port,
				Remark:           inbound.Remark,
				SubscriptionURL:  "",
				Note:             inbound.Remark,
				ProxySettings:    make(map[string]interface{}),
				GroupIDs:         []int{},
			}
			users = append(users, user)
			userID++
		}
	}

	output := models.PasarGuardUsersExportFile{
		ExportDate: time.Now().Format(time.RFC3339),
		PanelType:  "3X-UI",
		TotalUsers: len(users),
		Users:      users,
	}

	fileData, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return fmt.Errorf("error creating output JSON: %v", err)
	}
	if err := os.WriteFile(filename, fileData, 0644); err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}

// SavePasarGuardUsersToJSON saves PasarGuard users to a JSON file and prints stats.
func SavePasarGuardUsersToJSON(users []models.PasarGuardUser, filename string) error {
	output := models.PasarGuardUsersExportFile{
		ExportDate: time.Now().Format(time.RFC3339),
		PanelType:  "PasarGuard",
		TotalUsers: len(users),
		Users:      users,
	}
	fileData, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return fmt.Errorf("error creating output JSON: %v", err)
	}
	if err := os.WriteFile(filename, fileData, 0644); err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	fmt.Println("\n" + utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightYellow+"📊 EXPORT STATISTICS (PasarGuard)"+utils.ColorReset, 70) + utils.ColorBrightCyan + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)

	activeUsers := 0
	var totalTrafficUsed, totalTrafficLimit, totalTrafficRemaining int64
	for _, user := range users {
		if user.Enable {
			activeUsers++
		}
		totalTrafficUsed += user.UsedTraffic
		if user.TotalGB > 0 {
			totalTrafficLimit += user.TotalGB
		}
		if user.RemainingTraffic > 0 {
			totalTrafficRemaining += user.RemainingTraffic
		}
	}

	fmt.Println("\n " + utils.ColorBrightBlue + "┌─ GENERAL INFORMATION" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorGreen+"📁 Export File: "+utils.ColorReset+"%s\n", filename)
	fmt.Printf(" │ "+utils.ColorCyan+"🕐 Export Date: "+utils.ColorReset+"%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf(" │ " + utils.ColorMagenta + "🌐 Panel Type: " + utils.ColorReset + "PasarGuard\n")
	fmt.Println(" " + utils.ColorBrightBlue + "└" + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightMagenta + "┌─ USER STATISTICS" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorMagenta+"👥 Total Users: "+utils.ColorReset+"%d\n", len(users))
	fmt.Printf(" │ "+utils.ColorGreen+"✅ Active Users: "+utils.ColorReset+"%d\n", activeUsers)
	fmt.Printf(" │ "+utils.ColorRed+"❌ Inactive Users: "+utils.ColorReset+"%d\n", len(users)-activeUsers)
	fmt.Println(" " + utils.ColorBrightMagenta + "└" + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightYellow + "┌─ TRAFFIC STATISTICS" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorCyan+"📤 Total Traffic Used: "+utils.ColorReset+"%s\n", utils.FormatBytes(totalTrafficUsed))
	fmt.Printf(" │ "+utils.ColorBlue+"📊 Total Traffic Allocated: "+utils.ColorReset+"%s\n", utils.FormatBytes(totalTrafficLimit))
	fmt.Printf(" │ "+utils.ColorGreen+"📥 Total Traffic Remaining: "+utils.ColorReset+"%s\n", utils.FormatBytes(totalTrafficRemaining))
	if totalTrafficLimit > 0 {
		usagePercent := (float64(totalTrafficUsed) / float64(totalTrafficLimit)) * 100
		bar := utils.GenerateProgressBar(int(usagePercent))
		fmt.Printf(" │ "+utils.ColorYellow+"📈 Overall Usage: "+utils.ColorReset+"%s (%.1f%%)\n", bar, usagePercent)
	}
	fmt.Println(" " + utils.ColorBrightYellow + "└" + utils.ColorReset)

	sort.Slice(users, func(i, j int) bool {
		return users[i].UsedTraffic > users[j].UsedTraffic
	})
	fmt.Println("\n " + utils.ColorBrightMagenta + "┌─ TOP TRAFFIC CONSUMERS" + utils.ColorReset)
	topN := 5
	if len(users) < topN {
		topN = len(users)
	}
	if topN > 0 {
		for i := 0; i < topN; i++ {
			user := users[i]
			email := user.Email
			if email == "" {
				email = user.Username
			}
			if email == "" {
				email = "No Email/Username"
			}
			used := utils.FormatBytes(user.UsedTraffic)
			if user.TotalGB <= 0 {
				fmt.Printf(" │ "+utils.ColorBrightCyan+"🥇 %d. "+utils.ColorReset+"%-25s │ %s (Unlimited)\n", i+1, email[:utils.Min(25, len(email))], used)
			} else {
				remaining := utils.FormatBytes(user.RemainingTraffic)
				fmt.Printf(" │ "+utils.ColorBrightCyan+"🥇 %d. "+utils.ColorReset+"%-25s │ %s / %s\n", i+1, email[:utils.Min(25, len(email))], used, remaining)
			}
		}
	} else {
		fmt.Println(" │ " + utils.ColorDim + "No users found" + utils.ColorReset)
	}
	fmt.Println(" " + utils.ColorBrightMagenta + "└" + utils.ColorReset)
	fmt.Println("\n" + utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"✅ EXPORT COMPLETED SUCCESSFULLY"+utils.ColorReset, 70) + utils.ColorBrightGreen + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset + "\n")
	return nil
}
