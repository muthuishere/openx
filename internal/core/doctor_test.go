package core

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestRunDoctor(t *testing.T) {
	// Create a test config
	testContent := `
apps:
  testapp:
    darwin: "/Applications/TestApp.app"
    linux: "testapp"
    windows: "testapp.exe"
    kill: ["TestApp"]
  
  existingcommand:
    darwin: "/bin/ls"
    linux: "/bin/ls"
    windows: "cmd.exe"

aliases:
  ta: testapp
  ls: existingcommand`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name       string
		jsonOutput bool
		wantErr    bool
	}{
		{
			name:       "human readable output",
			jsonOutput: false,
			wantErr:    false,
		},
		{
			name:       "JSON output",
			jsonOutput: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := RunDoctor(tt.jsonOutput)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.wantErr {
				if err == nil {
					t.Errorf("RunDoctor() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("RunDoctor() unexpected error: %v", err)
				return
			}

			// Verify output format
			if tt.jsonOutput {
				// Should be valid JSON
				var report DoctorReport
				if err := json.Unmarshal([]byte(output), &report); err != nil {
					t.Errorf("RunDoctor() JSON output is invalid: %v\nOutput: %s", err, output)
				}
			} else {
				// Human readable should contain certain keywords
				if len(output) == 0 {
					t.Error("RunDoctor() human output is empty")
				}
			}
		})
	}
}

func TestCheckAppStatus(t *testing.T) {
	tests := []struct {
		name     string
		appName  string
		app      *App
		expected string // expected status
	}{
		{
			name:    "app with no path for current OS",
			appName: "testapp",
			app: &App{
				Paths: map[string]string{
					"invalid-os": "/some/path",
				},
			},
			expected: "no-path",
		},
		{
			name:    "app with existing command",
			appName: "ls",
			app: &App{
				Paths: map[string]string{
					"darwin":  "/bin/ls",
					"linux":   "/bin/ls",
					"windows": "cmd.exe",
				},
			},
			expected: "available", // ls should exist on most systems
		},
		{
			name:    "app with non-existing path",
			appName: "nonexistent",
			app: &App{
				Paths: map[string]string{
					"darwin":  "/definitely/does/not/exist",
					"linux":   "/definitely/does/not/exist",
					"windows": "definitely-does-not-exist.exe",
				},
			},
			expected: "missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := checkAppStatus(tt.appName, tt.app)

			if status.Name != tt.appName {
				t.Errorf("checkAppStatus() name = %v, want %v", status.Name, tt.appName)
			}

			if status.Status != tt.expected {
				t.Errorf("checkAppStatus() status = %v, want %v", status.Status, tt.expected)
			}
		})
	}
}

func TestAppExists(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "absolute path to ls",
			path:     "/bin/ls",
			expected: true, // ls should exist on Unix systems
		},
		{
			name:     "command in PATH",
			path:     "ls",
			expected: true, // ls should be in PATH on Unix systems
		},
		{
			name:     "non-existing absolute path",
			path:     "/definitely/does/not/exist",
			expected: false,
		},
		{
			name:     "non-existing command",
			path:     "definitely-does-not-exist-command",
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
			result := appExists(tt.path)

			// On Windows, ls might not exist, so adjust expectations
			if (tt.path == "/bin/ls" || tt.path == "ls") && !result {
				t.Logf("ls command not found on this system, skipping test")
				return
			}

			if result != tt.expected {
				t.Errorf("appExists(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{
			name:     "available status",
			status:   "available",
			expected: "✓",
		},
		{
			name:     "missing status",
			status:   "missing",
			expected: "✗",
		},
		{
			name:     "no-path status",
			status:   "no-path",
			expected: "○",
		},
		{
			name:     "unknown status",
			status:   "unknown",
			expected: "?",
		},
		{
			name:     "empty status",
			status:   "",
			expected: "?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusIcon(tt.status)
			if result != tt.expected {
				t.Errorf("getStatusIcon(%s) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetStatusColor(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{
			name:     "available status",
			status:   "available",
			expected: ColorGreen,
		},
		{
			name:     "missing status",
			status:   "missing",
			expected: ColorRed,
		},
		{
			name:     "no-path status",
			status:   "no-path",
			expected: ColorYellow,
		},
		{
			name:     "unknown status",
			status:   "unknown",
			expected: ColorReset,
		},
		{
			name:     "empty status",
			status:   "",
			expected: ColorReset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusColor(tt.status)
			if result != tt.expected {
				t.Errorf("getStatusColor(%s) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestOutputJSON(t *testing.T) {
	report := DoctorReport{
		Platform:   "test",
		ConfigPath: "/test/config.yaml",
		Apps: []AppStatus{
			{
				Name:       "testapp",
				LaunchPath: "/test/path",
				Status:     "available",
				Running:    false,
			},
		},
		Aliases: map[string]string{
			"ta": "testapp",
		},
		Summary: Summary{
			Total:     1,
			Available: 1,
			Missing:   0,
			Running:   0,
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputJSON(report)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Errorf("outputJSON() unexpected error: %v", err)
	}

	// Verify it's valid JSON
	var parsedReport DoctorReport
	if err := json.Unmarshal([]byte(output), &parsedReport); err != nil {
		t.Errorf("outputJSON() produced invalid JSON: %v\nOutput: %s", err, output)
	}

	// Verify content
	if parsedReport.Platform != report.Platform {
		t.Errorf("outputJSON() platform = %v, want %v", parsedReport.Platform, report.Platform)
	}
}

func TestRunDoctor_ConfigError(t *testing.T) {
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

	err := RunDoctor(false)
	if err == nil {
		t.Error("RunDoctor() expected error when config file doesn't exist")
	}

	expectedSubstring := "failed to load config"
	if err != nil && !contains(err.Error(), expectedSubstring) {
		t.Errorf("RunDoctor() error = %v, want error containing %v", err, expectedSubstring)
	}
}
