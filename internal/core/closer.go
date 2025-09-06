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

	// Try each kill pattern until one works
	killed := false
	for _, pattern := range killPatterns {
		if err := killByPattern(pattern); err == nil {
			fmt.Printf("Killed processes matching: %s\n", pattern)
			killed = true
			break
		}
	}

	if !killed {
		fmt.Printf("No running processes found for: %s\n", alias)
	}

	return nil
}

// killByPattern kills processes matching the given pattern
func killByPattern(pattern string) error {
	switch runtime.GOOS {
	case "darwin":
		return killMacOS(pattern)
	case "linux":
		return killLinux(pattern)
	case "windows":
		return killWindows(pattern)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// killMacOS kills processes on macOS
func killMacOS(pattern string) error {
	// Try graceful quit via AppleScript first
	if err := quitMacOSApp(pattern); err == nil {
		return nil
	}

	// Fallback to pkill
	return exec.Command("pkill", "-f", pattern).Run()
}

// quitMacOSApp tries to quit an app gracefully via AppleScript
func quitMacOSApp(appName string) error {
	script := fmt.Sprintf(`tell application "%s" to quit`, appName)
	return exec.Command("osascript", "-e", script).Run()
}

// killLinux kills processes on Linux
func killLinux(pattern string) error {
	return exec.Command("pkill", "-f", pattern).Run()
}

// killWindows kills processes on Windows
func killWindows(pattern string) error {
	// Try with .exe extension first
	if err := exec.Command("taskkill", "/F", "/IM", pattern+".exe").Run(); err == nil {
		return nil
	}

	// Try without .exe extension
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
