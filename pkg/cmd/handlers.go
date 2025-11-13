package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"panels_user_manager/pkg/clients"
	"panels_user_manager/pkg/exporters"
	"panels_user_manager/pkg/importers"
	"panels_user_manager/pkg/utils"
)

// ShowMenu displays the main menu with attractive styling.
func ShowMenu() {
	utils.ClearScreen()
	fmt.Println("\n" + utils.ColorBrightYellow + utils.ColorBold + "🌐 MULTI-PANEL DATA MANAGER 🌐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "Professional Inbound & User Management" + utils.ColorReset)
	fmt.Println()
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBold + utils.ColorBrightWhite + "                              📋 MAIN MENU" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightCyan + "  🔧 PANEL SELECTION" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[1] 3X-UI Panel Operations" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "├─ Export all inbounds and users" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "├─ Export users (PasarGuard format)" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Import inbounds and users from file" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[2] PasarGuard Panel Operations" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "├─ Export users" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Import users from file" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightRed + "  🚪 APPLICATION CONTROL" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[3] Exit Application" + utils.ColorReset + utils.ColorDim + " (close and return to system)" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Print(utils.ColorBrightMagenta + "  ➜ Select an option (1-3): " + utils.ColorReset)
}

// Show3XUIMenu displays the 3X-UI panel menu.
func Show3XUIMenu() {
	utils.ClearScreen()
	fmt.Println("\n" + utils.ColorBrightYellow + utils.ColorBold + "🌐 3X-UI PANEL OPERATIONS 🌐" + utils.ColorReset)
	fmt.Println()
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBold + utils.ColorBrightWhite + "                          📋 3X-UI OPERATIONS MENU" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightGreen + "  📤 EXPORT OPERATIONS" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[1] Export all inbounds and users to JSON" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Backup all configurations and user data" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[2] Export users only (PasarGuard format)" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Export users in PasarGuard-compatible format" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightYellow + "  📥 IMPORT OPERATIONS" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[3] Import inbounds and users from JSON" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Restore from backup file" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightRed + "  🔙 NAVIGATION" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[4] Return to main menu" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Print(utils.ColorBrightMagenta + "  ➜ Select an option (1-4): " + utils.ColorReset)
}

// ShowPasarGuardMenu displays the PasarGuard panel menu.
func ShowPasarGuardMenu() {
	utils.ClearScreen()
	fmt.Println("\n" + utils.ColorBrightYellow + utils.ColorBold + "🛡️ PASARGUARD PANEL OPERATIONS 🛡️" + utils.ColorReset)
	fmt.Println()
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBold + utils.ColorBrightWhite + "                      📋 PASARGUARD OPERATIONS MENU" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightGreen + "  📤 EXPORT OPERATIONS" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[1] Export users to JSON" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Backup all user data" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightYellow + "  📥 IMPORT OPERATIONS" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[2] Import users from JSON" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + "     " + utils.ColorDim + "└─ Restore users from backup file" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Println(utils.ColorBrightRed + "  🔙 NAVIGATION" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  ┌─────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  │" + utils.ColorReset + " " + utils.ColorBrightWhite + "[3] Return to main menu" + utils.ColorReset)
	fmt.Println(utils.ColorBrightBlue + "  └─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
	fmt.Println()

	fmt.Print(utils.ColorBrightMagenta + "  ➜ Select an option (1-3): " + utils.ColorReset)
}

// GetLoginSettings prompts the user for panel connection details.
func GetLoginSettings() (string, string, string) {
	fmt.Println("\n" + utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"🔐 PANEL CONNECTION SETTINGS"+utils.ColorReset, 70) + utils.ColorBrightMagenta + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightMagenta + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightBlue + "┌─ Panel Configuration" + utils.ColorReset)
	baseURL := PromptForInputStyled("Panel Address", " │ (e.g., http://127.0.0.1:2053)", utils.ColorBrightGreen)
	fmt.Println("\n " + utils.ColorBrightBlue + "┌─ Authentication Credentials" + utils.ColorReset)
	username := PromptForInputStyled("Username", " │", utils.ColorBrightYellow)
	password := PromptForInputStyled("Password", " └", utils.ColorBrightRed)
	fmt.Println()
	fmt.Println(" " + utils.ColorBrightCyan + "┌─ Connection Summary" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorGreen+"🔗 URL: "+utils.ColorReset+"%s\n", baseURL)
	fmt.Printf(" │ "+utils.ColorYellow+"👤 User: "+utils.ColorReset+"%s\n", username)
	fmt.Println(" " + utils.ColorBrightCyan + "└─ Ready to connect " + utils.ColorGreen + "✓" + utils.ColorReset)
	return baseURL, username, password
}

