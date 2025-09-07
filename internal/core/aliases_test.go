package core

import (
	"runtime"
	"testing"

	"openx/shared/config"
)

func TestNewAliasResolver(t *testing.T) {
	// Create a mock config for testing
	mockConfig := &config.Config{
		Apps: map[string]*config.App{
			"vscode": {
				Paths: map[string]string{
					"darwin":  "Visual Studio Code.app",
					"linux":   "code",
					"windows": "Code.exe",
				},
			},
		},
	}

	resolver := newAliasResolver(mockConfig)

	if resolver == nil {
		t.Fatal("newAliasResolver() returned nil")
	}

	if resolver.config == nil {
		t.Fatal("config is nil")
	}

	if resolver.synonyms == nil {
		t.Fatal("synonyms map is nil")
	}

	// Check that some synonyms are initialized
	if len(resolver.synonyms) == 0 {
		t.Error("synonyms map is empty, expected some entries")
	}
}

func TestAliasResolver_Resolve(t *testing.T) {
	// Create a mock config for testing
	mockConfig := &config.Config{
		Apps: map[string]*config.App{
			"vscode": {
				Paths: map[string]string{
					"darwin":  "Visual Studio Code.app",
					"linux":   "code",
					"windows": "Code.exe",
				},
			},
		},
	}

	resolver := newAliasResolver(mockConfig)

	tests := []struct {
		name     string
		alias    string
		expected string
		wantOk   bool
	}{
		{
			name:     "vscode alias",
			alias:    "vscode",
			expected: getExpectedVSCodePath(),
			wantOk:   true,
		},
		{
			name:     "code synonym for vscode",
			alias:    "code",
			expected: getExpectedVSCodePath(),
			wantOk:   true,
		},
		{
			name:     "unknown alias",
			alias:    "nonexistent",
			expected: "",
			wantOk:   false,
		},
		{
			name:     "empty alias",
			alias:    "",
			expected: "",
			wantOk:   false,
		},
		{
			name:     "case insensitive",
			alias:    "VSCODE",
			expected: getExpectedVSCodePath(),
			wantOk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := resolver.Resolve(tt.alias)

			if ok != tt.wantOk {
				t.Errorf("Resolve() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if result != tt.expected {
				t.Errorf("Resolve() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAliasResolver_InitializeSynonyms(t *testing.T) {
	mockConfig := &config.Config{
		Apps: map[string]*config.App{
			"vscode": {
				Paths: map[string]string{
					"darwin":  "Visual Studio Code.app",
					"linux":   "code",
					"windows": "Code.exe",
				},
			},
		},
	}

	resolver := newAliasResolver(mockConfig)

	// Should be populated after initialization
	if len(resolver.synonyms) == 0 {
		t.Error("synonyms should not be empty after initialization")
	}

	// Check synonyms
	if synonym, exists := resolver.synonyms["code"]; !exists || synonym != "vscode" {
		t.Error("'code' should be a synonym for 'vscode'")
	}
}

// Helper function to get expected VS Code path based on OS
func getExpectedVSCodePath() string {
	switch runtime.GOOS {
	case "darwin":
		return "Visual Studio Code.app"
	case "windows":
		return "Code.exe"
	default: // linux and others
		return "code"
	}
}
