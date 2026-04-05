package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	AppVersion = "14.6.0"
	UserAgent  = "Dalvik/2.1.0 (Linux; U; Android 13; Pixel 7 Build/TQ3A.230805.001)"
)

// Config holds user-editable credentials loaded from config.json.
type Config struct {
	ClientKey     string `json:"client_key"`
	ClientSecret  string `json:"client_secret"`
	SigningSecret string `json:"signing_secret"`
}

// ValidationError reports missing config fields.
type ValidationError struct {
	Path   string
	Fields []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"config %s: required fields are not set (%s) - edit %s",
		e.Path,
		strings.Join(e.Fields, ", "),
		e.Path,
	)
}

// Load loads the user config from the default config location, creating it
// with empty defaults when it does not yet exist.
func Load() (*Config, error) {
	path, err := FilePath()
	if err != nil {
		return nil, err
	}
	return LoadPath(path)
}

// LoadPath loads config from path, creating a default file when missing.
func LoadPath(path string) (*Config, error) {
	if err := ensureConfigFile(path); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg := defaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if err := validate(cfg, path); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// FilePath returns the default config file path.
func FilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "idealista-mcp", "config.json"), nil
}

func ensureConfigFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat config %s: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir %s: %w", filepath.Dir(path), err)
	}

	data, err := json.MarshalIndent(defaultConfig(), "", "  ")
	if err != nil {
		return fmt.Errorf("marshal default config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write default config %s: %w", path, err)
	}

	return nil
}

func defaultConfig() *Config {
	return &Config{}
}

func validate(cfg *Config, path string) error {
	var fields []string

	if isUnset(cfg.ClientKey) {
		fields = append(fields, "client_key")
	}
	if isUnset(cfg.ClientSecret) {
		fields = append(fields, "client_secret")
	}
	if isUnset(cfg.SigningSecret) {
		fields = append(fields, "signing_secret")
	}

	if len(fields) > 0 {
		return &ValidationError{Path: path, Fields: fields}
	}

	return nil
}

func isUnset(value string) bool {
	return strings.TrimSpace(value) == ""
}
