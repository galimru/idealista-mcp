package auth

import (
	"testing"

	"github.com/galimru/idealista-mcp/internal/session"
)

func TestNewRequiresCredentials(t *testing.T) {
	_, err := New(&session.Session{}, "", "secret")
	if err == nil {
		t.Fatal("expected error for missing client key")
	}

	_, err = New(&session.Session{}, "key", "")
	if err == nil {
		t.Fatal("expected error for missing client secret")
	}
}
