package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "loopr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create .loopr directory
	looprDir := filepath.Join(tmpDir, ".loopr")
	if err := os.Mkdir(looprDir, 0755); err != nil {
		t.Fatalf("Failed to create .loopr directory: %v", err)
	}

	// Test case 1: Valid config with docker sandbox
	t.Run("ValidDockerConfig", func(t *testing.T) {
		configPath := filepath.Join(looprDir, "config.json")
		validConfig := `{
  "$schema": "../schema/config.schema.json",
  "sandbox": "docker",
  "looprDir": ".loopr"
}`
		if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}

		if cfg.Sandbox != "docker" {
			t.Errorf("Expected sandbox to be 'docker', got '%s'", cfg.Sandbox)
		}
	})

	// Test case 2: Invalid sandbox type
	t.Run("InvalidSandbox", func(t *testing.T) {
		configPath := filepath.Join(looprDir, "config.json")
		invalidConfig := `{
  "$schema": "../schema/config.schema.json",
  "sandbox": "e2b",
  "looprDir": ".loopr"
}`
		if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := Load()
		if err == nil {
			t.Error("Expected error for non-docker sandbox, got nil")
		}
		if err.Error() != "only 'docker' sandbox supported in v1" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 3: Config file not found
	t.Run("ConfigNotFound", func(t *testing.T) {
		configPath := filepath.Join(looprDir, "config.json")
		if err := os.Remove(configPath); err != nil {
			t.Fatalf("Failed to remove config file: %v", err)
		}

		_, err := Load()
		if err == nil {
			t.Error("Expected error for missing config file, got nil")
		}
	})
}

func TestLooprDirExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "loopr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test case 1: .loopr directory does not exist
	t.Run("DirectoryDoesNotExist", func(t *testing.T) {
		if LooprDirExists() {
			t.Error("Expected LooprDirExists() to return false, got true")
		}
	})

	// Test case 2: .loopr directory exists
	t.Run("DirectoryExists", func(t *testing.T) {
		looprDir := filepath.Join(tmpDir, ".loopr")
		if err := os.Mkdir(looprDir, 0755); err != nil {
			t.Fatalf("Failed to create .loopr directory: %v", err)
		}

		if !LooprDirExists() {
			t.Error("Expected LooprDirExists() to return true, got false")
		}
	})

	// Test case 3: .loopr exists but is a file, not a directory
	t.Run("FileNotDirectory", func(t *testing.T) {
		looprPath := filepath.Join(tmpDir, ".loopr")
		if err := os.Remove(looprPath); err != nil {
			t.Fatalf("Failed to remove .loopr directory: %v", err)
		}

		if err := os.WriteFile(looprPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create .loopr file: %v", err)
		}

		if LooprDirExists() {
			t.Error("Expected LooprDirExists() to return false for file, got true")
		}
	})
}
