package utils

import (
	"fmt"
	"strings"
)

// --- COLOR CODES FOR ANSI TERMINAL ---
const (
	// Foreground Colors
	ColorReset     = "\033[0m"
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
	ColorItalic    = "\033[3m"
	ColorUnderline = "\033[4m"
	// Standard Colors
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	// Bright Colors
	ColorBrightBlack   = "\033[90m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"
	// Background Colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
	// Bright Background Colors
	BgBrightBlack   = "\033[100m"
	BgBrightRed     = "\033[101m"
	BgBrightGreen   = "\033[102m"
	BgBrightYellow  = "\033[103m"
	BgBrightBlue    = "\033[104m"
	BgBrightMagenta = "\033[105m"
	BgBrightCyan    = "\033[106m"
	BgBrightWhite   = "\033[107m"
)

// Global verbose mode flag
var VerboseMode = false

// VerboseLog prints only when verbose mode is enabled
func VerboseLog(format string, args ...interface{}) {
	if VerboseMode {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// ClearScreen clears the console screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

// PrintBox prints a centered box with styling
func PrintBox(title, content string) {
	fmt.Println("\n " + ColorBrightMagenta + "┏" + strings.Repeat("━", 66) + "┓" + ColorReset)
	fmt.Println(" " + ColorBrightMagenta + "┃" + ColorReset + CenterText(ColorBold+ColorBrightCyan+title+ColorReset, 66) + ColorBrightMagenta + "┃" + ColorReset)
	fmt.Println(" " + ColorBrightMagenta + "┣" + strings.Repeat("━", 66) + "┫" + ColorReset)
	fmt.Println(" " + ColorBrightMagenta + "┃" + ColorReset + strings.Repeat(" ", 66) + ColorBrightMagenta + "┃" + ColorReset)
	for _, line := range strings.Split(content, "\n") {
		if len(line) > 66 {
			line = line[:63] + "..."
		}
		fmt.Println(" " + ColorBrightMagenta + "┃" + ColorReset + " " + ColorCyan + line + ColorReset + strings.Repeat(" ", 63-len(line)) + ColorBrightMagenta + "┃" + ColorReset)
	}
	fmt.Println(" " + ColorBrightMagenta + "┃" + ColorReset + strings.Repeat(" ", 66) + ColorBrightMagenta + "┃" + ColorReset)
	fmt.Println(" " + ColorBrightMagenta + "┗" + strings.Repeat("━", 66) + "┛" + ColorReset + "\n")
}

// PrintSuccess prints a success message with styling
func PrintSuccess(message string) {
	fmt.Printf("\n "+ColorBrightGreen+"✅ %s"+ColorReset+"\n", message)
}

// PrintError prints an error message with styling
func PrintError(message string) {
	fmt.Printf("\n "+ColorBrightRed+"❌ %s"+ColorReset+"\n", message)
}

// PrintWarning prints a warning message with styling
func PrintWarning(message string) {
	fmt.Printf("\n "+ColorBrightYellow+"⚠️ %s"+ColorReset+"\n", message)
}

// PrintInfo prints an info message with styling
func PrintInfo(message string) {
	fmt.Printf("\n "+ColorBrightBlue+"ℹ️ %s"+ColorReset+"\n", message)
}

// FormatBytes converts bytes into a human-readable format.
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	index := 0
	value := float64(bytes)
	for value >= 1024 && index < len(units)-1 {
		value /= 1024
		index++
	}
	return fmt.Sprintf("%.2f %s", value, units[index])
}

// GenerateProgressBar creates a visual progress bar
func GenerateProgressBar(percent int) string {
	if percent > 100 {
		percent = 100
	}
	filled := percent / 5
	empty := 20 - filled
	// Color the bar based on percentage
	var color string
	if percent >= 80 {
		color = ColorBrightRed
	} else if percent >= 60 {
		color = ColorBrightYellow
	} else if percent >= 40 {
		color = ColorBrightCyan
	} else {
		color = ColorBrightGreen
	}
	return "[" + color + strings.Repeat("█", filled) + ColorReset + ColorBrightBlack + strings.Repeat("░", empty) + ColorReset + "]"
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