// GetExportSettings prompts for all details required for exporting 3X-UI.
func GetExportSettings() (string, string, string, string) {
	baseURL, username, password := GetLoginSettings()
	fmt.Println("\n" + utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightYellow+"📤 EXPORT CONFIGURATION"+utils.ColorReset, 70) + utils.ColorBrightGreen + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
	defaultFilename := "3xui_users_data.json"
	fmt.Printf("\n " + utils.ColorBrightCyan + "📁 Output File Configuration\n" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorCyan+"Default filename: "+utils.ColorReset+"%s\n", defaultFilename)
	fmt.Printf(" │\n")
	filename := PromptForInputStyled("Enter custom filename (or press Enter for default)", " └", utils.ColorBrightMagenta)
	if filename == "" {
		filename = defaultFilename
		fmt.Printf(" "+utils.ColorGreen+"✓ Using default: "+utils.ColorReset+"%s\n", filename)
	}
	return baseURL, username, password, filename
}

// GetPasarGuardExportSettings prompts for all details required for exporting PasarGuard users.
func GetPasarGuardExportSettings() (string, string, string, string) {
	baseURL, username, password := GetLoginSettings()
	fmt.Println("\n" + utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightYellow+"📤 EXPORT CONFIGURATION (PasarGuard)"+utils.ColorReset, 70) + utils.ColorBrightGreen + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
	defaultFilename := "pasarguard_users_data.json"
	fmt.Printf("\n " + utils.ColorBrightCyan + "📁 Output File Configuration\n" + utils.ColorReset)
	fmt.Printf(" │ "+utils.ColorCyan+"Default filename: "+utils.ColorReset+"%s\n", defaultFilename)
	fmt.Printf(" │\n")
	filename := PromptForInputStyled("Enter custom filename (or press Enter for default)", " └", utils.ColorBrightMagenta)
	if filename == "" {
		filename = defaultFilename
		fmt.Printf(" "+utils.ColorGreen+"✓ Using default: "+utils.ColorReset+"%s\n", filename)
	}
	return baseURL, username, password, filename
}

// PromptForInputStyled displays a styled prompt and returns the user's input.
func PromptForInputStyled(label, prefix, color string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s %s%s%s: ", prefix, color, label, utils.ColorReset)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Handle3XUIMenu handles the 3X-UI panel menu operations.
func Handle3XUIMenu(reader *bufio.Reader) {
	for {
		Show3XUIMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			baseURL, username, password, filename := GetExportSettings()
			RunExporter(baseURL, username, password, filename)
			fmt.Println("\nPress Enter to return to the menu...")
			reader.ReadString('\n')
		case "2":
			baseURL, username, password, filename := GetExportSettings()
			RunUsersExporter(baseURL, username, password, filename)
			fmt.Println("\nPress Enter to return to the menu...")
			reader.ReadString('\n')
		case "3":
			baseURL, username, password := GetLoginSettings()
			client := clients.NewThreeXUIClient(baseURL, username, password)
			if err := client.Login(); err != nil {
				fmt.Printf("\n✗ Login failed, cannot proceed with import: %v\n", err)
				fmt.Println("\nPress Enter to return to the menu...")
				reader.ReadString('\n')
				continue
			}
			importers.ImportFromJSON(client)
			fmt.Println("\nPress Enter to return to the menu...")
			reader.ReadString('\n')
		case "4":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			fmt.Println("\nPress Enter to continue...")
			reader.ReadString('\n')
		}
	}
}

