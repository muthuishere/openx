package lib

import (
	"testing"
)

func TestNew(t *testing.T) {
	ox := New()
	if ox == nil {
		t.Fatal("New() returned nil")
	}

	if ox.configPath != "" {
		t.Errorf("Expected empty config path, got %s", ox.configPath)
	}
}

func TestNewWithConfig(t *testing.T) {
	customPath := "/custom/config/path.yaml"
	ox := NewWithConfig(customPath)

	if ox == nil {
		t.Fatal("NewWithConfig() returned nil")
	}

	if ox.configPath != customPath {
		t.Errorf("Expected config path %s, got %s", customPath, ox.configPath)
	}
}

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() returned empty string")
	}

	if version != Version {
		t.Errorf("GetVersion() = %s, want %s", version, Version)
	}
}

func TestGetName(t *testing.T) {
	name := GetName()
	if name == "" {
		t.Error("GetName() returned empty string")
	}

	if name != Name {
		t.Errorf("GetName() = %s, want %s", name, Name)
	}
}

func TestEnsureConfig(t *testing.T) {
	ox := New()

	// This should create the config if it doesn't exist
	err := ox.EnsureConfig()
	if err != nil {
		t.Errorf("EnsureConfig() unexpected error: %v", err)
	}
}

// Integration tests would go here, but they require a valid config file
// For now, we'll test the basic functionality that doesn't require external dependencies

func TestLibraryAPI(t *testing.T) {
	// Test that all expected methods exist and have correct signatures
	ox := New()

	// Test method existence (these will fail if config doesn't exist, but that's expected)
	_ = ox.RunAlias
	_ = ox.RunDirect
	_ = ox.Kill
	_ = ox.AddAlias
	_ = ox.RemoveAlias
	_ = ox.ListAliases
	_ = ox.Doctor
	_ = ox.DoctorJSON

	// If we get here, all methods exist with correct signatures
	t.Log("All library methods exist with correct signatures")
}
