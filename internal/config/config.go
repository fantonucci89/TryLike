package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	DefaultEditor = "vim"
	configDir     = ".config/tryl"
	configFile    = "config.toml"
)

// Config holds the application configuration.
type Config struct {
	BasePath string `toml:"base_path"`
	Editor   string `toml:"editor"`
}

// Load reads (or creates) the config file from ~/.config/tryl/config.toml.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, configDir)
	path := filepath.Join(dir, configFile)

	// Provide defaults.
	cfg := &Config{
		BasePath: filepath.Join(home, "Work"),
		Editor:   DefaultEditor,
	}

	// If the file doesn't exist, create it with defaults.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		enc := toml.NewEncoder(f)
		if err := enc.Encode(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	// Parse the existing file.
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}

	// Expand ~ in BasePath.
	if len(cfg.BasePath) > 1 && cfg.BasePath[:2] == "~/" {
		cfg.BasePath = filepath.Join(home, cfg.BasePath[2:])
	}

	return cfg, nil
}

// ConfigFilePath returns the absolute path to the config file.
func ConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}
