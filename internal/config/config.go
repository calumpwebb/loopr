package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config represents the loopr configuration
type Config struct {
	Schema     string            `json:"$schema"`
	Sandbox    string            `json:"sandbox"`
	LooprDir   string            `json:"looprDir,omitempty"`
	Model      *ModelConfig      `json:"model,omitempty"`
	Git        *GitConfig        `json:"git,omitempty"`
	Iterations *IterationsConfig `json:"iterations,omitempty"`
}

// ModelConfig holds model settings for different phases
type ModelConfig struct {
	Plan  string `json:"plan,omitempty"`
	Build string `json:"build,omitempty"`
}

// GitConfig holds git-related settings
type GitConfig struct {
	AutoPush   bool `json:"autoPush,omitempty"`
	AutoCommit bool `json:"autoCommit,omitempty"`
}

// IterationsConfig holds default iteration counts
type IterationsConfig struct {
	Plan  int `json:"plan,omitempty"`
	Build int `json:"build,omitempty"`
}

// Load reads and parses the .loopr/config.json file
func Load() (*Config, error) {
	looprDir := filepath.Join(".", ".loopr")
	configPath := filepath.Join(looprDir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Validate that sandbox is "docker"
	if cfg.Sandbox != "docker" {
		return nil, errors.New("only 'docker' sandbox supported in v1")
	}

	return &cfg, nil
}

// LooprDirExists checks if the .loopr directory exists
func LooprDirExists() bool {
	info, err := os.Stat(".loopr")
	return err == nil && info.IsDir()
}
