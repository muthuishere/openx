package core

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// setupTestConfig creates a temporary config file for testing
func setupTestConfig(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create directory
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create test config directory: %v", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	return configPath
}

// setTempConfigPath temporarily overrides the config path for testing
func setTempConfigPath(t *testing.T, path string) func() {
	// Override getConfigPath for testing by setting environment variable
	oldHome := os.Getenv("HOME")
	oldXDG := os.Getenv("XDG_CONFIG_HOME")

	// Set XDG_CONFIG_HOME to control where config is looked for
	tempDir := filepath.Dir(filepath.Dir(path))
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	return func() {
		if oldXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
		os.Setenv("HOME", oldHome)
	}
}

func TestConfig_Success(t *testing.T) {
	testContent := `apps:
  code:
    darwin: /Applications/Visual Studio Code.app
    linux: /usr/bin/code
    windows: C:\Users\%USERNAME%\AppData\Local\Programs\Microsoft VS Code\Code.exe
    kill:
      - "Visual Studio Code"
      - "code"
  chrome:
    darwin: /Applications/Google Chrome.app
    linux: /usr/bin/google-chrome
    windows: C:\Program Files\Google\Chrome\Application\chrome.exe
  firefox:
    darwin: /Applications/Firefox.app
    linux: /usr/bin/firefox
    kill:
      - "firefox"
      - "Firefox"
aliases:
  vs: code
  gc: chrome
  ff: firefox`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	// Test apps are loaded correctly
	if config.Apps == nil {
		t.Fatal("Apps map is nil")
	}

	if len(config.Apps) != 3 {
		t.Errorf("Expected 3 apps, got %d", len(config.Apps))
	}

	// Test code app
	codeApp, exists := config.Apps["code"]
	if !exists {
		t.Fatal("Code app not found")
	}

	expectedCodePath := ""
	switch runtime.GOOS {
	case "darwin":
		expectedCodePath = "/Applications/Visual Studio Code.app"
	case "linux":
		expectedCodePath = "/usr/bin/code"
	case "windows":
		expectedCodePath = "C:\\Users\\%USERNAME%\\AppData\\Local\\Programs\\Microsoft VS Code\\Code.exe"
	}

	if codeApp.GetLaunchPath() != expectedCodePath {
		t.Errorf("Expected code path %s, got %s", expectedCodePath, codeApp.GetLaunchPath())
	}

	// Test explicit kill patterns
	codeKillPatterns := codeApp.GetKillPatterns()
	expectedKillPatterns := []string{"Visual Studio Code", "code"}
	if len(codeKillPatterns) != len(expectedKillPatterns) {
		t.Errorf("Expected %d kill patterns, got %d", len(expectedKillPatterns), len(codeKillPatterns))
	}
	for i, pattern := range expectedKillPatterns {
		if i >= len(codeKillPatterns) || codeKillPatterns[i] != pattern {
			t.Errorf("Expected kill pattern %s at index %d, got %s", pattern, i, codeKillPatterns[i])
		}
	}

	// Test aliases
	if config.Aliases == nil {
		t.Fatal("Aliases map is nil")
	}

	if len(config.Aliases) != 3 {
		t.Errorf("Expected 3 aliases, got %d", len(config.Aliases))
	}

	if config.Aliases["vs"] != "code" {
		t.Errorf("Expected alias 'vs' to point to 'code', got %s", config.Aliases["vs"])
	}
}

