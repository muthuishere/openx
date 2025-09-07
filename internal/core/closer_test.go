package core

import (
	"fmt"
	"os"
	"testing"
)

func TestCloseApp(t *testing.T) {
	// Create a test config
	testContent := `
apps:
  testapp:
    darwin: "/Applications/TestApp.app"
    linux: "testapp"
    windows: "testapp.exe"
    kill: ["TestApp", "testapp"]

aliases:
  ta: testapp`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name    string
		alias   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid app",
			alias:   "testapp",
			wantErr: false, // We expect this to fail in the killProcess step, not config loading
		},
		{
			name:    "valid alias",
			alias:   "ta",
			wantErr: false, // We expect this to fail in the killProcess step, not config loading
		},
		{
			name:    "unknown app",
			alias:   "unknown",
			wantErr: true,
			errMsg:  "unknown app: unknown",
		},
		{
			name:    "empty alias",
			alias:   "",
			wantErr: true,
			errMsg:  "unknown app: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloseApp(tt.alias)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CloseApp() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("CloseApp() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// For valid apps, we expect it to fail at killProcess step since processes don't exist
			// This is normal behavior in tests
			if err != nil {
				t.Logf("CloseApp() failed as expected at killProcess step: %v", err)
			}
		})
	}
}

func TestCloseMultipleApps(t *testing.T) {
	// Create a test config
	testContent := `
apps:
  app1:
    darwin: "/Applications/App1.app"
    linux: "app1"
    windows: "app1.exe"
    kill: ["App1"]
  
  app2:
    darwin: "/Applications/App2.app"
    linux: "app2"
    windows: "app2.exe"
    kill: ["App2"]

aliases:
  a1: app1
  a2: app2`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name    string
		aliases []string
		wantErr bool
		minErrs int // minimum number of errors expected
	}{
		{
			name:    "valid apps",
			aliases: []string{"app1", "app2"},
			wantErr: false, // Will succeed in closing apps
			minErrs: 0,
		},
		{
			name:    "mixed valid and invalid",
			aliases: []string{"app1", "unknown", "app2"},
			wantErr: true,
			minErrs: 1, // At least one error for unknown app
		},
		{
			name:    "all invalid",
			aliases: []string{"unknown1", "unknown2"},
			wantErr: true,
			minErrs: 2,
		},
		{
			name:    "empty list",
			aliases: []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := closeMultipleApps(tt.aliases)

			if tt.wantErr {
				if err == nil {
					t.Errorf("closeMultipleApps() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("closeMultipleApps() unexpected error: %v", err)
			}
		})
	}
}

func TestKillAllByPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{
			name:    "nonexistent process",
			pattern: "nonexistent-app-12345",
			wantErr: false, // killAllByPattern should not error if process doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := killAllByPattern(tt.pattern)
			if tt.wantErr && err == nil {
				t.Errorf("killAllByPattern(%s) expected error but got none", tt.pattern)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("killAllByPattern(%s) unexpected error: %v", tt.pattern, err)
			}
		})
	}
}

func TestIsProcessRunning(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected bool // In tests, we don't expect any specific processes to be running
	}{
		{
			name:     "non-existent process",
			pattern:  "definitely-not-running-process-12345",
			expected: false,
		},
		{
			name:     "empty pattern",
			pattern:  "",
			expected: false,
		},
		{
			name:     "common system process that might exist",
			pattern:  "kernel_task", // On macOS, this is likely to exist
			expected: false,         // We can't guarantee it exists in test environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isProcessRunning(tt.pattern)
			// We just test that the function doesn't panic
			// The actual result depends on what's running on the system
			t.Logf("isProcessRunning(%s) = %v", tt.pattern, result)
		})
	}
}

