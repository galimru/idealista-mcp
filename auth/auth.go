package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/galimru/idealista-mcp/internal/api"
	"github.com/galimru/idealista-mcp/internal/config"
	"github.com/galimru/idealista-mcp/internal/debug"
	"github.com/galimru/idealista-mcp/internal/session"
)

// Client manages an Idealista OAuth access token obtained via client_credentials grant.
type Client struct {
	mu         sync.Mutex
	basicAuth  string
	sess       *session.Session
	httpClient *http.Client
}

// New creates a Client using config-provided credentials. The provided session
// is used to cache the access token across process restarts.
func New(sess *session.Session, key, secret string) (*Client, error) {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(secret) == "" {
		return nil, fmt.Errorf("client key and client secret are required")
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(key) + ":" + url.QueryEscape(secret)))

	var transport http.RoundTripper = http.DefaultTransport
	if os.Getenv("IDEALISTA_DEBUG") != "" {
		transport = debug.NewTransport(transport)
	}

	return &Client{
		basicAuth:  basicAuth,
		sess:       sess,
		httpClient: &http.Client{Timeout: 30 * time.Second, Transport: transport},
	}, nil
}

// Token returns a valid Bearer access token, using the session cache when possible
// and refreshing only when the token has expired or is about to expire.
func (c *Client) Token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sess.AccessToken != "" && time.Now().Add(60*time.Second).Before(c.sess.ExpiresAt) {
		return c.sess.AccessToken, nil
	}
	return c.fetch(ctx)
}

// Invalidate clears the cached token, forcing a re-fetch on the next call.
func (c *Client) Invalidate() {
	c.mu.Lock()
	c.sess.AccessToken = ""
	c.mu.Unlock()
}

// fetch must be called with c.mu held.
func (c *Client) fetch(ctx context.Context) (string, error) {
	body := url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"write"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api.TokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Basic "+c.basicAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("app_version", config.AppVersion)
	req.Header.Set("device_identifier", c.sess.DeviceIdentifier)
	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch token: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, respBody)
	}

	var tok api.TokenResponse
	if err := json.Unmarshal(respBody, &tok); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("empty access token in response")
	}

	ttl := time.Duration(tok.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 55 * time.Minute
	}
	c.sess.AccessToken = tok.AccessToken
	c.sess.ExpiresAt = time.Now().Add(ttl)
	if err := c.sess.Save(); err != nil {
		return "", fmt.Errorf("save session: %w", err)
	}
	return c.sess.AccessToken, nil
}
