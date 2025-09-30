// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	TargetPath      string `yaml:"target_path"`
	BackupPath      string `yaml:"backup_path"`
	DestinationPath string `yaml:"destination_path"`
	TmpDir          string `yaml:"tmp_dir"`
}

// Default returns the default configuration
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	tmpDir := filepath.Join(homeDir, "Pictures", "tmp")

	return &Config{
		TargetPath:      "/Volumes/1-1",
		BackupPath:      "/Volumes/1-2",
		DestinationPath: "/Volumes/a7iii",
		TmpDir:          tmpDir,
	}
}

// Load loads configuration from a YAML file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to a YAML file
func Save(config *Config, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetPath returns the configuration file path.
// It looks for config in the following order:
// 1. ./config.yaml (current directory)
// 2. ~/.config/rename-sony-photos/config.yaml (XDG config directory)
func GetPath() string {
	// Check current directory first
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	// Check ~/.config directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".config", "rename-sony-photos", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	return ""
}

// LoadOrDefault loads existing config or returns default configuration
func LoadOrDefault() (*Config, error) {
	configPath := GetPath()

	// If config file exists, load it
	if configPath != "" {
		config, err := Load(configPath)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// Return default config
	return Default(), nil
}
