package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the entire configuration
type Config struct {
	Apps    map[string]*App   `yaml:"apps"`
	Aliases map[string]string `yaml:"aliases"`
}

// App represents a single application configuration
type App struct {
	Paths map[string]string `yaml:",inline"`
	Kill  []string          `yaml:"kill,omitempty"`
}

// GetLaunchPath returns the launch path for the current OS
func (a *App) GetLaunchPath() string {
	osKey := runtime.GOOS

	// Check direct OS key first
	if path, ok := a.Paths[osKey]; ok && path != "" {
		return expandTilde(path)
	}

	return ""
}

// GetKillPatterns returns the kill patterns for this app
func (a *App) GetKillPatterns() []string {
	// If explicitly specified, use those
	if len(a.Kill) > 0 {
		return a.Kill
	}

	// Otherwise, derive from launch path
	return a.deriveKillPatterns()
}

// deriveKillPatterns derives kill patterns from the launch path
func (a *App) deriveKillPatterns() []string {
	launchPath := a.GetLaunchPath()
	if launchPath == "" {
		return []string{}
	}

	baseName := filepath.Base(launchPath)

	switch runtime.GOOS {
	case "darwin":
		if strings.HasSuffix(baseName, ".app") {
			appName := strings.TrimSuffix(baseName, ".app")
			// Handle known exceptions
			if mapped := processNameExceptions[appName]; mapped != "" {
				return []string{mapped}
			}
			return []string{appName}
		}
		return []string{baseName}
	case "windows":
		if strings.HasSuffix(baseName, ".exe") {
			return []string{strings.TrimSuffix(baseName, ".exe")}
		}
		return []string{baseName}
	case "linux":
		return []string{baseName}
	default:
		return []string{baseName}
	}
}

// processNameExceptions maps app bundle names to actual process names
var processNameExceptions = map[string]string{
	"Visual Studio Code": "Code",
	"Android Studio":     "studio",
	"IntelliJ IDEA":      "idea",
}

// loadConfig loads the configuration from file
func loadConfig() (*Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found at %s (run 'openx doctor' to create it)", configPath)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize empty maps if not present
	if config.Apps == nil {
		config.Apps = make(map[string]*App)
	}
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}

	return &config, nil
}

// saveConfig saves the configuration to file
func saveConfig(config *Config) error {
	configPath := getConfigPath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