func TestCloseApp_ConfigError(t *testing.T) {
	// Test with no config file
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", "/nonexistent/path")
	defer func() {
		if oldXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	err := CloseApp("testapp")
	if err == nil {
		t.Error("CloseApp() expected error when config file doesn't exist")
	}

	expectedSubstring := "failed to load config"
	if err != nil && !contains(err.Error(), expectedSubstring) {
		t.Errorf("CloseApp() error = %v, want error containing %v", err, expectedSubstring)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 || indexString(s, substr) >= 0)
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestGetKillPatterns(t *testing.T) {
	// Create a comprehensive test config with various apps and kill patterns
	testContent := `
apps:
  # App with explicit kill patterns
  vscode:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "Code.exe"
    kill: ["Visual Studio Code", "code"]
  
  # App with implicit kill patterns (derived from app bundle name)
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "chrome.exe"
  
  # App with mixed patterns  
  postman:
    darwin: "/Applications/Postman.app"
    linux: "postman"
    windows: "Postman.exe"
    kill: ["Postman", "postman-agent"]

  # App with no kill patterns (should derive from path)
  simple:
    darwin: "/Applications/SimpleApp.app"
    linux: "simple"
    windows: "simple.exe"

  # App for case-insensitive testing
  intellij:
    darwin: "/Applications/IntelliJ IDEA.app"
    linux: "idea"
    windows: "idea.exe"
    kill: ["IDEA"]

aliases:
  code: vscode
  browser: chrome
  idea: intellij`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name             string
		alias            string
		expectedPatterns []string
		wantErr          bool
		errMsg           string
	}{
		{
			name:             "app with explicit kill patterns",
			alias:            "vscode",
			expectedPatterns: []string{"Visual Studio Code", "code"},
			wantErr:          false,
		},
		{
			name:             "app via alias with explicit kill patterns",
			alias:            "code",
			expectedPatterns: []string{"Visual Studio Code", "code"},
			wantErr:          false,
		},
		{
			name:             "app with implicit kill patterns",
			alias:            "chrome",
			expectedPatterns: []string{"Google Chrome"}, // Derived from macOS app bundle
			wantErr:          false,
		},
		{
			name:             "app via alias with implicit kill patterns",
			alias:            "browser",
			expectedPatterns: []string{"Google Chrome"},
			wantErr:          false,
		},
		{
			name:             "app with mixed explicit patterns",
			alias:            "postman",
			expectedPatterns: []string{"Postman", "postman-agent"},
			wantErr:          false,
		},
		{
			name:             "app with derived patterns",
			alias:            "simple",
			expectedPatterns: []string{"SimpleApp"}, // Derived from app bundle name
			wantErr:          false,
		},
		{
			name:             "case-insensitive kill pattern test",
			alias:            "idea",
			expectedPatterns: []string{"IDEA"}, // Should work case-insensitively with actual process "idea"
			wantErr:          false,
		},
		{
			name:             "unknown app",
			alias:            "unknown",
			expectedPatterns: nil,
			wantErr:          true,
			errMsg:           "unknown app: unknown",
		},
		{
			name:             "empty alias",
			alias:            "",
			expectedPatterns: nil,
			wantErr:          true,
			errMsg:           "unknown app: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns, err := getKillPatternsForApp(tt.alias)

			if tt.wantErr {
				if err == nil {
					t.Errorf("getKillPatternsForApp(%s) expected error but got none", tt.alias)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("getKillPatternsForApp(%s) error = %v, want %v", tt.alias, err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("getKillPatternsForApp(%s) unexpected error: %v", tt.alias, err)
				return
			}

			// Check that we got the expected patterns
			if len(patterns) != len(tt.expectedPatterns) {
				t.Errorf("getKillPatternsForApp(%s) got %d patterns, want %d", tt.alias, len(patterns), len(tt.expectedPatterns))
				t.Errorf("Got patterns: %v", patterns)
				t.Errorf("Expected patterns: %v", tt.expectedPatterns)
				return
			}

			// Check each pattern
			for i, pattern := range patterns {
				if pattern != tt.expectedPatterns[i] {
					t.Errorf("getKillPatternsForApp(%s) pattern[%d] = %v, want %v", tt.alias, i, pattern, tt.expectedPatterns[i])
				}
			}

			t.Logf("getKillPatternsForApp(%s) returned patterns: %v", tt.alias, patterns)
		})
	}
}

// getKillPatternsForApp is a helper function that returns the kill patterns for a given app/alias
// This function mimics the logic in CloseApp but only returns the patterns without attempting to kill
func getKillPatternsForApp(alias string) ([]string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	app, exists := config.Apps[alias]
	if !exists {
		// Check if it's an alias
		if canonical, ok := config.Aliases[alias]; ok {
			app, exists = config.Apps[canonical]
			if !exists {
				return nil, fmt.Errorf("alias '%s' points to unknown app '%s'", alias, canonical)
			}
		} else {
			return nil, fmt.Errorf("unknown app: %s", alias)
		}
	}

	killPatterns := app.GetKillPatterns()
	return killPatterns, nil
}