func TestApp_GetLaunchPath_Success(t *testing.T) {
	tests := []struct {
		name     string
		app      *App
		expected string
	}{
		{
			name: "current OS path exists",
			app: &App{
				Paths: map[string]string{
					runtime.GOOS: "/test/app/path",
					"other":      "/other/path",
				},
			},
			expected: "/test/app/path",
		},
		{
			name: "tilde expansion",
			app: &App{
				Paths: map[string]string{
					runtime.GOOS: "~/Applications/Test.app",
				},
			},
			expected: func() string {
				home, _ := os.UserHomeDir()
				return filepath.Join(home, "Applications", "Test.app")
			}(),
		},
		{
			name: "no path for current OS",
			app: &App{
				Paths: map[string]string{
					"other": "/other/path",
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.app.GetLaunchPath()
			if result != tt.expected {
				t.Errorf("GetLaunchPath() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestApp_GetKillPatterns_Success(t *testing.T) {
	tests := []struct {
		name     string
		app      *App
		expected []string
	}{
		{
			name: "explicit kill patterns",
			app: &App{
				Paths: map[string]string{
					runtime.GOOS: "/Applications/Test.app",
				},
				Kill: []string{"Test App", "test"},
			},
			expected: []string{"Test App", "test"},
		},
		{
			name: "derived from macOS app bundle",
			app: &App{
				Paths: map[string]string{
					"darwin": "/Applications/Visual Studio Code.app",
				},
			},
			expected: func() []string {
				if runtime.GOOS == "darwin" {
					return []string{"Code"} // Uses processNameExceptions
				}
				return []string{}
			}(),
		},
		{
			name: "derived from regular executable",
			app: &App{
				Paths: map[string]string{
					runtime.GOOS: "/usr/bin/firefox",
				},
			},
			expected: []string{"firefox"},
		},
		{
			name: "no launch path",
			app: &App{
				Paths: map[string]string{
					"other": "/other/path",
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.app.GetKillPatterns()
			if len(result) != len(tt.expected) {
				t.Errorf("GetKillPatterns() returned %d patterns, expected %d", len(result), len(tt.expected))
			}
			for i, pattern := range tt.expected {
				if i >= len(result) || result[i] != pattern {
					t.Errorf("GetKillPatterns()[%d] = %s, expected %s", i, result[i], pattern)
				}
			}
		})
	}
}

func TestSaveConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "openx", "config.yaml")

	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	config := &Config{
		Apps: map[string]*App{
			"test": {
				Paths: map[string]string{
					"darwin": "/Applications/Test.app",
					"linux":  "/usr/bin/test",
				},
				Kill: []string{"test", "Test"},
			},
		},
		Aliases: map[string]string{
			"t": "test",
		},
	}

	err := saveConfig(config)
	if err != nil {
		t.Fatalf("saveConfig() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Verify content by loading it back
	loadedConfig, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// Check apps
	if len(loadedConfig.Apps) != 1 {
		t.Errorf("Expected 1 app, got %d", len(loadedConfig.Apps))
	}

	testApp, exists := loadedConfig.Apps["test"]
	if !exists {
		t.Fatal("Test app not found in loaded config")
	}

	if testApp.Paths["darwin"] != "/Applications/Test.app" {
		t.Errorf("Expected darwin path '/Applications/Test.app', got %s", testApp.Paths["darwin"])
	}

	// Check aliases
	if len(loadedConfig.Aliases) != 1 {
		t.Errorf("Expected 1 alias, got %d", len(loadedConfig.Aliases))
	}

	if loadedConfig.Aliases["t"] != "test" {
		t.Errorf("Expected alias 't' to point to 'test', got %s", loadedConfig.Aliases["t"])
	}
}

func TestProcessNameExceptions_Success(t *testing.T) {
	tests := []struct {
		appName  string
		expected string
	}{
		{"Visual Studio Code", "Code"},
		{"Android Studio", "studio"},
		{"IntelliJ IDEA", "idea"},
		{"Unknown App", ""}, // Not in exceptions map
	}

	for _, tt := range tests {
		t.Run(tt.appName, func(t *testing.T) {
			result := processNameExceptions[tt.appName]
			if result != tt.expected {
				t.Errorf("processNameExceptions[%s] = %s, expected %s", tt.appName, result, tt.expected)
			}
		})
	}
}

// TestLoadConfig_E2E tests the entire config loading process end-to-end
func TestLoadConfig_E2E_Success(t *testing.T) {
	// Create a realistic config file
	realisticConfig := `apps:
  code:
    darwin: /Applications/Visual Studio Code.app
    linux: /usr/bin/code
    windows: C:\Users\%USERNAME%\AppData\Local\Programs\Microsoft VS Code\Code.exe
    kill:
      - "Visual Studio Code"
      - "code"
  chrome:
    darwin: /Applications/Google Chrome.app
    linux: /usr/bin/google-chrome
    windows: C:\Program Files\Google\Chrome\Application\chrome.exe
  firefox:
    darwin: /Applications/Firefox.app
    linux: /usr/bin/firefox
    windows: C:\Program Files\Mozilla Firefox\firefox.exe
    kill:
      - "firefox"
  postman:
    darwin: /Applications/Postman.app
    linux: /snap/bin/postman
    windows: C:\Users\%USERNAME%\AppData\Local\Postman\Postman.exe
  docker:
    darwin: /Applications/Docker.app
    linux: /usr/bin/docker-desktop
    windows: C:\Program Files\Docker\Docker\Docker Desktop.exe
aliases:
  vs: code
  vscode: code
  gc: chrome
  ff: firefox
  pm: postman`

	configPath := setupTestConfig(t, realisticConfig)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	// Load the config
	config, err := loadConfig()
	if err != nil {
		t.Fatalf("E2E loadConfig() failed: %v", err)
	}

	// Test that all expected apps are present
	expectedApps := []string{"code", "chrome", "firefox", "postman", "docker"}
	for _, appName := range expectedApps {
		if _, exists := config.Apps[appName]; !exists {
			t.Errorf("Expected app %s not found", appName)
		}
	}

	// Test launch paths for current OS work
	for appName, app := range config.Apps {
		launchPath := app.GetLaunchPath()
		if launchPath != "" {
			// Should not be empty if configured for current OS
			t.Logf("App %s launch path: %s", appName, launchPath)
		}
	}

	// Test kill patterns are generated
	for appName, app := range config.Apps {
		killPatterns := app.GetKillPatterns()
		if len(killPatterns) == 0 && app.GetLaunchPath() != "" {
			t.Errorf("App %s has launch path but no kill patterns", appName)
		}
		if len(killPatterns) > 0 {
			t.Logf("App %s kill patterns: %v", appName, killPatterns)
		}
	}

	// Test aliases work
	expectedAliases := map[string]string{
		"vs":     "code",
		"vscode": "code",
		"gc":     "chrome",
		"ff":     "firefox",
		"pm":     "postman",
	}

	for alias, expectedTarget := range expectedAliases {
		if target, exists := config.Aliases[alias]; !exists {
			t.Errorf("Expected alias %s not found", alias)
		} else if target != expectedTarget {
			t.Errorf("Alias %s points to %s, expected %s", alias, target, expectedTarget)
		}
	}

	// Test that we can save and reload the same config
	err = saveConfig(config)
	if err != nil {
		t.Fatalf("E2E saveConfig() failed: %v", err)
	}

	reloadedConfig, err := loadConfig()
	if err != nil {
		t.Fatalf("E2E reload config failed: %v", err)
	}

	// Compare key aspects
	if len(reloadedConfig.Apps) != len(config.Apps) {
		t.Errorf("Reloaded config has %d apps, original had %d", len(reloadedConfig.Apps), len(config.Apps))
	}

	if len(reloadedConfig.Aliases) != len(config.Aliases) {
		t.Errorf("Reloaded config has %d aliases, original had %d", len(reloadedConfig.Aliases), len(config.Aliases))
	}
}

// TestDeriveKillPatterns_E2E tests the kill pattern derivation for different OS scenarios
func TestDeriveKillPatterns_E2E_Success(t *testing.T) {
	tests := []struct {
		name        string
		osType      string
		launchPath  string
		expected    []string
		description string
	}{
		{
			name:        "macOS app bundle with exception",
			osType:      "darwin",
			launchPath:  "/Applications/Visual Studio Code.app",
			expected:    []string{"Code"},
			description: "Uses processNameExceptions mapping",
		},
		{
			name:        "macOS app bundle without exception",
			osType:      "darwin",
			launchPath:  "/Applications/MyApp.app",
			expected:    []string{"MyApp"},
			description: "Strips .app extension",
		},
		{
			name:        "Windows executable",
			osType:      "windows",
			launchPath:  "C:\\Program Files\\App\\app.exe",
			expected:    []string{"app"},
			description: "Strips .exe extension",
		},
		{
			name:        "Linux executable",
			osType:      "linux",
			launchPath:  "/usr/bin/firefox",
			expected:    []string{"firefox"},
			description: "Uses basename as-is",
		},
	}

	originalGOOS := runtime.GOOS

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't actually change runtime.GOOS in tests,
			// but we can test the logic with the current OS

			app := &App{
				Paths: map[string]string{
					tt.osType: tt.launchPath,
				},
			}

			// If this is the current OS, test the actual derivation
			if tt.osType == originalGOOS {
				patterns := app.deriveKillPatterns()
				t.Logf("OS: %s, Path: %s, Patterns: %v", tt.osType, tt.launchPath, patterns)

				// For the current OS, verify the pattern makes sense
				if len(patterns) == 0 {
					t.Errorf("Expected non-empty kill patterns for %s", tt.launchPath)
				}
			}
		})
	}
}
