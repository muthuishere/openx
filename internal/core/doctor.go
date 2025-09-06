package core

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorGray   = "\033[90m"
)

// DoctorReport represents the status of all configured applications
type DoctorReport struct {
	Platform   string            `json:"platform"`
	ConfigPath string            `json:"configPath"`
	Apps       []AppStatus       `json:"apps"`
	Aliases    map[string]string `json:"aliases"`
	Summary    Summary           `json:"summary"`
}

// AppStatus represents the status of a single application
type AppStatus struct {
	Name        string `json:"name"`
	LaunchPath  string `json:"launchPath"`
	Status      string `json:"status"` // "available", "missing", "no-path"
	KillPattern string `json:"killPattern"`
	Running     bool   `json:"running"`
}

// Summary provides aggregate statistics
type Summary struct {
	Total     int `json:"total"`
	Available int `json:"available"`
	Missing   int `json:"missing"`
	Running   int `json:"running"`
}

// RunDoctor performs a health check of all configured applications
func RunDoctor(jsonOutput bool) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath := getConfigPath()
	report := DoctorReport{
		Platform:   runtime.GOOS,
		ConfigPath: configPath,
		Apps:       []AppStatus{},
		Aliases:    config.Aliases,
		Summary:    Summary{},
	}

	// Check each application
	appNames := make([]string, 0, len(config.Apps))
	for name := range config.Apps {
		appNames = append(appNames, name)
	}
	sort.Strings(appNames)

	for _, name := range appNames {
		app := config.Apps[name]
		status := checkAppStatus(name, app)
		report.Apps = append(report.Apps, status)

		// Update summary
		report.Summary.Total++
		switch status.Status {
		case "available":
			report.Summary.Available++
		case "missing":
			report.Summary.Missing++
		}
		if status.Running {
			report.Summary.Running++
		}
	}

	if jsonOutput {
		return outputJSON(report)
	}

	return outputHuman(report)
}

// checkAppStatus checks the status of a single application
func checkAppStatus(name string, app *App) AppStatus {
	status := AppStatus{
		Name:        name,
		KillPattern: strings.Join(app.GetKillPatterns(), ", "),
	}

	// Check if we have a launch path for this platform
	launchPath := app.GetLaunchPath()
	if launchPath == "" {
		status.Status = "no-path"
		status.LaunchPath = fmt.Sprintf("(no path for %s)", runtime.GOOS)
		return status
	}

	status.LaunchPath = launchPath

	// Check if the application exists
	if appExists(launchPath) {
		status.Status = "available"
	} else {
		status.Status = "missing"
	}

	// Check if the application is running
	killPatterns := app.GetKillPatterns()
	for _, pattern := range killPatterns {
		if isProcessRunning(pattern) {
			status.Running = true
			break
		}
	}

	return status
}

// appExists checks if an application exists at the given path
func appExists(path string) bool {
	if strings.ContainsAny(path, `/\`) {
		// Absolute or relative path
		return exists(path)
	}

	// Command in PATH
	_, err := exec.LookPath(path)
	return err == nil
}

// outputJSON outputs the doctor report in JSON format
func outputJSON(report DoctorReport) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// outputHuman outputs the doctor report in human-readable format
func outputHuman(report DoctorReport) error {
	fmt.Printf("openx doctor (%s)\n", report.Platform)
	fmt.Printf("Config: %s\n\n", report.ConfigPath)

	// Applications status
	fmt.Println("Applications:")
	for _, app := range report.Apps {
		status := getStatusIcon(app.Status)
		statusColor := getStatusColor(app.Status)
		running := ""
		if app.Running {
			running = ColorGreen + " (running)" + ColorReset
		}

		fmt.Printf("  %s%s%s %-15s %s%s\n", statusColor, status, ColorReset, app.Name, app.LaunchPath, running)
		if app.KillPattern != "" {
			fmt.Printf("    %s└─ kill: %s%s\n", ColorGray, app.KillPattern, ColorReset)
		}
	}

	// Aliases
	if len(report.Aliases) > 0 {
		fmt.Println("\nAliases:")
		aliasNames := make([]string, 0, len(report.Aliases))
		for alias := range report.Aliases {
			aliasNames = append(aliasNames, alias)
		}
		sort.Strings(aliasNames)

		for _, alias := range aliasNames {
			target := report.Aliases[alias]
			fmt.Printf("  %-10s → %s\n", alias, target)
		}
	}

	// Summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total: %d apps\n", report.Summary.Total)
	fmt.Printf("  %sAvailable: %d%s\n", ColorGreen, report.Summary.Available, ColorReset)
	if report.Summary.Missing > 0 {
		fmt.Printf("  %sMissing: %d%s\n", ColorRed, report.Summary.Missing, ColorReset)
	} else {
		fmt.Printf("  Missing: %d\n", report.Summary.Missing)
	}
	if report.Summary.Running > 0 {
		fmt.Printf("  %sRunning: %d%s\n", ColorGreen, report.Summary.Running, ColorReset)
	} else {
		fmt.Printf("  Running: %d\n", report.Summary.Running)
	}

	if report.Summary.Missing > 0 {
		fmt.Printf("\n%sNote: Missing apps may need to be installed or paths updated in config.%s\n", ColorYellow, ColorReset)
	}

	return nil
}

// getStatusIcon returns an icon for the given status
func getStatusIcon(status string) string {
	switch status {
	case "available":
		return "✓"
	case "missing":
		return "✗"
	case "no-path":
		return "○"
	default:
		return "?"
	}
}

// getStatusColor returns the color code for the given status
func getStatusColor(status string) string {
	switch status {
	case "available":
		return ColorGreen
	case "missing":
		return ColorRed
	case "no-path":
		return ColorYellow
	default:
		return ColorReset
	}
}
