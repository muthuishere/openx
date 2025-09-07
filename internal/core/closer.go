package core

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CloseApp closes an application by killing its processes
func CloseApp(alias string) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	app, exists := config.Apps[alias]
	if !exists {
		// Check if it's an alias
		if canonical, ok := config.Aliases[alias]; ok {
			app, exists = config.Apps[canonical]
			if !exists {
				return fmt.Errorf("alias '%s' points to unknown app '%s'", alias, canonical)
			}
		} else {
			return fmt.Errorf("unknown app: %s", alias)
		}
	}

	killPatterns := app.GetKillPatterns()
	if len(killPatterns) == 0 {
		return fmt.Errorf("no kill patterns available for %s", alias)
	}

	// Try each kill pattern and kill all matching processes
	killed := false
	for _, pattern := range killPatterns {
		if err := killAllByPattern(pattern); err == nil {
			fmt.Printf("Killed all processes matching: %s\n", pattern)
			killed = true
		}
	}

	if !killed {
		fmt.Printf("No running processes found for: %s\n", alias)
	}

	return nil
}

// killAllByPattern kills all processes matching the given pattern
func killAllByPattern(pattern string) error {
	switch runtime.GOOS {
	case "darwin":
		return killAllMacOS(pattern)
	case "linux":
		return killAllLinux(pattern)
	case "windows":
		return killAllWindows(pattern)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// killAllMacOS kills all processes on macOS matching the pattern
func killAllMacOS(pattern string) error {
	// For macOS apps, try graceful quit first for GUI apps
	if err := quitMacOSApp(pattern); err == nil {
		// After graceful quit, check if any processes are still running
		// and force kill them if needed
		if isProcessRunning(pattern) {
			return exec.Command("pkill", "-f", pattern).Run()
		}
		return nil
	}

	// If graceful quit failed, force kill all matching processes
	return exec.Command("pkill", "-f", pattern).Run()
}

// quitMacOSApp tries to quit an app gracefully via AppleScript
func quitMacOSApp(appName string) error {
	// First try to quit all instances of the app gracefully
	script := fmt.Sprintf(`
		tell application "System Events"
			set appList to (name of every application process whose name contains "%s")
			repeat with appProcess in appList
				try
					tell application appProcess to quit
				end try
			end repeat
		end tell`, appName)
	return exec.Command("osascript", "-e", script).Run()
}

// killAllLinux kills all processes on Linux matching the pattern
func killAllLinux(pattern string) error {
	return exec.Command("pkill", "-f", pattern).Run()
}

// killAllWindows kills all processes on Windows matching the pattern
func killAllWindows(pattern string) error {
	// Try with .exe extension first - use /F to force kill all processes
	if err := exec.Command("taskkill", "/F", "/IM", pattern+".exe").Run(); err == nil {
		return nil
	}

	// Try without .exe extension - use /F to force kill all processes
	return exec.Command("taskkill", "/F", "/IM", pattern).Run()
}

// closeMultipleApps closes multiple applications
func closeMultipleApps(aliases []string) error {
	errors := 0
	for _, alias := range aliases {
		if err := CloseApp(alias); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing %s: %v\n", alias, err)
			errors++
		}
	}

	if errors > 0 {
		return fmt.Errorf("%d apps failed to close", errors)
	}

	return nil
}

// isProcessRunning checks if a process matching the pattern is running
func isProcessRunning(pattern string) bool {
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd := exec.Command("pgrep", "-f", pattern)
		return cmd.Run() == nil
	case "windows":
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s*", pattern))
		output, err := cmd.Output()
		return err == nil && strings.Contains(string(output), pattern)
	default:
		return false
	}
}
