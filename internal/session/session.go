package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Session holds persistent per-installation state: a stable device identifier
// and the cached OAuth access token so restarts reuse a valid token.
type Session struct {
	DeviceIdentifier string    `json:"device_identifier"`
	AccessToken      string    `json:"access_token,omitempty"`
	ExpiresAt        time.Time `json:"expires_at,omitempty"`

	path string // not serialised
}

// Load reads the session file. If no file exists a new device_identifier is
// generated and the session is saved before returning.
func Load() (*Session, error) {
	p, err := FilePath()
	if err != nil {
		return nil, err
	}
	return LoadPath(p)
}

// LoadPath reads the session file from path. If no file exists a new
// device_identifier is generated and the session is saved before returning.
func LoadPath(path string) (*Session, error) {
	s := &Session{path: path}

	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, s)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read session %s: %w", path, err)
	}

	if s.DeviceIdentifier == "" {
		s.DeviceIdentifier = uuid.New().String()
		if err := s.Save(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Save persists the session to disk.
func (s *Session) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// FilePath returns the default session file path.
func FilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "idealista-mcp", "session.json"), nil
}
