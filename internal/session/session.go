package session

import (
	"encoding/json"
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
	p, err := filePath()
	if err != nil {
		return nil, err
	}

	s := &Session{path: p}

	data, err := os.ReadFile(p)
	if err == nil {
		_ = json.Unmarshal(data, s)
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

func filePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "idealista-mcp", "session.json"), nil
}