// HandlePasarGuardMenu handles the PasarGuard panel menu operations.
func HandlePasarGuardMenu(reader *bufio.Reader) {
	for {
		ShowPasarGuardMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			baseURL, username, password, filename := GetPasarGuardExportSettings()
			RunPasarGuardExporter(baseURL, username, password, filename)
			fmt.Println("\nPress Enter to return to the menu...")
			reader.ReadString('\n')
		case "2":
			baseURL, username, password := GetLoginSettings()
			client := clients.NewPasarGuardClient(baseURL, username, password)
			if err := client.Login(); err != nil {
				fmt.Printf("\n✗ Login failed, cannot proceed with import: %v\n", err)
				fmt.Println("\nPress Enter to return to the menu...")
				reader.ReadString('\n')
				continue
			}
			importers.ImportPasarGuardUsersFromJSON(client)
			fmt.Println("\nPress Enter to return to the menu...")
			reader.ReadString('\n')
		case "3":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			fmt.Println("\nPress Enter to continue...")
			reader.ReadString('\n')
		}
	}
}

// RunExporter executes the main export logic for 3X-UI.
func RunExporter(baseURL, username, password, filename string) {
	fmt.Println("\n" + utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightYellow+"📤 EXPORT PROCESS STARTED"+utils.ColorReset, 70) + utils.ColorBrightCyan + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightBlue + "Processing stages:" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "┌─────────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	client := clients.NewThreeXUIClient(baseURL, username, password)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [1/4] " + utils.ColorBrightGreen + "Authenticating with panel..." + utils.ColorReset)
	if err := client.Login(); err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Failed to log in: %v", err))
		return
	}
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " " + utils.ColorGreen + "✓ Authentication successful" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [2/4] " + utils.ColorBrightGreen + "Fetching inbounds list..." + utils.ColorReset)
	inbounds, err := client.GetAllInbounds()
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error fetching inbounds: %v", err))
		return
	}
	if len(inbounds) == 0 {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No inbounds found")
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Found %d inbound(s)\n", len(inbounds))
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [3/4] " + utils.ColorBrightGreen + "Extracting client data..." + utils.ColorReset)
	inboundsData, totalUsers, err := client.ExtractClientsFromInbounds(inbounds)
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error extracting clients: %v", err))
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Extracted data for %d users\n", totalUsers)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [4/4] " + utils.ColorBrightGreen + "Saving to JSON file..." + utils.ColorReset)
	if len(inboundsData) > 0 {
		if err := exporters.SaveToJSON(inboundsData, totalUsers, filename); err != nil {
			fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
			utils.PrintError(fmt.Sprintf("Error saving file: %v", err))
		} else {
			fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
			utils.PrintSuccess(fmt.Sprintf("Export completed successfully! Saved to: %s", filename))
		}
	} else {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No inbounds were processed")
	}
}

