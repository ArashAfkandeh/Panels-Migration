package main

import (
	"bufio"
	"flag"
	"os"

	"panels_user_manager/pkg/cmd"
	"panels_user_manager/pkg/utils"
)

func main() {
	// Parse command-line flags
	flag.BoolVar(&utils.VerboseMode, "v", false, "Enable verbose logging (use: -v)")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	for {
		cmd.ShowMenu()
		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1] // Remove newline

		switch choice {
		case "1":
			// 3X-UI Panel Operations
			cmd.Handle3XUIMenu(reader)
		case "2":
			// PasarGuard Panel Operations
			cmd.HandlePasarGuardMenu(reader)
		case "3":
			os.Exit(0)
		default:
			println("Invalid option. Please try again.\n")
			println("Press Enter to continue...")
			reader.ReadString('\n')
		}
	}
}
