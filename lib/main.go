package lib

import (
	_ "embed"
	"fmt"
	"openx/internal/core"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed versions.txt
var versionData string

// OpenX represents the main library interface for managing applications
type OpenX struct {
	configPath string
}

// New creates a new OpenX instance with the default config location
func New() *OpenX {
	return &OpenX{}
}

// NewWithConfig creates a new OpenX instance with a custom config file path
func NewWithConfig(configPath string) *OpenX {
	return &OpenX{
		configPath: configPath,
	}
}

// EnsureConfig ensures that the configuration file exists and is properly set up
func (ox *OpenX) EnsureConfig() error {
	return core.EnsureConfig()
}

// RunAlias runs an application by alias with optional arguments
func (ox *OpenX) RunAlias(alias string, args ...string) error {
	return core.LaunchApp(alias, args)
}

// RunDirect runs an application by direct path with optional arguments
func (ox *OpenX) RunDirect(path string, args ...string) error {
	return ox.executeDirectPath(path, args...)
}

// Kill terminates an application by alias
func (ox *OpenX) Kill(alias string) error {
	return core.CloseApp(alias)
}

// AddAlias adds a new alias to the configuration
func (ox *OpenX) AddAlias(alias, appName string) error {
	config, err := ox.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if the app exists in the configuration
	if _, exists := config.Apps[appName]; !exists {
		return fmt.Errorf("application '%s' is not configured", appName)
	}

	// Add the alias
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}
	config.Aliases[alias] = appName

	return ox.saveConfig(config)
}

// RemoveAlias removes an alias from the configuration
func (ox *OpenX) RemoveAlias(alias string) error {
	config, err := ox.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.Aliases == nil {
		return fmt.Errorf("alias '%s' not found", alias)
	}

	if _, exists := config.Aliases[alias]; !exists {
		return fmt.Errorf("alias '%s' not found", alias)
	}

	delete(config.Aliases, alias)

	return ox.saveConfig(config)
}

// ListAliases returns a map of all configured aliases
func (ox *OpenX) ListAliases() (map[string]string, error) {
	config, err := ox.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if config.Aliases == nil {
		return make(map[string]string), nil
	}

	// Return a copy to prevent external modification
	aliases := make(map[string]string)
	for k, v := range config.Aliases {
		aliases[k] = v
	}

	return aliases, nil
}

// Doctor performs a health check on all configured applications
func (ox *OpenX) Doctor() error {
	return core.RunDoctor(false)
}

// DoctorJSON performs a health check and returns results in JSON format
func (ox *OpenX) DoctorJSON() error {
	return core.RunDoctor(true)
}

// Helper methods for internal use

// loadConfig loads the configuration from the default location
func (ox *OpenX) loadConfig() (*core.Config, error) {
	// Use the core package's internal loadConfig through EnsureConfig
	if err := core.EnsureConfig(); err != nil {
		return nil, err
	}

	// Read the config file directly
	configPath := ox.getConfigPath()

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config core.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// saveConfig saves the configuration to the default location
func (ox *OpenX) saveConfig(config *core.Config) error {
	configPath := ox.getConfigPath()

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(config)
}

// getConfigPath returns the configuration file path
func (ox *OpenX) getConfigPath() string {
	if ox.configPath != "" {
		return ox.configPath
	}

	// Use XDG config directory or fallback to home directory
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "openx", "config.yaml")
}

// executeDirectPath executes an application by direct path
func (ox *OpenX) executeDirectPath(appPath string, args ...string) error {
	// Expand path if it starts with ~
	if strings.HasPrefix(appPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		appPath = filepath.Join(homeDir, appPath[2:])
	}

	// Check if the path exists
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return fmt.Errorf("application not found: %s", appPath)
	}

	// For macOS .app bundles, we need special handling
	if runtime.GOOS == "darwin" && strings.HasSuffix(appPath, ".app") {
		return ox.launchMacOSApp(appPath, args)
	}

	// For regular executables
	cmd := exec.Command(appPath, args...)
	return cmd.Start()
}

// launchMacOSApp launches a macOS .app bundle
func (ox *OpenX) launchMacOSApp(appPath string, args []string) error {
	// Try to find the executable inside the .app bundle
	executablePath := filepath.Join(appPath, "Contents", "MacOS")

	entries, err := os.ReadDir(executablePath)
	if err != nil {
		// Fallback to using 'open' command
		return ox.launchWithOpen(appPath, args)
	}

	// Find the main executable
	for _, entry := range entries {
		if !entry.IsDir() {
			execPath := filepath.Join(executablePath, entry.Name())
			if info, err := entry.Info(); err == nil && info.Mode()&0111 != 0 {
				cmd := exec.Command(execPath, args...)
				return cmd.Start()
			}
		}
	}

	// Fallback to using 'open' command
	return ox.launchWithOpen(appPath, args)
}

// launchWithOpen uses macOS 'open' command to launch an application
func (ox *OpenX) launchWithOpen(appPath string, args []string) error {
	openArgs := []string{appPath}
	if len(args) > 0 {
		openArgs = append(openArgs, "--args")
		openArgs = append(openArgs, args...)
	}

	cmd := exec.Command("open", openArgs...)
	return cmd.Start()
}

// Version information
const (
	Name = "OpenX"
)

// GetVersion returns the library version from embedded versions.txt
func GetVersion() string {
	return strings.TrimSpace(versionData)
}

// GetName returns the library name
func GetName() string {
	return Name
}
