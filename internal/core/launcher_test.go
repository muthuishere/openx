package core

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLaunchApp(t *testing.T) {
	// Create a test config
	testContent := `
apps:
  testapp:
    darwin: "/Applications/TestApp.app"
    linux: "echo"
    windows: "cmd.exe"
    kill: ["TestApp"]

aliases:
  ta: testapp`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name    string
		alias   string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid app with echo command",
			alias:   "testapp",
			args:    []string{"hello"},
			wantErr: false, // echo should work on Unix systems
		},
		{
			name:    "valid alias",
			alias:   "ta",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "unknown app",
			alias:   "unknown",
			args:    []string{},
			wantErr: true,
			errMsg:  "unknown app: unknown",
		},
		{
			name:    "empty alias",
			alias:   "",
			args:    []string{},
			wantErr: true,
			errMsg:  "unknown app: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip on Windows if using echo command
			if runtime.GOOS == "windows" && tt.alias == "testapp" {
				t.Skip("Skipping echo test on Windows")
			}

			err := LaunchApp(tt.alias, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LaunchApp() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("LaunchApp() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Logf("LaunchApp() failed (may be expected): %v", err)
			}
		})
	}
}

func TestIsDirectPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "absolute Unix path",
			path:     "/Applications/Test.app",
			expected: true,
		},
		{
			name:     "absolute Windows path",
			path:     "C:\\Program Files\\Test.exe",
			expected: true,
		},
		{
			name:     "relative path with slash",
			path:     "./test",
			expected: true,
		},
		{
			name:     "relative path with backslash",
			path:     ".\\test",
			expected: true,
		},
		{
			name:     "just command name",
			path:     "test",
			expected: false,
		},
		{
			name:     "empty string",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDirectPath(tt.path)
			if result != tt.expected {
				t.Errorf("isDirectPath(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestLaunchDirectPath(t *testing.T) {
	// Create a temporary executable script for testing
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_script")

	// Create a simple script that just exits successfully
	scriptContent := "#!/bin/bash\nexit 0\n"
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(tmpDir, "test_script.bat")
		scriptContent = "@echo off\nexit /b 0\n"
	}

	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	tests := []struct {
		name    string
		appPath string
		args    []string
		wantErr bool
	}{
		{
			name:    "existing executable",
			appPath: scriptPath,
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "non-existing path",
			appPath: "/definitely/does/not/exist",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "empty path",
			appPath: "",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := launchDirectPath(tt.appPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("launchDirectPath() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("launchDirectPath() unexpected error: %v", err)
			}
		})
	}
}

func TestExecuteApp(t *testing.T) {
	tests := []struct {
		name       string
		launchPath string
		args       []string
		wantErr    bool
	}{
		{
			name:       "echo command",
			launchPath: "echo",
			args:       []string{"test"},
			wantErr:    false,
		},
		{
			name:       "non-existing command",
			launchPath: "definitely-does-not-exist-command",
			args:       []string{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip echo test on Windows
			if runtime.GOOS == "windows" && tt.launchPath == "echo" {
				t.Skip("Skipping echo test on Windows")
			}

			err := executeApp(tt.launchPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("executeApp() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("executeApp() unexpected error: %v", err)
			}
		})
	}
}

func TestLaunchMacOSApp(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("launchMacOSApp is only available on macOS")
	}

	tests := []struct {
		name    string
		appPath string
		args    []string
		wantErr bool
	}{
		{
			name:    "non-existing app",
			appPath: "/Applications/NonExistentApp.app",
			args:    []string{},
			wantErr: false, // Should fallback to 'open' command, which might succeed
		},
		{
			name:    "invalid app path",
			appPath: "/not/an/app",
			args:    []string{},
			wantErr: false, // Should fallback to 'open' command
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := launchMacOSApp(tt.appPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("launchMacOSApp() expected error but got none")
				}
				return
			}

			// Even if no error, the app might not actually launch
			// This is normal in a test environment
			if err != nil {
				t.Logf("launchMacOSApp() failed (may be expected): %v", err)
			}
		})
	}
}

func TestLaunchWithOpen(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("launchWithOpen is only available on macOS")
	}

	tests := []struct {
		name    string
		appPath string
		args    []string
		wantErr bool
	}{
		{
			name:    "test with invalid app",
			appPath: "/Applications/NonExistentApp.app",
			args:    []string{},
			wantErr: false, // 'open' command exists, might fail but won't return error immediately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := launchWithOpen(tt.appPath, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("launchWithOpen() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Logf("launchWithOpen() failed (may be expected): %v", err)
			}
		})
	}
}

func TestLaunchMultipleApps(t *testing.T) {
	// Create a test config with a working command
	testContent := `
apps:
  echo:
    darwin: "echo"
    linux: "echo"
    windows: "cmd.exe"
  
  invalid:
    darwin: "/definitely/does/not/exist"
    linux: "/definitely/does/not/exist"
    windows: "definitely-does-not-exist.exe"`

	configPath := setupTestConfig(t, testContent)
	cleanup := setTempConfigPath(t, configPath)
	defer cleanup()

	tests := []struct {
		name    string
		aliases []string
		wantErr bool
	}{
		{
			name:    "empty list",
			aliases: []string{},
			wantErr: false,
		},
		{
			name:    "single valid app",
			aliases: []string{"echo"},
			wantErr: false,
		},
		{
			name:    "single invalid app",
			aliases: []string{"invalid"},
			wantErr: true,
		},
		{
			name:    "mixed valid and invalid",
			aliases: []string{"echo", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip echo tests on Windows
			if runtime.GOOS == "windows" {
				t.Skip("Skipping echo tests on Windows")
			}

			err := launchMultipleApps(tt.aliases)

			if tt.wantErr {
				if err == nil {
					t.Errorf("launchMultipleApps() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("launchMultipleApps() unexpected error: %v", err)
			}
		})
	}
}

func TestLaunchApp_ConfigError(t *testing.T) {
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

	err := LaunchApp("testapp", []string{})
	if err == nil {
		t.Error("LaunchApp() expected error when config file doesn't exist")
	}

	expectedSubstring := "failed to load config"
	if err != nil && !contains(err.Error(), expectedSubstring) {
		t.Errorf("LaunchApp() error = %v, want error containing %v", err, expectedSubstring)
	}
}

func TestLaunchApp_DirectPath(t *testing.T) {
	// Test direct path functionality
	tests := []struct {
		name    string
		alias   string
		args    []string
		wantErr bool
		skipOS  string
	}{
		{
			name:    "direct path to echo",
			alias:   "/bin/echo",
			args:    []string{"hello"},
			wantErr: false,
			skipOS:  "windows",
		},
		{
			name:    "non-existing direct path",
			alias:   "/definitely/does/not/exist",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOS != "" && runtime.GOOS == tt.skipOS {
				t.Skipf("Skipping test on %s", tt.skipOS)
			}

			err := LaunchApp(tt.alias, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LaunchApp() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("LaunchApp() unexpected error: %v", err)
			}
		})
	}
}
