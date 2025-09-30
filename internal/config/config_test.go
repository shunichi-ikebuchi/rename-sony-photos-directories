package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := Default()

	if config.TargetPath != "/Volumes/1-1" {
		t.Errorf("Expected TargetPath to be /Volumes/1-1, got %s", config.TargetPath)
	}
	if config.BackupPath != "/Volumes/1-2" {
		t.Errorf("Expected BackupPath to be /Volumes/1-2, got %s", config.BackupPath)
	}
	if config.DestinationPath != "/Volumes/a7iii" {
		t.Errorf("Expected DestinationPath to be /Volumes/a7iii, got %s", config.DestinationPath)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temporary config file
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `target_path: /test/path
backup_path: /test/backup
destination_path: /test/dest
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config
	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify values
	if config.TargetPath != "/test/path" {
		t.Errorf("Expected TargetPath to be /test/path, got %s", config.TargetPath)
	}
	if config.BackupPath != "/test/backup" {
		t.Errorf("Expected BackupPath to be /test/backup, got %s", config.BackupPath)
	}
	if config.DestinationPath != "/test/dest" {
		t.Errorf("Expected DestinationPath to be /test/dest, got %s", config.DestinationPath)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent config file, got nil")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	// Create temporary config file with invalid YAML
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	invalidYAML := `target_path: /test/path
	invalid yaml content
	no proper structure
`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestSaveConfig(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create and save config
	config := &Config{
		TargetPath:      "/test/target",
		BackupPath:      "/test/backup",
		DestinationPath: "/test/destination",
	}

	if err := Save(config, configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load and verify content
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedConfig.TargetPath != config.TargetPath {
		t.Errorf("TargetPath mismatch: got %s, want %s", loadedConfig.TargetPath, config.TargetPath)
	}
	if loadedConfig.BackupPath != config.BackupPath {
		t.Errorf("BackupPath mismatch: got %s, want %s", loadedConfig.BackupPath, config.BackupPath)
	}
	if loadedConfig.DestinationPath != config.DestinationPath {
		t.Errorf("DestinationPath mismatch: got %s, want %s", loadedConfig.DestinationPath, config.DestinationPath)
	}
}

func TestGetConfigPath(t *testing.T) {
	// Save current directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test 1: No config file exists
	configPath := GetPath()
	if configPath != "" {
		t.Errorf("Expected empty path when no config exists, got %s", configPath)
	}

	// Test 2: Config file in current directory
	localConfigPath := "config.yaml"
	if err := os.WriteFile(localConfigPath, []byte("target_path: /test"), 0644); err != nil {
		t.Fatalf("Failed to create local config: %v", err)
	}

	configPath = GetPath()
	if configPath != localConfigPath {
		t.Errorf("Expected %s, got %s", localConfigPath, configPath)
	}
}

func TestLoadOrCreateConfig(t *testing.T) {
	// Save current directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test: No config file exists (should return default)
	config, err := LoadOrDefault()
	if err != nil {
		t.Fatalf("LoadOrDefault failed: %v", err)
	}

	defaultConfig := Default()
	if config.TargetPath != defaultConfig.TargetPath {
		t.Errorf("Expected default TargetPath %s, got %s", defaultConfig.TargetPath, config.TargetPath)
	}
}
