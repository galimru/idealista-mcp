package tools

import (
	"fmt"
	"sync"

	"github.com/galimru/idealista-mcp/auth"
	"github.com/galimru/idealista-mcp/client"
	"github.com/galimru/idealista-mcp/internal/config"
	"github.com/galimru/idealista-mcp/internal/session"
	"github.com/galimru/idealista-mcp/internal/signing"
)

// RuntimeProvider lazily initializes expensive runtime dependencies and caches
// successful results for reuse across tool calls.
type RuntimeProvider struct {
	mu sync.Mutex

	loadConfig    func() (*config.Config, error)
	loadSession   func() (*session.Session, error)
	newSigner     func(string) (*signing.Signer, error)
	newAuthClient func(*session.Session, string, string) (*auth.Client, error)
	newAPIClient  func(*auth.Client, *signing.Signer, string) client.APIClient

	cfg        *config.Config
	sess       *session.Session
	signer     *signing.Signer
	authClient *auth.Client
	apiClient  client.APIClient
}

func NewRuntimeProvider() *RuntimeProvider {
	return &RuntimeProvider{
		loadConfig:    config.Load,
		loadSession:   session.Load,
		newSigner:     signing.New,
		newAuthClient: auth.New,
		newAPIClient: func(a *auth.Client, s *signing.Signer, deviceIdentifier string) client.APIClient {
			return client.New(a, s, deviceIdentifier)
		},
	}
}

func (p *RuntimeProvider) Config() (*config.Config, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cfg != nil {
		return p.cfg, nil
	}

	cfg, err := p.loadConfig()
	if err != nil {
		return nil, err
	}
	p.cfg = cfg
	return cfg, nil
}

func (p *RuntimeProvider) Session() (*session.Session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.sess != nil {
		return p.sess, nil
	}

	sess, err := p.loadSession()
	if err != nil {
		return nil, err
	}
	p.sess = sess
	return sess, nil
}

func (p *RuntimeProvider) Signer() (*signing.Signer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.signer != nil {
		return p.signer, nil
	}

	cfg, err := p.getConfigLocked()
	if err != nil {
		return nil, err
	}

	signerClient, err := p.newSigner(cfg.SigningSecret)
	if err != nil {
		return nil, fmt.Errorf("create signer: %w", err)
	}
	p.signer = signerClient
	return signerClient, nil
}

func (p *RuntimeProvider) APIClient() (client.APIClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.apiClient != nil {
		return p.apiClient, nil
	}

	cfg, err := p.getConfigLocked()
	if err != nil {
		return nil, err
	}

	sess, err := p.getSessionLocked()
	if err != nil {
		return nil, err
	}

	signerClient, err := p.getSignerLocked(cfg)
	if err != nil {
		return nil, err
	}

	authClient, err := p.getAuthClientLocked(cfg, sess)
	if err != nil {
		return nil, err
	}

	p.apiClient = p.newAPIClient(authClient, signerClient, sess.DeviceIdentifier)
	return p.apiClient, nil
}

func (p *RuntimeProvider) getConfigLocked() (*config.Config, error) {
	if p.cfg != nil {
		return p.cfg, nil
	}

	cfg, err := p.loadConfig()
	if err != nil {
		return nil, err
	}
	p.cfg = cfg
	return cfg, nil
}

func (p *RuntimeProvider) getSessionLocked() (*session.Session, error) {
	if p.sess != nil {
		return p.sess, nil
	}

	sess, err := p.loadSession()
	if err != nil {
		return nil, err
	}
	p.sess = sess
	return sess, nil
}

func (p *RuntimeProvider) getSignerLocked(cfg *config.Config) (*signing.Signer, error) {
	if p.signer != nil {
		return p.signer, nil
	}

	signerClient, err := p.newSigner(cfg.SigningSecret)
	if err != nil {
		return nil, fmt.Errorf("create signer: %w", err)
	}
	p.signer = signerClient
	return signerClient, nil
}

func (p *RuntimeProvider) getAuthClientLocked(cfg *config.Config, sess *session.Session) (*auth.Client, error) {
	if p.authClient != nil {
		return p.authClient, nil
	}

	authClient, err := p.newAuthClient(sess, cfg.ClientKey, cfg.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("create auth client: %w", err)
	}
	p.authClient = authClient
	return authClient, nil
}
