package importers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"panels_user_manager/pkg/clients"
	"panels_user_manager/pkg/models"
	"panels_user_manager/pkg/utils"
)

// ImportFromJSON handles the process of importing 3X-UI inbounds from a file.
// Similar to PasarGuard import: checks for UUID conflicts and updates existing users.
func ImportFromJSON(client *clients.ThreeXUIClient) {
	fmt.Println("\n" + utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"📥 IMPORT PROCESS STARTED (3X-UI)"+utils.ColorReset, 70) + utils.ColorBrightMagenta + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	filePath := PromptForInputStyled("Enter the path to the JSON file", "\n ➜", utils.ColorBrightYellow)
	fmt.Println("\n " + utils.ColorBrightBlue + "Processing stages:" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "┌─────────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	// 1. Read the file content
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [1/4] " + utils.ColorBrightGreen + "Reading JSON file..." + utils.ColorReset)
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error reading file '%s': %v", filePath, err))
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ File loaded (%d bytes)\n", len(fileBytes))
	// 2. Parse the JSON
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [2/4] " + utils.ColorBrightGreen + "Parsing JSON content..." + utils.ColorReset)
	var dataToImport models.OutputFile
	if err := json.Unmarshal(fileBytes, &dataToImport); err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error parsing JSON file. Make sure it's a valid export file: %v", err))
		return
	}
	if len(dataToImport.Inbounds) == 0 {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No inbounds found in the file to import")
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Found %d inbound(s) to import\n", len(dataToImport.Inbounds))

	// 3. Fetch existing inbounds to check for conflicts
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [3/4] " + utils.ColorBrightGreen + "Checking for conflicts..." + utils.ColorReset)
	existingInbounds, err := client.GetAllInbounds()
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error fetching existing inbounds: %v", err))
		return
	}

	// Build map of existing inbound ports and tags with their IDs
	existingPorts := make(map[int]int)   // port -> inbound ID
	existingTags := make(map[string]int) // tag -> inbound ID
	for _, ib := range existingInbounds {
		existingPorts[ib.Port] = ib.ID
		existingTags[ib.Tag] = ib.ID
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Found %d existing inbound(s)\n", len(existingInbounds))

	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [4/4] " + utils.ColorBrightGreen + "Importing inbounds to panel..." + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	successCount := 0
	failureCount := 0
	updateCount := 0

	// 4. حجم کاربران را بر اساس traffic_remaining تنظیم کنید
	// منطق: حجم کاربر = حجم باقیمانده (traffic_remaining)
	for idx := range dataToImport.Inbounds {
		for jdx := range dataToImport.Inbounds[idx].Clients {
			client := &dataToImport.Inbounds[idx].Clients[jdx]
			// اگر traffic_remaining موجود و مثبت باشد، آن را به عنوان حجم کل استفاده کنید
			if client.TrafficRemaining > 0 {
				client.ClientTotalGB = client.TrafficRemaining
			} else if client.TrafficRemaining == 0 {
				// اگر باقی نیست، حجم را 0 (unlimited) سیٹ کنید
				client.ClientTotalGB = 0
			}
			// اگر TrafficRemaining منفی باشد، ClientTotalGB بدون تغییر باقی می‌ماند
		}
	}

	// 5. Loop through each inbound and create/update it on the panel
	for idx, inbound := range dataToImport.Inbounds {
		fmt.Printf("\n " + utils.ColorBrightBlue + "════════════════════════════════════════════════════════════════════════\n" + utils.ColorReset)
		fmt.Printf(" "+utils.ColorBrightYellow+"[%d/%d] Processing: %s (Port: %d)\n"+utils.ColorReset, idx+1, len(dataToImport.Inbounds), inbound.Remark, inbound.Port)

		// Check if port or tag already exists
		existingID, portExists := existingPorts[inbound.Port]
		existingTagID, tagExists := existingTags[inbound.Tag]

		if portExists || tagExists {
			fmt.Printf(" " + utils.ColorBrightYellow + "⚠️  Conflict detected (Port/Tag already exists)\n" + utils.ColorReset)
			if utils.VerboseMode {
				fmt.Printf(" "+utils.ColorCyan+"Port exists: %v (ID: %d) | Tag exists: %v (ID: %d)\n"+utils.ColorReset, portExists, existingID, tagExists, existingTagID)
				fmt.Printf(" "+utils.ColorCyan+"Protocol: %s | Clients: %d\n"+utils.ColorReset, inbound.Protocol, len(inbound.Clients))
			}

			// Use the existing ID to update
			updateID := existingID
			if tagExists {
				updateID = existingTagID
			}

			fmt.Printf(" " + utils.ColorBrightYellow + "↻ Attempting to update existing inbound...\n" + utils.ColorReset)
			err := client.UpdateInbound(updateID, inbound)
			if err != nil {
				fmt.Printf(" "+utils.ColorBrightRed+"❌ UPDATE FAILED: %v\n"+utils.ColorReset, err)
				failureCount++
			} else {
				fmt.Printf(" " + utils.ColorBrightCyan + "✅ SUCCESS (Updated)\n" + utils.ColorReset)
				updateCount++
			}
		} else {
			if utils.VerboseMode {
				fmt.Printf(" "+utils.ColorCyan+"Protocol: %s | Clients: %d\n"+utils.ColorReset, inbound.Protocol, len(inbound.Clients))
			}
			err := client.AddInbound(inbound)
			if err != nil {
				fmt.Printf(" "+utils.ColorBrightRed+"❌ FAILED: %v\n"+utils.ColorReset, err)
				failureCount++
			} else {
				fmt.Printf(" " + utils.ColorBrightGreen + "✅ SUCCESS (Created)\n" + utils.ColorReset)
				successCount++
				existingPorts[inbound.Port] = inbound.ID
				existingTags[inbound.Tag] = inbound.ID
			}
		}
		fmt.Println(" " + utils.ColorBrightBlue + "════════════════════════════════════════════════════════════════════════" + utils.ColorReset)
	}

	fmt.Println("\n" + utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"📥 IMPORT SUMMARY"+utils.ColorReset, 70) + utils.ColorBrightMagenta + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Printf("\n "+utils.ColorGreen+"✓ Successful imports: %d\n"+utils.ColorReset, successCount)
	if failureCount > 0 {
		fmt.Printf(" "+utils.ColorRed+"✗ Failed/Skipped: %d\n"+utils.ColorReset, failureCount)
	}
	if updateCount > 0 {
		fmt.Printf(" "+utils.ColorYellow+"↻ Updated: %d\n"+utils.ColorReset, updateCount)
	}
	fmt.Printf(" "+utils.ColorCyan+"📊 Total inbounds: %d\n"+utils.ColorReset, len(dataToImport.Inbounds))
	fmt.Printf(" "+utils.ColorBlue+"👥 Total users: %d\n\n"+utils.ColorReset, dataToImport.TotalUsers)
}

// ImportPasarGuardUsersFromJSON handles the process of importing PasarGuard users from a file.
func ImportPasarGuardUsersFromJSON(client *clients.PasarGuardClient) {
	fmt.Println("\n" + utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"📥 IMPORT PROCESS STARTED (PasarGuard)"+utils.ColorReset, 70) + utils.ColorBrightMagenta + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	filePath := PromptForInputStyled("Enter the path to the JSON file", "\n ➜", utils.ColorBrightYellow)
	fmt.Println("\n " + utils.ColorBrightBlue + "Processing stages:" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "┌─────────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	// 1. Read the file content
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [1/3] " + utils.ColorBrightGreen + "Reading JSON file..." + utils.ColorReset)
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error reading file '%s': %v", filePath, err))
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ File loaded (%d bytes)\n", len(fileBytes))
	// 2. Parse the JSON
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [2/3] " + utils.ColorBrightGreen + "Parsing JSON content..." + utils.ColorReset)
	var dataToImport models.PasarGuardUsersExportFile
	if err := json.Unmarshal(fileBytes, &dataToImport); err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error parsing JSON file. Make sure it's a valid PasarGuard export file: %v", err))
		return
	}
	if len(dataToImport.Users) == 0 {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No users found in the file to import")
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Found %d user(s) to import\n", len(dataToImport.Users))

	// Ask user which Groups to assign imported users to
	selectedGroupIDs := []int{}
	groups, gErr := client.GetAllGroups()
	if gErr == nil && len(groups) > 0 {
		fmt.Println("\n " + utils.ColorBrightBlue + "│" + utils.ColorReset + " " + utils.ColorBrightCyan + "Available Groups:" + utils.ColorReset)
		for i, g := range groups {
			fmt.Printf(" \t[%d] %s (id=%d)\n", i+1, g.Name, g.ID)
		}
		sel := PromptForInputStyled("Select group number(s) to assign to imported users (comma-separated), or press Enter to skip", "\n ➜", utils.ColorBrightYellow)
		sel = strings.TrimSpace(sel)
		if sel != "" {
			parts := strings.Split(sel, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				idx, err := strconv.Atoi(p)
				if err != nil {
					fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Invalid number '%s', skipping\n"+utils.ColorReset, p)
					continue
				}
				if idx < 1 || idx > len(groups) {
					fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Selection '%d' out of range, skipping\n"+utils.ColorReset, idx)
					continue
				}
				selectedGroupIDs = append(selectedGroupIDs, groups[idx-1].ID)
			}
			fmt.Printf(" "+utils.ColorBrightGreen+"✓ Selected group IDs: %v\n"+utils.ColorReset, selectedGroupIDs)
		} else {
			fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " No group assignment selected (skipping)")
		}
	} else if gErr != nil {
		fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Could not fetch groups: %v\n"+utils.ColorReset, gErr)
	}

	// Assign groups to all users to import
	for i := range dataToImport.Users {
		dataToImport.Users[i].GroupIDs = selectedGroupIDs
	}
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [3/3] " + utils.ColorBrightGreen + "Importing users to panel..." + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	successCount := 0
	failureCount := 0

	// Prefetch existing users
	existingUsers, prefetchErr := client.GetAllUsers()
	usersByUUID := make(map[string]models.PasarGuardUser)
	usersByUsername := make(map[string]models.PasarGuardUser)
	allUUIDsMap := make(map[string]models.PasarGuardUser)

	if prefetchErr != nil {
		fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Could not prefetch existing users: %v\n"+utils.ColorReset, prefetchErr)
	} else {
		if utils.VerboseMode {
			fmt.Printf(" "+utils.ColorBrightCyan+"📋 Prefetched %d existing users, building UUID maps...\n"+utils.ColorReset, len(existingUsers))
		}
		for idx, existing := range existingUsers {
			uuidKey := strings.ToLower(strings.TrimSpace(existing.UUID))
			if uuidKey != "" {
				usersByUUID[uuidKey] = existing
				allUUIDsMap[uuidKey] = existing
				if utils.VerboseMode {
					fmt.Printf(" "+utils.ColorBrightCyan+"  [%d] Mapped PRIMARY UUID: '%s' → User: '%s'\n"+utils.ColorReset, idx+1, uuidKey, existing.Username)
				}
			}
			if existing.ProxySettings != nil {
				allUUIDs := extractAllUUIDsFromProxySettings(existing.ProxySettings)
				for _, uuid := range allUUIDs {
					if uuid != "" && uuid != uuidKey {
						allUUIDsMap[uuid] = existing
						if utils.VerboseMode {
							fmt.Printf(" "+utils.ColorBrightCyan+"  [%d] Mapped PROXY UUID: '%s' → User: '%s'\n"+utils.ColorReset, idx+1, uuid, existing.Username)
						}
					}
				}
			}
			usernameKey := strings.ToLower(strings.TrimSpace(existing.Username))
			if usernameKey != "" {
				usersByUsername[usernameKey] = existing
			}
		}
		if utils.VerboseMode {
			fmt.Printf(" "+utils.ColorBrightGreen+"✓ UUID map ready with %d entries\n"+utils.ColorReset, len(allUUIDsMap))
		}
	}

	refreshUserLookups := func() {
		fallbackUsers, err := client.GetAllUsers()
		if err != nil {
			if utils.VerboseMode {
				fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Unable to refresh users list: %v\n"+utils.ColorReset, err)
			}
			return
		}
		usersByUUID = make(map[string]models.PasarGuardUser)
		usersByUsername = make(map[string]models.PasarGuardUser)
		allUUIDsMap = make(map[string]models.PasarGuardUser)
		for _, existing := range fallbackUsers {
			uuidKey := strings.ToLower(strings.TrimSpace(existing.UUID))
			if uuidKey != "" {
				usersByUUID[uuidKey] = existing
				allUUIDsMap[uuidKey] = existing
			}
			if existing.ProxySettings != nil {
				allUUIDs := extractAllUUIDsFromProxySettings(existing.ProxySettings)
				for _, uuid := range allUUIDs {
					if uuid != "" {
						allUUIDsMap[uuid] = existing
					}
				}
			}
			usernameKey := strings.ToLower(strings.TrimSpace(existing.Username))
			if usernameKey != "" {
				usersByUsername[usernameKey] = existing
			}
		}
		prefetchErr = nil
	}

	// 3. Loop through each user and create/update it on the panel
	for idx, user := range dataToImport.Users {
		fmt.Printf("\n " + utils.ColorBrightBlue + "════════════════════════════════════════════════════════════════════════\n" + utils.ColorReset)

		if utils.VerboseMode {
			fmt.Printf(utils.ColorBrightMagenta+"[DEBUG] Original - TotalGB: %d, RemainingTraffic: %d\n"+utils.ColorReset, user.TotalGB, user.RemainingTraffic)
		}

		if user.RemainingTraffic > 0 {
			user.TotalGB = user.RemainingTraffic
			if utils.VerboseMode {
				fmt.Printf(utils.ColorBrightMagenta+"[DEBUG] Updated - TotalGB set to RemainingTraffic: %d\n"+utils.ColorReset, user.TotalGB)
			}
		}
		user.UsedTraffic = 0

		if user.ExpiryTime > 1e11 {
			user.ExpiryTime = user.ExpiryTime / 1000
		}

		originalUsername := user.Username
		sanitizedUsername := strings.ToLower(user.Username)
		sanitizedUsername = strings.ReplaceAll(sanitizedUsername, " ", "_")
		sanitizedUsername = strings.TrimSpace(sanitizedUsername)
		user.Username = sanitizedUsername

		email := user.Email
		if email == "" {
			email = originalUsername
		}
		if email == "" {
			email = fmt.Sprintf("User_%d", idx+1)
		}
		fmt.Printf(" "+utils.ColorBrightYellow+"[%d/%d] Processing: %s (Original: %s)\n"+utils.ColorReset, idx+1, len(dataToImport.Users), sanitizedUsername, originalUsername)
		if utils.VerboseMode {
			fmt.Printf(" "+utils.ColorCyan+"Protocol: %s | Port: %d\n"+utils.ColorReset, user.Protocol, user.Port)
		}
		quotaInGB := float64(user.TotalGB) / (1024 * 1024 * 1024)
		if utils.VerboseMode {
			fmt.Printf(" "+utils.ColorBrightGreen+"📊 Final Quota: %.2f GB (%d bytes) | Used Traffic reset to: 0\n"+utils.ColorReset, quotaInGB, user.TotalGB)
		}

		newUUID := strings.ToLower(strings.TrimSpace(user.UUID))
		if newUUID == "" {
			fmt.Printf(" " + utils.ColorBrightRed + "❌ FAILED: User UUID is empty\n" + utils.ColorReset)
			failureCount++
			continue
		}

		if utils.VerboseMode {
			fmt.Printf(" "+utils.ColorBrightCyan+"Checking for existing user with UUID '%s'...\n"+utils.ColorReset, newUUID)
		}
		if utils.VerboseMode {
			fmt.Printf(" "+utils.ColorBrightCyan+"  (Available UUIDs in map: %d)\n"+utils.ColorReset, len(allUUIDsMap))
			for mapUUID := range allUUIDsMap {
				fmt.Printf(" "+utils.ColorBrightCyan+"    - '%s'\n"+utils.ColorReset, mapUUID)
			}
		}
		existingEntry, uuidExists := allUUIDsMap[newUUID]
		if !uuidExists && prefetchErr != nil {
			if utils.VerboseMode {
				fmt.Printf(" " + utils.ColorBrightYellow + "⚠️ UUID not found, attempting refresh...\n" + utils.ColorReset)
			}
			refreshUserLookups()
			existingEntry, uuidExists = allUUIDsMap[newUUID]
		}

		if uuidExists {
			if utils.VerboseMode {
				fmt.Printf(" "+utils.ColorBrightYellow+"✓ Found existing user with matching UUID (Old username: '%s')\n"+utils.ColorReset, existingEntry.Username)
				fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Updating user with UUID '%s' to use new username: '%s'\n"+utils.ColorReset,
					newUUID, user.Username)
			}

			updateErr := client.UpdateUserByIdentifier(existingEntry.Username, user)
			if updateErr != nil {
				fmt.Printf(" "+utils.ColorBrightRed+"❌ UPDATE FAILED: %v\n"+utils.ColorReset, updateErr)
				failureCount++
			} else {
				if utils.VerboseMode {
					fmt.Printf(" "+utils.ColorBrightGreen+"✅ UPDATED SUCCESSFULLY (new username: '%s')\n"+utils.ColorReset, user.Username)
				}

				if len(user.GroupIDs) == 0 {
					if utils.VerboseMode {
						fmt.Printf(" "+utils.ColorBrightCyan+"🗑️ Clearing group assignments for user '%s'...\n"+utils.ColorReset, user.Username)
					}
					clearErr := client.ClearUserGroups(user.Username)
					if clearErr != nil {
						if utils.VerboseMode {
							fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Warning: Could not clear groups: %v\n"+utils.ColorReset, clearErr)
						}
					} else {
						if utils.VerboseMode {
							fmt.Printf(" " + utils.ColorBrightGreen + "✅ Groups cleared successfully\n" + utils.ColorReset)
						}
					}
				}

				successCount++

				updatedEntry := existingEntry
				updatedEntry.Username = user.Username
				updatedEntry.TotalGB = user.TotalGB
				updatedEntry.ExpiryTime = user.ExpiryTime
				updatedEntry.Enable = user.Enable
				updatedEntry.Note = user.Note
				updatedEntry.LimitIP = user.LimitIP
				updatedEntry.UsedTraffic = user.UsedTraffic
				updatedEntry.RemainingTraffic = user.RemainingTraffic
				updatedEntry.GroupIDs = user.GroupIDs
				usersByUUID[newUUID] = updatedEntry
				allUUIDsMap[newUUID] = updatedEntry
				usersByUsername[strings.ToLower(strings.TrimSpace(updatedEntry.Username))] = updatedEntry
			}
			continue
		}

		if originalUsername == "" {
			originalUsername = fmt.Sprintf("user_%d", idx+1)
		}

		maxAttempts := 10
		addSuccess := false
		failureRecorded := false

		for attemptCount := 0; attemptCount < maxAttempts; attemptCount++ {
			candidateUsername := sanitizedUsername
			if attemptCount > 0 {
				candidateUsername = fmt.Sprintf("%s_%d", sanitizedUsername, attemptCount)
				fmt.Printf(" "+utils.ColorBrightYellow+"Trying username '%s'...\n"+utils.ColorReset, candidateUsername)
			}

			usernameKey := strings.ToLower(strings.TrimSpace(candidateUsername))
			if usernameKey == "" {
				continue
			}

			if _, exists := usersByUsername[usernameKey]; exists {
				if attemptCount == 0 {
					fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Username '%s' already exists, searching for free variant...\n"+utils.ColorReset, candidateUsername)
				}
				continue
			}

			user.Username = candidateUsername
			addErr := client.AddUser(user)
			if addErr == nil {
				if attemptCount == 0 {
					fmt.Printf(" " + utils.ColorBrightGreen + "✅ SUCCESS\n" + utils.ColorReset)
				} else {
					fmt.Printf(" "+utils.ColorBrightGreen+"✅ SUCCESS (created as '%s')\n"+utils.ColorReset, candidateUsername)
				}
				successCount++
				addSuccess = true

				stored := user
				stored.Username = candidateUsername
				usersByUUID[newUUID] = stored
				allUUIDsMap[newUUID] = stored
				usersByUsername[usernameKey] = stored
				break
			}

			errMsg := strings.ToLower(addErr.Error())
			if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "409") || strings.Contains(errMsg, "user already exists") {
				if attemptCount == 0 {
					fmt.Printf(" "+utils.ColorBrightYellow+"⚠️ Username '%s' already exists with different UUID, trying alternatives...\n"+utils.ColorReset, candidateUsername)
				}
				usersByUsername[usernameKey] = models.PasarGuardUser{Username: candidateUsername}
				continue
			}

			fmt.Printf(" "+utils.ColorBrightRed+"❌ FAILED: %v\n"+utils.ColorReset, addErr)
			failureCount++
			failureRecorded = true
			break
		}

		if !addSuccess && !failureRecorded {
			fmt.Printf(" "+utils.ColorBrightRed+"❌ FAILED: Could not create user after %d attempts\n"+utils.ColorReset, maxAttempts)
			failureCount++
		}
	}

	fmt.Println(" " + utils.ColorBrightBlue + "════════════════════════════════════════════════════════════════════════" + utils.ColorReset)
	fmt.Println("\n" + utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"📥 IMPORT SUMMARY (PasarGuard)"+utils.ColorReset, 70) + utils.ColorBrightMagenta + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Printf("\n "+utils.ColorGreen+"✓ Successful imports/updates: %d\n"+utils.ColorReset, successCount)
	if failureCount > 0 {
		fmt.Printf(" "+utils.ColorRed+"✗ Failed imports: %d\n"+utils.ColorReset, failureCount)
	}
	fmt.Printf(" "+utils.ColorCyan+"📊 Total users: %d\n\n"+utils.ColorReset, len(dataToImport.Users))
}

// PromptForInputStyled displays a styled prompt and returns the user's input.
func PromptForInputStyled(label, prefix, color string) string {
	fmt.Printf("%s %s%s%s: ", prefix, color, label, utils.ColorReset)
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
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
