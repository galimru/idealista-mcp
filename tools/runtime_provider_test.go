package tools

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"testing"

	"github.com/galimru/idealista-mcp/auth"
	"github.com/galimru/idealista-mcp/client"
	"github.com/galimru/idealista-mcp/internal/config"
	"github.com/galimru/idealista-mcp/internal/session"
	"github.com/galimru/idealista-mcp/internal/signing"
)

type stubAPIClient struct{}

func (stubAPIClient) Get(_ context.Context, _ string, _ url.Values) ([]byte, error) {
	return nil, nil
}

func (stubAPIClient) Post(_ context.Context, _ string, _ url.Values) ([]byte, error) {
	return nil, nil
}

func TestRuntimeProviderCachesSuccessfulAPIClientInitialization(t *testing.T) {
	var loadConfigCalls, loadSessionCalls, newSignerCalls, newAuthCalls, newClientCalls int

	provider := &RuntimeProvider{
		loadConfig: func() (*config.Config, error) {
			loadConfigCalls++
			return &config.Config{ClientKey: "key", ClientSecret: "secret", SigningSecret: "signing"}, nil
		},
		loadSession: func() (*session.Session, error) {
			loadSessionCalls++
			return &session.Session{DeviceIdentifier: "device-id"}, nil
		},
		newSigner: func(secret string) (*signing.Signer, error) {
			newSignerCalls++
			return signing.New(secret)
		},
		newAuthClient: func(sess *session.Session, key, secret string) (*auth.Client, error) {
			newAuthCalls++
			return auth.New(sess, key, secret)
		},
		newAPIClient: func(_ *auth.Client, _ *signing.Signer, _ string) client.APIClient {
			newClientCalls++
			return stubAPIClient{}
		},
	}

	one, err := provider.APIClient()
	if err != nil {
		t.Fatalf("first APIClient() error = %v", err)
	}
	two, err := provider.APIClient()
	if err != nil {
		t.Fatalf("second APIClient() error = %v", err)
	}

	if one != two {
		t.Fatal("expected API client instance to be cached")
	}

	if loadConfigCalls != 1 || loadSessionCalls != 1 || newSignerCalls != 1 || newAuthCalls != 1 || newClientCalls != 1 {
		t.Fatalf(
			"unexpected call counts: config=%d session=%d signer=%d auth=%d client=%d",
			loadConfigCalls, loadSessionCalls, newSignerCalls, newAuthCalls, newClientCalls,
		)
	}
}

func TestRuntimeProviderRetriesAfterLoadConfigFailure(t *testing.T) {
	var loadConfigCalls int

	provider := &RuntimeProvider{
		loadConfig: func() (*config.Config, error) {
			loadConfigCalls++
			if loadConfigCalls == 1 {
				return nil, errors.New("boom")
			}
			return &config.Config{ClientKey: "key", ClientSecret: "secret", SigningSecret: "signing"}, nil
		},
		loadSession: func() (*session.Session, error) {
			return &session.Session{DeviceIdentifier: "device-id"}, nil
		},
		newSigner: func(secret string) (*signing.Signer, error) {
			return signing.New(secret)
		},
		newAuthClient: func(sess *session.Session, key, secret string) (*auth.Client, error) {
			return auth.New(sess, key, secret)
		},
		newAPIClient: func(_ *auth.Client, _ *signing.Signer, _ string) client.APIClient {
			return stubAPIClient{}
		},
	}

	if _, err := provider.APIClient(); err == nil {
		t.Fatal("expected first APIClient() call to fail")
	}

	if _, err := provider.APIClient(); err != nil {
		t.Fatalf("expected second APIClient() call to succeed, got %v", err)
	}

	if loadConfigCalls != 2 {
		t.Fatalf("loadConfig calls = %d, want 2", loadConfigCalls)
	}
}

func TestRuntimeProviderConcurrentConfigInitialization(t *testing.T) {
	var (
		mu              sync.Mutex
		loadConfigCalls int
	)

	provider := &RuntimeProvider{
		loadConfig: func() (*config.Config, error) {
			mu.Lock()
			loadConfigCalls++
			mu.Unlock()
			return &config.Config{ClientKey: "key", ClientSecret: "secret", SigningSecret: "signing"}, nil
		},
	}

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := provider.Config(); err != nil {
				t.Errorf("Config() error = %v", err)
			}
		}()
	}
	wg.Wait()

	if loadConfigCalls != 1 {
		t.Fatalf("loadConfig calls = %d, want 1", loadConfigCalls)
	}
}
