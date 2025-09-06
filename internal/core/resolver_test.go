package core

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "http URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "https URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "ftp URL",
			input:    "ftp://example.com",
			expected: true,
		},
		{
			name:     "file URL",
			input:    "file:///path/to/file",
			expected: true,
		},
		{
			name:     "custom protocol",
			input:    "custom://something",
			expected: true,
		},
		{
			name:     "regular file path",
			input:    "/path/to/file",
			expected: false,
		},
		{
			name:     "relative path",
			input:    "./file.txt",
			expected: false,
		},
		{
			name:     "just a name",
			input:    "filename",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "windows path",
			input:    "C:\\Users\\test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isURL(tt.input)
			if result != tt.expected {
				t.Errorf("isURL(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExpandTilde(t *testing.T) {
	// Get user home directory for comparison
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde at start",
			input:    "~/Documents",
			expected: filepath.Join(homeDir, "Documents"),
		},
		{
			name:     "just tilde",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "tilde with slash",
			input:    "~/",
			expected: homeDir, // expandTilde removes trailing slash
		},
		{
			name:     "no tilde",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "tilde not at start",
			input:    "path/~/file",
			expected: "path/~/file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandTilde(tt.input)
			if result != tt.expected {
				t.Errorf("expandTilde(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExists(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     tmpFile.Name(),
			expected: true,
		},
		{
			name:     "existing directory",
			path:     tmpDir,
			expected: true,
		},
		{
			name:     "non-existing file",
			path:     "/path/that/does/not/exist",
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := exists(tt.path)
			if result != tt.expected {
				t.Errorf("exists(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestFindAppExecutable(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("findAppExecutable is only available on macOS")
	}

	tests := []struct {
		name    string
		appPath string
		wantErr bool
	}{
		{
			name:    "non-existing app",
			appPath: "/Applications/NonExistentApp.app",
			wantErr: true,
		},
		{
			name:    "invalid path",
			appPath: "/not/an/app/path",
			wantErr: true,
		},
		{
			name:    "empty path",
			appPath: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := findAppExecutable(tt.appPath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("findAppExecutable(%s) expected error but got none", tt.appPath)
				}
			} else {
				if err != nil {
					t.Errorf("findAppExecutable(%s) unexpected error: %v", tt.appPath, err)
				}
			}
		})
	}
}

func TestResolveTarget(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name   string
		target string
		// We can't predict exact results for all cases, so we'll check basic behavior
	}{
		{
			name:   "URL target",
			target: "https://example.com",
		},
		{
			name:   "existing file",
			target: tmpFile.Name(),
		},
		{
			name:   "tilde path",
			target: "~/Documents",
		},
		{
			name:   "relative path",
			target: "./file.txt",
		},
		{
			name:   "absolute path",
			target: "/absolute/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveTarget(tt.target)

			// Basic sanity checks
			if result == "" && tt.target != "" {
				t.Errorf("resolveTarget(%s) returned empty string", tt.target)
			}

			// For tilde paths, check that tilde was expanded
			if tt.target == "~/Documents" {
				expected := filepath.Join(homeDir, "Documents")
				if result != expected {
					t.Errorf("resolveTarget(%s) = %v, want %v", tt.target, result, expected)
				}
			}

			// For URLs, should return as-is
			if tt.target == "https://example.com" {
				if result != tt.target {
					t.Errorf("resolveTarget(%s) = %v, want %v", tt.target, result, tt.target)
				}
			}
		})
	}
}

func TestResolveTargets(t *testing.T) {
	targets := []string{
		"https://example.com",
		"~/Documents",
		"./file.txt",
	}

	results := resolveTargets(targets)

	if len(results) != len(targets) {
		t.Errorf("resolveTargets() returned %d results, want %d", len(results), len(targets))
	}

	// Check that each target was processed
	for i, target := range targets {
		if i >= len(results) {
			t.Errorf("Missing result for target %s", target)
			continue
		}

		// URLs should be unchanged
		if target == "https://example.com" && results[i] != target {
			t.Errorf("URL target %s was modified to %s", target, results[i])
		}
	}
}

func TestValidateTarget(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name:    "valid URL",
			target:  "https://example.com",
			wantErr: false,
		},
		{
			name:    "invalid URL format",
			target:  "not-a-url",
			wantErr: true, // This should be treated as a file path and fail since it doesn't exist
		},
		{
			name:    "existing file",
			target:  tmpFile.Name(),
			wantErr: false,
		},
		{
			name:    "non-existing file",
			target:  "/path/that/does/not/exist",
			wantErr: true,
		},
		{
			name:    "empty target",
			target:  "",
			wantErr: false, // Empty string might be allowed by validateTarget
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTarget(tt.target)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateTarget(%s) expected error but got none", tt.target)
				}
			} else {
				if err != nil {
					t.Errorf("validateTarget(%s) unexpected error: %v", tt.target, err)
				}
			}
		})
	}
}
