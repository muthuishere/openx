package core

import (
	"runtime"
	"testing"
)

func TestNewAliasResolver(t *testing.T) {
	resolver := newAliasResolver()

	if resolver == nil {
		t.Fatal("newAliasResolver() returned nil")
	}

	if resolver.canonicals == nil {
		t.Fatal("canonicals map is nil")
	}

	if resolver.synonyms == nil {
		t.Fatal("synonyms map is nil")
	}

	// Check that some canonical apps are initialized
	if len(resolver.canonicals) == 0 {
		t.Error("canonicals map is empty, expected some entries")
	}
}

func TestAliasResolver_Resolve(t *testing.T) {
	resolver := newAliasResolver()

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

func TestAliasResolver_InitializeCanonicals(t *testing.T) {
	resolver := &AliasResolver{
		canonicals: map[string]map[string]string{},
		synonyms:   map[string]string{},
	}

	// Initially empty
	if len(resolver.canonicals) != 0 {
		t.Error("canonicals should be empty before initialization")
	}

	if len(resolver.synonyms) != 0 {
		t.Error("synonyms should be empty before initialization")
	}

	resolver.initializeCanonicals()

	// Should be populated after initialization
	if len(resolver.canonicals) == 0 {
		t.Error("canonicals should not be empty after initialization")
	}

	if len(resolver.synonyms) == 0 {
		t.Error("synonyms should not be empty after initialization")
	}

	// Check specific entries
	vscode, exists := resolver.canonicals["vscode"]
	if !exists {
		t.Error("vscode should exist in canonicals")
	}

	if len(vscode) == 0 {
		t.Error("vscode should have platform-specific entries")
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