// RunUsersExporter exports 3X-UI users in PasarGuard-compatible format.
func RunUsersExporter(baseURL, username, password, filename string) {
	fmt.Println("\n" + utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightYellow+"📤 USERS EXPORT (PasarGuard Format)"+utils.ColorReset, 70) + utils.ColorBrightCyan + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightBlue + "Processing stages:" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "┌─────────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	client := clients.NewThreeXUIClient(baseURL, username, password)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [1/4] " + utils.ColorBrightGreen + "Authenticating with panel..." + utils.ColorReset)
	if err := client.Login(); err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Failed to log in: %v", err))
		return
	}
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " " + utils.ColorGreen + "✓ Authentication successful" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [2/4] " + utils.ColorBrightGreen + "Fetching inbounds list..." + utils.ColorReset)
	inbounds, err := client.GetAllInbounds()
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error fetching inbounds: %v", err))
		return
	}
	if len(inbounds) == 0 {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No inbounds found")
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Found %d inbound(s)\n", len(inbounds))
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [3/4] " + utils.ColorBrightGreen + "Extracting client data..." + utils.ColorReset)
	inboundsData, totalUsers, err := client.ExtractClientsFromInbounds(inbounds)
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error extracting clients: %v", err))
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Extracted data for %d users\n", totalUsers)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [4/4] " + utils.ColorBrightGreen + "Saving to JSON file..." + utils.ColorReset)
	if len(inboundsData) > 0 {
		if err := exporters.SaveThreeXUIUsersToJSON(inboundsData, filename); err != nil {
			fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
			utils.PrintError(fmt.Sprintf("Error saving file: %v", err))
		} else {
			fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
			fmt.Println("\n" + utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
			fmt.Println(utils.ColorBrightGreen + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightCyan+"✅ EXPORT COMPLETED SUCCESSFULLY"+utils.ColorReset, 70) + utils.ColorBrightGreen + "║" + utils.ColorReset)
			fmt.Println(utils.ColorBrightGreen + strings.Repeat("═", 72) + utils.ColorReset)
			fmt.Printf("\n "+utils.ColorGreen+"📁 Export File: "+utils.ColorReset+"%s\n", filename)
			fmt.Printf(" "+utils.ColorCyan+"👥 Total Users: "+utils.ColorReset+"%d\n\n", totalUsers)
		}
	} else {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No inbounds were processed")
	}
}

// RunPasarGuardExporter executes the export logic for PasarGuard panel (users only).
func RunPasarGuardExporter(baseURL, username, password, filename string) {
	fmt.Println("\n" + utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + "║" + utils.ColorReset + utils.CenterText(utils.ColorBold+utils.ColorBrightYellow+"📤 EXPORT PROCESS STARTED (PasarGuard)"+utils.ColorReset, 70) + utils.ColorBrightCyan + "║" + utils.ColorReset)
	fmt.Println(utils.ColorBrightCyan + strings.Repeat("═", 72) + utils.ColorReset)
	fmt.Println("\n " + utils.ColorBrightBlue + "Processing stages:" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "┌─────────────────────────────────────────────────────────────────┐" + utils.ColorReset)
	client := clients.NewPasarGuardClient(baseURL, username, password)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [1/3] " + utils.ColorBrightGreen + "Authenticating with PasarGuard panel..." + utils.ColorReset)
	if err := client.Login(); err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Failed to log in: %v", err))
		return
	}
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " " + utils.ColorGreen + "✓ Authentication successful" + utils.ColorReset)
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [2/3] " + utils.ColorBrightGreen + "Fetching users list..." + utils.ColorReset)
	users, err := client.GetAllUsers()
	if err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error fetching users: %v", err))
		return
	}
	if len(users) == 0 {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintWarning("No users found")
		return
	}
	fmt.Printf(" "+utils.ColorBrightBlue+"│"+utils.ColorReset+" "+utils.ColorGreen+"✓ Found %d user(s)\n", len(users))
	fmt.Println(" " + utils.ColorBrightBlue + "│" + utils.ColorReset + " [3/3] " + utils.ColorBrightGreen + "Saving to JSON file..." + utils.ColorReset)
	if err := exporters.SavePasarGuardUsersToJSON(users, filename); err != nil {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintError(fmt.Sprintf("Error saving file: %v", err))
	} else {
		fmt.Println(" " + utils.ColorBrightBlue + "└─────────────────────────────────────────────────────────────────┘" + utils.ColorReset)
		utils.PrintSuccess(fmt.Sprintf("Export completed successfully! Saved to: %s", filename))
	}
}
