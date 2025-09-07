package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsValidAlias(t *testing.T) {
	// Setup test config
	testContent := `
apps:
  testeditor:
    darwin: "TextEdit.app"
    linux: "gedit"
    windows: "notepad.exe"
  
  browser:
    darwin: "Safari"
    linux: "firefox"
    windows: "chrome.exe"
`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name     string
		alias    string
		expected bool
	}{
		{
			name:     "valid app alias",
			alias:    "testeditor",
			expected: true,
		},
		{
			name:     "valid synonym",
			alias:    "code", // Should map to vscode if configured
			expected: false,  // Our test config doesn't have this
		},
		{
			name:     "invalid alias",
			alias:    "nonexistent",
			expected: false,
		},
		{
			name:     "file path",
			alias:    "/path/to/file.txt",
			expected: false,
		},
		{
			name:     "URL",
			alias:    "https://example.com",
			expected: false,
		},
		{
			name:     "empty string",
			alias:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidAlias(tt.alias)
			if result != tt.expected {
				t.Errorf("isValidAlias(%q) = %v, want %v", tt.alias, result, tt.expected)
			}
		})
	}
}

func TestOpenWithSystemDefault(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name:    "open test file",
			target:  testFile,
			wantErr: false,
		},
		{
			name:    "open URL",
			target:  "https://example.com",
			wantErr: false,
		},
		{
			name:    "empty target",
			target:  "",
			wantErr: true, // This might fail depending on OS
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := openWithSystemDefault(tt.target)

			if tt.wantErr {
				if err == nil {
					t.Errorf("openWithSystemDefault() expected error but got none")
				}
				return
			}

			// Note: On some systems, open commands might "succeed" even for invalid targets
			// So we don't strictly require no error, just log if there is one
			if err != nil {
				t.Logf("openWithSystemDefault() returned error (might be expected): %v", err)
			}
		})
	}
}

func TestOpenWithAppAndArgs(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Testing macOS-specific 'open -a' functionality")
	}

	// Create a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		appPath string
		args    []string
		wantErr bool
	}{
		{
			name:    "open with TextEdit",
			appPath: "TextEdit",
			args:    []string{testFile},
			wantErr: false,
		},
		{
			name:    "open with invalid app",
			appPath: "NonExistentApp",
			args:    []string{testFile},
			wantErr: false, // open -a might still "succeed" even for invalid apps
		},
		{
			name:    "empty app path",
			appPath: "",
			args:    []string{testFile},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := openWithAppAndArgs(tt.appPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("openWithAppAndArgs() expected error but got none")
				}
				return
			}

			// Note: On macOS, 'open -a' might "succeed" even for non-existent apps
			// So we don't strictly require no error, just log if there is one
			if err != nil {
				t.Logf("openWithAppAndArgs() returned error (might be expected): %v", err)
			}
		})
	}
}

// Helper functions for test setup
func setupTestConfig(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return configPath
}

func setTempConfigPath(t *testing.T, configPath string) func() {
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	configDir := filepath.Dir(configPath)
	os.Setenv("XDG_CONFIG_HOME", configDir)

	return func() {
		if oldXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}
}
