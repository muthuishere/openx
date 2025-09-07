package core

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// LaunchApp launches an application with the given arguments
func LaunchApp(alias string, args []string) error {
	// Check if it's a direct path to an application
	if isDirectPath(alias) {
		return launchDirectPath(alias, args)
	}

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

	launchPath := app.GetLaunchPath()
	if launchPath == "" {
		return fmt.Errorf("no launch path configured for %s on %s", alias, runtime.GOOS)
	}

	// Resolve and prepare arguments
	resolvedArgs := resolveTargets(args)

	// Launch the application
	if err := executeApp(launchPath, resolvedArgs); err != nil {
		return fmt.Errorf("failed to launch %s: %w", alias, err)
	}

	fmt.Printf("Launched: %s\n", alias)
	if len(args) > 0 {
		fmt.Printf("Arguments: %v\n", args)
	}

	return nil
}

// executeApp handles the actual launching of the application
func executeApp(launchPath string, args []string) error {
	// Handle macOS .app bundles
	if runtime.GOOS == "darwin" {

		return launchMacOSApp(launchPath, args)
	}

	// Handle regular executables
	cmd := exec.Command(launchPath, args...)
	return cmd.Start()
}

// launchMacOSApp launches a macOS .app bundle
func launchMacOSApp(appPath string, args []string) error {
	// Find the actual executable inside the .app bundle
	execPath, err := findAppExecutable(appPath)
	if err != nil {
		// Fallback to using 'open' command
		return launchWithOpen(appPath, args)
	}

	// Launch the executable directly
	cmd := exec.Command(execPath, args...)
	return cmd.Start()
}

// launchWithOpen uses macOS 'open' command as fallback
func launchWithOpen(appPath string, args []string) error {
	openArgs := []string{"-a", appPath}
	if len(args) > 0 {
		// openArgs = append(openArgs, "--args")
		openArgs = append(openArgs, args...)
	}
	fmt.Printf("Using 'open' command: open %s\n", strings.Join(openArgs, " "))

	cmd := exec.Command("open", openArgs...)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error with 'open -a %s': %v\n", appPath, err)
		return fmt.Errorf("failed to launch %s with 'open' command: %w", appPath, err)
	}

	fmt.Printf("Successfully launched with 'open -a %s'\n", appPath)
	return nil
}

// launchMultipleApps launches multiple applications
func launchMultipleApps(aliases []string) error {
	errors := 0
	for _, alias := range aliases {
		if err := LaunchApp(alias, []string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Error launching %s: %v\n", alias, err)
			errors++
		}
	}

	if errors > 0 {
		return fmt.Errorf("%d apps failed to launch", errors)
	}

	return nil
}

// isDirectPath checks if the given string is a direct path to an application
func isDirectPath(path string) bool {
	// Check if it contains path separators
	if strings.Contains(path, "/") || strings.Contains(path, "\\") {
		return true
	}
	return false
}

// launchDirectPath launches an application using a direct path
func launchDirectPath(appPath string, args []string) error {
	// Check if the application exists
	if !exists(appPath) {
		return fmt.Errorf("application not found: %s", appPath)
	}

	// Resolve and prepare arguments
	resolvedArgs := resolveTargets(args)

	// Launch the application
	if err := executeApp(appPath, resolvedArgs); err != nil {
		return fmt.Errorf("failed to launch %s: %w", appPath, err)
	}

	fmt.Printf("Launched: %s\n", appPath)
	if len(args) > 0 {
		fmt.Printf("Arguments: %v\n", args)
	}

	return nil
}
