package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPathCreatesDefaultFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	cfg, err := LoadPath(path)
	if err == nil {
		t.Fatal("expected validation error for default config")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config to be returned alongside validation error")
	}

	if _, statErr := os.Stat(path); statErr != nil {
		t.Fatalf("expected config file to be created: %v", statErr)
	}

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read created config: %v", readErr)
	}

	text := string(data)
	for _, field := range []string{"client_key", "client_secret", "signing_secret"} {
		if !strings.Contains(text, field) {
			t.Fatalf("expected created config to contain %q, got %s", field, text)
		}
	}
}

func TestLoadPathValidConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
  "client_key": "key",
  "client_secret": "secret",
  "signing_secret": "signing"
}`

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadPath(path)
	if err != nil {
		t.Fatalf("LoadPath() error = %v", err)
	}

	if cfg.ClientKey != "key" || cfg.ClientSecret != "secret" || cfg.SigningSecret != "signing" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadPathValidationErrorIncludesPathAndFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
  "client_key": "",
  "client_secret": "secret",
  "signing_secret": "replace-with-signing-secret"
}`

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadPath(path)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if cfg == nil {
		t.Fatal("expected config to be returned")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}

	if validationErr.Path != path {
		t.Fatalf("validation path = %q, want %q", validationErr.Path, path)
	}

	if !strings.Contains(err.Error(), "client_key") || !strings.Contains(err.Error(), "signing_secret") {
		t.Fatalf("unexpected validation error text: %v", err)
	}
}
